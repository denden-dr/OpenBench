package config

import "time"

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func LoadAuthConfig() (AuthConfig, error) {
	var err error
	var cfg AuthConfig

	if cfg.AccessSecret, err = requireEnv("JWT_ACCESS_SECRET"); err != nil {
		return cfg, err
	}
	if cfg.RefreshSecret, err = requireEnv("JWT_REFRESH_SECRET"); err != nil {
		return cfg, err
	}
	if cfg.AccessExpiry, err = requireEnvDuration("JWT_ACCESS_EXPIRATION"); err != nil {
		return cfg, err
	}
	if cfg.RefreshExpiry, err = requireEnvDuration("JWT_REFRESH_EXPIRATION"); err != nil {
		return cfg, err
	}

	return cfg, nil
}
