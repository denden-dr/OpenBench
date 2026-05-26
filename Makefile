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
	migrate -path apps/backend/migrations -database "$(DB_URL)" down $(or $(STEPS),1)

migrate-create:
	migrate create -ext sql -dir apps/backend/migrations -seq $(NAME)

mock-backend:
	cd apps/backend && mockery

backend-tidy:
	cd apps/backend && go mod tidy

backend-fmt:
	cd apps/backend && go fmt ./... && goimports -w ./...

test-backend-unit:
	cd apps/backend && go test ./... -v

test-backend-integration:
	cd apps/backend && go test -tags=integration ./... -v

BACKEND_PID := /tmp/openbench-backend.pid
FRONTEND_PID := /tmp/openbench-frontend.pid

.PHONY: run-backend run-frontend up down start stop

run-backend:
	cd apps/backend && go build -o bin/tmp-server main.go && ./bin/tmp-server & echo $$! > $(BACKEND_PID) && wait

run-frontend-mock:
	cd apps/frontend && npm run dev:mock

run-frontend:
	cd apps/frontend && npm run dev & echo $$! > $(FRONTEND_PID) && wait

up start:
	$(MAKE) compose-up
	$(MAKE) -j 2 run-backend run-frontend

down stop:
	@echo "Stopping all services..."
	$(COMPOSE_TOOL) down 2>/dev/null || true
	@if [ -f $(BACKEND_PID) ]; then kill $$(cat $(BACKEND_PID)) 2>/dev/null && rm $(BACKEND_PID) || true; fi
	@if [ -f $(FRONTEND_PID) ]; then kill $$(cat $(FRONTEND_PID)) 2>/dev/null && rm $(FRONTEND_PID) || true; fi
	rm -f apps/backend/bin/tmp-server
	@echo "All services stopped"
