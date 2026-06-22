package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/auth"
	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/health"
	"github.com/denden-dr/openbench/apps/backend/internal/ticket"
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

	// Seed Development/Testing Users
	if cfg.Env == "development" || cfg.Env == "test" {
		if err := database.SeedDevUsers(db); err != nil {
			log.Fatalf("Database seeding failed: %v", err)
		}
	}

	// Initialize Auth Repository and Service
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, db, cfg.JWTSecret)

	// Initialize Ticket Repository and Services
	ticketRepo := ticket.NewRepository(db)
	publicTicketService := ticket.NewService(ticketRepo, db)
	adminTicketService := ticket.NewAdminService(ticketRepo, db)
	ticketHandler := ticket.NewHandler(adminTicketService, publicTicketService)

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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Security Headers (Helmet)
	app.Use(helmet.New())

	// Health Probes (TD-008)
	healthHandler := health.NewHandler(db, cfg.Env)
	app.Get("/health/liveness", healthHandler.Liveness)
	app.Get("/health/readiness", healthHandler.Readiness)
	app.Get("/health/metrics", healthHandler.Metrics)
	app.Get("/health", healthHandler.LegacyHealth)

	// Auth Routes
	api := app.Group("/api/v1")
	isDev := cfg.Env == "development" || cfg.Env == "test"
	authHandler := auth.NewHandler(authService, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry, isDev)
	api.Post("/auth/signin", authHandler.SignIn)
	api.Post("/auth/refresh", authHandler.Refresh)
	api.Post("/auth/signout", authHandler.SignOut)
	api.Get("/auth/me", auth.RequireAuth(cfg.JWTSecret), authHandler.Me)

	// Public Ticket Tracker Route
	api.Get("/tracker/:id", ticketHandler.GetPublicTrackerTicket)

	// Protected Admin Routes
	admin := api.Group("/admin", auth.RequireAuth(cfg.JWTSecret), auth.RequireRole("admin"))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		return c.JSON(fiber.Map{
			"message": "Welcome to the admin dashboard",
			"user_id": userID,
		})
	})
	admin.Get("/tickets", ticketHandler.ListTickets)
	admin.Post("/tickets", ticketHandler.CreateTicket)
	admin.Get("/tickets/:id", ticketHandler.GetTicket)
	admin.Patch("/tickets/:id", ticketHandler.UpdateTicket)
	admin.Post("/tickets/:id/emergency", ticketHandler.EmergencyUpdateTicket)
	admin.Get("/warranties", ticketHandler.ListWarranties)

	// Periodic expired refresh token cleanup (TD-003)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := authRepo.PurgeExpiredTokens(context.Background()); err != nil {
				log.Printf("[Cleanup] Failed to purge expired tokens: %v", err)
			} else {
				log.Println("[Cleanup] Expired/revoked refresh tokens purged successfully")
			}
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
