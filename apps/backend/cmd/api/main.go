package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Database Connection
	db, err := database.NewConnection(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		AppName:      "OpenBench API",
	})

	// Apply Middlewares
	app.Use(recover.New()) // Capture panics
	app.Use(logger.New())  // Standard request logging

	// Secure CORS settings (TD-002)
	allowOrigins := strings.Join(cfg.AllowedOrigins, ",")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:5173"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Security Headers (Helmet)
	app.Use(helmet.New())

	// Liveness Probe (TD-008)
	app.Get("/health/liveness", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "OK",
			"environment": cfg.Env,
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	})

	// Readiness Probe (TD-008)
	app.Get("/health/readiness", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		start := time.Now()
		err := db.DB.PingContext(ctx)
		duration := time.Since(start)

		// Slow operation warning (TD-009)
		if duration > 100*time.Millisecond {
			log.Printf("Warning: Slow database ping took %v (threshold: 100ms)", duration)
		}

		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "UNREADY",
				"error":  "database is not reachable",
			})
		}

		// Retrieve database stats for observability (TD-009)
		stats := db.Stats()

		return c.JSON(fiber.Map{
			"status":  "READY",
			"db":      "CONNECTED",
			"latency": duration.String(),
			"pool": fiber.Map{
				"open_connections": stats.OpenConnections,
				"in_use":           stats.InUse,
				"idle":             stats.Idle,
				"wait_count":       stats.WaitCount,
				"wait_duration_ms": stats.WaitDuration.Milliseconds(),
			},
		})
	})

	// Legacy Health route redirects to readiness
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Redirect("/health/readiness", fiber.StatusMovedPermanently)
	})

	// Periodic DB Connection Pool Stats Logging (TD-009)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			stats := db.Stats()
			log.Printf("[DB Stats] OpenConnections=%d, InUse=%d, Idle=%d, WaitCount=%d, WaitDuration=%v, MaxIdleClosed=%d, MaxLifetimeClosed=%d",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration, stats.MaxIdleClosed, stats.MaxLifetimeClosed)
		}
	}()

	// Channel to capture server startup errors
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s in %s mode...", cfg.Port, cfg.Env)
		serverErrors <- app.Listen(":" + cfg.Port)
	}()

	// Graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until signal or startup error
	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)
	case <-quit:
		log.Println("Shutting down server gracefully...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}

	log.Println("Server exited cleanly")
}
