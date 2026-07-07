include .env
export

.PHONY: up down run build test mod-tidy migrate-up migrate-down

up:
	podman compose up -d

down:
	podman compose down

run:
	cd apps/backend && go run cmd/api/main.go

mod-tidy:
	cd apps/backend && go mod tidy

build:
	cd apps/backend && go build -o bin/api cmd/api/main.go

test:
	cd apps/backend && go test -v ./...

migrate-up:
	migrate -path apps/backend/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	migrate -path apps/backend/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down
