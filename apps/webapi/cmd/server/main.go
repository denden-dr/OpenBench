package main

import (
	"log/slog"
	"os"

	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/auth"
	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/events"
	"github.com/denden-dr/OpenBench/internal/inventory"
	"github.com/denden-dr/OpenBench/internal/logger"
	"github.com/denden-dr/OpenBench/internal/pos"
	"github.com/denden-dr/OpenBench/internal/ticket"
	"github.com/denden-dr/OpenBench/internal/warranty"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.InitLogger(cfg.App.Env)

	slog.Info("Application starting",
		slog.String("app_name", cfg.App.AppName),
		slog.String("env", cfg.App.Env),
		slog.String("port", cfg.App.Port),
		slog.String("allowed_origins", cfg.App.AllowedOrigins),
		slog.String("db_host", cfg.DB.Host),
		slog.String("db_name", cfg.DB.Name),
	)

	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Database connection pool established",
		slog.Int("max_open_conns", int(cfg.DB.MaxConns)),
		slog.Int("max_idle_conns", int(cfg.DB.MinConns)),
		slog.Duration("max_conn_lifetime", cfg.DB.MaxConnLifetime),
		slog.Duration("max_conn_idle_time", cfg.DB.MaxConnIdleTime),
	)

	eventBus := events.NewAsyncEventBus(100)
	defer eventBus.Close()

	txManager := database.NewTxManager(db)

	// Wire Auth Layers
	authMod, stopAuth := auth.NewModule(db, cfg)
	defer stopAuth()

	warrantyMod := warranty.NewModule(db, txManager)
	ticketMod := ticket.NewModule(db, txManager, warrantyMod.Service, eventBus, cfg.App.EncryptionKey)
	inventoryMod := inventory.NewModule(db)
	posMod := pos.NewModule(db, txManager, inventoryMod.QueryRepo, inventoryMod.CommandRepo)

	app := newFiberApp(cfg)

	registerRoutes(app, cfg, db, authMod, warrantyMod, ticketMod, inventoryMod, posMod)
	listenAndServe(app, cfg.App.Port)
}
