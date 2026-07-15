package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/logger"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/google/uuid"
	"github.com/samber/hot"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.InitLogger(cfg.App.Env)

	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	cache := hot.NewHotCache[string, bool](hot.WTinyLFU, 100).Build()
	defer cache.StopJanitor()

	authQueryRepo := auth.NewQueryRepository(db, cache)
	authCommandRepo := auth.NewCommandRepository(db, cache)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Running database seeder...")
	if err := SeedDefaultAdmin(ctx, authQueryRepo, authCommandRepo, cfg); err != nil {
		slog.Error("Failed to seed default admin", "error", err)
		os.Exit(1)
	}

	slog.Info("Database seeding completed.")
}

func SeedDefaultAdmin(ctx context.Context, queryRepo auth.QueryRepository, commandRepo auth.CommandRepository, cfg *config.Config) error {
	if cfg.App.Env != "development" {
		return nil
	}

	email := "admin@openbench.com"
	password := "secretpassword123"

	// Cek apakah user dengan email tersebut sudah ada
	existing, err := queryRepo.GetUserByEmail(ctx, email)
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

	err = commandRepo.CreateUser(ctx, admin)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully seeded default admin", slog.String("email", email))
	return nil
}
