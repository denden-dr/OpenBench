package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigins string
}

func LoadConfig() *Config {
	// Attempt to load .env file, but do not fail if it's missing
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/openbench?sslmode=disable"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
