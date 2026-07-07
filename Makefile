.PHONY: up down run build test mod-tidy

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
