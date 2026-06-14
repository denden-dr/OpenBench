---
name: configuring-postgres-compose
description: Use when setting up or debugging PostgreSQL containers, connection pooling, or docker-compose database services. Do not use for MySQL, SQLite, or other databases.
version: 1.0.0
---

# Configuring Postgres Compose & Connection Pooling

## Overview
PostgreSQL in Docker Compose and Go should be configured for security, persistence, isolation, service readiness, and reliable pool management.

## When to Use
- Adding or configuring a PostgreSQL service using Docker/Podman Compose.
- Writing Makefile cleanup targets that spawn and destroy test containers.
- Setting up or tuning the database connection package in Go (e.g. `sqlx` parameters).
- Implementing database retry pings and connection pool stats observability.

## Step-by-Step Instructions

1. **Verify Docker Compose configuration**: Read `assets/docker-compose.yml.template` and apply the configurations for development/testing environments, ensuring container health checks are enabled via `pg_isready`. Audit environment declarations in all compose files to eliminate dead or unused variables (such as unused API URLs).
2. **Add Makefile teardown targets**: Read `references/makefile-targets.md` and add necessary commands to stop and remove containers cleanly, avoiding network bridge locks under Podman.
3. **Configure connection pooling in Go**: Read `assets/db-connection.go.template` and implement connection pooling in your database setup, ensuring connection health is verified with individual context timeouts during startup retries.
4. **Expose database connection statistics**: Expose `DB.Stats()` to track pool exhaustion and monitor database health under loaded states.

## Common Mistakes
- **No Individual Context timeouts**: Creating a single context outside the retry loop, causing later ping retries to fail immediately if the connection takes time.
- **Hanging Networks**: Relying only on `compose down` under Podman. Always stop and force-remove containers by name first.
- **No Observability**: Storing the DB connection globally without exposing `db.DB.Stats()` to track pool exhaustion.
- **No Ping Verification**: Running `sqlx.Open` without verifying connection viability using `.PingContext()`.
- **No Application-Level Readiness Wait**: Using `pg_isready` to verify database health but not waiting for the backend API (`/api/health`) and frontend dev server to be ready before running E2E tests. Makefile targets that orchestrate test environments must include wait steps for all services, not just the database.
