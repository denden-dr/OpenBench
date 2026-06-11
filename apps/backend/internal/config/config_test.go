package config_test

import (
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Defaults(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("PORT", "9999")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "openbench_dev")
	t.Setenv("DB_SSLMODE", "disable")

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "9999", cfg.Port)
	assert.Equal(t, "127.0.0.1", cfg.DB.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, "postgres", cfg.DB.User)
	assert.Equal(t, "secret", cfg.DB.Password)
	assert.Equal(t, "openbench_dev", cfg.DB.DBName)
	assert.Equal(t, "disable", cfg.DB.SSLMode)
}

func TestLoadConfig_Validation(t *testing.T) {
	t.Run("production rejects empty password", func(t *testing.T) {
		t.Setenv("APP_ENV", "production")
		t.Setenv("DB_PASSWORD", "")
		t.Setenv("DB_SSLMODE", "require")

		cfg, err := config.LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "database password")
	})

	t.Run("production rejects disabled SSLMode", func(t *testing.T) {
		t.Setenv("APP_ENV", "production")
		t.Setenv("DB_PASSWORD", "secure_pass")
		t.Setenv("DB_SSLMODE", "disable")

		cfg, err := config.LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "SSL mode")
	})
}

func TestLoadConfig_PoolSettings(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")
	t.Setenv("DB_CONN_MAX_LIFETIME", "1h")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "30m")

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 50, cfg.DB.MaxOpenConns)
	assert.Equal(t, 10, cfg.DB.MaxIdleConns)
	assert.Equal(t, 1*time.Hour, cfg.DB.ConnMaxLifetime)
	assert.Equal(t, 30*time.Minute, cfg.DB.ConnMaxIdleTime)
}
