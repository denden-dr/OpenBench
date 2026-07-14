package main

import (
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/health"
	"github.com/denden-dr/OpenBench/apps/backend/internal/inventory"
	"github.com/denden-dr/OpenBench/apps/backend/internal/pos"
	"github.com/denden-dr/OpenBench/apps/backend/internal/ticket"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/denden-dr/OpenBench/apps/backend/internal/warranty"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
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
	authGroup.Post("/login", authLimiter, authMod.Handler.Login)
	authGroup.Post("/refresh", authLimiter, authMod.Handler.Refresh)
	authGroup.Post("/logout", authMod.Handler.Logout)

	// Protected Admin Routes
	adminGroup := app.Group("/api/v1/admin", auth.RequireAuth(cfg, authMod.QueryRepo), auth.RequireRole("ADMIN"))
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
	ticketGroup.Post("/", ticketMod.Handler.CreateTicket)
	ticketGroup.Get("/", ticketMod.Handler.GetTickets)
	ticketGroup.Add([]string{"QUERY"}, "/search", ticketMod.Handler.SearchTickets)
	ticketGroup.Get("/:ticket_id", ticketMod.Handler.GetTicketByID)
	ticketGroup.Patch("/:ticket_id/status", ticketMod.Handler.UpdateTicketStatus)
	ticketGroup.Put("/:ticket_id", ticketMod.Handler.UpdateTicketDetails)
	ticketGroup.Put("/:ticket_id/emergency", ticketMod.Handler.EmergencyUpdateTicket)

	// Warranty Routes
	warrGroup := adminGroup.Group("/warranties")
	warrGroup.Get("/by-ticket/:ticket_id", warrantyMod.Handler.GetWarrantyByTicketID)
	warrGroup.Patch("/:warranty_id/status", warrantyMod.Handler.UpdateWarrantyStatus)

	// Claim Routes
	claimGroup := adminGroup.Group("/claims")
	claimGroup.Post("/", warrantyMod.Handler.CreateClaim)
	claimGroup.Get("/", warrantyMod.Handler.GetClaims)
	claimGroup.Get("/:claim_id", warrantyMod.Handler.GetClaimByID)
	claimGroup.Patch("/:claim_id/status", warrantyMod.Handler.UpdateClaimStatus)
	claimGroup.Put("/:claim_id", warrantyMod.Handler.UpdateClaim)
	claimGroup.Post("/:claim_id/evaluate", warrantyMod.Handler.EvaluateClaim)

	// Inventory Routes
	invGroup := adminGroup.Group("/products")
	invGroup.Post("/", inventoryMod.Handler.CreateProduct)
	invGroup.Get("/", inventoryMod.Handler.GetProducts)
	invGroup.Get("/:id", inventoryMod.Handler.GetProductByID)
	invGroup.Put("/:id", inventoryMod.Handler.UpdateProduct)
	invGroup.Patch("/:id/stock", inventoryMod.Handler.AdjustStock)
	invGroup.Delete("/:id", inventoryMod.Handler.DeleteProduct)

	// POS Routes
	posGroup := adminGroup.Group("/pos")
	posGroup.Post("/checkout", posMod.Handler.Checkout)
	posGroup.Get("/transactions", posMod.Handler.GetTransactions)
	posGroup.Get("/transactions/:id", posMod.Handler.GetTransactionByID)
}
