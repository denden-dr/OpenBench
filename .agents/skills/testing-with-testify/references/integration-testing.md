# Integration Testing Configurations and Build Tags

Unit tests should have zero external dependencies and run instantly. Integration tests that require databases, networks, or containers should be isolated using Go build tags.

### 1. Add Build Tags
Place `//go:build integration` at the top of the test file:
```go
//go:build integration

package database_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
)
```

### 2. Makefile Configuration
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
