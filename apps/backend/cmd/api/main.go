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
	"github.com/denden-dr/OpenBench/apps/backend/internal/events"
	"github.com/denden-dr/OpenBench/apps/backend/internal/health"
	"github.com/denden-dr/OpenBench/apps/backend/internal/ticket"
	"github.com/denden-dr/OpenBench/apps/backend/internal/warranty"

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

	// Initialize Event Bus
	eventBus := events.NewSyncEventBus()

	// Wire Auth Layers
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService, cfg)

	// Wire Ticket Layers
	ticketRepo := ticket.NewRepository(db)
	ticketService := ticket.NewService(ticketRepo, eventBus)
	ticketHandler := ticket.NewHandler(ticketService)

	// Wire Warranty & Claims Layers
	warrantyRepo := warranty.NewRepository(db)
	warrantyService := warranty.NewService(warrantyRepo)
	warrantyHandler := warranty.NewHandler(warrantyService)
	warrantyEventHandler := warranty.NewEventHandler(warrantyService)

	// Register Domain Event Subscribers
	eventBus.Subscribe(events.TicketCompletedType, warrantyEventHandler.HandleTicketCompleted)

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

	// Ticket Routes
	ticketGroup := adminGroup.Group("/services")
	ticketGroup.Post("/", ticketHandler.CreateTicket)
	ticketGroup.Get("/", ticketHandler.GetTickets)
	ticketGroup.Add([]string{"QUERY"}, "/search", ticketHandler.SearchTickets)
	ticketGroup.Get("/:ticket_id", ticketHandler.GetTicketByID)
	ticketGroup.Patch("/:ticket_id/status", ticketHandler.UpdateTicketStatus)
	ticketGroup.Put("/:ticket_id", ticketHandler.UpdateTicketDetails)
	ticketGroup.Put("/:ticket_id/emergency", ticketHandler.EmergencyUpdateTicket)

	// Warranty Routes
	warrGroup := adminGroup.Group("/warranties")
	warrGroup.Get("/by-ticket/:ticket_id", warrantyHandler.GetWarrantyByTicketID)
	warrGroup.Patch("/:warranty_id/status", warrantyHandler.UpdateWarrantyStatus)

	// Claim Routes
	claimGroup := adminGroup.Group("/claims")
	claimGroup.Post("/", warrantyHandler.CreateClaim)
	claimGroup.Get("/", warrantyHandler.GetClaims)
	claimGroup.Get("/:claim_id", warrantyHandler.GetClaimByID)
	claimGroup.Patch("/:claim_id/status", warrantyHandler.UpdateClaimStatus)
	claimGroup.Put("/:claim_id", warrantyHandler.UpdateClaim)
	claimGroup.Post("/:claim_id/evaluate", warrantyHandler.EvaluateClaim)

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
