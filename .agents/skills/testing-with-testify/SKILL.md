---
name: testing-with-testify
description: Use when Go backend unit or database integration tests are failing, slow, resource-leaking, require mock dependencies, or need mock auto-generation with mockery.
version: 2.0.0
---

# Go Testing with Testify, Testcontainers & Mockery

## Overview
Go backend tests must be readable, isolated, and hermetic. We use `stretchr/testify` for assertions/suites, `testcontainers-go` (via `testutil`) for Postgres integration tests, and `mockery` for auto-generating interface mocks.

## When to Use
- Writing unit or integration tests for Go handlers, services, and repositories.
- Verifying database schema changes or SQL query behaviors.
- Auto-generating or regenerating interface mocks after signature changes.

### When NOT to Use
- Testing frontend UI interactions (use Playwright instead).

---

## Core Pattern: Database Integration Test

### Before (Anti-Pattern: Hardcoded Ports & Unmanaged State)
```go
// ❌ BAD: Depends on pre-existing db, uses static port 5433, no automatic teardown
func TestSaveUser(t *testing.T) {
	db, _ := sql.Open("postgres", "postgres://postgres:pass@localhost:5433/test_db")
	err := db.QueryRow("INSERT INTO users ...")
}
```

### After (Best Practice: Isolated Testcontainers Suite)
```go
//go:build integration

package database_test

import (
	"log"
	"os"
	"testing"
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
)

type UserRepoTestSuite struct {
	testutil.IntegrationSuite
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}

func (s *UserRepoTestSuite) TestSaveUser() {
	ctx := context.Background()
	_, err := s.DB.DB.ExecContext(ctx, "INSERT INTO users...")
	s.Require().NoError(err)
}

// MANDATORY: Teardown container once for all tests in this package.
func TestMain(m *testing.M) {
	tdb, err := testutil.SetupTestDB()
	if err != nil {
		log.Fatalf("Failed to setup integration test database: %v", err)
	}

	code := m.Run()
	tdb.Terminate()
	os.Exit(code)
}
```

---

## Mock Generation with Mockery

For auto-generating interface mocks, read `references/mockery-generation.md`. Key steps:
1. Add `//go:generate mockery --name=<Interface> --output=mocks --outpkg=mocks` directive.
2. Run `go generate ./...` to regenerate.
3. Use `mocks.NewRepository(t)` in tests — cleanup assertions are automatic.

---

## Quick Reference: Assertions

| Scenario | Assertion Method | Behavior on Failure |
|---|---|---|
| **Critical Setup** (DB connects, error checks) | `require.NoError(t, err)` | Stops test execution immediately |
| **Object Lifecycle** (Pointers) | `require.NotNil(t, obj)` | Stops test execution (prevents nil pointer panics) |
| **Value Verification** (Fields, count) | `assert.Equal(t, expected, actual)` | Logs failure, continues executing test |
| **Boolean State** (Validity checks) | `assert.True(t, condition)` | Logs failure, continues executing test |

## Common Mistakes
- **Ryuk Should Stay Enabled by Default**: Ryuk is a safety net that cleans up containers when the test process crashes or is killed. Do not disable it globally **at any layer** — not in Go helper code, not in Makefile targets, and not in CI scripts. Only disable for runtimes that technically cannot run it (rootless Podman), and make it opt-in by **allowing the caller** to set `TESTCONTAINERS_RYUK_DISABLED=true` in their shell environment. Makefile targets should declare a default and pass it through:
  ```makefile
  # Top of Makefile: Ryuk enabled by default, developer overrides via env
  TESTCONTAINERS_RYUK_DISABLED ?= false

  test-integration:
  	cd apps/backend && TESTCONTAINERS_RYUK_DISABLED=$(TESTCONTAINERS_RYUK_DISABLED) go test -count=1 -tags=integration ./...
  ```
  If Ryuk must be disabled, also add a Makefile cleanup target as fallback:
  ```makefile
  clean-test-containers:
  	docker rm -f $$(docker ps -q --filter "label=org.testcontainers=true") 2>/dev/null || true
  ```
- **Container Orphans on Crash**: When Ryuk is disabled and the test process is killed before `TestMain` calls `tdb.Terminate()`, containers are left running indefinitely. Always keep Ryuk enabled as the primary cleanup mechanism; `Terminate()` in `TestMain` is the secondary cleanup.
- **Migration Pool Contention**: Running migrations using the same pool as the test code (`pgxMigrate.WithInstance(db.DB.DB, ...)`) holds a connection throughout the migration, reducing effective pool capacity for tests. Always create a separate database connection for migrations and close it before tests begin:
  ```go
  // ✓ Good: separate connection for migrations
  migrationDB, _ := sql.Open("postgres", dsn)
  defer migrationDB.Close()
  driver, _ := pgxMigrate.WithInstance(migrationDB, &pgxMigrate.Config{})
  m, _ := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
  m.Up()
  // main pool remains fully available for tests
  ```
- **Unmet Mock Expectations**: Forgetting to assert mock expectations. Use `mocks.New<Interface>(t)` which registers `t.Cleanup` automatically, or manually add `t.Cleanup(func() { mockSQL.ExpectationsWereMet() })`.
- **Out-of-Sync Mocks**: Modifying an interface and forgetting to run `go generate ./...`.
- **Nil Pointer Panic**: Using `assert.NoError` on an error return and immediately dereferencing the return variable. Always use `require.NoError`.
- **State Pollution**: Mutating the system environment via `os.Setenv` (use `t.Setenv`) or modifying database rows without table truncation.
