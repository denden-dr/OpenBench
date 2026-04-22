.PHONY: dev-server dev-client dev test-server test-client test

dev-server:
	cd server && go run cmd/api/main.go

dev-client:
	cd client && npm run dev

dev:
	make -j2 dev-server dev-client

test-server:
	cd server && go test ./...

test-client:
	cd client && npm run test

test:
	make test-server test-client

setup:
	@echo "Setting up project..."
	cd server && go mod tidy
	cd client && npm install
