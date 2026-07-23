package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// AppConfig holds core application configurations.
type AppConfig struct {
	Env            string `mapstructure:"env" validate:"required,oneof=development staging production testing"`
	AppName        string `mapstructure:"name" validate:"required"`
	Port           string `mapstructure:"port" validate:"required,numeric"`
	AllowedOrigins string `mapstructure:"allowed_origins" validate:"required"`
	EncryptionKey  string `mapstructure:"encryption_key" validate:"required,len=32"`
}

// Config wraps all application configuration groups.
type Config struct {
	App  AppConfig  `mapstructure:"app" validate:"required"`
	DB   DBConfig   `mapstructure:"db" validate:"required"`
	Auth AuthConfig `mapstructure:"auth" validate:"required"`
}

// findProjectRoot walks up the directory tree starting from the current working
// directory to find the directory containing go.mod. It falls back to the current
// directory if go.mod is not found.
func findProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}

	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return cwd
}

// Load reads config from settings.json and environment variables.
func Load() (*Config, error) {
	root := findProjectRoot()

	// Load .env first if it exists (so OS environment variables are populated)
	if os.Getenv("TEST_NO_ENV_FILE") != "true" {
		_ = godotenv.Load(filepath.Join(root, ".env"))
	}

	v := viper.New()

	v.SetConfigName("settings")
	v.SetConfigType("json")
	v.AddConfigPath(root) // Look in project root
	v.AddConfigPath(".")  // Fallback to current directory (e.g. inside Docker)

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Environment variables setup
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables that don't match the mapstructure structure directly
	envBindings := map[string]string{
		"app.env":             "APP_ENV",
		"app.name":            "APP_NAME",
		"app.port":            "PORT",
		"app.allowed_origins": "CORS_ALLOWED_ORIGINS",
		"app.encryption_key":  "APP_ENCRYPTION_KEY",
		"db.host":             "DB_HOST",
		"db.port":             "DB_PORT",
		"db.user":             "DB_USER",
		"db.name":             "DB_NAME",
		"db.sslmode":          "DB_SSLMODE",
		"db.sslrootcert":      "DB_SSLROOTCERT",
		"db.sslcert":          "DB_SSLCERT",
		"db.sslkey":           "DB_SSLKEY",
		"db.password":         "DB_PASSWORD",
		"auth.access_secret":  "JWT_ACCESS_SECRET",
		"auth.refresh_secret": "JWT_REFRESH_SECRET",
		"auth.access_expiry":  "JWT_ACCESS_EXPIRATION",
		"auth.refresh_expiry": "JWT_REFRESH_EXPIRATION",
	}
	for key, env := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			return nil, fmt.Errorf("failed to bind env %s to config key %s: %w", env, key, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate configuration using struct tags
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Database SSL enforcement in production
	if cfg.App.Env == "production" && (cfg.DB.SSLMode == "disable" || cfg.DB.SSLMode == "") {
		return nil, errors.New("database SSL mode cannot be disabled in production environment")
	}

	return &cfg, nil
}
