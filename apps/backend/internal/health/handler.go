package health

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/gofiber/fiber/v2"
)

// HealthHandler holds dependencies for health probes
type HealthHandler struct {
	db  *database.Database
	env string
}

// NewHandler creates a new health Handler instance
func NewHandler(db *database.Database, env string) *HealthHandler {
	return &HealthHandler{
		db:  db,
		env: env,
	}
}

// Liveness returns the liveness status of the API service
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":      "OK",
		"environment": h.env,
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}

// Readiness verifies that the database is reachable and active
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
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
func (h *HealthHandler) LegacyHealth(c *fiber.Ctx) error {
	return c.Redirect("/health/readiness", fiber.StatusMovedPermanently)
}

// Metrics returns database connection pool stats in Prometheus format
func (h *HealthHandler) Metrics(c *fiber.Ctx) error {
	stats := h.db.Stats()

	c.Set(fiber.HeaderContentType, "text/plain; version=0.0.4; charset=utf-8")

	metrics := []struct {
		name string
		help string
		val  interface{}
	}{
		{"go_db_open_connections", "Number of established connections both in use and idle.", stats.OpenConnections},
		{"go_db_in_use_connections", "Number of connections currently in use.", stats.InUse},
		{"go_db_idle_connections", "Number of idle connections.", stats.Idle},
		{"go_db_wait_count", "Total number of connections waited for.", stats.WaitCount},
		{"go_db_wait_duration_seconds", "Total time blocked waiting for a new connection.", stats.WaitDuration.Seconds()},
		{"go_db_max_idle_closed", "Total number of connections closed due to SetMaxIdleConns.", stats.MaxIdleClosed},
		{"go_db_max_lifetime_closed", "Total number of connections closed due to SetConnMaxLifetime.", stats.MaxLifetimeClosed},
	}

	var sb strings.Builder
	for _, m := range metrics {
		sb.WriteString(fmt.Sprintf("# HELP %s %s\n", m.name, m.help))
		sb.WriteString(fmt.Sprintf("# TYPE %s gauge\n", m.name))
		sb.WriteString(fmt.Sprintf("%s %v\n\n", m.name, m.val))
	}

	return c.SendString(sb.String())
}
