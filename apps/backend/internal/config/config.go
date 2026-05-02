package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	return &Config{
		Port:  getEnv("PORT", "3000"),
		DBURL: dbURL,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
