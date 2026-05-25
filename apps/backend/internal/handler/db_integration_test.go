//go:build integration

package handler_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	code, err := runTests(m)
	if err != nil {
		log.Fatalf("Test setup failed: %v", err)
	}
	os.Exit(code)
}

func runTests(m *testing.M) (int, error) {
	setupContainerEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("openbench_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to start container: %w", err)
	}

	defer func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable&timezone=UTC")
	if err != nil {
		return 0, fmt.Errorf("failed to get connection string: %w", err)
	}

	testDB, err = database.NewDB(config.DefaultDatabaseConfig(connStr))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer testDB.Close()

	migrationsPath, err := findMigrationsPath()
	if err != nil {
		return 0, fmt.Errorf("failed to find migrations path: %w", err)
	}

	migrator, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		return 0, fmt.Errorf("failed to initialize migrator: %w", err)
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return 0, fmt.Errorf("failed to run migrations: %w", err)
	}

	return m.Run(), nil
}

func SetupTestDB() *sqlx.DB {
	if testDB == nil {
		log.Fatal("testDB is not initialized")
	}
	return testDB
}

func CleanTestDB(t testing.TB, db *sqlx.DB) {
	_, err := db.Exec("TRUNCATE TABLE tickets, idempotency_keys RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func setupContainerEnv() {
	uid := os.Getuid()
	podmanSock := fmt.Sprintf("/run/user/%d/podman/podman.sock", uid)
	if _, err := os.Stat(podmanSock); err == nil {
		// Verify read/write socket access permissions before exporting DOCKER_HOST
		if err := syscall.Access(podmanSock, 06); err == nil {
			if os.Getenv("DOCKER_HOST") == "" {
				os.Setenv("DOCKER_HOST", "unix://"+podmanSock)
			}
			os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		}
	}
}

func findMigrationsPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		target := filepath.Join(dir, "migrations")
		if info, err := os.Stat(target); err == nil && info.IsDir() {
			return target, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("migrations directory not found")
}
