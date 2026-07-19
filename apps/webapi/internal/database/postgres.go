package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/denden-dr/OpenBench/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
)

// slogAdapter implements sqldblogger.Logger interface using standard log/slog.
type slogAdapter struct{}

func (s *slogAdapter) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	var slogLevel slog.Level
	switch level {
	case sqldblogger.LevelError:
		slogLevel = slog.LevelError
	case sqldblogger.LevelInfo:
		slogLevel = slog.LevelInfo
	default:
		slogLevel = slog.LevelDebug
	}

	slog.LogAttrs(ctx, slogLevel, msg, attrs...)
}

// NewPostgresDB creates a connection pool to the PostgreSQL database with retry logic.
func NewPostgresDB(cfg config.DBConfig) (*sqlx.DB, error) {
	// First, obtain the raw pgx driver.
	rawDB, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %w", err)
	}
	driver := rawDB.Driver()
	rawDB.Close()

	// Wrap the driver with sqldb-logger for query logging.
	loggerAdapter := &slogAdapter{}
	wrappedDB := sqldblogger.OpenDriver(cfg.DSN(), driver, loggerAdapter,
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
	)

	// Apply connection pool configurations.
	wrappedDB.SetMaxOpenConns(int(cfg.MaxConns))
	wrappedDB.SetMaxIdleConns(int(cfg.MinConns))
	wrappedDB.SetConnMaxLifetime(cfg.MaxConnLifetime)
	wrappedDB.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	var lastErr error
	for attempt := 0; attempt < cfg.MaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		lastErr = wrappedDB.PingContext(ctx)
		cancel()

		if lastErr == nil {
			break
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
		time.Sleep(delay)
	}

	if lastErr != nil {
		wrappedDB.Close()
		return nil, fmt.Errorf("could not connect to database after %d attempts: %w", cfg.MaxRetries, lastErr)
	}

	sqlxDB := sqlx.NewDb(wrappedDB, "pgx")
	return sqlxDB, nil
}
