package middleware

import (
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// NewRateLimiters creates public and admin rate limiters based on the config.
func NewRateLimiters(cfg config.RateLimitConfig) (fiber.Handler, fiber.Handler) {
	publicLimiter := limiter.New(limiter.Config{
		Max:        cfg.MaxPublic,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Next: func(c *fiber.Ctx) bool {
			return cfg.Disable
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many requests, please try again later",
			})
		},
	})

	adminLimiter := limiter.New(limiter.Config{
		Max:        cfg.MaxAdmin,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Next: func(c *fiber.Ctx) bool {
			return cfg.Disable
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many admin requests, please try again later",
			})
		},
	})

	return publicLimiter, adminLimiter
}
