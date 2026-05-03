package test

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// SetupTestDB loads configuration and initializes a database connection.
// It relies on a TEST_DB_URL env var or defaults to a standard local test database.
func SetupTestDB() *sqlx.DB {
	// Attempt to load .env if present
	_ = godotenv.Load("../.env")

	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		// Fallback to a standard local test db instance if not provided
		dbURL = "postgres://postgres:postgres@localhost:5432/openbench_test?sslmode=disable"
	}

	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

// CleanTestDB truncates the relevant tables to ensure a clean state between tests.
func CleanTestDB(db *sqlx.DB) {
	_, err := db.Exec("TRUNCATE TABLE tickets CASCADE;")
	if err != nil {
		log.Printf("Failed to clean test DB: %v", err)
	}
}
