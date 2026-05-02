package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port  string
    DBURL string
}

func Load() *Config {
    _ = godotenv.Load()
    return &Config{
        Port:  getEnv("PORT", "8080"),
        DBURL: getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/openbench?sslmode=disable"),
    }
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
