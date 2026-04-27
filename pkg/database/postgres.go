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
