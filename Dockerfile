FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install make, nodejs, npm for templ and tailwind
RUN apk add --no-cache make nodejs npm

# Install templ and tailwind globally
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN npm install -g tailwindcss@3

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build templ and tailwind css
RUN make templ
RUN make tailwind-build

# Build the Go binary
RUN go build -o /app/bin/server ./cmd/server
RUN go build -o /app/bin/seed ./cmd/seed

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/bin/seed .
COPY --from=builder /app/ui/static ./ui/static
COPY --from=builder /app/settings.json .

EXPOSE 3000
CMD ["./server"]
