---
name: configuring-postgres-compose
description: Use when configuring PostgreSQL 16 services in docker-compose.yml or docker-compose-test.yml, ensuring persistent volumes, environment variable configuration, port mapping to avoid conflicts, and safe teardown mechanics.
---

# Configuring Postgres Compose & Container Teardown

## Overview
PostgreSQL in Docker Compose should be configured with security, persistence, isolation, service readiness, and reliable teardown in mind. This means avoiding hardcoded credentials, configuring persistent volumes for development, utilizing different ports for test environments, implementing container health checks, and writing bulletproof cleanup scripts to prevent container, pod, or network lockups (especially under Podman).

## When to Use
- Adding a new database service using Docker Compose.
- Setting up a test database environment alongside a development database.
- Defining volumes and environment variables for PostgreSQL in containers.
- Writing Makefile cleanup targets (`clean`, `test-integration`) that spawn and destroy container resources.

## Core Pattern

### 1. Development `docker-compose.yml` Example
```yaml
version: '3.8'

services:
  postgres-dev:
    image: postgres:16-alpine
    container_name: openbench-postgres-dev
    environment:
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-openbench_secure_password}
      POSTGRES_DB: ${DB_NAME:-openbench_db}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-postgres} -d ${DB_NAME:-openbench_db}"]
      interval: 3s
      timeout: 3s
      retries: 5
    restart: unless-stopped

volumes:
  postgres_data:
```

### 2. Test `docker-compose-test.yml` Example
```yaml
version: '3.8'

services:
  postgres-test:
    image: postgres:16-alpine
    container_name: openbench-postgres-test
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: test_password
      POSTGRES_DB: openbench_test_db
    ports:
      - "5433:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d openbench_test_db"]
      interval: 3s
      timeout: 3s
      retries: 5
    restart: "no"
```

### 3. Safe Teardown & Makefile Integration (Podman / Docker)
Podman clusters compose services inside a single shared Pod (e.g. `pod_openbench`). Standard `down` commands can fail or refuse to remove the bridge network if any container inside the pod is still running or in an improper state.

Always write Makefile targets that explicitly stop and remove containers by name using direct container commands as a safety layer before running compose down.

```make
test-integration:
	# Start database container
	podman-compose -f docker-compose-test.yml up -d postgres-test
	
	# Wait for postgres-test container health
	@until [ "$$(podman inspect --format='{{.State.Health.Status}}' openbench-postgres-test 2>/dev/null)" = "healthy" ]; do \
		sleep 1; \
	done
	
	# Execute tests
	cd apps/backend && APP_ENV=test go test -count=1 -tags=integration ./...
	
	# Targeted teardown to prevent pod/network locks
	podman stop openbench-postgres-test || true
	podman rm -f -v openbench-postgres-test || true

clean:
	@echo "Stopping database containers..."
	@podman stop openbench-postgres-dev || true
	@podman rm -f -v openbench-postgres-dev || true
	@podman stop openbench-postgres-test || true
	@podman rm -f -v openbench-postgres-test || true
	
	@echo "Stopping compose stacks..."
	podman-compose --env-file apps/backend/.env down || true
	podman-compose -f docker-compose-test.yml down || true
	
	@echo "Pruning remaining pods and volumes..."
	@podman pod rm -f pod_openbench pod_openbench-test || true
	podman volume rm openbench_postgres_data || true
```

## Common Mistakes
- **Hanging Networks & Pods**: Relying only on `podman-compose down` when some containers are running or crashed, which causes network cleanup to fail with `"associated containers... network is being used"`. Always stop and force-remove containers by name first.
- **No Healthcheck Gating**: Relying on simple `depends_on` without `condition: service_healthy` check, which causes client applications to crash on start because the database container is active but PostgreSQL inside it is still booting up.
- **Hardcoding Passwords**: Never hardcode database passwords inside `docker-compose.yml`. Use environment variables instead.
- **Port Conflicts**: Mapping both development and test databases to the same external port `5432`. Always use distinct ports (e.g. `5432` for dev, `5433` for test).
