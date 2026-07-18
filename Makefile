-include .env
export

.PHONY: up down run build test mod-tidy migrate-up migrate-down seed lint test-integration bench templ tailwind-build tailwind-watch test-env-up test-env-down test-e2e test-e2e-ui

up:
	podman compose up -d

down:
	podman compose down

templ:
	templ generate

tailwind-build:
	npx -y tailwindcss@3 -i ./ui/static/css/input.css -o ./ui/static/css/style.css --minify

tailwind-watch:
	npx -y tailwindcss@3 -i ./ui/static/css/input.css -o ./ui/static/css/style.css --watch

run: templ tailwind-build
	go run ./cmd/server

# Development with hot-reload (Go, Templ, Tailwind)
dev:
	air

seed:
	go run ./cmd/seed

mod-tidy:
	go mod tidy

fmt:
	go fmt ./...

lint:
	go vet ./...

build: templ
	go build -o bin/server ./cmd/server

test:
	go test -v ./...

test-integration:
	go test -v -tags=integration ./...

migrate-up:
	migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down

bench:
	go test -bench=. -benchmem ./...

test-env-up:
	podman build -t openbench-test-base:latest .
	podman-compose -f docker-compose.test.yml up -d

test-env-down:
	@echo "Tearing down test environment..."
	podman-compose -f docker-compose.test.yml down -v

test-e2e:
	cd e2e && PLAYWRIGHT_TEST_BASE_URL=http://localhost:3000 npx playwright test

test-e2e-ui:
	cd e2e && PLAYWRIGHT_TEST_BASE_URL=http://localhost:3000 npx playwright test --ui
