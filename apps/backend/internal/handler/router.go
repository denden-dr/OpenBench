package handler

import (
	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes configures all routes for the Fiber application.
func RegisterRoutes(
	app *fiber.App,
	cfg *config.Config,
	ticketHandler *TicketHandler,
	warrantyClaimHandler *WarrantyClaimHandler,
	healthHandler *HealthHandler,
) {
	api := app.Group("/api/v1")

	// Rate Limiting Middlewares
	publicLimiter, adminLimiter := middleware.NewRateLimiters(cfg.RateLimit)

	public := api.Group("/public", publicLimiter)
	public.Get("/tickets/:id", ticketHandler.GetPublicByID)
	public.Post("/track", ticketHandler.TrackPublic)

	tickets := api.Group("/tickets", adminLimiter)
	tickets.Post("/", ticketHandler.Create)
	tickets.Get("/", ticketHandler.List)
	tickets.Get("/:id", ticketHandler.GetByID)
	tickets.Patch("/:id", ticketHandler.Update)
	tickets.Delete("/:id", ticketHandler.Delete)

	warrantyClaims := api.Group("/warranty-claims", adminLimiter)
	warrantyClaims.Post("/", warrantyClaimHandler.Create)
	warrantyClaims.Get("/", warrantyClaimHandler.List)
	warrantyClaims.Post("/:id/approve", warrantyClaimHandler.Approve)
	warrantyClaims.Post("/:id/void", warrantyClaimHandler.Void)

	app.Get("/health", healthHandler.Check)
}
