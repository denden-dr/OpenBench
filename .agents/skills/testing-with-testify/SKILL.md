---
name: testing-with-testify
description: Use when writing or structuring Go unit tests, creating mocks with testify, using assert vs require, or setting up test suites with stretchr/testify.
---

# Testing with Testify

## Overview
Go's built-in testing library is powerful, but writing assertions manually can lead to verbose boilerplate. The `github.com/stretchr/testify` library provides friendly assertion functions, mock helpers, and test suite structures that make unit and integration tests cleaner, more readable, and easier to maintain.

To keep tests deterministic and isolated, avoid mutating global OS process variables like `os.Setenv`, use build tags to separate unit and integration suites, and use robust teardown commands to prevent container state conflicts.

## When to Use
- Writing unit and integration tests in Go.
- Defining assertions (comparisons, errors, type assertions).
- Mocking internal services or database interfaces for unit tests.
- Implementing setup/teardown lifecycles for test environments using suites.
- Organizing and separating unit and integration test executions.

## Core Patterns

### 1. Assert vs Require
- **`assert`** (`github.com/stretchr/testify/assert`): Checks the condition and logs the error, but **does not stop** execution of the current test. Use this for independent assertions where subsequent checks still make sense.
- **`require`** (`github.com/stretchr/testify/require`): Checks the condition and **terminates** the test run immediately (calls `t.FailNow()`). Use this when continuing the test makes no sense if this assertion fails (e.g., checking if an error is nil before accessing the returned value).

```go
func TestGetUser(t *testing.T) {
	user, err := GetUserByID(123)
	
	// Terminate immediately if error occurs to avoid nil-pointer panic on user fields
	require.NoError(t, err)
	require.NotNil(t, user)

	// Continue testing other fields if one assertion fails
	assert.Equal(t, 123, user.ID)
	assert.Equal(t, "Alice", user.Name)
}
```

### 2. Environment Isolation in Tests (TD-006)
Never use `os.Setenv` or `os.Unsetenv` inside tests. This mutates the environment variables globally for the whole Go test process, causing race conditions and order dependency when tests run concurrently.
Instead:
- Pass configurations explicitly to constructors rather than reading env files inside modules.
- Use `t.Setenv(...)` when environment variables must be changed. Go automatically cleans up and restores original environment values when the test finishes.

```go
func TestConfigLoader(t *testing.T) {
	// Automatically restored to original value after this test finishes
	t.Setenv("APP_ENV", "test")
	t.Setenv("PORT", "9999")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "9999", cfg.Port)
}
```

### 3. Test Suites and Separating Integration Tests (BE-005)
Unit tests should have zero external dependencies and run instantly. Integration tests that require databases, networks, or containers should be isolated using Go build tags.

1. **Add Build Tags**: Place `//go:build integration` at the top of the test file:
```go
//go:build integration

package database_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
)
```

2. **Makefile Configuration**:
Configure targets to run unit tests and integration tests separately, appending `-count=1` to bypass Go's test result caching. Make sure the teardown section stops and removes containers explicitly to prevent pod/network locks in Podman/Docker.

```make
# Run unit tests only (excludes integration-tagged tests)
test-unit:
	@echo "Running unit tests..."
	cd apps/backend && go test -count=1 ./...

# Provision environment and run integration tests with safe teardown
test-integration:
	podman-compose -f docker-compose-test.yml up -d postgres-test
	@echo "Waiting for database readiness..."
	@until [ "$$(podman inspect --format='{{.State.Health.Status}}' openbench-postgres-test 2>/dev/null)" = "healthy" ]; do \
		sleep 1; \
	done
	cd apps/backend && APP_ENV=test go test -count=1 -tags=integration ./...
	@echo "Tearing down test database..."
	podman stop openbench-postgres-test || true
	podman rm -f -v openbench-postgres-test || true
```

## Common Mistakes
- **Mutating Global Process Environment (TD-006)**: Modifying env vars with `os.Setenv` inside tests, which pollutes the global test runner state and breaks parallel execution. Use `t.Setenv` instead.
- **Polluting Unit Tests with Integration Logic**: Letting integration tests run as part of standard `go test ./...` without build tags, which breaks local developer pipelines if the external DB is off (BE-005).
- **Relying on Cached Test Results**: Forgetting `-count=1`, which allows Go to report success from cached results even if the database configuration has drifted or failed since the last run.
- **Failing to use `.AssertExpectations(t)` on mock**: Forgetting to verify that the mock methods were actually called. Always call this at the end of the test.
- **Using Assert instead of Require for err checks**: Using `assert.NoError` before checking fields of a returned struct pointer. If the pointer is nil, subsequent calls will panic. Always use `require.NoError` first.
