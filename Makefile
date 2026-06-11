.PHONY: dev-db-up dev-db-down test-env-build test-env-up test-env-down run-backend test-backend test-unit test-integration clean

# --- Development Environment (Database only) ---
dev-db-up:
	podman-compose --env-file apps/backend/.env up -d

dev-db-down:
	podman-compose --env-file apps/backend/.env down

# --- Test Environment (Database, Backend, Frontend Containerized) ---
test-env-build:
	podman-compose -f docker-compose-test.yml build

test-env-up:
	podman-compose -f docker-compose-test.yml up -d

test-env-down:
	podman-compose -f docker-compose-test.yml down

# --- Local Go Backend Commands ---
run-backend:
	cd apps/backend && go run ./cmd/api/main.go

test-backend: test-unit test-integration

test-unit:
	@echo "Verifying code formatting..."
	@if [ -n "$$(gofmt -l apps/backend/cmd apps/backend/internal)" ]; then \
		echo "Files not formatted with gofmt:"; \
		gofmt -l apps/backend/cmd apps/backend/internal; \
		exit 1; \
	fi
	@echo "Running unit tests..."
	cd apps/backend && go test -count=1 ./...

test-integration:
	podman-compose -f docker-compose-test.yml up -d postgres-test
	@echo "Waiting for test database to be ready..."
	@until [ "$$(podman inspect --format='{{.State.Health.Status}}' openbench-postgres-test 2>/dev/null)" = "healthy" ]; do \
		sleep 1; \
	done
	@echo "Database is ready. Running integration tests..."
	cd apps/backend && APP_ENV=test go test -count=1 -tags=integration ./...
	@echo "Tearing down test database..."
	podman stop openbench-postgres-test || true
	podman rm -f -v openbench-postgres-test || true

# --- Cleanup ---
clean:
	@echo "Stopping and removing dev database..."
	@podman stop openbench-postgres-dev || true
	@podman rm -f -v openbench-postgres-dev || true
	@echo "Stopping and removing test database..."
	@podman stop openbench-postgres-test || true
	@podman rm -f -v openbench-postgres-test || true
	@echo "Tearing down compose environments..."
	podman-compose --env-file apps/backend/.env down || true
	podman-compose -f docker-compose-test.yml down || true
	@podman pod rm -f pod_openbench || true
	@podman pod rm -f pod_openbench-test || true
	podman volume rm openbench_postgres_data || true

fmt:
	cd apps/backend && go fmt ./...
