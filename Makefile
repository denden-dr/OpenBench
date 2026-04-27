.PHONY: build run fmt tidy db-up db-down db-reset db-logs db-shell db-check migrate-up migrate-down migrate-create

build:
	go build -o bin/api cmd/api/main.go

run:
	air

fmt:
	go fmt ./...

tidy:
	go mod tidy

# --- Database Setup ---

# Detection for podman vs docker
DOCKER_CMD := $(shell which docker 2>/dev/null || which podman 2>/dev/null)
COMPOSE_CMD := $(shell if $(DOCKER_CMD) compose version >/dev/null 2>&1; then echo "$(DOCKER_CMD) compose"; else which docker-compose 2>/dev/null || which podman-compose 2>/dev/null; fi)

db-up:
	@echo "Starting PostgreSQL container using $(COMPOSE_CMD)..."
	@$(COMPOSE_CMD) up -d
	@echo "Waiting for database to be healthy..."
	@until $(COMPOSE_CMD) ps | grep -q "healthy"; do \
		sleep 1; \
	done
	@echo "Database is ready!"

db-down:
	@echo "Stopping PostgreSQL container..."
	@$(COMPOSE_CMD) down

db-reset:
	@echo "Resetting PostgreSQL container and volume..."
	@$(COMPOSE_CMD) down -v
	@$(MAKE) db-up

db-logs:
	@$(COMPOSE_CMD) logs -f postgres

db-shell:
	@$(COMPOSE_CMD) exec postgres psql -U postgres -d openbench

db-check:
	@if ! $(COMPOSE_CMD) ps | grep -q "healthy"; then \
		echo "Error: Database container is not running or not healthy. Run 'make db-up' first."; \
		exit 1; \
	fi

# --- Migrations ---

# Extract DATABASE_URL from .env file
DB_URL=$(shell grep "^DATABASE_URL=" .env | cut -d '=' -f 2- | tr -d '"')

migrate-up: db-check
	@echo "Running migrations..."
	@migrate -path migrations -database "$(DB_URL)" up

migrate-down: db-check
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "$(DB_URL)" down -all

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
