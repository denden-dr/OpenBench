-include apps/webapi/.env
export

.PHONY: up down install-api install-user install-admin dev-api dev-user dev-admin build-all build-api build-user build-admin test-api test-integration migrate-up migrate-down seed

# --- Database / Infrastructure ---
up:
	podman compose up -d

down:
	podman compose down

migrate-up:
	cd apps/webapi && migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	cd apps/webapi && migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down

seed:
	cd apps/webapi && go run ./cmd/seed

# --- Installation ---
install-api:
	cd apps/webapi && go mod tidy

install-user:
	cd apps/web-user && pnpm install

install-admin:
	cd apps/web-admin && pnpm install

# --- Development Servers ---
dev-api:
	cd apps/webapi && air

dev-user:
	cd apps/web-user && pnpm dev

dev-admin:
	cd apps/web-admin && pnpm dev

# --- Testing & Building ---
test-api:
	cd apps/webapi && go test -v ./...

test-integration:
	cd apps/webapi && go test -v -tags=integration ./...

build-api:
	cd apps/webapi && go build -o bin/server ./cmd/server

build-user:
	cd apps/web-user && pnpm build

build-admin:
	cd apps/web-admin && pnpm build

build-all: build-api build-user build-admin
