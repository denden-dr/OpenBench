package main

import (
	"github.com/denden-dr/OpenBench/internal/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func registerWebRoutes(app *fiber.App, authMod auth.Module) {
	// Static Files
	app.Get("/static/*", static.New("./ui/static"))

	// Web UI Home
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("OpenBench Web UI (HTMX ready)")
	})

	// Web Auth Routes
	app.Get("/login", authMod.WebHandler.LoginPage)
	app.Post("/login", authMod.WebHandler.LoginPost)
}
