package logger

import (
	"log"

	"go.uber.org/zap"
)

// NewLogger creates and configures a new zap Logger instance
func NewLogger() *zap.Logger {
	// Using NewProduction provides structured JSON logging by default
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return logger
}
