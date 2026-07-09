package database

import (
	"context"
	"log/slog"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SeedDefaultAdmin(ctx context.Context, repo auth.Repository, cfg *config.Config) error {
	if cfg.App.Env != "development" {
		return nil
	}

	email := "admin@openbench.com"
	password := "secretpassword123"

	// Cek apakah user dengan email tersebut sudah ada
	existing, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if existing != nil {
		slog.InfoContext(ctx, "Admin user already exists. Skipping...", slog.String("email", email))
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashed),
		FullName:     "Super Admin",
		Role:         "ADMIN",
	}

	err = repo.CreateUser(ctx, admin)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully seeded default admin", slog.String("email", email))
	return nil
}
