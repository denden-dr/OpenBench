//go:build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
)

func TestDatabaseConnection(t *testing.T) {
	// Set the environment to test (TD-006)
	t.Setenv("APP_ENV", "test")

	// Load configuration (it will automatically find .env.test by traversing up)
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Attempt database connection
	db, err := database.NewConnection(&cfg.DB)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v. Make sure the podman test db container is running.", err)
	}
	defer db.Close()

	// Verify query execution
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var result int
	err = db.DB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Failed to execute test query: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected query to return 1, got %d", result)
	}
}
