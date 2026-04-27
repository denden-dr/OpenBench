# Implementation Reference: Setup Local PostgreSQL via Docker

This document provides the exact code changes and file additions required to implement the local development PostgreSQL setup as defined in `plan.md`.

## 1. `docker-compose.yml`

Create a new file `docker-compose.yml` in the project root:

```yaml
services:
  postgres:
    image: postgres:17-alpine
    container_name: openbench_pg
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: openbench
    ports:
      # Binds to localhost only to prevent accidental external exposure.
      # Allows overriding host port via DB_PORT environment variable.
      - "127.0.0.1:${DB_PORT:-5432}:5432"
    volumes:
      - openbench_pg_data:/var/lib/postgresql/data
    healthcheck:
      # Ensures the database is accepting connections
      test: ["CMD-SHELL", "pg_isready -U postgres -d openbench"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  openbench_pg_data:
```

## 2. `Makefile`

Add the database setup targets and update the migration targets. Here are the additions and changes to make to the `Makefile`:

```makefile
# --- Database Setup ---

# Detection for podman vs docker
DOCKER_CMD := $(shell which docker 2>/dev/null || which podman 2>/dev/null)
COMPOSE_CMD := $(shell if $$(DOCKER_CMD) compose version >/dev/null 2>&1; then echo "$$(DOCKER_CMD) compose"; else which docker-compose 2>/dev/null || which podman-compose 2>/dev/null; fi)

.PHONY: db-up db-down db-reset db-logs db-shell db-check

db-up:
	@echo "Starting PostgreSQL container using $$(COMPOSE_CMD)..."
	@$$(COMPOSE_CMD) up -d
	@echo "Waiting for database to be healthy..."
	@timeout=60; while ! $$(COMPOSE_CMD) ps | grep -q "healthy"; do \
		timeout=$$((timeout - 1)); \
		if [ $$timeout -le 0 ]; then echo "Error: Timeout waiting for database to become healthy"; exit 1; fi; \
		sleep 1; \
	done
	@echo "Database is ready!"

db-down:
	@echo "Stopping PostgreSQL container..."
	@$$(COMPOSE_CMD) down

db-reset:
	@echo "Resetting PostgreSQL container and volume..."
	@$$(COMPOSE_CMD) down -v
	@$(MAKE) db-up

db-logs:
	@$$(COMPOSE_CMD) logs -f postgres

db-shell:
	@$$(COMPOSE_CMD) exec postgres psql -U postgres -d openbench

db-check:
	@if ! $$(COMPOSE_CMD) ps | grep -q "healthy"; then \
		echo "Error: Database container is not running or not healthy. Run 'make db-up' first."; \
		exit 1; \
	fi

# --- Update existing Migration targets ---
# Make sure DB_URL extraction handles potential quotes correctly
DB_URL=$(shell grep "^DATABASE_URL=" .env | cut -d '=' -f 2- | tr -d '"')

# Add db-check prerequisite
migrate-up: db-check
	@echo "Running migrations..."
	@migrate -path migrations -database "$(DB_URL)" up

# Add db-check prerequisite
migrate-down: db-check
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "$(DB_URL)" down -all

# Add a target to quickly scaffold new migrations
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
```

## 3. `migrations/.gitkeep`

Ensure the migrations directory exists for `golang-migrate` to work properly on fresh checkouts.

Run this command in the project root:

```bash
mkdir -p migrations && touch migrations/.gitkeep
```

## 4. `.env.example`

Update the comment above `DATABASE_URL` to clarify its relationship with Docker:

```env
# Database Configuration
# NOTE: These default credentials match the local docker-compose.yml configuration.
# Change them ONLY if you have customized your setup via a docker-compose.override.yml file.
# Format: postgres://user:password@host:port/dbname?sslmode=disable
DATABASE_URL=postgres://postgres:postgres@localhost:5432/openbench?sslmode=disable
```

## 5. `.gitignore`

Append the Docker override file to the `.gitignore`:

```gitignore
# Docker
docker-compose.override.yml
```

## 6. `README.md`

Add a Local Development section to clearly instruct developers:

```markdown
## Local Development

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Go](https://go.dev/doc/install) (1.21+)
- Make
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI

### Quick Start
1. Setup environment variables:
   ```bash
   cp .env.example .env
   ```
2. Start the local database:
   ```bash
   make db-up
   ```
3. Run the migrations:
   ```bash
   make migrate-up
   ```
4. Start the application with hot-reloading:
   ```bash
   make run
   ```

### Database Management
We use a Docker container for local development to ensure a consistent PostgreSQL environment.

| Command | Action |
|---------|--------|
| `make db-up` | Starts the PostgreSQL database container in the background. |
| `make db-down` | Stops and removes the container (data is preserved in a volume). |
| `make db-reset` | **Destroys the database volume** and recreates a fresh, empty container. |
| `make db-logs` | Streams the database container logs for debugging. |
| `make db-shell` | Opens an interactive `psql` shell as the `postgres` user. |
```
