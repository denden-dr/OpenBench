#!/usr/bin/env bash
set -euo pipefail

COMPOSE_FILE="docker-compose.test.yml"
CONTAINERS=(
  openbench-postgres-test
  openbench-migrate-test
  openbench-seed-test
  openbench-webapi-test
  openbench-web-admin-test
  openbench-web-user-test
)

cleanup() {
  echo "Tearing down E2E test stack..."
  podman compose -f "$COMPOSE_FILE" down -v >/dev/null 2>&1 || true
  podman rm -f "${CONTAINERS[@]}" >/dev/null 2>&1 || true
}

wait_for_postgres() {
  echo "Waiting for postgres-test database to be ready..."
  until podman exec openbench-postgres-test pg_isready -h localhost -U postgres -d openbench_test >/dev/null 2>&1; do
    sleep 1
  done
  sleep 2
}

wait_for_webapi() {
  echo "Waiting for webapi-test server to be ready..."
  until podman exec openbench-webapi-test wget -qO- http://127.0.0.1:3000/health >/dev/null 2>&1; do
    sleep 1
  done
}

trap cleanup EXIT

echo "Cleaning previous E2E test containers..."
podman rm -f "${CONTAINERS[@]}" >/dev/null 2>&1 || true
podman compose -f "$COMPOSE_FILE" down -v >/dev/null 2>&1 || true

echo "Rebuilding webapi image fresh (seed + server)..."
podman build -t localhost/openbench_webapi-test -f apps/webapi/Dockerfile apps/webapi

echo "Spinning up E2E test database..."
podman compose -f "$COMPOSE_FILE" up -d postgres-test
wait_for_postgres

echo "Running database migrations..."
podman compose -f "$COMPOSE_FILE" run --rm migrate-test

echo "Seeding database data..."
podman compose -f "$COMPOSE_FILE" run --rm seed-test

echo "Spinning up webapi, web-admin, and web-user services..."
podman compose -f "$COMPOSE_FILE" up -d --build webapi-test web-admin-test web-user-test
wait_for_webapi

echo "Running Playwright E2E tests..."
cd apps/e2e-testing && pnpm test
