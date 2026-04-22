package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"openbench/server/internal/config"
	"openbench/server/internal/handler"
	"openbench/server/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "OpenBench API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.CORSOrigins},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	// Initialize Services
	healthService := service.NewHealthService("0.1.0")

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler(healthService)

	// Routes
	api := app.Group("/api")
	api.Get("/health", healthHandler.GetHealth)

	// Graceful shutdown setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
