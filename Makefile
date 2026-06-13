.PHONY: dev dev-db-up dev-db-down test-env-build test-env-up test-env-down run-backend run-frontend run-frontend-mock test-frontend test-backend test-unit test-integration clean

# --- Integrated Development Launch ---
dev:
	@echo "Starting development database..."
	$(MAKE) dev-db-up
	@echo "Waiting for database health check..."
	@until [ "$$(podman inspect --format='{{.State.Health.Status}}' openbench-postgres-dev 2>/dev/null)" = "healthy" ]; do \
		sleep 1; \
	done
	@echo "Database is ready. Running migrations..."
	$(MAKE) migrate-up
	@echo "Starting Go backend and SvelteKit frontend..."
	@trap 'kill 0' EXIT; \
	(cd apps/backend && go run ./cmd/api/main.go) & \
	(cd apps/frontend && npm run dev) & \
	wait

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

# --- Local SvelteKit Frontend Commands ---
run-frontend:
	cd apps/frontend && npm run dev

run-frontend-mock:
	cd apps/frontend && npm run dev:mock

test-frontend:
	cd apps/frontend && npm run test

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
	@sleep 2
	@echo "Database is ready. Running migrations..."
	$(MAKE) migrate-up-test
	@echo "Migrations applied. Running integration tests..."
	cd apps/backend && APP_ENV=test go test -count=1 -tags=integration ./...
	@echo "Tearing down test database..."
	podman stop openbench-postgres-test || true
	podman rm -f -v openbench-postgres-test || true

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

migrate-up-test:
	@if [ -f apps/backend/.env.test ]; then \
		export $$(cat apps/backend/.env.test | grep -v '^#' | xargs); \
	fi; \
	migrate -path apps/backend/migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" up

migrate-down-test:
	@if [ -f apps/backend/.env.test ]; then \
		export $$(cat apps/backend/.env.test | grep -v '^#' | xargs); \
	fi; \
	migrate -path apps/backend/migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" down 1

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
