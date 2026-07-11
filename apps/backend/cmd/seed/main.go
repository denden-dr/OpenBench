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

	authRepo := auth.NewRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Running database seeder...")
	if err := database.SeedDefaultAdmin(ctx, authRepo, cfg); err != nil {
		slog.Error("Failed to seed default admin", "error", err)
		os.Exit(1)
	}

	slog.Info("Database seeding completed.")
}
