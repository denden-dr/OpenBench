package health

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

// HealthHandler holds the database pool to perform health checks.
type HealthHandler struct {
	db *sqlx.DB
}

// NewHealthHandler initializes the HealthHandler.
func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheckPublic verifies that the server is running without leaking database status.
func (h *HealthHandler) HealthCheckPublic(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "up",
	})
}

// HealthCheckDetail verifies the server and database status in detail.
func (h *HealthHandler) HealthCheckDetail(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbStatus := "up"
	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "down"
		slog.ErrorContext(ctx, "Health check: database is unreachable", slog.Any("error", err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "up",
		"database":  dbStatus,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
