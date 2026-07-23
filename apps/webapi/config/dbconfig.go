package config

import (
	"fmt"
	"net/url"
	"time"
)

type DBConfig struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            string        `mapstructure:"port" validate:"required"`
	User            string        `mapstructure:"user" validate:"required"`
	Password        string        `mapstructure:"password" validate:"required"`
	Name            string        `mapstructure:"name" validate:"required"`
	SSLMode         string        `mapstructure:"sslmode"`
	SSLRootCert     string        `mapstructure:"sslrootcert"`
	SSLCert         string        `mapstructure:"sslcert"`
	SSLKey          string        `mapstructure:"sslkey"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxRetries      int           `mapstructure:"max_retries"`
	RetryBaseDelay  time.Duration `mapstructure:"retry_base_delay"`
	RetryMaxDelay   time.Duration `mapstructure:"retry_max_delay"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

func (db DBConfig) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(db.User, db.Password),
		Host:   fmt.Sprintf("%s:%s", db.Host, db.Port),
		Path:   db.Name,
	}
	q := u.Query()
	if db.SSLMode != "" {
		q.Set("sslmode", db.SSLMode)
	}
	if db.SSLRootCert != "" {
		q.Set("sslrootcert", db.SSLRootCert)
	}
	if db.SSLCert != "" {
		q.Set("sslcert", db.SSLCert)
	}
	if db.SSLKey != "" {
		q.Set("sslkey", db.SSLKey)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
