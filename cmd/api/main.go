package main

import (
	"github.com/denden-dr/OpenBench/internal/handlers"
	"github.com/denden-dr/OpenBench/internal/middleware"
	"github.com/denden-dr/OpenBench/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "OpenBench API",
	})

	// Use zap-based middleware for all HTTP requests
	app.Use(middleware.ZapLogger(log))

	// Define routes
	app.Get("/health", handlers.HealthCheck)

	// Log server start
	log.Info("Starting OpenBench API server on port 3000")

	// Start server
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
