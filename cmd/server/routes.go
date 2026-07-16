package main

import (
	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/auth"
	"github.com/denden-dr/OpenBench/internal/health"
	"github.com/denden-dr/OpenBench/internal/inventory"
	"github.com/denden-dr/OpenBench/internal/pos"
	"github.com/denden-dr/OpenBench/internal/ticket"
	"github.com/denden-dr/OpenBench/internal/warranty"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func registerRoutes(
	app *fiber.App,
	cfg *config.Config,
	db *sqlx.DB,
	authMod auth.Module,
	warrantyMod warranty.Module,
	ticketMod ticket.Module,
	inventoryMod inventory.Module,
	posMod pos.Module,
) {
	healthHandler := health.NewHealthHandler(db)

	// Global Public Health Route
	app.Get("/health", healthHandler.HealthCheckPublic)

	// Register Web UI Routes
	registerWebRoutes(app, cfg, authMod, ticketMod, posMod, warrantyMod)

	// Register JSON API Routes
	registerAPIRoutes(app, cfg, healthHandler, authMod, warrantyMod, ticketMod, inventoryMod, posMod)
}
