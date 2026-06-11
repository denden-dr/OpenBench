---
name: secure-containerization
description: Use when containerizing applications (Go, SvelteKit, Node.js) with Docker or Podman, ensuring non-root execution, pinned image tags, multi-stage builds, and dockerignore configurations.
---

# Secure Containerization

## Overview
Container security is crucial for production deployments. Running containers as `root` or baking sensitive credentials into image layers creates severe security risks. Safe containerization relies on multi-stage builds, pinned stable base tags, non-root user execution, and strict ignore patterns.

## When to Use
- Writing or updating a Dockerfile for backend or frontend applications.
- Setting up or reviewing container build files.
- Configuring `.dockerignore` files to prevent secret leakage.

## Core Patterns

### 1. The `.dockerignore` File
To prevent local `.env` files, build caches, git history, or local `node_modules` from leaking into docker image layers, always place a `.dockerignore` file in the service context root.

Example backend `.dockerignore`:
```ignore
.env
.env.test
.git
tmp
dist
node_modules
```

### 2. Multi-Stage Builds & Non-Root User Execution
Separate compilation dependencies from runtime dependencies. In the final stage, configure the container to run as a restricted, non-root user.

#### Go Backend Dockerfile Example
```dockerfile
# Stage 1: Build binary using official pinned stable Go image
FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a fully static, optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Stage 2: Minimal runtime using stable Alpine Linux
FROM alpine:3.20

# Create dedicated non-root user and group
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/main .

# Switch context to the non-root user (BE-007)
USER appuser

EXPOSE 8080

CMD ["./main"]
```

#### SvelteKit / Node.js Frontend Dockerfile Example
```dockerfile
# Stage 1: Build application assets
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Stage 2: Pinned production runtime
FROM node:20-alpine

WORKDIR /app
RUN chown node:node /app

# Switch context to node's pre-existing non-root user
USER node

# Copy manifests and install production-only dependencies
COPY --chown=node:node package*.json ./
RUN npm ci --omit=dev

# Copy built artifacts with correct node user ownership
COPY --chown=node:node --from=builder /app/build ./build

EXPOSE 3000
ENV PORT=3000
ENV NODE_ENV=production

CMD ["node", "build"]
```

## Common Mistakes
- **No `.dockerignore`**: Leaving local `.env` files in the build folder, causing them to be copied by `COPY . .` and permanently baked into the image history.
- **Running as Root**: Leaving the default `root` user context. If a remote code execution vulnerability occurs, the attacker gains full root privileges inside the container, easing breakouts to the host.
- **Mutable Base Tags**: Using `latest` or generic tags (e.g., `FROM golang:latest` or `FROM node:alpine`). When base images update upstream, your build pipeline can break or inherit unexpected security bugs.
- **Baking Node Modules**: Copying local development `node_modules` instead of running `npm ci --omit=dev` inside the builder context, leading to bloated images and environment mismatch issues.
