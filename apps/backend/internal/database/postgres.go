package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/jackc/pgx/v5"
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
	configPool.ConnConfig.Tracer = &slogQueryTracer{}

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

		// Calculate exponential delay: base * (2^attempt) capped at RetryMaxDelay
		factor := 1 << attempt
		if attempt > 30 {
			factor = 1 << 30
		}
		delay := cfg.RetryBaseDelay * time.Duration(factor)
		if delay > cfg.RetryMaxDelay {
			delay = cfg.RetryMaxDelay
		}

		slog.WarnContext(ctx, "DB Connect Attempt failed",
			slog.Int("attempt", attempt+1),
			slog.Int("max_retries", cfg.MaxRetries),
			slog.Any("error", lastErr),
			slog.Duration("retry_delay", delay),
		)
		cancel()
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", cfg.MaxRetries, lastErr)
}

type queryKey string

const queryStartKey queryKey = "query_start_time"

type slogQueryTracer struct{}

func (slogQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	slog.DebugContext(ctx, "SQL Query Start", slog.String("sql", data.SQL), slog.Any("args", data.Args))
	return context.WithValue(ctx, queryStartKey, time.Now())
}

func (slogQueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	duration := time.Duration(0)
	if startTime, ok := ctx.Value(queryStartKey).(time.Time); ok {
		duration = time.Since(startTime)
	}

	if data.Err != nil {
		slog.DebugContext(ctx, "SQL Query End with Error", slog.Duration("duration", duration), slog.Any("error", data.Err))
	} else {
		slog.DebugContext(ctx, "SQL Query End", slog.Duration("duration", duration))
	}
}
