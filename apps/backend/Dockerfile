# Build stage
FROM golang:1.25-alpine AS builder

# Install ca-certificates and git in case they are needed for fetching dependencies
RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app

# Copy module files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Production stage
FROM alpine:3.20

# Create a non-root group and user
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Switch to non-root user
USER appuser

# Expose Fiber port
EXPOSE 8080

# Run the app
CMD ["./main"]
