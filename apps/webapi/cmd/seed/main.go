package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/auth"
	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/inventory"
	"github.com/denden-dr/OpenBench/internal/logger"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

	if err := run(cfg); err != nil {
		slog.Error("Seed failed", "error", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		return err
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
		return err
	}

	if err := SeedDefaultProducts(ctx, db, cfg); err != nil {
		return err
	}

	slog.Info("Database seeding completed.")
	return nil
}

func SeedDefaultAdmin(ctx context.Context, queryRepo auth.QueryRepository, commandRepo auth.CommandRepository, cfg *config.Config) error {
	if cfg.App.Env != "development" && cfg.App.Env != "testing" {
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

func SeedDefaultProducts(ctx context.Context, db *sqlx.DB, cfg *config.Config) error {
	slog.Debug("SeedDefaultProducts: starting", slog.String("env", cfg.App.Env))

	if cfg.App.Env != "development" && cfg.App.Env != "testing" {
		slog.Debug("SeedDefaultProducts: skipped (env not dev/test)")
		return nil
	}

	commandRepo := inventory.NewCommandRepository(db)

	products := []struct {
		Name  string
		Price int64
		Stock int
	}{
		{Name: "Tempered Glass iPhone 15 Pro Max", Price: 75000, Stock: 18},
		{Name: "Silicon Case iPhone 15", Price: 120000, Stock: 4},
		{Name: "USB-C Charger Adapter 20W", Price: 299000, Stock: 2},
		{Name: "MicroUSB Cable 1m", Price: 35000, Stock: 25},
		{Name: "Lightning to USB-C Cable 2m", Price: 150000, Stock: 0},
	}

	slog.Debug("SeedDefaultProducts: products to seed", slog.Int("count", len(products)))

	for _, p := range products {
		product := &models.Product{
			ID:    uuid.New().String(),
			Name:  p.Name,
			Price: p.Price,
			Stock: p.Stock,
		}
		if err := commandRepo.Create(ctx, product); err != nil {
			slog.Error("SeedDefaultProducts: create failed", slog.String("name", p.Name), slog.Any("error", err))
			return err
		}
		slog.InfoContext(ctx, "Seeded product", slog.String("name", p.Name))
	}

	slog.Debug("SeedDefaultProducts: completed successfully")
	return nil
}
