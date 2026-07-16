-include .env
export

.PHONY: up down run build test mod-tidy migrate-up migrate-down seed lint test-integration bench templ tailwind-build tailwind-watch

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
	golangci-lint run ./...

build: templ
	go build -o bin/server cmd/server/main.go

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
