package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/handlers"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// 1. Load application configurations
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Initialize database connection pool
	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection pool established successfully")

	// 3. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: cfg.App.AppName,
	})

	// 4. Register handlers
	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)

	// 5. Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down API server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("error during server shutdown: %v", err)
		}
	}()

	// 6. Start server
	log.Printf("Starting API server on port %s...", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Printf("server exited: %v", err)
	}

	log.Println("Server stopped")
}
