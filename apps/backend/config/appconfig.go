package config

import (
	"github.com/joho/godotenv"
)

// AppConfig holds core application configurations.
type AppConfig struct {
	Env     string
	AppName string
	Port    string
}

// Config wraps all application configuration groups.
type Config struct {
	App  AppConfig
	DB   DBConfig
	Auth AuthConfig
}

func LoadAppConfig() (AppConfig, error) {
	var err error
	var cfg AppConfig

	if cfg.Env, err = requireEnv("APP_ENV"); err != nil {
		return cfg, err
	}
	if cfg.AppName, err = requireEnv("APP_NAME"); err != nil {
		return cfg, err
	}
	if cfg.Port, err = requireEnv("PORT"); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Load reads config from environment variables. If any variable is missing or malformed, it returns an error.
func Load() (*Config, error) {
	// Load .env from workspace root or current directory if it exists
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	appCfg, err := LoadAppConfig()
	if err != nil {
		return nil, err
	}

	dbCfg, err := LoadDBConfig()
	if err != nil {
		return nil, err
	}

	authCfg, err := LoadAuthConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App:  appCfg,
		DB:   dbCfg,
		Auth: authCfg,
	}

	return cfg, nil
}
