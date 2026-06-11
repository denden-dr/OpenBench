---
name: initializing-go-fiber-api
description: Use when initializing a Go web application using the Fiber framework, configuring security middleware (CORS, Recover, Helmet), and implementing graceful shutdown.
---

# Initializing Go Fiber API

## Overview
Go Fiber is an Express-inspired web framework for Go. When setting up a Fiber API, security middleware, fail-fast server startup, and graceful shutdown must be implemented from the start.

## When to Use
- Initializing a new Go web server using Fiber.
- Hardening an existing Fiber application with security headers, recovery logic, or CORS.
- Adding graceful shutdown handling to a Go service.

## Core Pattern
Always capture the server startup errors on a channel and block the main execution using a `select` statement. This ensures the application crashes immediately if the port bind fails (Fail-Fast), rather than running silently in a hang state.

### Graceful Shutdown and Middleware Pattern
```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Security Middleware
	app.Use(recover.New()) // Capture panics
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // Restrict origins
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Channel to capture server startup errors (BE-001)
	serverErrors := make(chan error, 1)

	// Listen in goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		serverErrors <- app.Listen(":" + port)
	}()

	// Graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until a signal or startup error occurs (Fail-Fast)
	select {
	case err := <-serverErrors:
		log.Fatalf("Server startup failed: %v", err)
	case <-quit:
		log.Println("Shutting down server gracefully...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		log.Println("Server exited cleanly")
	}
}
```

## Common Mistakes
- **Hanging Startup**: Running `app.Listen` inside a goroutine without returning its errors, which leaves the main process alive but unable to serve traffic when port bind fails (BE-001).
- **No Recover Middleware**: Not using the `recover` middleware. Any unhandled panic in a route handler will crash the entire server.
- **Permissive CORS**: Using `AllowOrigins: "*"` in production. Be restrictive and define exact domain origins.
- **No Read/Write Timeouts**: Leaving timeouts undefined in `fiber.Config`. This leaves the server vulnerable to slowloris attacks.
- **Abrupt Termination**: Killing the process with `os.Exit` or allowing the main function to end immediately without calling `app.Shutdown()`, which terminates in-flight client requests mid-transaction.
