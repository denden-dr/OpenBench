package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/handler"
	"github.com/denden-dr/openbench/apps/backend/internal/middleware"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	// Initialize Database
	db, err := database.NewDB(cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize Layers (Manual DI)
	ticketRepo := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	warrantyClaimRepo := repository.NewWarrantyClaimRepository(db)
	warrantyClaimService := service.NewWarrantyClaimService(warrantyClaimRepo, ticketRepo)
	warrantyClaimHandler := handler.NewWarrantyClaimHandler(warrantyClaimService)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowMethods: "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,Idempotency-Key",
	}))

	// Idempotency Middleware
	// Fiber's default idempotency lock is in-memory. This is acceptable while
	// OpenBench runs as a single backend instance. Use a DB-backed lock before
	// moving to multi-instance or rolling deployments.
	idempotencyStore := database.NewPostgresStorage(db)
	defer idempotencyStore.Close()
	app.Use(middleware.ScopeIdempotencyKey(idempotencyStore))
	app.Use(middleware.NewIdempotency(idempotencyStore))

	// Routes
	api := app.Group("/api/v1")

	tickets := api.Group("/tickets")
	tickets.Post("/", ticketHandler.Create)
	tickets.Get("/", ticketHandler.List)
	tickets.Get("/:id", ticketHandler.GetByID)
	tickets.Patch("/:id", ticketHandler.Update)
	tickets.Delete("/:id", ticketHandler.Delete)

	warrantyClaims := api.Group("/warranty-claims")
	warrantyClaims.Post("/", warrantyClaimHandler.Create)
	warrantyClaims.Get("/", warrantyClaimHandler.List)
	warrantyClaims.Post("/:id/approve", warrantyClaimHandler.Approve)
	warrantyClaims.Post("/:id/void", warrantyClaimHandler.Void)

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := db.PingContext(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"error":   "Database connection lost",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"status":  "ok",
				"message": "Hello from OpenBench Backend!",
			},
		})
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		slog.Info("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
	}()

	slog.Info("Server starting", "port", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil && ctx.Err() == nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
