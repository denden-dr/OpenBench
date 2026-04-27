package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type HealthHandler struct {
	db          *sqlx.DB
	pingTimeout time.Duration
}

func NewHealthHandler(db *sqlx.DB, pingTimeout time.Duration) *HealthHandler {
	return &HealthHandler{
		db:          db,
		pingTimeout: pingTimeout,
	}
}

func (h *HealthHandler) HealthCheck(c fiber.Ctx) error {
	// Configurable timeout for the ping
	ctx, cancel := context.WithTimeout(c.Context(), h.pingTimeout)
	defer cancel()

	status := "healthy"
	dbStatus := "up"
	dbMessage := ""
	httpStatus := fiber.StatusOK

	if err := h.db.PingContext(ctx); err != nil {
		status = "degraded"
		dbStatus = "down"
		dbMessage = err.Error()
		httpStatus = fiber.StatusServiceUnavailable
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"status": status,
		"checks": fiber.Map{
			"database": fiber.Map{
				"status":  dbStatus,
				"message": dbMessage,
			},
		},
	})
}
