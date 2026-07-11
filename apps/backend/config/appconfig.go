package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// AppConfig holds core application configurations.
type AppConfig struct {
	Env            string `mapstructure:"env"`
	AppName        string `mapstructure:"name"`
	Port           string `mapstructure:"port"`
	AllowedOrigins string `mapstructure:"allowed_origins"`
	EncryptionKey  string `mapstructure:"encryption_key"`
}

// Config wraps all application configuration groups.
type Config struct {
	App  AppConfig  `mapstructure:"app"`
	DB   DBConfig   `mapstructure:"db"`
	Auth AuthConfig `mapstructure:"auth"`
}

// Load reads config from settings.json and environment variables.
func Load() (*Config, error) {
	// Load .env first if it exists (so OS environment variables are populated)
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	v := viper.New()

	v.SetConfigName("settings")
	v.SetConfigType("json")
	v.AddConfigPath(".")              // Look in current directory
	v.AddConfigPath("..")             // Look in parent directory (e.g., if run from cmd/api or config)
	v.AddConfigPath("../..")          // Look in grandparent directory (e.g., if run from sub-packages)
	v.AddConfigPath("./apps/backend") // Look in apps/backend if run from workspace root

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Environment variables setup
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables that don't match the mapstructure structure directly
	_ = v.BindEnv("app.port", "PORT")
	_ = v.BindEnv("app.allowed_origins", "CORS_ALLOWED_ORIGINS")
	_ = v.BindEnv("app.encryption_key", "APP_ENCRYPTION_KEY")
	_ = v.BindEnv("db.host", "DB_HOST")
	_ = v.BindEnv("db.port", "DB_PORT")
	_ = v.BindEnv("db.user", "DB_USER")
	_ = v.BindEnv("db.name", "DB_NAME")
	_ = v.BindEnv("db.sslmode", "DB_SSLMODE")
	_ = v.BindEnv("db.password", "DB_PASSWORD")
	_ = v.BindEnv("auth.access_secret", "JWT_ACCESS_SECRET")
	_ = v.BindEnv("auth.refresh_secret", "JWT_REFRESH_SECRET")
	_ = v.BindEnv("auth.access_expiry", "JWT_ACCESS_EXPIRATION")
	_ = v.BindEnv("auth.refresh_expiry", "JWT_REFRESH_EXPIRATION")

	// Set defaults
	v.SetDefault("app.allowed_origins", "http://localhost:5173,http://localhost:3000")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate required variables
	if cfg.App.Env == "" {
		return nil, fmt.Errorf("app.env (APP_ENV) is required")
	}
	if cfg.App.AppName == "" {
		return nil, fmt.Errorf("app.name (APP_NAME) is required")
	}
	if cfg.App.Port == "" {
		return nil, fmt.Errorf("app.port (PORT) is required")
	}
	if cfg.App.EncryptionKey == "" {
		return nil, fmt.Errorf("app.encryption_key (APP_ENCRYPTION_KEY) is required")
	}
	if len(cfg.App.EncryptionKey) != 32 {
		return nil, fmt.Errorf("app.encryption_key (APP_ENCRYPTION_KEY) must be exactly 32 characters long")
	}
	
	// DB Validation
	if cfg.DB.Host == "" {
		return nil, fmt.Errorf("db.host (DB_HOST) is required")
	}
	if cfg.DB.Port == "" {
		return nil, fmt.Errorf("db.port (DB_PORT) is required")
	}
	if cfg.DB.User == "" {
		return nil, fmt.Errorf("db.user (DB_USER) is required")
	}
	if cfg.DB.Name == "" {
		return nil, fmt.Errorf("db.name (DB_NAME) is required")
	}
	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("db.password (DB_PASSWORD) is required")
	}
	if cfg.Auth.AccessSecret == "" {
		return nil, fmt.Errorf("auth.access_secret (JWT_ACCESS_SECRET) is required")
	}
	if cfg.Auth.RefreshSecret == "" {
		return nil, fmt.Errorf("auth.refresh_secret (JWT_REFRESH_SECRET) is required")
	}

	// Validate JWT Access Secret
	if len(cfg.Auth.AccessSecret) < 32 {
		return nil, fmt.Errorf("auth.access_secret (JWT_ACCESS_SECRET) must be at least 32 bytes long, got %d", len(cfg.Auth.AccessSecret))
	}
	if strings.Contains(cfg.Auth.AccessSecret, "your_super_secret") {
		return nil, fmt.Errorf("auth.access_secret (JWT_ACCESS_SECRET) cannot use default placeholder value")
	}

	// Validate JWT Refresh Secret
	if len(cfg.Auth.RefreshSecret) < 32 {
		return nil, fmt.Errorf("auth.refresh_secret (JWT_REFRESH_SECRET) must be at least 32 bytes long, got %d", len(cfg.Auth.RefreshSecret))
	}
	if strings.Contains(cfg.Auth.RefreshSecret, "your_super_secret") {
		return nil, fmt.Errorf("auth.refresh_secret (JWT_REFRESH_SECRET) cannot use default placeholder value")
	}

	return &cfg, nil
}
