package health

import (
	"context"
	"log"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/gofiber/fiber/v2"
)

// Handler holds dependencies for health probes
type Handler struct {
	db  *database.Database
	env string
}

// NewHandler creates a new health Handler instance
func NewHandler(db *database.Database, env string) *Handler {
	return &Handler{
		db:  db,
		env: env,
	}
}

// Liveness returns the liveness status of the API service
func (h *Handler) Liveness(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":      "OK",
		"environment": h.env,
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}

// Readiness verifies that the database is reachable and active
func (h *Handler) Readiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	start := time.Now()
	err := h.db.DB.PingContext(ctx)
	duration := time.Since(start)

	// Slow database ping check
	if duration > 100*time.Millisecond {
		log.Printf("Warning: Slow database ping took %v (threshold: 100ms)", duration)
	}

	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "UNREADY",
			"error":  "database is not reachable",
		})
	}

	// Retrieve database connection pool stats
	stats := h.db.Stats()

	return c.JSON(fiber.Map{
		"status":  "READY",
		"db":      "CONNECTED",
		"latency": duration.String(),
		"pool": fiber.Map{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
			"wait_count":       stats.WaitCount,
			"wait_duration_ms": stats.WaitDuration.Milliseconds(),
		},
	})
}

// LegacyHealth redirects legacy health requests to readiness
func (h *Handler) LegacyHealth(c *fiber.Ctx) error {
	return c.Redirect("/health/readiness", fiber.StatusMovedPermanently)
}
