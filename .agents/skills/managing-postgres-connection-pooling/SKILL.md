---
name: managing-postgres-connection-pooling
description: Use when setting up PostgreSQL connection management in Go, configuring connection pool parameters, verifying connection health, and using dependency injection.
---

# Managing Postgres Connection Pooling

## Overview
Proper PostgreSQL connection management in Go using a connection pool prevents database exhaustion and resource leaks. Avoid hardcoding connection pool limits, connection lifetime parameters, or retry settings. Instead, load them from environment variables with safe defaults. Expose pool observability metrics to diagnose starvation or latency regressions.

## When to Use
- Setting up a database connection package in Go.
- Tuning database pool parameters (maximum open connections, idle connections, lifespans).
- Implementing health check/ping mechanics and retry logic during application startup.
- Integrating database pool stats and observability.

## Core Pattern
Always configure connection pool limits and retry metrics dynamically on the `*sqlx.DB` instance. Create a fresh context for each ping attempt inside the retry loop. Expose a method to retrieve pool statistics for logging or telemetry.

### Database Pool & Observability Example
```go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	RetryCount      int
	RetryDelay      time.Duration
	PingTimeout     time.Duration
}

// DSN returns the connection string for pgx driver
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

type Database struct {
	DB *sqlx.DB
}

func NewConnection(cfg *Config) (*Database, error) {
	dsn := cfg.DSN()

	// Apply safe defaults if not provided (TD-004)
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 5
	}
	if cfg.ConnMaxLifetime <= 0 {
		cfg.ConnMaxLifetime = 30 * time.Minute
	}
	if cfg.ConnMaxIdleTime <= 0 {
		cfg.ConnMaxIdleTime = 15 * time.Minute
	}
	if cfg.RetryCount <= 0 {
		cfg.RetryCount = 5
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = 2 * time.Second
	}
	if cfg.PingTimeout <= 0 {
		cfg.PingTimeout = 2 * time.Second
	}

	var db *sqlx.DB
	var err error

	// Retry connection logic
	for i := 0; i < cfg.RetryCount; i++ {
		db, err = sqlx.Open("pgx", dsn)
		if err == nil {
			// Apply Pool Settings
			db.SetMaxOpenConns(cfg.MaxOpenConns)
			db.SetMaxIdleConns(cfg.MaxIdleConns)
			db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
			db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

			// Ping database using fresh context timeout (BE-002)
			pingCtx, pingCancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
			err = db.PingContext(pingCtx)
			pingCancel()

			if err == nil {
				log.Printf("Successfully connected to the database with sqlx (Pool: OpenConns=%d, IdleConns=%d)", cfg.MaxOpenConns, cfg.MaxIdleConns)
				return &Database{DB: db}, nil
			}
			db.Close()
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...", i+1, cfg.RetryCount, err, cfg.RetryDelay)
		time.Sleep(cfg.RetryDelay)
	}

	return nil, fmt.Errorf("could not establish database connection after %d attempts: %w", cfg.RetryCount, err)
}

// Stats returns connection pool statistics for observability (TD-009)
func (db *Database) Stats() sql.DBStats {
	if db.DB != nil {
		return db.DB.Stats()
	}
	return sql.DBStats{}
}

func (db *Database) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}
```

## Common Mistakes
- **Hardcoded Pool & Retry Limits (TD-004, BE-002)**: Storing numbers like `25` (connections) or `5` (retry count) directly in the code. This prevents runtime tuning per environment (test vs production).
- **No Pool Observability (TD-009)**: Failing to expose pool telemetry. Always expose pool stats (`db.DB.Stats()`) so monitoring tools or logs can track connection pool saturation.
- **Shared Context Timeout in Loop**: Creating a single 10s context outside the retry loop. If early attempts are slow, the context expires and blocks subsequent attempts from establishing connection even if PostgreSQL boots up (BE-002).
- **Global DB Variables**: Storing the database instance in a package-level global variable. This creates side-effects and prevents testing with mock databases.
- **No Ping Verification**: Running `sqlx.Open` without verifying the connection with `.PingContext(ctx)`. `Open` only validates the arguments and returns a handle; it does not check if the server is reachable.
