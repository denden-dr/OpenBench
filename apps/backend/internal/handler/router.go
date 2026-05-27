package handler

import (
	"database/sql"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes configures all routes for the Fiber application.
func RegisterRoutes(
	app *fiber.App,
	db *sql.DB,
	cfg *config.Config,
	ticketHandler *TicketHandler,
	warrantyClaimHandler *WarrantyClaimHandler,
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

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := db.PingContext(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "healthy",
		})
	})
}
