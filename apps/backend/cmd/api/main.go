package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/events"
	"github.com/denden-dr/OpenBench/apps/backend/internal/health"
	"github.com/denden-dr/OpenBench/apps/backend/internal/inventory"
	"github.com/denden-dr/OpenBench/apps/backend/internal/logger"
	"github.com/denden-dr/OpenBench/apps/backend/internal/pos"
	"github.com/denden-dr/OpenBench/apps/backend/internal/ticket"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/denden-dr/OpenBench/apps/backend/internal/warranty"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// 1. Load application configurations
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize Logger
	logger.InitLogger(cfg.App.Env)

	// 2. Initialize database connection pool
	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Database connection pool established successfully")

	// Initialize Event Bus
	eventBus := events.NewAsyncEventBus(100)
	defer eventBus.Close()

	// TxManager
	txManager := database.NewTxManager(db)

	// Wire Auth Layers
	authQueryRepo := auth.NewQueryRepository(db)
	authCommandRepo := auth.NewCommandRepository(db)
	authService := auth.NewService(authQueryRepo, authCommandRepo, cfg)
	authHandler := auth.NewHandler(authService, cfg)

	// Token Cleanup Worker
	authCleanupWorker := auth.NewCleanupWorker(authCommandRepo, 24*time.Hour)
	authCleanupWorker.Start()
	defer authCleanupWorker.Stop()

	// Wire Warranty & Claims Layers
	warrantyQueryRepo := warranty.NewQueryRepository(db)
	warrantyCommandRepo := warranty.NewCommandRepository(db)
	warrantyService := warranty.NewService(warrantyQueryRepo, warrantyCommandRepo, txManager)
	warrantyHandler := warranty.NewHandler(warrantyService)

	// Wire Ticket Layers
	ticketQueryRepo := ticket.NewQueryRepository(db)
	ticketCommandRepo := ticket.NewCommandRepository(db)
	ticketService := ticket.NewService(ticketQueryRepo, ticketCommandRepo, txManager, warrantyService, eventBus, cfg.App.EncryptionKey)
	ticketHandler := ticket.NewHandler(ticketService)

	// Wire Inventory Layers
	inventoryQueryRepo := inventory.NewQueryRepository(db)
	inventoryCommandRepo := inventory.NewCommandRepository(db)
	inventoryService := inventory.NewService(inventoryQueryRepo, inventoryCommandRepo)
	inventoryHandler := inventory.NewHandler(inventoryService)

	// Wire POS Layers
	posQueryRepo := pos.NewQueryRepository(db)
	posCommandRepo := pos.NewCommandRepository(db)
	posService := pos.NewService(posQueryRepo, posCommandRepo, inventoryQueryRepo, inventoryCommandRepo, txManager)
	posHandler := pos.NewHandler(posService)

	// 3. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.AppName,
		ErrorHandler: globalErrorHandler,
	})

	// Register recover middleware
	app.Use(recover.New())

	// Register structured logging middleware
	app.Use(logger.NewMiddleware())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.App.AllowedOrigins, ","),
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	}))

	// 4. Register handlers
	healthHandler := health.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheckPublic)

	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return utils.SendProblem(c, fiber.StatusTooManyRequests, "/errors/too-many-requests", "Too Many Requests", "Terlalu banyak percobaan masuk. Silakan coba lagi dalam 1 menit.")
		},
	})

	// Auth Public Routes
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/login", authLimiter, authHandler.Login)
	authGroup.Post("/refresh", authLimiter, authHandler.Refresh)
	authGroup.Post("/logout", authHandler.Logout)

	// Protected Admin Routes
	adminGroup := app.Group("/api/v1/admin", auth.RequireAuth(cfg, authQueryRepo), auth.RequireRole("ADMIN"))
	adminGroup.Get("/health/detail", healthHandler.HealthCheckDetail)
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

	// Inventory Routes
	invGroup := adminGroup.Group("/products")
	invGroup.Post("/", inventoryHandler.CreateProduct)
	invGroup.Get("/", inventoryHandler.GetProducts)
	invGroup.Get("/:id", inventoryHandler.GetProductByID)
	invGroup.Put("/:id", inventoryHandler.UpdateProduct)
	invGroup.Patch("/:id/stock", inventoryHandler.AdjustStock)
	invGroup.Delete("/:id", inventoryHandler.DeleteProduct)

	// POS Routes
	posGroup := adminGroup.Group("/pos")
	posGroup.Post("/checkout", posHandler.Checkout)
	posGroup.Get("/transactions", posHandler.GetTransactions)
	posGroup.Get("/transactions/:id", posHandler.GetTransactionByID)

	// 5. Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Shutting down API server...")
		if err := app.Shutdown(); err != nil {
			slog.Error("error during server shutdown", "error", err)
		}
	}()

	// 6. Start server
	slog.Info("Starting API server", "port", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		slog.Error("server exited", "error", err)
	}

	slog.Info("Server stopped")
}
