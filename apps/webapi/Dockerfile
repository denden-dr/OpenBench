FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go binaries
RUN go build -o /app/bin/server ./cmd/server && go build -o /app/bin/seed ./cmd/seed

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/bin/seed .
COPY --from=builder /app/settings.json .

EXPOSE 3000
CMD ["./server"]
