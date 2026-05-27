package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultCORSAllowOrigins(t *testing.T) {
	t.Setenv("DB_URL", "postgres://localhost:5432/openbench_test?sslmode=disable")
	if err := os.Unsetenv("CORS_ALLOW_ORIGINS"); err != nil {
		t.Fatalf("unset cors env: %v", err)
	}

	cfg := Load()

	const expected = "http://localhost:5173,http://127.0.0.1:5173"
	if cfg.CORSAllowOrigins != expected {
		t.Fatalf("expected default cors origins %q, got %q", expected, cfg.CORSAllowOrigins)
	}
}

func TestLoad_UsesConfiguredCORSAllowOrigins(t *testing.T) {
	t.Setenv("DB_URL", "postgres://localhost:5432/openbench_test?sslmode=disable")
	t.Setenv("CORS_ALLOW_ORIGINS", "https://admin.example.com")

	cfg := Load()

	if cfg.CORSAllowOrigins != "https://admin.example.com" {
		t.Fatalf("expected configured cors origins to be used, got %q", cfg.CORSAllowOrigins)
	}
}

func TestLoad_RateLimitConfig(t *testing.T) {
	t.Setenv("DB_URL", "postgres://localhost:5432/openbench_test?sslmode=disable")
	t.Setenv("APP_ENV", "production")
	t.Setenv("RATE_LIMIT_DISABLE", "true")
	t.Setenv("RATE_LIMIT_MAX", "50")
	t.Setenv("ADMIN_RATE_LIMIT_MAX", "150")

	cfg := Load()

	if !cfg.RateLimit.Disable {
		t.Error("expected rate limit to be disabled")
	}
	if cfg.RateLimit.MaxPublic != 50 {
		t.Errorf("expected public max rate limit to be 50, got %d", cfg.RateLimit.MaxPublic)
	}
	if cfg.RateLimit.MaxAdmin != 150 {
		t.Errorf("expected admin max rate limit to be 150, got %d", cfg.RateLimit.MaxAdmin)
	}

	// Test fallbacks when not set
	os.Unsetenv("RATE_LIMIT_DISABLE")
	os.Unsetenv("RATE_LIMIT_MAX")
	os.Unsetenv("ADMIN_RATE_LIMIT_MAX")
	t.Setenv("APP_ENV", "production")

	cfgDefault := Load()
	if cfgDefault.RateLimit.Disable {
		t.Error("expected default rate limit to not be disabled")
	}
	if cfgDefault.RateLimit.MaxPublic != 20 {
		t.Errorf("expected default public max to be 20, got %d", cfgDefault.RateLimit.MaxPublic)
	}
	if cfgDefault.RateLimit.MaxAdmin != 100 {
		t.Errorf("expected default admin max to be 100, got %d", cfgDefault.RateLimit.MaxAdmin)
	}

	// Test test defaults
	t.Setenv("APP_ENV", "test")
	cfgTest := Load()
	if cfgTest.RateLimit.MaxPublic != 1000 {
		t.Errorf("expected test public max to be 1000, got %d", cfgTest.RateLimit.MaxPublic)
	}
	if cfgTest.RateLimit.MaxAdmin != 1000 {
		t.Errorf("expected test admin max to be 1000, got %d", cfgTest.RateLimit.MaxAdmin)
	}
}
