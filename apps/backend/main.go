package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/denden-dr/openbench/apps/backend/internal/config"
)

func main() {
    cfg := config.Load()

    app := fiber.New()

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
