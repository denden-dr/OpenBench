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
