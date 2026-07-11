package health

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler holds the database pool to perform health checks.
type HealthHandler struct {
	db *pgxpool.Pool
}

// NewHealthHandler initializes the HealthHandler.
func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
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
	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "down"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "up",
		"database":  dbStatus,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
