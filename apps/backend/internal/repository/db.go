package repository

import (
	"context"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(12) // Half of MaxOpenConns to avoid wasting resources
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close() // Explicitly close to avoid connection leak
		return nil, err
	}

	slog.Info("Database connection established")
	return db, nil
}
