---
name: managing-multi-environment-config
description: Use when setting up environment-based configurations in a Go application, managing different environments (development, testing, production) using .env files, and loading variables cleanly.
---

# Managing Multi-Environment Config

## Overview
Go applications should support running in different environments (development, testing, production) seamlessly. The configuration should be loaded dynamically based on the current environment, controlled by a standardized `APP_ENV` environment variable, ensuring that test runs do not affect development databases or systems.

To keep configurations deterministic and secure, avoid implicit directory walking to discover `.env` files, avoid hardcoding service origins (e.g., CORS), and validate all production/non-local settings.

## When to Use
- Designing the configuration loader for a Go service.
- Defining environment files (`.env`, `.env.test`, `.env.example`).
- Isolating development and testing parameters (e.g., database host, port, credentials).

## Core Pattern
Use a central configuration package that reads `APP_ENV` to determine which dotenv file to load at startup. Provide an `.env.example` template as documentation. Add a validation gate to reject empty passwords or insecure SSL configs when running in a non-local environment. Expose CORS origins dynamically from the configuration.

### Multi-Environment Config Loader Example
```go
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AppConfig struct {
	Env          string
	Port         string
	AllowedOrigins []string
	DB           DatabaseConfig
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// LoadConfig loads configuration from env files or system environment variables.
// It searches for env files in the current working directory or a path explicitly passed.
func LoadConfig() (*AppConfig, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // default fallback
	}

	var envFile string
	switch env {
	case "test":
		envFile = ".env.test"
	default:
		envFile = ".env"
	}

	// Read directly from current working directory or look at parent directory only if necessary,
	// but avoid infinite upward directory traversal (TD-003).
	envFilePath := envFile
	if _, err := os.Stat(envFilePath); err != nil {
		// Look 1 level up only as a fallback for test directories
		fallbackPath := filepath.Join("..", envFile)
		if _, err := os.Stat(fallbackPath); err == nil {
			envFilePath = fallbackPath
		} else {
			envFilePath = ""
		}
	}

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
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "openbench"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	// Validate required variables for non-dev and non-test environments (BE-004)
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
```

### Environment Files Template Structure

#### 1. `.env.example` (Committed to Git)
```ini
PORT=8080
APP_ENV=development
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=openbench_user
DB_PASSWORD=openbench_secure_password
DB_NAME=openbench_db
DB_SSLMODE=disable
```

## Common Mistakes
- **Implicit Env File Discovery (TD-003)**: Walking up parent directories indefinitely to find `.env` files. This makes configuration load-dependent on the current working directory. Keep it to a single level or explicit path parameter.
- **Hardcoded CORS Origins (TD-002)**: Hardcoding allowed origins inside the router or application startup code. Always pull allowed origins from the configuration.
- **ENV/APP_ENV Inconsistency**: Mixing up variable names like `ENV` in docker-compose or .env files but reading `APP_ENV` in the application loader (BE-003).
- **No Production Safety Gates**: Allowing the application to start up in a production environment with insecure local defaults like an empty database password or `DB_SSLMODE=disable` (BE-004).
- **Committing Sensitive Secrets**: Committing actual `.env` files with production or local personal passwords to Git. Always add `.env` to `.gitignore`.
