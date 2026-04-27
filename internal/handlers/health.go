package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type HealthHandler struct {
	db *sqlx.DB
}

func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c fiber.Ctx) error {
	// 2-second timeout for the ping
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	status := "healthy"
	dbStatus := "up"
	dbMessage := ""

	if err := h.db.PingContext(ctx); err != nil {
		status = "degraded"
		dbStatus = "down"
		dbMessage = err.Error()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": status,
		"checks": fiber.Map{
			"database": fiber.Map{
				"status":  dbStatus,
				"message": dbMessage,
			},
		},
	})
}
