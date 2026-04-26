.PHONY: build run fmt tidy

build:
	go build -o bin/api cmd/api/main.go

run:
	air

fmt:
	go fmt ./...

tidy:
	go mod tidy
