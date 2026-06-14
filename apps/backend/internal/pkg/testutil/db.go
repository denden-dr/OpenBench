package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	migrate "github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	postgresTC "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB wraps the test database container and connection pool
type TestDB struct {
	Container *postgresTC.PostgresContainer
	DB        *database.Database
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
}

var (
	testDBInstance *TestDB
	once           sync.Once
)

// SetupTestDB starts a postgres container, runs migrations, and returns the TestDB instance.
// It uses a singleton pattern to ensure only one container is active per package execution process.
func SetupTestDB() (*TestDB, error) {
	var initErr error
	once.Do(func() {
		ctx := context.Background()

		dbName := "openbench_test"
		dbUser := "postgres"
		dbPass := "postgres"

		// Find migration path dynamically
		migrationPath, err := findMigrationPath()
		if err != nil {
			initErr = fmt.Errorf("failed to find migration path: %w", err)
			return
		}

		// Spin up PostgreSQL container using testcontainers-go
		pgContainer, err := postgresTC.Run(ctx,
			"postgres:16-alpine",
			postgresTC.WithDatabase(dbName),
			postgresTC.WithUsername(dbUser),
			postgresTC.WithPassword(dbPass),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(30*time.Second),
			),
		)
		if err != nil {
			initErr = fmt.Errorf("failed to start postgres container: %w", err)
			return
		}

		// Get host and port
		host, err := pgContainer.Host(ctx)
		if err != nil {
			pgContainer.Terminate(ctx)
			initErr = fmt.Errorf("failed to get container host: %w", err)
			return
		}

		port, err := pgContainer.MappedPort(ctx, "5432")
		if err != nil {
			pgContainer.Terminate(ctx)
			initErr = fmt.Errorf("failed to get mapped port: %w", err)
			return
		}

		// Configure DatabaseConfig
		dbCfg := &config.DatabaseConfig{
			Host:            host,
			Port:            port.Port(),
			User:            dbUser,
			Password:        dbPass,
			DBName:          dbName,
			SSLMode:         "disable",
			MaxOpenConns:    10,
			MaxIdleConns:    2,
			ConnMaxLifetime: 10 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
		}

		// Connect to the test database for migrations using a separate, temporary pool
		migrationDBCfg := *dbCfg
		migrationDBCfg.MaxOpenConns = 1
		migrationDBCfg.MaxIdleConns = 1
		migrationDB, err := database.NewConnection(&migrationDBCfg)
		if err != nil {
			pgContainer.Terminate(ctx)
			initErr = fmt.Errorf("failed to connect to test database for migrations: %w", err)
			return
		}

		// Run migrations using golang-migrate
		driver, err := pgxMigrate.WithInstance(migrationDB.DB.DB, &pgxMigrate.Config{})
		if err != nil {
			migrationDB.Close()
			pgContainer.Terminate(ctx)
			initErr = fmt.Errorf("failed to create migration driver: %w", err)
			return
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://"+migrationPath,
			"pgx", driver,
		)
		if err != nil {
			migrationDB.Close()
			pgContainer.Terminate(ctx)
			initErr = fmt.Errorf("failed to create migrate instance: %w", err)
			return
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			_, _ = m.Close()
			migrationDB.Close()
			pgContainer.Terminate(ctx)
			initErr = fmt.Errorf("failed to run migrations up: %w", err)
			return
		}

		// Close migration tools and connections cleanly to return connection to database
		_, _ = m.Close()
		migrationDB.Close()

		// Connect using the production connection logic to maintain driver parity for the test suite
		db, err := database.NewConnection(dbCfg)
		if err != nil {
			pgContainer.Terminate(ctx)
			initErr = fmt.Errorf("failed to connect to test database: %w", err)
			return
		}

		testDBInstance = &TestDB{
			Container: pgContainer,
			DB:        db,
			DBHost:    host,
			DBPort:    port.Port(),
			DBUser:    dbUser,
			DBPass:    dbPass,
			DBName:    dbName,
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return testDBInstance, nil
}

// Terminate cleans up the container and DB connection
func (t *TestDB) Terminate() {
	if t.DB != nil {
		t.DB.Close()
	}
	if t.Container != nil {
		t.Container.Terminate(context.Background())
	}
	testDBInstance = nil
}

// findMigrationPath traverses up the directory tree starting from the caller's directory to locate the migrations folder
func findMigrationPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		// Traverse up from the location of db.go to find backend/migrations
		dir := filepath.Dir(filename) // testutil
		dir = filepath.Dir(dir)       // pkg
		dir = filepath.Dir(dir)       // internal
		dir = filepath.Dir(dir)       // backend
		migrationDir := filepath.Join(dir, "migrations")
		if _, err := os.Stat(migrationDir); err == nil {
			return migrationDir, nil
		}
	}

	// Fallback using working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path1 := filepath.Join(dir, "apps", "backend", "migrations")
		if _, err := os.Stat(path1); err == nil {
			return path1, nil
		}

		path2 := filepath.Join(dir, "migrations")
		if _, err := os.Stat(path2); err == nil {
			return path2, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find migrations directory")
}

// IntegrationSuite is the base testify suite for database integration tests.
// It automatically handles setting up the container and clean database state for each test.
type IntegrationSuite struct {
	suite.Suite
	TestDB *TestDB
	DB     *database.Database
}

func (s *IntegrationSuite) SetupSuite() {
	tdb, err := SetupTestDB()
	s.Require().NoError(err)
	s.TestDB = tdb
	s.DB = tdb.DB
}

func (s *IntegrationSuite) SetupTest() {
	// Clean all tables in the public schema except the schema migrations tracker table
	ctx := context.Background()
	rows, err := s.DB.DB.QueryContext(ctx,
		"SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'schema_migrations'",
	)
	s.Require().NoError(err)
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		s.Require().NoError(rows.Scan(&table))
		tables = append(tables, table)
	}
	s.Require().NoError(rows.Err())

	if len(tables) > 0 {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(tables, ", "))
		_, err = s.DB.DB.ExecContext(ctx, query)
		s.Require().NoError(err)
	}
}

// CleanupAll is a convenience helper to clean up database state manually when needed
func (s *IntegrationSuite) CleanupAll() {
	s.SetupTest()
}
