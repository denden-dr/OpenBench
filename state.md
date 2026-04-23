# Project State — 2026-04-22

## Overview
Initial scaffold for OpenBench completed. Monorepo structure with Go Fiber v3 backend and SvelteKit frontend. Docker Compose setup integrated for containerized development and testing.

## Components

### Backend (`server/`)
- **Status**: Ready / Containerized
- **Framework**: Fiber v3
- **Port**: 3000
- **Features**:
    - Health Check API (`/api/health`)
    - Layered Architecture (Handler, Service, DTO)
    - Environment variable loading with `godotenv`
    - CORS configured for frontend dev server
    - Graceful shutdown
    - **Dockerized**: Multi-stage `Dockerfile` (Go 1.25 + Alpine)

### Frontend (`client/`)
- **Status**: Ready / Home & Authentication UI
- **Framework**: SvelteKit (Svelte 5)
- **Adapter**: Switched to `@sveltejs/adapter-node` for container support
- **Styling**: TailwindCSS v4 with custom brand theme
- **Testing**: Vitest
- **Port**: 5173 (Mapped from 3000 in container)
- **Features**:
	- **High-Fidelity Home Page**: Industrial Minimalism design for mobile repair service.
	- **Authentication Portal**: Polished Login and Registration interfaces following brand aesthetics.
	- **Connectivity Monitoring**: Real-time backend status integrated into diagnostic feed.
	- **Transparency Engine**: Detailed feature grid and documentation mockups.
	- **Tech Data Layout**: custom fonts (Inter, Space Grotesk) and Material Symbols integration.
	- **Dockerized**: Multi-stage `Dockerfile` (Node 22 + Alpine)
		- Optimized to reuse `node_modules` from build stage for consistency.

### Infrastructure
- **Orchestration**: Docker Compose
- **Database**: PostgreSQL 16 (Image: `postgres:16-alpine`)
    - **DB Name**: `pg-openbench`
    - **External Port**: `5433`
- **Automation**: Updated root `Makefile` with `docker-up`, `docker-down`, and `docker-logs`.

## Connectivity
- Frontend-Backend connectivity in Docker is handled by a server-side proxy in SvelteKit (`hooks.server.ts`).
- `BACKEND_URL` points to `http://server:3000` within the Docker network.
- Health check verified: `curl http://localhost:5173/api/health` returns valid backend data.
- Full stack verified: Dashboard correctly displays backend status in production build.
