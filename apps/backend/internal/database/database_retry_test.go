package database_test

import (
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnectionRetry(t *testing.T) {
	// Configure database config to fail and retry quickly on a closed port
	cfg := &config.DatabaseConfig{
		Host:        "127.0.0.1",
		Port:        "5439", // Closed port
		User:        "postgres",
		Password:    "password",
		DBName:      "test_db",
		SSLMode:     "disable",
		RetryCount:  3,
		RetryDelay:  10 * time.Millisecond,
		PingTimeout: 10 * time.Millisecond,
	}

	start := time.Now()
	db, err := database.NewConnection(cfg)
	elapsed := time.Since(start)

	// Connection must fail
	assert.Nil(t, db)
	assert.Error(t, err)

	// Elapsed time should be at least 20ms (since it slept 2 times for 10ms each: retry 1, sleep 10ms, retry 2, sleep 10ms, retry 3, fail)
	assert.GreaterOrEqual(t, elapsed, 20*time.Millisecond, "Should have retried and slept at least 20ms")
}
