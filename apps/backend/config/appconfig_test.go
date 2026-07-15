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

func TestLoad_Validation(t *testing.T) {
	// Helper to set valid env vars first
	setValidEnv := func() {
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
	}

	clearEnv := func() {
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
	}

	t.Run("valid config", func(t *testing.T) {
		setValidEnv()
		defer clearEnv()
		cfg, err := Load()
		require.NoError(t, err)
		require.NotNil(t, cfg)
	})

	t.Run("missing DB_HOST", func(t *testing.T) {
		setValidEnv()
		defer clearEnv()
		os.Unsetenv("DB_HOST")
		_, err := Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Host")
	})

	t.Run("short encryption key", func(t *testing.T) {
		setValidEnv()
		defer clearEnv()
		os.Setenv("APP_ENCRYPTION_KEY", "short")
		_, err := Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "EncryptionKey")
	})

	t.Run("short JWT access secret", func(t *testing.T) {
		setValidEnv()
		defer clearEnv()
		os.Setenv("JWT_ACCESS_SECRET", "short_secret")
		_, err := Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "AccessSecret")
	})

	t.Run("placeholder JWT refresh secret", func(t *testing.T) {
		setValidEnv()
		defer clearEnv()
		os.Setenv("JWT_REFRESH_SECRET", "your_super_secret_for_jwt_refresh")
		_, err := Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RefreshSecret")
	})

	t.Run("invalid APP_ENV", func(t *testing.T) {
		setValidEnv()
		defer clearEnv()
		os.Setenv("APP_ENV", "invalid_env")
		_, err := Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Env")
	})

	t.Run("non-numeric PORT", func(t *testing.T) {
		setValidEnv()
		defer clearEnv()
		os.Setenv("PORT", "abc")
		_, err := Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Port")
	})
}

func TestDBConfig_DSN(t *testing.T) {
	db := DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "myuser@admin",
		Password: "p@ss:word#123",
		Name:     "mydb",
		SSLMode:  "verify-full",
	}

	dsn := db.DSN()
	// User and Password should be URL-encoded:
	// 'myuser@admin' -> 'myuser%40admin'
	// 'p@ss:word#123' -> 'p%40ss%3Aword%23123'
	assert.Contains(t, dsn, "myuser%40admin")
	assert.Contains(t, dsn, "p%40ss%3Aword%23123")
	assert.Contains(t, dsn, "sslmode=verify-full")
}
