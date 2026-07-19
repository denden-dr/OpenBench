package logger

import (
	"context"
	"time"

	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// NewMiddleware returns a Fiber middleware that injects a Request ID into the context
// and logs HTTP requests using slog.
func NewMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Get or generate request ID
		reqID := c.Get(fiber.HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// Set request ID in HTTP Response Header
		c.Set(fiber.HeaderXRequestID, reqID)

		// Inject Request ID into Go context
		ctx := context.WithValue(c.Context(), RequestIDKey, reqID)
		c.SetContext(ctx)

		// Process request
		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		switch {
		case err != nil:
			slog.ErrorContext(ctx, "HTTP Request Error",
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("method", method),
				slog.String("path", path),
				slog.Any("error", err),
			)
		case status >= 400:
			slog.WarnContext(ctx, "HTTP Request Warning",
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("method", method),
				slog.String("path", path),
			)
		default:
			slog.InfoContext(ctx, "HTTP Request Success",
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("method", method),
				slog.String("path", path),
			)
		}

		return err
	}
}
