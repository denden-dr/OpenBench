package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL             string `envconfig:"DATABASE_URL" required:"true"`
	DBMaxOpenConns          int    `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	DBMaxIdleConns          int    `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	DBConnMaxLifetimeSecs   int    `envconfig:"DB_CONN_MAX_LIFETIME_SECS" default:"300"`
	DBConnMaxIdleTimeSecs   int    `envconfig:"DB_CONN_MAX_IDLE_TIME_SECS" default:"60"`
	DBHealthPingTimeoutSecs int    `envconfig:"DB_HEALTH_PING_TIMEOUT_SECS" default:"2"`
}

// LoadConfig reads environment variables from .env (if present) and populates the Config struct.
func LoadConfig() (*Config, error) {
	// Best-effort .env loading — not an error if file is missing (e.g., in containers)
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	// Validate Database Connection Pool settings
	if cfg.DBMaxOpenConns < 1 {
		return nil, fmt.Errorf("invalid config: DB_MAX_OPEN_CONNS must be at least 1")
	}
	if cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		return nil, fmt.Errorf("invalid config: DB_MAX_IDLE_CONNS (%d) cannot exceed DB_MAX_OPEN_CONNS (%d)", cfg.DBMaxIdleConns, cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxIdleConns < 1 {
		return nil, fmt.Errorf("invalid config: DB_MAX_IDLE_CONNS must be at least 1")
	}
	if cfg.DBConnMaxLifetimeSecs <= 0 {
		return nil, fmt.Errorf("invalid config: DB_CONN_MAX_LIFETIME_SECS must be greater than 0")
	}
	if cfg.DBConnMaxIdleTimeSecs <= 0 {
		return nil, fmt.Errorf("invalid config: DB_CONN_MAX_IDLE_TIME_SECS must be greater than 0")
	}
	if cfg.DBHealthPingTimeoutSecs <= 0 {
		return nil, fmt.Errorf("invalid config: DB_HEALTH_PING_TIMEOUT_SECS must be greater than 0")
	}

	return &cfg, nil
}
