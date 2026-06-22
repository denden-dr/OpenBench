# ==============================================================================
# OpenBench Monorepo Makefile
# ==============================================================================
# Coordinates dev services, migrations, unit tests, and E2E browser environments.

.PHONY: dev dev-db-up dev-db-down compose-test-build compose-test-up compose-test-down \
        run-backend run-frontend run-frontend-mock \
        test-backend test-unit test-integration test-frontend \
        test-frontend-mock test-frontend-dev-env test-frontend-test-env \
        test-e2e-mock test-e2e-dev test-e2e-env \
        migrate-up migrate-down migrate-test-up migrate-test-down \
        generate-api-types generate-api-go generate-api-ts \
        fmt clean

# Default: Ryuk enabled by default, developer overrides via env (opt-out)
TESTCONTAINERS_RYUK_DISABLED ?= false
TEST_ENV_COMMAND_TIMEOUT ?= 180
TEST_ENV_READINESS_TIMEOUT ?= 60

# Container tooling defaults to Podman when available, with Docker fallback.
CONTAINER_RUNTIME ?= $(shell if command -v podman >/dev/null 2>&1; then echo podman; elif command -v docker >/dev/null 2>&1; then echo docker; else echo podman; fi)

# Derive COMPOSE from CONTAINER_RUNTIME to prevent mixing Docker and Podman commands.
ifeq ($(CONTAINER_RUNTIME),podman)
  COMPOSE ?= $(shell if command -v podman-compose >/dev/null 2>&1; then echo podman-compose; else echo podman-compose; fi)
else
  COMPOSE ?= $(shell if docker compose version >/dev/null 2>&1; then echo "docker compose"; elif command -v docker-compose >/dev/null 2>&1; then echo docker-compose; else echo docker-compose; fi)
endif

# Fail-fast check to ensure CONTAINER_RUNTIME and COMPOSE use the same engine.
ifneq ($(findstring podman,$(CONTAINER_RUNTIME)),$(findstring podman,$(COMPOSE)))
  ifneq ($(findstring docker,$(CONTAINER_RUNTIME)),$(findstring docker,$(COMPOSE)))
    $(error CONTAINER_RUNTIME ($(CONTAINER_RUNTIME)) and COMPOSE ($(COMPOSE)) must point to the same container engine!)
  endif
endif
COMPOSE_DEV = $(COMPOSE) --env-file apps/backend/.env
COMPOSE_TEST = $(COMPOSE) -f docker-compose-test.yml --profile frontend
DEV_DB_CONTAINER = openbench-postgres-dev
TEST_DB_CONTAINER = openbench-postgres-test

# --- Integrated Development Launch ---
dev:
	@echo "==> Starting development database..."
	$(MAKE) dev-db-up
	@echo "==> Waiting for database health check..."
	@until [ "$$($(CONTAINER_RUNTIME) inspect --format='{{.State.Health.Status}}' $(DEV_DB_CONTAINER) 2>/dev/null)" = "healthy" ]; do \
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
	$(COMPOSE_DEV) up -d

dev-db-down:
	$(COMPOSE_DEV) down

run-backend:
	cd apps/backend && go run ./cmd/api/main.go

run-frontend:
	cd apps/frontend && npm run dev

run-frontend-mock:
	cd apps/frontend && npm run dev:mock

# --- Containerized Test Environment ---
compose-test-build:
	$(COMPOSE_TEST) build

compose-test-up:
	@echo "==> Starting containerized test environment (database, migration, api, frontend)..."
	timeout $(TEST_ENV_COMMAND_TIMEOUT) $(COMPOSE_TEST) up -d

compose-test-down:
	$(COMPOSE_TEST) down

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
	cd apps/backend && TESTCONTAINERS_RYUK_DISABLED=$(TESTCONTAINERS_RYUK_DISABLED) go test -count=1 -tags=integration ./...

# --- Frontend Testing Suite ---
test-frontend:
	cd apps/frontend && npm run test

# --- E2E Browser Testing Suite ---
test-e2e-mock:
	@echo "==> Running Playwright E2E tests in Mock mode..."
	cd apps/e2e && npm run test:e2e:mock

test-e2e-dev:
	@echo "==> Running Playwright E2E tests against Dev environment..."
	cd apps/e2e && npm run test:e2e:dev

test-e2e-env:
	@set -e; \
	wait_for() { \
		name="$$1"; \
		url="$$2"; \
		timeout="$(TEST_ENV_READINESS_TIMEOUT)"; \
		elapsed=0; \
		echo "==> Waiting for $$name to be ready..."; \
		until curl -sf "$$url" >/dev/null 2>&1; do \
			if [ "$$elapsed" -ge "$$timeout" ]; then \
				echo "ERROR: Timed out after $${timeout}s waiting for $$name at $$url" >&2; \
				return 1; \
			fi; \
			sleep 1; \
			elapsed=$$((elapsed + 1)); \
		done; \
	}; \
	echo "==> Booting containerized test environment..."; \
	trap '$(MAKE) -C $(CURDIR) compose-test-down' EXIT INT TERM; \
	$(MAKE) compose-test-up; \
	wait_for backend-api-test http://127.0.0.1:8081/health/readiness; \
	wait_for frontend-web-test http://127.0.0.1:3001/auth/signin; \
	echo "==> Running Playwright E2E tests against test environment..."; \
	cd apps/e2e && npm run test:e2e:env

# --- Backward Compatibility Aliases ---
test-frontend-mock: test-e2e-mock
test-frontend-dev-env: test-e2e-dev
test-frontend-test-env: test-e2e-env

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

# --- API Contract Code Generation ---
generate-api-types: generate-api-go generate-api-ts

generate-api-go:
	@echo "==> Generating Go types from OpenAPI spec..."
	cd apps/backend && go generate ./internal/pkg/api/

generate-api-ts:
	@echo "==> Generating TypeScript types from OpenAPI spec..."
	cd apps/frontend && npm run generate-api-types

# --- Utilities & Cleanup ---
fmt:
	cd apps/backend && go fmt ./...

clean:
	@echo "==> Stopping and cleaning up dev databases..."
	@$(CONTAINER_RUNTIME) stop $(DEV_DB_CONTAINER) || true
	@$(CONTAINER_RUNTIME) rm -f -v $(DEV_DB_CONTAINER) || true
	@echo "==> Stopping and cleaning up test databases..."
	@$(CONTAINER_RUNTIME) stop $(TEST_DB_CONTAINER) || true
	@$(CONTAINER_RUNTIME) rm -f -v $(TEST_DB_CONTAINER) || true
	@echo "==> Tearing down active compose sessions..."
	$(COMPOSE_DEV) down || true
	$(COMPOSE_TEST) down || true
	@if [ "$(CONTAINER_RUNTIME)" = "podman" ]; then \
		$(CONTAINER_RUNTIME) pod rm -f pod_openbench || true; \
		$(CONTAINER_RUNTIME) pod rm -f pod_openbench-test || true; \
	fi
	$(CONTAINER_RUNTIME) volume rm openbench_postgres_data || true
