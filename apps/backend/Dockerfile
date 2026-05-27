FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go

FROM alpine:3.21
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
USER appuser
EXPOSE 3000
CMD ["./main"]
