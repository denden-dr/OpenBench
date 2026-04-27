# Implementation Reference: PostgreSQL Connection & Health Integration

This document provides the specific code changes required to implement the database connection pooling and health-check integration as defined in `plan.md`. Details unrelated to the database connection have been omitted for brevity.

## 1. `pkg/config/config.go`

Create the full configuration package from scratch:

```go
package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL           string `envconfig:"DATABASE_URL" required:"true"`
	DBMaxOpenConns        int    `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	DBMaxIdleConns        int    `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	DBConnMaxLifetimeSecs int    `envconfig:"DB_CONN_MAX_LIFETIME_SECS" default:"300"`
	DBConnMaxIdleTimeSecs int    `envconfig:"DB_CONN_MAX_IDLE_TIME_SECS" default:"60"`
}

// LoadConfig reads environment variables from .env (if present) and populates the Config struct.
func LoadConfig() (*Config, error) {
	// Best-effort .env loading — not an error if file is missing (e.g., in containers)
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	// Validate Database Connection Pool settings
	if cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		return nil, fmt.Errorf("invalid config: DB_MAX_IDLE_CONNS (%d) cannot exceed DB_MAX_OPEN_CONNS (%d)", cfg.DBMaxIdleConns, cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxIdleConns < 1 {
		return nil, fmt.Errorf("invalid config: DB_MAX_IDLE_CONNS must be at least 1")
	}
	if cfg.DBConnMaxLifetimeSecs <= 0 {
		return nil, fmt.Errorf("invalid config: DB_CONN_MAX_LIFETIME_SECS must be greater than 0")
	}
	if cfg.DBConnMaxIdleTimeSecs <= 0 {
		return nil, fmt.Errorf("invalid config: DB_CONN_MAX_IDLE_TIME_SECS must be greater than 0")
	}

	return &cfg, nil
}
```

## 2. `pkg/database/postgres.go`

Create the database connection package from scratch:

```go
package database

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// NewPostgresDB creates a new PostgreSQL database connection with configurable pool settings
func NewPostgresDB(dsn string, maxOpenConns, maxIdleConns, connMaxLifetimeSecs, connMaxIdleTimeSecs int) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool settings dynamically
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetimeSecs) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(connMaxIdleTimeSecs) * time.Second)

	return db, nil
}
```

## 3. `internal/handlers/health.go`

Update the health handler to check the database status via `PingContext`:

```go
package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type HealthHandler struct {
	db *sqlx.DB
}

func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c fiber.Ctx) error {
	// 2-second timeout for the ping
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	status := "healthy"
	dbStatus := "up"
	dbMessage := ""

	if err := h.db.PingContext(ctx); err != nil {
		status = "degraded"
		dbStatus = "down"
		dbMessage = err.Error()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": status,
		"checks": fiber.Map{
			"database": fiber.Map{
				"status":  dbStatus,
				"message": dbMessage,
			},
		},
	})
}
```

## 4. `cmd/api/main.go`

Update the existing `main.go` to wire in config loading, database initialisation, and health handler dependency injection. This is the **complete** updated file:

```go
package main

import (
	"github.com/denden-dr/OpenBench/internal/handlers"
	"github.com/denden-dr/OpenBench/internal/middleware"
	"github.com/denden-dr/OpenBench/pkg/config"
	"github.com/denden-dr/OpenBench/pkg/database"
	"github.com/denden-dr/OpenBench/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Database with configurable pool settings
	db, err := database.NewPostgresDB(
		cfg.DatabaseURL,
		cfg.DBMaxOpenConns,
		cfg.DBMaxIdleConns,
		cfg.DBConnMaxLifetimeSecs,
		cfg.DBConnMaxIdleTimeSecs,
	)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Log applied settings for observability
	log.Info("Database connection established",
		zap.Int("max_open_conns", cfg.DBMaxOpenConns),
		zap.Int("max_idle_conns", cfg.DBMaxIdleConns),
		zap.Int("conn_max_lifetime_secs", cfg.DBConnMaxLifetimeSecs),
		zap.Int("conn_max_idle_time_secs", cfg.DBConnMaxIdleTimeSecs),
	)

	// Dependency Injection for Health
	healthHandler := handlers.NewHealthHandler(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "OpenBench API",
	})

	// Use zap-based middleware for all HTTP requests
	app.Use(middleware.ZapLogger(log))

	// Define routes
	app.Get("/health", healthHandler.HealthCheck)

	// Log server start
	log.Info("Starting OpenBench API server on port 3000")

	// Start server
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
```

## 5. `.env` and `.env.example` additions

Append these lines under your existing database configuration block:

```env
# Connection Pool Configuration
# DB_MAX_OPEN_CONNS: Maximum number of open connections to the database.
DB_MAX_OPEN_CONNS=25
# DB_MAX_IDLE_CONNS: Maximum number of connections in the idle connection pool.
DB_MAX_IDLE_CONNS=5
# DB_CONN_MAX_LIFETIME_SECS: Maximum amount of time a connection may be reused.
DB_CONN_MAX_LIFETIME_SECS=300
# DB_CONN_MAX_IDLE_TIME_SECS: Maximum amount of time a connection may be idle.
DB_CONN_MAX_IDLE_TIME_SECS=60
```
