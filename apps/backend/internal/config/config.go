package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	RetryCount      int
	RetryDelay      time.Duration
	PingTimeout     time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type AppConfig struct {
	Env            string
	Port           string
	AllowedOrigins []string
	DB             DatabaseConfig
}

// DSN returns the connection string for pgx/stdlib driver
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// LoadConfig loads the environment files depending on APP_ENV
func LoadConfig() (*AppConfig, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var envFile string
	switch env {
	case "test":
		envFile = ".env.test"
	default:
		envFile = ".env"
	}

	// Find the environment file in current directory, one, or two levels up (TD-003)
	var envFilePath string
	if _, err := os.Stat(envFile); err == nil {
		envFilePath = envFile
	} else if _, err := os.Stat(filepath.Join("..", envFile)); err == nil {
		envFilePath = filepath.Join("..", envFile)
	} else if _, err := os.Stat(filepath.Join("..", "..", envFile)); err == nil {
		envFilePath = filepath.Join("..", "..", envFile)
	}

	// Load dotenv file if it exists, proceed with system env variables if not found
	if envFilePath != "" {
		if err := godotenv.Load(envFilePath); err != nil {
			log.Printf("Warning: Failed to load %s: %v", envFilePath, err)
		}
	} else {
		log.Printf("Info: No %s file found, relying on system environment variables", envFile)
	}

	// Parse CORS allowed origins (TD-002)
	originsRaw := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	var origins []string
	for _, o := range strings.Split(originsRaw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}

	cfg := &AppConfig{
		Env:            env,
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: origins,
		DB: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "openbench"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			RetryCount:      getEnvAsInt("DB_RETRY_COUNT", 5),
			RetryDelay:      getEnvAsDuration("DB_RETRY_DELAY", 2*time.Second),
			PingTimeout:     getEnvAsDuration("DB_PING_TIMEOUT", 2*time.Second),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 15*time.Minute),
		},
	}

	// Validate Configuration (BE-004)
	if cfg.Env != "development" && cfg.Env != "test" {
		if cfg.DB.Password == "" {
			return nil, fmt.Errorf("database password (DB_PASSWORD) must not be empty in environment %q", cfg.Env)
		}
		if cfg.DB.SSLMode == "disable" {
			return nil, fmt.Errorf("insecure database SSL mode (DB_SSLMODE=disable) is not allowed in environment %q", cfg.Env)
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valStr, exists := os.LookupEnv(key); exists {
		var val int
		if _, err := fmt.Sscanf(valStr, "%d", &val); err == nil {
			return val
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if valStr, exists := os.LookupEnv(key); exists {
		if d, err := time.ParseDuration(valStr); err == nil {
			return d
		}
	}
	return fallback
}
