package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds core application configurations.
type AppConfig struct {
	Env     string
	AppName string
	Port    string
}

// DBConfig holds PostgreSQL connection settings.
type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxRetries      int
	RetryBaseDelay  time.Duration
	RetryMaxDelay   time.Duration
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// Config wraps all application configuration groups.
type Config struct {
	App AppConfig
	DB  DBConfig
}

// DSN returns the connection string formatted for the pgxpool driver.
func (db DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

// Load reads config from environment variables. If any variable is missing or malformed, it returns an error.
func Load() (*Config, error) {
	// Load .env from workspace root or current directory if it exists
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	// APP variables
	appEnv, err := requireEnv("APP_ENV")
	if err != nil {
		return nil, err
	}

	appName, err := requireEnv("APP_NAME")
	if err != nil {
		return nil, err
	}

	appPort, err := requireEnv("PORT")
	if err != nil {
		return nil, err
	}

	// DB variables
	dbHost, err := requireEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbPort, err := requireEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbUser, err := requireEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPassword, err := requireEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := requireEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	dbSSLMode, err := requireEnv("DB_SSLMODE")
	if err != nil {
		return nil, err
	}

	dbMaxConnsStr, err := requireEnv("DB_MAX_CONNS")
	if err != nil {
		return nil, err
	}
	dbMaxConns, err := strconv.Atoi(dbMaxConnsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONNS: %w", err)
	}

	dbMinConnsStr, err := requireEnv("DB_MIN_CONNS")
	if err != nil {
		return nil, err
	}
	dbMinConns, err := strconv.Atoi(dbMinConnsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MIN_CONNS: %w", err)
	}

	dbMaxRetriesStr, err := requireEnv("DB_MAX_RETRIES")
	if err != nil {
		return nil, err
	}
	dbMaxRetries, err := strconv.Atoi(dbMaxRetriesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_RETRIES: %w", err)
	}

	dbRetryBaseDelayStr, err := requireEnv("DB_RETRY_BASE_DELAY")
	if err != nil {
		return nil, err
	}
	dbRetryBaseDelay, err := time.ParseDuration(dbRetryBaseDelayStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_RETRY_BASE_DELAY: %w", err)
	}

	dbRetryMaxDelayStr, err := requireEnv("DB_RETRY_MAX_DELAY")
	if err != nil {
		return nil, err
	}
	dbRetryMaxDelay, err := time.ParseDuration(dbRetryMaxDelayStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_RETRY_MAX_DELAY: %w", err)
	}

	dbMaxConnLifetimeStr, err := requireEnv("DB_MAX_CONN_LIFETIME")
	if err != nil {
		return nil, err
	}
	dbMaxConnLifetime, err := time.ParseDuration(dbMaxConnLifetimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONN_LIFETIME: %w", err)
	}

	dbMaxConnIdleTimeStr, err := requireEnv("DB_MAX_CONN_IDLE_TIME")
	if err != nil {
		return nil, err
	}
	dbMaxConnIdleTime, err := time.ParseDuration(dbMaxConnIdleTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONN_IDLE_TIME: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Env:     appEnv,
			AppName: appName,
			Port:    appPort,
		},
		DB: DBConfig{
			Host:            dbHost,
			Port:            dbPort,
			User:            dbUser,
			Password:        dbPassword,
			Name:            dbName,
			SSLMode:         dbSSLMode,
			MaxConns:        int32(dbMaxConns),
			MinConns:        int32(dbMinConns),
			MaxRetries:      dbMaxRetries,
			RetryBaseDelay:  dbRetryBaseDelay,
			RetryMaxDelay:   dbRetryMaxDelay,
			MaxConnLifetime: dbMaxConnLifetime,
			MaxConnIdleTime: dbMaxConnIdleTime,
		},
	}

	return cfg, nil
}

func requireEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return "", fmt.Errorf("environment variable %s is required", key)
	}
	return value, nil
}
