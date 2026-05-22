CONTAINER_TOOL := $(shell which podman 2>/dev/null || which docker 2>/dev/null)
COMPOSE_TOOL := $(shell which podman-compose 2>/dev/null || which docker-compose 2>/dev/null || echo "$(CONTAINER_TOOL) compose")
DB_URL=postgres://postgres:postgres@localhost:5432/openbench?sslmode=disable

# Supabase/Podman configuration
ifeq ($(shell basename $(CONTAINER_TOOL)),podman)
	export DOCKER_HOST ?= unix:///run/user/$(shell id -u)/podman/podman.sock
endif

db-up:
	$(COMPOSE_TOOL) up -d db

db-down:
	$(COMPOSE_TOOL) stop db && $(COMPOSE_TOOL) rm -f db

supabase-up:
	supabase start

supabase-down:
	supabase stop

supabase-status:
	supabase status

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

run-backend:
	cd apps/backend && go run main.go

test-backend:
	cd apps/backend && go test ./... -v
