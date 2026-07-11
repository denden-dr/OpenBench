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
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "openbench")
	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("JWT_ACCESS_SECRET", "testing_access_secret_longer_than_32_bytes_1234567890")
	os.Setenv("JWT_REFRESH_SECRET", "testing_refresh_secret_longer_than_32_bytes_1234567890")
	os.Setenv("APP_ENCRYPTION_KEY", "this_is_a_secret_key_32_chars_ok")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
		os.Unsetenv("APP_ENCRYPTION_KEY")
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
