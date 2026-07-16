package utils

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

// Render is a helper to render templ.Component to a fiber.Ctx
func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
