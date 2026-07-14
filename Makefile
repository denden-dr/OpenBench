-include .env
export

.PHONY: up down run build test mod-tidy migrate-up migrate-down seed lint test-integration

up:
	podman compose up -d

down:
	podman compose down

run:
	cd apps/backend && go run ./cmd/api

seed:
	cd apps/backend && go run ./cmd/seed

mod-tidy:
	cd apps/backend && go mod tidy

fmt:
	cd apps/backend && go fmt ./...

lint:
	cd apps/backend && golangci-lint run ./...

build:
	cd apps/backend && go build -o bin/api cmd/api/main.go

test:
	cd apps/backend && go test -v ./...

test-integration:
	cd apps/backend && go test -v -tags=integration ./...

migrate-up:
	migrate -path apps/backend/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	migrate -path apps/backend/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down

