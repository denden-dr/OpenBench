package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/health"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// 1. Load application configurations
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Initialize database connection pool
	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection pool established successfully")

	// Wire Auth Layers
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService, cfg)

	// Run Seeder if APP_ENV == development
	ctxSeed, cancelSeed := context.WithTimeout(context.Background(), 10*time.Second)
	if err := database.SeedDefaultAdmin(ctxSeed, authRepo, cfg); err != nil {
		log.Printf("Warning: Failed to seed default admin: %v", err)
	}
	cancelSeed()

	// 3. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: cfg.App.AppName,
	})

	// 4. Register handlers
	healthHandler := health.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)

	// Auth Public Routes
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.Refresh)
	authGroup.Post("/logout", authHandler.Logout)

	// Protected Admin Routes
	adminGroup := app.Group("/api/v1/admin", auth.RequireAuth(cfg))
	adminGroup.Get("/profile", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": fiber.Map{
				"user_id": c.Locals("userID"),
				"role":    c.Locals("userRole"),
			},
		})
	})

	// 5. Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down API server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("error during server shutdown: %v", err)
		}
	}()

	// 6. Start server
	log.Printf("Starting API server on port %s...", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Printf("server exited: %v", err)
	}

	log.Println("Server stopped")
}
