package handlers

import (
	"github.com/gofiber/fiber/v3"
)

// HealthCheck returns a 200 OK status to indicate the server is running
func HealthCheck(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "OpenBench API is healthy",
	})
}
