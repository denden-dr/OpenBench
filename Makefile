CONTAINER_TOOL := $(shell which podman 2>/dev/null || which docker 2>/dev/null)
COMPOSE_TOOL := $(shell which podman-compose 2>/dev/null || which docker-compose 2>/dev/null || echo "$(CONTAINER_TOOL) compose")
DB_URL=postgres://postgres:postgres@localhost:5432/openbench?sslmode=disable

# Supabase/Podman configuration
ifeq ($(shell basename $(CONTAINER_TOOL)),podman)
	export DOCKER_HOST ?= unix:///run/user/$(shell id -u)/podman/podman.sock
endif

compose-up:
	$(COMPOSE_TOOL) up -d --build

compose-down:
	$(COMPOSE_TOOL) down

compose-test-up:
	$(COMPOSE_TOOL) -f compose.test.yaml up -d --build

compose-test-down:
	$(COMPOSE_TOOL) -f compose.test.yaml down

migrate-up:
	migrate -path apps/backend/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path apps/backend/migrations -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir apps/backend/migrations -seq $(NAME)

mock-backend:
	cd apps/backend && go run github.com/vektra/mockery/v2

backend-tidy:
	cd apps/backend && go mod tidy

backend-fmt:
	cd apps/backend && go fmt ./... && goimports -w ./...

test-backend-unit:
	cd apps/backend && go test ./... -v

test-backend-integration:
	cd apps/backend && go test -tags=integration ./... -v

.PHONY: run-db run-backend run-frontend up down

run-db: compose-up

run-backend:
	cd apps/backend && go run main.go

run-frontend-mock:
	cd apps/frontend && npm run dev:mock

run-frontend:
	cd apps/frontend && npm run dev

up: run-db
	$(MAKE) -j 2 run-backend run-frontend

down: compose-down
	@echo "Stopping backend and frontend processes..."
	@pkill -f "go run main.go" || true
	@pkill -f "vite dev" || true
