package main

import (
	"time"

	"github.com/denden-dr/OpenBench/internal/handlers"
	"github.com/denden-dr/OpenBench/internal/middleware"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/denden-dr/OpenBench/internal/service"
	"github.com/denden-dr/OpenBench/pkg/config"
	"github.com/denden-dr/OpenBench/pkg/database"
	"github.com/denden-dr/OpenBench/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Database with configurable pool settings
	db, err := database.NewPostgresDB(
		cfg.DatabaseURL,
		cfg.DBMaxOpenConns,
		cfg.DBMaxIdleConns,
		cfg.DBConnMaxLifetimeSecs,
		cfg.DBConnMaxIdleTimeSecs,
	)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Log applied settings for observability
	log.Info("Database connection established",
		zap.Int("max_open_conns", cfg.DBMaxOpenConns),
		zap.Int("max_idle_conns", cfg.DBMaxIdleConns),
		zap.Int("conn_max_lifetime_secs", cfg.DBConnMaxLifetimeSecs),
		zap.Int("conn_max_idle_time_secs", cfg.DBConnMaxIdleTimeSecs),
	)

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	techRepo := repository.NewTechnicianRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	techBookingService := service.NewTechBookingService(bookingRepo)
	techHandler := handlers.NewTechHandler(techBookingService)

	bookingService := service.NewBookingService(bookingRepo)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	healthHandler := handlers.NewHealthHandler(db, time.Duration(cfg.DBHealthPingTimeoutSecs)*time.Second)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "OpenBench API",
	})

	// Use zap-based middleware for all HTTP requests
	app.Use(middleware.ZapLogger(log))

	// Define routes
	app.Get("/health", healthHandler.HealthCheck)

	// User routes
	v1 := app.Group("/api/v1")
	users := v1.Group("/users")
	// TODO: Add AuthMiddleware here once implemented
	users.Get("/me", userHandler.GetMe)

	// Booking routes (Customer)
	bookings := v1.Group("/bookings")
	{
		bookings.Post("/", bookingHandler.Create())
		bookings.Get("/:id", bookingHandler.GetByID())
		bookings.Post("/:id/approve", bookingHandler.Approve())
		bookings.Post("/:id/cancel", bookingHandler.Cancel())
	}

	// Technician routes
	tech := v1.Group("/tech", middleware.RequireTech(techRepo))
	{
		tech.Get("/bookings/available", techHandler.GetAvailableBookings())
		tech.Get("/bookings/mine", techHandler.GetMyBookings())
		tech.Post("/bookings/:id/assign", techHandler.AssignBooking())
		tech.Post("/bookings/:id/diagnose", techHandler.DiagnoseBooking())
		tech.Post("/bookings/:id/status", techHandler.UpdateBookingStatus())
	}

	// Log server start
	log.Info("Starting OpenBench API server on port 3000")

	// Start server
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
