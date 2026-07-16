package main

import (
	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/auth"
	"github.com/denden-dr/OpenBench/internal/pos"
	"github.com/denden-dr/OpenBench/internal/ticket"
	"github.com/denden-dr/OpenBench/internal/utils"
	"github.com/denden-dr/OpenBench/internal/warranty"
	admin_pages "github.com/denden-dr/OpenBench/ui/views/pages/admin"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func registerWebRoutes(app *fiber.App, cfg *config.Config, authMod auth.Module, ticketMod ticket.Module, posMod pos.Module, warrantyMod warranty.Module) {
	// Static Files
	app.Get("/static/*", static.New("./ui/static"))

	// Dashboard (Protected Web Route)
	webAuth := auth.RequireWebAuth(cfg, authMod.QueryRepo)
	app.Get("/", webAuth, func(c fiber.Ctx) error {
		return utils.Render(c, admin_pages.DashboardPage())
	})

	// Ticket Routes (Protected)
	ticketsGroup := app.Group("/tickets", webAuth)
	ticketsGroup.Get("/", ticketMod.WebHandler.TicketsPage)
	ticketsGroup.Get("/new", ticketMod.WebHandler.NewTicketPage)
	ticketsGroup.Get("/:id", ticketMod.WebHandler.TicketDetailPage)

	// POS & Inventory Routes (Protected)
	posGroup := app.Group("/pos", webAuth)
	posGroup.Get("/", posMod.WebHandler.CheckoutPage)
	posGroup.Get("/inventory", posMod.WebHandler.InventoryPage)
	posGroup.Get("/inventory/new", posMod.WebHandler.NewProductPage)
	posGroup.Get("/history", posMod.WebHandler.HistoryPage)

	// Warranty Routes (Protected)
	warrantiesGroup := app.Group("/warranties", webAuth)
	warrantiesGroup.Get("/", warrantyMod.WebHandler.WarrantiesPage)
	warrantiesGroup.Get("/claims/new", warrantyMod.WebHandler.NewClaimPage)
	warrantiesGroup.Get("/claims/:id", warrantyMod.WebHandler.ClaimDetailPage)

	// Web Auth Routes
	app.Get("/login", authMod.WebHandler.LoginPage)
	app.Post("/login", authMod.WebHandler.LoginPost)
}
