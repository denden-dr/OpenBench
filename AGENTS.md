# OpenBench — Agent Guide

Monorepo with 3 apps under `apps/`:

| app | stack | entrypoint |
|-----|-------|------------|
| `webapi` | Go 1.26, Fiber v3, sqlx+pgx, golang-migrate | `cmd/server/main.go` |
| `web-user` | React 19, Vite, Tailwind v4, oxlint | `src/main.tsx` |
| `web-admin` | React 19, Vite, Tailwind v4, shadcn/ui (base-nova), zustand, react-router, oxlint | `src/main.tsx` |

## Commands

All via `make` from repo root:

| command | what it does |
|---------|-------------|
| `make dev-api` | Air hot-reload on `apps/webapi` (requires `go install github.com/air-verse/air@latest`) |
| `make dev-user` | `pnpm dev` on `web-user` (Vite, port 5173) |
| `make dev-admin` | `pnpm dev` on `web-admin` (Vite, port 5173, proxies `/api` → `localhost:3000`) |
| `make test-api` | `go test -v ./...` (unit only) |
| `make test-integration` | `go test -v -tags=integration ./...` (needs Docker/Podman, uses testcontainers) |
| `make test-admin` | `pnpm test` (vitest, `apps/web-admin`) |
| `make build-all` | Go binary + both Vite builds |
| `make lint-api` | `golangci-lint run` |
| `make lint-user` / `make lint-admin` | `oxlint` |
| `make migrate-up` / `make migrate-down` | golang-migrate via `DB_*` env vars |
| `make seed` | `go run ./cmd/seed` |
| `make up` / `make down` | `podman compose` to manage PostgreSQL |

Frontend build requires `tsc -b` typecheck before `vite build` (`pnpm build` does both).

Use `pnpm` (not npm/yarn) for both frontends. Lint with oxlint, not eslint.

## Config

Two sources merged via viper:
1. `apps/webapi/settings.json` — static defaults (env, port, origins, pool sizes, auth expiry)
2. `apps/webapi/.env` — secrets (copy from `.env.example`)

`APP_ENCRYPTION_KEY` must be exactly 32 characters.

Skip `.env` loading in tests: set `TEST_NO_ENV_FILE=true`.

## Backend Architecture

Server setup: `cmd/server/server.go` → Fiber app with recover, request logger, CORS middleware.

DI wiring: `cmd/server/main.go` — each domain is a `Module` (auth, ticket, warranty, inventory, pos) with handler/repository/service sub-components.

Cross-cutting:
- `internal/database/tx.go` — `TxManager` for atomic multi-repo operations
- `internal/events/bus.go` — `AsyncEventBus` for domain events (used by ticket → warranty flow)
- `internal/apierrors` — RFC 7807 Problem Details error format
- `internal/auth` — JWT (access + refresh) via cookie `access_token` + `Authorization: Bearer` header; rate limiter on login (5/min per IP)

### API routes (`cmd/server/routes_api.go`)

| prefix | auth | notes |
|--------|------|-------|
| `GET /health` | none | public |
| `POST /api/v1/auth/login` | none | rate-limited |
| `POST /api/v1/auth/refresh` | none | rate-limited |
| `POST /api/v1/auth/logout` | none | |
| `GET /api/v1/auth/me` | JWT | |
| `/api/v1/admin/*` | JWT + ADMIN role | all admin routes |

Routes under `admin`: `/services`, `/warranties`, `/claims`, `/products`, `/pos`.

## Go Testing

- Unit tests: no tag needed
- Integration tests: `//go:build integration` build tag + testcontainers-postgres module
- VSCode `gopls` already configured with `-tags=integration` in `.vscode/settings.json`
- Integration test files: `*_integration_test.go` (in `auth/`, `ticket/`)

## Frontend

- `@/` path alias → `src/` in web-admin (both vite and vitest config)
- `web-admin` vitest: jsdom environment, `@testing-library/jest-dom` setup
- `web-user` has no tests or router; `web-admin` has full routing + `services/` API layer
- CSS framework: Tailwind v4 (no `tailwind.config.js` — CSS-based config in `index.css`)
- Icon library: lucide-react
- Auth stores: zustand (`src/stores/`)

## Repo conventions

- Commits are in Indonesian (bahasa)
- Go module path: `github.com/denden-dr/OpenBench`
- Fiber v3 error handler returns `application/problem+json` (RFC 7807)
- `.env.example` has the canonical list of required env vars
- No root `package.json` or workspace — each frontend is standalone
- Integration tests need Docker/Podman running

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

When the user types `/graphify`, invoke the `skill` tool with `skill: "graphify"` before doing anything else.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- Dirty graphify-out/ files are expected after hooks or incremental updates; dirty graph files are not a reason to skip graphify. Only skip graphify if the task is about stale or incorrect graph output, or the user explicitly says not to use it.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
