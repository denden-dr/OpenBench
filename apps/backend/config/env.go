package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func requireEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return "", fmt.Errorf("environment variable %s is required", key)
	}
	return value, nil
}

func requireEnvInt(key string) (int, error) {
	strValue, err := requireEnv(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return value, nil
}

func requireEnvDuration(key string) (time.Duration, error) {
	strValue, err := requireEnv(key)
	if err != nil {
		return 0, err
	}
	value, err := time.ParseDuration(strValue)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return value, nil
}
