package config

import (
	"fmt"
	"time"
)

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

func LoadDBConfig() (DBConfig, error) {
	var err error
	var cfg DBConfig

	if cfg.Host, err = requireEnv("DB_HOST"); err != nil {
		return cfg, err
	}
	if cfg.Port, err = requireEnv("DB_PORT"); err != nil {
		return cfg, err
	}
	if cfg.User, err = requireEnv("DB_USER"); err != nil {
		return cfg, err
	}
	if cfg.Password, err = requireEnv("DB_PASSWORD"); err != nil {
		return cfg, err
	}
	if cfg.Name, err = requireEnv("DB_NAME"); err != nil {
		return cfg, err
	}
	if cfg.SSLMode, err = requireEnv("DB_SSLMODE"); err != nil {
		return cfg, err
	}

	maxConns, err := requireEnvInt("DB_MAX_CONNS")
	if err != nil {
		return cfg, err
	}
	cfg.MaxConns = int32(maxConns)

	minConns, err := requireEnvInt("DB_MIN_CONNS")
	if err != nil {
		return cfg, err
	}
	cfg.MinConns = int32(minConns)

	if cfg.MaxRetries, err = requireEnvInt("DB_MAX_RETRIES"); err != nil {
		return cfg, err
	}
	if cfg.RetryBaseDelay, err = requireEnvDuration("DB_RETRY_BASE_DELAY"); err != nil {
		return cfg, err
	}
	if cfg.RetryMaxDelay, err = requireEnvDuration("DB_RETRY_MAX_DELAY"); err != nil {
		return cfg, err
	}
	if cfg.MaxConnLifetime, err = requireEnvDuration("DB_MAX_CONN_LIFETIME"); err != nil {
		return cfg, err
	}
	if cfg.MaxConnIdleTime, err = requireEnvDuration("DB_MAX_CONN_IDLE_TIME"); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (db DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}
