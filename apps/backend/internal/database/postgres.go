package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresDB creates a connection pool to the PostgreSQL database with retry logic.
func NewPostgresDB(cfg config.DBConfig) (*pgxpool.Pool, error) {
	configPool, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Apply pool configurations from config
	configPool.MaxConns = cfg.MaxConns
	configPool.MinConns = cfg.MinConns
	configPool.MaxConnLifetime = cfg.MaxConnLifetime
	configPool.MaxConnIdleTime = cfg.MaxConnIdleTime

	var pool *pgxpool.Pool
	var lastErr error

	for attempt := 0; attempt < cfg.MaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		pool, lastErr = pgxpool.NewWithConfig(ctx, configPool)
		if lastErr == nil {
			// Try pinging the database to confirm it's ready
			lastErr = pool.Ping(ctx)
			if lastErr == nil {
				cancel()
				return pool, nil
			}
			pool.Close()
		}
		cancel()

		// Calculate exponential delay: base * (2^attempt) capped at RetryMaxDelay
		factor := 1 << attempt
		if attempt > 30 {
			factor = 1 << 30
		}
		delay := cfg.RetryBaseDelay * time.Duration(factor)
		if delay > cfg.RetryMaxDelay {
			delay = cfg.RetryMaxDelay
		}

		log.Printf("[DB Connect] Attempt %d/%d failed: %v. Retrying in %v...", attempt+1, cfg.MaxRetries, lastErr, delay)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", cfg.MaxRetries, lastErr)
}
