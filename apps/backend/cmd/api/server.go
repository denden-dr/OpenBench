package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func newFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.AppName,
		ErrorHandler: globalErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.NewMiddleware())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.App.AllowedOrigins, ","),
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	}))

	return app
}

func listenAndServe(app *fiber.App, port string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Shutting down API server...")
		if err := app.Shutdown(); err != nil {
			slog.Error("error during server shutdown", "error", err)
		}
	}()

	slog.Info("Starting API server", "port", port)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("server exited", "error", err)
	}

	slog.Info("Server stopped")
}
