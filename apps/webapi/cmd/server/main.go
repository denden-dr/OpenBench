package main

import (
	"log/slog"
	"os"

	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/auth"
	"github.com/denden-dr/OpenBench/internal/dashboard"
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

	// Wire Domains
	// 1. Repositories
	ticketQR := ticket.NewQueryRepository(db)
	ticketCR := ticket.NewCommandRepository(db)

	warrantyQR := warranty.NewQueryRepository(db)
	warrantyCR := warranty.NewCommandRepository(db)

	invQR := inventory.NewQueryRepository(db)
	invCR := inventory.NewCommandRepository(db)

	posQR := pos.NewQueryRepository(db)
	posCR := pos.NewCommandRepository(db)

	// 2. Cross-domain helpers (Creators / Generators)
	ticketCreator := ticket.NewCreator(ticketQR, ticketCR)
	warrantyGen := warranty.NewGenerator(warrantyCR)

	// 3. Services
	warrantySvc := warranty.NewService(warrantyQR, warrantyCR, txManager, ticketCreator)
	ticketSvc := ticket.NewService(ticketQR, ticketCR, txManager, warrantyGen, eventBus, cfg.App.EncryptionKey)
	invSvc := inventory.NewService(invQR, invCR)
	posSvc := pos.NewService(posQR, posCR, invQR, invCR, txManager)

	// 4. Modules
	warrantyMod := warranty.NewModule(warrantySvc)
	ticketMod := ticket.NewModule(ticketSvc)
	inventoryMod := inventory.NewModule(invSvc, invQR, invCR)
	posMod := pos.NewModule(posSvc)
	dashboardMod := dashboard.NewModule(db)

	app := newFiberApp(cfg)

	registerRoutes(app, cfg, db, authMod, warrantyMod, ticketMod, inventoryMod, posMod, dashboardMod)
	listenAndServe(app, cfg.App.Port)
}
