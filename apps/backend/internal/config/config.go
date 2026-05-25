package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	CORSAllowOrigins string
	Database         DatabaseConfig
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	PingTimeout     time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	return &Config{
		Port:             getEnv("PORT", "3000"),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173"),
		Database:         loadDatabaseConfig(dbURL),
	}
}

func DefaultDatabaseConfig(dbURL string) DatabaseConfig {
	return DatabaseConfig{
		URL:             dbURL,
		MaxOpenConns:    25,
		MaxIdleConns:    12,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		PingTimeout:     5 * time.Second,
	}
}

func loadDatabaseConfig(dbURL string) DatabaseConfig {
	cfg := DefaultDatabaseConfig(dbURL)
	cfg.MaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", cfg.MaxOpenConns)
	cfg.MaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", cfg.MaxIdleConns)
	cfg.ConnMaxLifetime = getEnvDuration("DB_CONN_MAX_LIFETIME", cfg.ConnMaxLifetime)
	cfg.ConnMaxIdleTime = getEnvDuration("DB_CONN_MAX_IDLE_TIME", cfg.ConnMaxIdleTime)
	cfg.PingTimeout = getEnvDuration("DB_PING_TIMEOUT", cfg.PingTimeout)
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("%s must be an integer: %v", key, err)
	}
	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("%s must be a duration like 5s or 5m: %v", key, err)
	}
	return parsed
}
