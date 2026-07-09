package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Temporarily set env vars
	os.Setenv("APP_ENV", "testing")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("JWT_ACCESS_SECRET", "test_access_secret")
	os.Setenv("JWT_REFRESH_SECRET", "test_refresh_secret")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
	}()

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Env was overridden by environment variable
	assert.Equal(t, "testing", cfg.App.Env)

	// Password loaded from env
	assert.Equal(t, "testpassword", cfg.DB.Password)

	// Loaded from settings.json
	assert.Equal(t, "OpenBench API", cfg.App.AppName)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, int32(25), cfg.DB.MaxConns)

	// Check duration loading from JSON
	assert.Equal(t, 15*time.Minute, cfg.Auth.AccessExpiry)
}
