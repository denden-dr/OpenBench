package main

import (
    "log"

    "github.com/denden-dr/openbench/apps/backend/internal/config"
    "github.com/denden-dr/openbench/apps/backend/internal/handler"
    "github.com/denden-dr/openbench/apps/backend/internal/repository"
    "github.com/denden-dr/openbench/apps/backend/internal/service"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    cfg := config.Load()

    // Initialize Database
    db, err := repository.NewDB(cfg.DBURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize Layers (Manual DI)
    ticketRepo := repository.NewTicketRepository(db)
    ticketService := service.NewTicketService(ticketRepo)
    ticketHandler := handler.NewTicketHandler(ticketService)

    app := fiber.New()

    // Middleware
    app.Use(logger.New())
    app.Use(recover.New())

    // Routes
    api := app.Group("/api/v1")

    tickets := api.Group("/tickets")
    tickets.Post("/", ticketHandler.Create)
    tickets.Get("/:id", ticketHandler.GetByID)

    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "success": true,
            "data": fiber.Map{
                "status":  "ok",
                "message": "Hello from OpenBench Backend!",
            },
        })
    })

    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(app.Listen(":" + cfg.Port))
}
