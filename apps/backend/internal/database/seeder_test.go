//go:build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func TestSeedDevUsers(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewConnection(&cfg.DB)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clean up any existing admin user to test creation
	_, err = db.DB.ExecContext(ctx, "DELETE FROM users WHERE email = $1", "admin@openbench.dev")
	if err != nil {
		t.Fatalf("Failed to clean up test user: %v", err)
	}

	// First seeding (Creation)
	err = database.SeedDevUsers(db)
	if err != nil {
		t.Fatalf("Failed to run seed: %v", err)
	}

	// Verify details
	var u struct {
		ID           string `db:"id"`
		Email        string `db:"email"`
		PasswordHash string `db:"password_hash"`
		Role         string `db:"role"`
	}
	err = db.DB.GetContext(ctx, &u, "SELECT id, email, password_hash, role FROM users WHERE email = $1", "admin@openbench.dev")
	if err != nil {
		t.Fatalf("Failed to retrieve seeded user: %v", err)
	}

	if u.Email != "admin@openbench.dev" {
		t.Errorf("Expected email to be admin@openbench.dev, got %s", u.Email)
	}
	if u.Role != "admin" {
		t.Errorf("Expected role to be admin, got %s", u.Role)
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("SecureAdminPassword123!"))
	if err != nil {
		t.Errorf("Password hash verification failed: %v", err)
	}

	// Second seeding (Idempotency check)
	err = database.SeedDevUsers(db)
	if err != nil {
		t.Errorf("SeedDevUsers should be idempotent and not return an error when user already exists, got: %v", err)
	}
}
