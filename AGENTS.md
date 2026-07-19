# AGENTS.md — OpenBench

## Project Identity
Phone/electronics repair shop management app.
Architecture: **Go Web API + React Micro Frontends**
- Backend: Go + Fiber v3 (JSON API only)
- Frontend: React 19 + TypeScript + Vite + Tailwind CSS v4
- Database: PostgreSQL 16 via pgx/v5

## Repository Structure
Monorepo using Go Workspaces (`go.work`) and `pnpm`.

```text
apps/
  webapi/       → Go Backend API (Fiber)
  web-user/     → Public User Portal (React + Vite)
  web-admin/    → Internal Admin Dashboard (React + Vite)
```

## Essential Commands (Run from root)

```bash
make install-api    # go mod tidy in webapi
make install-user   # pnpm install in web-user
make install-admin  # pnpm install in web-admin

make dev-api        # hot-reload (Air) for Go backend
make dev-user       # Vite dev server for public portal
make dev-admin      # Vite dev server for admin dashboard

make test-api       # all unit tests for backend
make build-all      # build Go binary and both React apps
make up / make down # Podman compose start/stop PostgreSQL
make seed           # seed the database (runs in webapi)
make migrate-up     # run migrations
```

## Build prerequisites
- `air` (hot-reload): `go install github.com/air-verse/air@latest`
- No external linter required (use `go vet` natively)
- `migrate` CLI (golang-migrate) if running migrations manually
- `pnpm` (Node.js package manager) for frontends

## Architecture (apps/webapi)

```text
cmd/server/main.go          → wires everything, creates Fiber app, starts API server
cmd/seed/main.go            → standalone seeder binary
config/                     → settings.json (static) + .env (secrets), loaded via viper + validator
internal/
  auth/                     → auth domain (JWT Bearer API auth)
  ticket/                   → service tickets domain
  inventory/                → product inventory domain
  pos/                      → point-of-sale domain
  warranty/                 → warranty claims domain
  models/                   → shared domain models + validation
  database/                 → postgres connection pool, re-entrant TxManager, seeder
  events/                   → async in-process event bus
  apierrors/                → RFC 7807 error handler + stack-wrapped errors
  health/                   → health check handler
  logger/                   → slog logger + middleware
  testutils/                → testcontainers DB setup (integration tests)
```

Every domain module (`internal/<domain>`) follows the same pattern:
- `module.go` → `Module` struct exporting `Handler` and (optional) `QueryRepo`
- `service.go` / `service_test.go` → business logic
- `repository.go` / `repository_integration_test.go` → DB access via `sqlx`
- `handler.go` / `handler_integration_test.go` → JSON API handlers (Fiber)

## Key Conventions

### Config
- `settings.json` holds non-secret defaults (conn pool, timeouts, app name)
- `.env` holds secrets (DB credentials, JWT keys, encryption key)
- All config keys validated with struct tags; `APP_ENCRYPTION_KEY` must be **exactly 32 chars**
- `config.Load()` auto-finds project root by walking up to `go.mod`
- Environment vars override settings.json values

### Database
- PostgreSQL 16 via `pgx/v5` + `sqlx`; connection pool with exponential backoff retry
- **Re-entrant TxManager**: `RunInTx(fn)` propagates tx via context; if a tx is already in context, it reuses it. Use `database.GetQuerier(ctx, db)` in repos to auto-detect tx vs pool.
- SQL builder: `squirrel` (github.com/Masterminds/squirrel)
- Query logging via `sqldblogger` wrapping the pgx driver
- Migrations in `apps/webapi/migrations/` (up/down pairs, `golang-migrate` format)

### HTTP API
- **Fiber v3** (not v2): `fiber.Ctx` is a value type, **not a pointer**
- All routes are JSON API under `/api/v1/`
- API auth: `auth.RequireAuth(cfg, queryRepo)` middleware (JWT Bearer)
- All API errors go through `apierrors.GlobalErrorHandler` → RFC 7807 Problem Details JSON
- Validation errors respect `Accept-Language` header (en/id) for translated messages

### UI Design (React Frontends)
- Both `web-user` and `web-admin` use Tailwind CSS v4 via `@tailwindcss/vite` plugin.
- Custom colors and theme variables are defined in `src/index.css` using `@theme`.
- **Icons**: MUST use Lucide React (`lucide-react`) for frontends.
- **Fonts**: Plus Jakarta Sans (body), Outfit (headings), JetBrains Mono (mono)
- **Colors**: palette defined in `DESIGN.md` — `slate` base, custom `primary`/`secondary`/`accent`/`tertiary`/`danger`
- **Aesthetic**: glassmorphism cards (`bg-white/70 backdrop-blur-xl border border-white/40 shadow-2xl`), rounded corners (`rounded-xl`/`rounded-2xl`), smooth transitions

### Error Handling
- Return domain-specific sentinel errors (e.g., `ticket.ErrTicketNotFound`, `pos.ErrInsufficientStock`) — the global error handler maps them to HTTP statuses
- Use `apierrors.Wrap(err, "context")` or `apierrors.New("message")` for stack-traced errors (logged on 500s)
- Never return bare `errors.New(...)` in handler/service code; use sentinel errors from the domain package

### Testing
- **Integration tests** use build tag `//go:build integration` — run with `make test-integration` (or `make test` inside `apps/webapi`)
- Integration tests spin up PostgreSQL via `testutils.SetupTestDatabase(ctx)` (testcontainers), auto-migrate, return `db` + `teardown`
- Between subtests, use `testutils.CleanTable(db, "table_name")` to reset
- Tests use `testify/assert` + `testify/require`

### Other
- Protobuf validation via `go-playground/validator/v10`
- Password hashing: `golang.org/x/crypto/bcrypt`
- JWT: `github.com/golang-jwt/jwt/v5` (access + refresh token pair)
- Token blacklist uses in-memory cache (`samber/hot`) with periodic DB-backed cleanup worker
- Event bus: in-process async, buffer size 100, used for side effects (e.g., ticket → warranty generation)
- This project uses **Podman**, not Docker. All `docker-compose.yml` files use `podman compose`
