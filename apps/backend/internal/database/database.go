package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB *sqlx.DB
}

// NewConnection initializes a new database connection pool with pgx stdlib and sqlx
func NewConnection(cfg *config.DatabaseConfig) (*Database, error) {
	var db *sqlx.DB
	var err error
	dsn := cfg.DSN()

	retryCount := cfg.RetryCount
	if retryCount <= 0 {
		retryCount = 5
	}
	retryDelay := cfg.RetryDelay
	if retryDelay <= 0 {
		retryDelay = 2 * time.Second
	}
	pingTimeout := cfg.PingTimeout
	if pingTimeout <= 0 {
		pingTimeout = 2 * time.Second
	}

	// Retry connection logic
	for i := 0; i < retryCount; i++ {
		// Use "pgx" driver provided by pgx/v5/stdlib
		db, err = sqlx.Open("pgx", dsn)
		if err == nil {
			// Connection Pool Settings (TD-004)
			db.SetMaxOpenConns(cfg.MaxOpenConns)
			db.SetMaxIdleConns(cfg.MaxIdleConns)
			db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
			db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

			// Ping database to verify connection health using a fresh per-attempt context
			pingCtx, pingCancel := context.WithTimeout(context.Background(), pingTimeout)
			err = db.PingContext(pingCtx)
			pingCancel()

			if err == nil {
				log.Printf("Successfully connected to the database with sqlx (Pool settings: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v)",
					cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime, cfg.ConnMaxIdleTime)
				return &Database{DB: db}, nil
			}
			db.Close()
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...", i+1, retryCount, err, retryDelay)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("could not establish database connection: %w", err)
}

// Stats returns connection pool statistics (TD-009)
func (db *Database) Stats() sql.DBStats {
	if db.DB != nil {
		return db.DB.Stats()
	}
	return sql.DBStats{}
}

// Close closes the database connection pool
func (db *Database) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}
