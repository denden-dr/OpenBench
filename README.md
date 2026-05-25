# OpenBench (PhoneFix Admin)

A single-user administrative dashboard for a phone repair business. Track the full intake-to-pickup lifecycle of customer devices with status workflow, payment tracking, and warranty management.

## Project Architecture

```
apps/
├── backend/          Go + Fiber REST API
│   ├── main.go       Entry point, DI wiring, route registration
│   ├── migrations/   PostgreSQL schema migrations (golang-migrate)
│   ├── mocks/        Generated Mockery mocks
│   └── internal/
│       ├── config/       Env-based configuration
│       ├── database/     DB connection + idempotency storage
│       ├── dto/          Request/response types
│       ├── handler/      HTTP handlers
│       ├── middleware/   Error handler + idempotency middleware
│       ├── model/        Domain model & validation
│       ├── repository/   SQL queries via sqlx
│       └── service/      Business logic & lifecycle invariants
└── frontend/         Svelte 5 + SvelteKit (adapter-node)
    └── src/
        ├── routes/       +page.svelte, +layout.svelte
        ├── lib/          Components, stores, mocks
        └── hooks/        API proxy hooks
```

## Core Features

- **Ticket Intake**: Log repairs with customer details, device specs, accessories, pricing, and warranty duration (default 30 days).
- **Status Workflow**: Six stages — `service_in` → `on_process` → `waiting_confirmation` → `fixed` → `picked_up`, plus `cancelled`.
- **Payment Lifecycle**: Moving to `picked_up` auto-records `exit_date`, sets `payment_status=paid`, and computes `warranty_expiry_date` (exit + warranty days).
- **Lifecycle Invariants**: Backend enforces rules — picked_up requires exit_date + paid; non-picked_up cannot have exit_date.
- **Idempotency**: Postgres-backed `X-Idempotency-Key` middleware prevents duplicate mutations. Different payload with same key returns `409 Conflict`.
- **Dashboard Stats**: KPI cards for revenue, active repairs, today's completions, unpaid repairs.
- **Search & Filters**: Filter tickets by customer name, brand, model, or issue.
- **Mock API Mode**: Run frontend without backend via `npm run dev:mock`.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Svelte 5 (Runes), SvelteKit, Tailwind CSS v4, Lucide icons |
| Backend | Go, Fiber v2, sqlx, shopspring/decimal |
| Database | PostgreSQL 16 |
| Testing | Go testing, testify, Mockery, Testcontainers |
| Container | Docker/Podman + Compose |

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker or Podman
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI

### Running Locally

```bash
# 1. Start database + services
make compose-up

# 2. Apply migrations
make migrate-up

# 3. Run backend on :3000
make run-backend

# 4. Run frontend on :5173
make run-frontend
```

Open http://localhost:5173.

### Makefile Targets

| Target | Description |
|--------|-------------|
| `compose-up` / `compose-down` | Start/stop local stack |
| `compose-test-up` / `compose-test-down` | Test container stack |
| `migrate-up` / `migrate-down` | Run/rollback DB migrations |
| `migrate-create NAME=foo` | Create new migration |
| `run-backend` | Start Go API server |
| `run-frontend` | Start SvelteKit dev server (proxies `/api` to backend) |
| `run-frontend-mock` | Start frontend with mock API (no backend needed) |
| `mock-backend` | Regenerate Mockery mocks |
| `test-backend-unit` | Run unit tests |
| `test-backend-integration` | Run integration tests (Testcontainers) |
| `backend-tidy` / `backend-fmt` | Tidy and format Go code |

## Project Graph

The repository maintains a knowledge graph at `graphify-out/`. After code changes:

```bash
graphify update .
```

Open `graphify-out/graph.html` in a browser for an interactive visualization of the codebase architecture.
