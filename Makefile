# ==============================================================================
# OpenBench Monorepo Makefile
# ==============================================================================
# Coordinates dev services, migrations, unit tests, and E2E browser environments.

.PHONY: dev dev-db-up dev-db-down compose-test-build compose-test-up compose-test-down \
        run-backend run-frontend run-frontend-mock \
        test-backend test-unit test-integration test-frontend \
        test-frontend-mock test-frontend-dev-env test-frontend-test-env \
        migrate-up migrate-down migrate-test-up migrate-test-down \
        fmt clean

# --- Integrated Development Launch ---
dev:
	@echo "==> Starting development database..."
	$(MAKE) dev-db-up
	@echo "==> Waiting for database health check..."
	@until [ "$$(podman inspect --format='{{.State.Health.Status}}' openbench-postgres-dev 2>/dev/null)" = "healthy" ]; do \
		sleep 1; \
	done
	@echo "==> Running database migrations..."
	$(MAKE) migrate-up
	@echo "==> Starting backend and frontend concurrently..."
	@trap 'kill 0' EXIT; \
	(cd apps/backend && go run ./cmd/api/main.go) & \
	(cd apps/frontend && npm run dev) & \
	wait

# --- Local Services (Dev) ---
dev-db-up:
	podman-compose --env-file apps/backend/.env up -d

dev-db-down:
	podman-compose --env-file apps/backend/.env down

run-backend:
	cd apps/backend && go run ./cmd/api/main.go

run-frontend:
	cd apps/frontend && npm run dev

run-frontend-mock:
	cd apps/frontend && npm run dev:mock

# --- Containerized Test Environment ---
compose-test-build:
	podman-compose -f docker-compose-test.yml --profile frontend build

compose-test-up:
	podman-compose -f docker-compose-test.yml --profile frontend up -d postgres-test
	@echo "==> Waiting for test database to be ready..."
	@until [ "$$(podman inspect --format='{{.State.Health.Status}}' openbench-postgres-test 2>/dev/null)" = "healthy" ]; do \
		sleep 1; \
	done
	@sleep 2
	@echo "==> Database ready. Running test migrations..."
	$(MAKE) migrate-test-up
	@echo "==> Starting test environment frontend services..."
	podman-compose -f docker-compose-test.yml --profile frontend up -d

compose-test-down:
	podman-compose -f docker-compose-test.yml --profile frontend down

# --- Backend Testing Suite ---
test-backend: test-unit test-integration

test-unit:
	@echo "==> Verifying code formatting..."
	@if [ -n "$$(gofmt -l apps/backend/cmd apps/backend/internal)" ]; then \
		echo "Files not formatted with gofmt:"; \
		gofmt -l apps/backend/cmd apps/backend/internal; \
		exit 1; \
	fi
	@echo "==> Running unit tests..."
	cd apps/backend && go test -count=1 ./...

test-integration:
	@echo "==> Running integration tests using Testcontainers..."
	cd apps/backend && go test -count=1 -tags=integration ./...

# --- Frontend Testing Suite ---
test-frontend:
	cd apps/frontend && npm run test

test-frontend-mock:
	@echo "==> Running Playwright E2E tests in Mock mode..."
	cd apps/frontend && npm run test:e2e:mock

test-frontend-dev-env:
	@echo "==> Running Playwright E2E tests against Dev environment..."
	cd apps/frontend && npm run test:e2e:dev

test-frontend-test-env:
	@echo "==> Booting containerized test environment..."
	$(MAKE) compose-test-up
	@echo "==> Running Playwright E2E tests against test environment..."
	cd apps/frontend && npm run test:e2e:env || ( $(MAKE) compose-test-down && exit 1 )
	@echo "==> Tearing down containerized test environment..."
	$(MAKE) compose-test-down

# --- Database Migrations ---
migrate-up:
	@if [ -f apps/backend/.env ]; then \
		export $$(cat apps/backend/.env | grep -v '^#' | xargs); \
	fi; \
	migrate -path apps/backend/migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" up

migrate-down:
	@if [ -f apps/backend/.env ]; then \
		export $$(cat apps/backend/.env | grep -v '^#' | xargs); \
	fi; \
	migrate -path apps/backend/migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" down 1

migrate-test-up:
	@if [ -f apps/backend/.env.test ]; then \
		export $$(cat apps/backend/.env.test | grep -v '^#' | xargs); \
	fi; \
	migrate -path apps/backend/migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" up

migrate-test-down:
	@if [ -f apps/backend/.env.test ]; then \
		export $$(cat apps/backend/.env.test | grep -v '^#' | xargs); \
	fi; \
	migrate -path apps/backend/migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" down 1

# --- Utilities & Cleanup ---
fmt:
	cd apps/backend && go fmt ./...

clean:
	@echo "==> Stopping and cleaning up dev databases..."
	@podman stop openbench-postgres-dev || true
	@podman rm -f -v openbench-postgres-dev || true
	@echo "==> Stopping and cleaning up test databases..."
	@podman stop openbench-postgres-test || true
	@podman rm -f -v openbench-postgres-test || true
	@echo "==> Tearing down active compose sessions..."
	podman-compose --env-file apps/backend/.env down || true
	podman-compose -f docker-compose-test.yml --profile frontend down || true
	@podman pod rm -f pod_openbench || true
	@podman pod rm -f pod_openbench-test || true
	podman volume rm openbench_postgres_data || true
