package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

type ContextHandler struct {
	slog.Handler
}

// Handle implements slog.Handler. It extracts request_id from context and adds it to the record.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
			r.AddAttrs(slog.String("request_id", reqID))
		}
	}
	return h.Handler.Handle(ctx, r)
}

// InitLogger initializes the global structured logger.
// If env is "production", it uses JSON format with Info level.
// Otherwise, it uses colorful text format (tint) with Debug level.
func InitLogger(env string) {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})
	}

	// Wrap handler with ContextHandler to support request correlation
	logger := slog.New(&ContextHandler{Handler: handler})
	slog.SetDefault(logger)
}
