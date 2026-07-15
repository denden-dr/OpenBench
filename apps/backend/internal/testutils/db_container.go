package testutils

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	testcontainerpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupTestDatabase spins up a PostgreSQL container, runs migrations, and returns a sqlx.DB connection pool along with a teardown function.
func SetupTestDatabase(ctx context.Context) (*sqlx.DB, func(), error) {
	dbName := "testdb"
	dbUser := "postgres"
	dbPassword := "postgres"

	// Spin up the Postgres container
	postgresContainer, err := testcontainerpostgres.Run(ctx,
		"docker.io/library/postgres:16-alpine",
		testcontainerpostgres.WithDatabase(dbName),
		testcontainerpostgres.WithUsername(dbUser),
		testcontainerpostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get container port: %w", err)
	}

	// Create DB config
	dbCfg := config.DBConfig{
		Host:            host,
		Port:            port.Port(),
		User:            dbUser,
		Password:        dbPassword,
		Name:            dbName,
		SSLMode:         "disable",
		MaxConns:        10,
		MinConns:        2,
		MaxRetries:      5,
		RetryBaseDelay:  100 * time.Millisecond,
		RetryMaxDelay:   1 * time.Second,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	// Establish connection pool
	db, err := database.NewPostgresDB(dbCfg)
	if err != nil {
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Run migrations
	_, filename, _, _ := runtime.Caller(0)
	testutilsDir := filepath.Dir(filename)
	migrationsDir := filepath.Join(testutilsDir, "..", "..", "migrations")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		db.Close()
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		db.Close()
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		db.Close()
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	teardown := func() {
		db.Close()
		postgresContainer.Terminate(context.Background())
	}

	return db, teardown, nil
}

// CleanTable clears all data from a table (useful between test runs)
func CleanTable(db *sqlx.DB, tableName string) error {
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName))
	return err
}
