package handler

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	if err := h.db.PingContext(c.Context()); err != nil {
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
}
