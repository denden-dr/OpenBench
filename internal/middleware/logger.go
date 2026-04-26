package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// ZapLogger logs HTTP requests and their relevant metadata using the provided zap logger
func ZapLogger(log *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Proceed to the next handler
		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.Duration("latency", latency),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			log.Error("Request failed", fields...)
		} else {
			log.Info("Request completed", fields...)
		}

		return err
	}
}
