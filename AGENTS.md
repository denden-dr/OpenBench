# AGENTS.md — OpenBench

## Project Identity
Phone/electronics repair shop management app. GOTTH stack: **Go + Templ + Tailwind CSS + HTMX**, plus Alpine.js for micro client-side state.

## Essential Commands

```bash
make dev            # hot-reload (Air): watches .go, .templ, .css — recompiles everything
make run            # one-shot: templ → tailwind → go run
make test           # all unit tests: go test -v ./...
make test-integration # integration tests (requires Docker/Podman for testcontainers): go test -v -tags=integration ./...
make lint           # go vet ./...
make fmt            # go fmt ./...
make templ          # templ generate only
make tailwind-build # npx tailwindcss@3 build only
make up / make down # Podman compose start/stop PostgreSQL
make seed           # seed the database
make build          # templ → go build -o bin/server

# E2E (needs test env running)
make test-env-up    # build Docker image + start test stack (migrate → seed → app)
make test-e2e       # npx playwright test (from e2e/ dir)
make test-e2e-ui    # Playwright UI mode
make test-env-down  # tear down test stack
```

## Build prerequisites
- `templ` CLI (codegen): `go install github.com/a-h/templ/cmd/templ@latest`
- `air` (hot-reload): `go install github.com/air-verse/air@latest`
- No external linter required (`make lint` uses `go vet`, which ships with Go)
- `migrate` CLI (golang-migrate) if running migrations manually
- Node.js + npm (for npx tailwindcss@3)

## Architecture

```
cmd/server/main.go          → wires everything, creates Fiber app, starts server
cmd/seed/main.go             → standalone seeder binary
config/                      → settings.json (static) + .env (secrets), loaded via viper + validator
internal/
  auth/                      → auth domain (JWT, cookie-based web auth, bearer API auth)
  ticket/                    → service tickets domain
  inventory/                 → product inventory domain
  pos/                       → point-of-sale domain
  warranty/                  → warranty claims domain
  models/                    → shared domain models + validation
  database/                  → postgres connection pool, re-entrant TxManager, seeder
  events/                    → async in-process event bus
  apierrors/                 → RFC 7807 error handler + stack-wrapped errors
  health/                    → health check handler
  logger/                    → slog logger + middleware
  testutils/                 → testcontainers DB setup (integration tests)
ui/
  static/css/{input,style}.css → tailwind input & build output (style.css is gitignored)
  static/fonts/              → self-hosted .woff2 (Plus Jakarta Sans, Outfit, JetBrains Mono)
  views/layouts/             → main.templ (base layout)
  views/pages/               → full-page Templ components
  views/components/          → reusable Templ fragments (HTMX responses)
```

Every domain module (`internal/<domain>`) follows the same pattern:
- `module.go` → `Module` struct exporting `Handler`, `WebHandler`, and (optional) `QueryRepo`
- `service.go` / `service_test.go` → business logic
- `repository.go` / `repository_integration_test.go` → DB access via `sqlx`
- `handler.go` / `handler_integration_test.go` → JSON API handlers (Fiber)
- `web_handler.go` → Templ/HTMX web handlers

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
- Migrations in `migrations/` (up/down pairs, `golang-migrate` format)

### HTTP & Rendering
- **Fiber v3** (not v2): `fiber.Ctx` is a value type, **not a pointer**
- **Two route groups**: JSON API (`/api/v1/`) and Web UI (Templ/HTMX, cookie auth)
- API auth: `auth.RequireAuth(cfg, queryRepo)` middleware (JWT Bearer)
- Web auth: `auth.RequireWebAuth(cfg, queryRepo)` middleware (cookie-based)
- HTML rendering: `utils.Render(c, component)` — wraps Templ component in base layout
- **HTMX redirects**: server returns `HX-Redirect` header, not HTTP 3xx
- **HTMX errors**: return only the fragment (e.g., `ErrorBanner("message")`), not full page
- All API errors go through `apierrors.GlobalErrorHandler` → RFC 7807 Problem Details JSON
- Validation errors respect `Accept-Language` header (en/id) for translated messages

### UI Design
- **Icons**: MUST use `github.com/dimmerz92/go-icons/lucide` (Lucide), never raw `<svg>`
- **Fonts**: self-hosted — Plus Jakarta Sans (body), Outfit (headings), JetBrains Mono (mono)
- **Colors**: palette defined in `DESIGN.md` — `slate` base, custom `primary`/`secondary`/`accent`/`tertiary`/`danger`
- **Aesthetic**: glassmorphism cards (`bg-white/70 backdrop-blur-xl`), rounded corners (`rounded-xl`/`rounded-2xl`), smooth transitions

### Error Handling
- Return domain-specific sentinel errors (e.g., `ticket.ErrTicketNotFound`, `pos.ErrInsufficientStock`) — the global error handler maps them to HTTP statuses
- Use `apierrors.Wrap(err, "context")` or `apierrors.New("message")` for stack-traced errors (logged on 500s)
- Never return bare `errors.New(...)` in handler/service code; use sentinel errors from the domain package

### Testing
- **Integration tests** use build tag `//go:build integration` — run with `make test-integration`
- Integration tests spin up PostgreSQL via `testutils.SetupTestDatabase(ctx)` (testcontainers), auto-migrate, return `db` + `teardown`
- Between subtests, use `testutils.CleanTable(db, "table_name")` to reset
- Tests use `testify/assert` + `testify/require`
- E2E tests in `e2e/` use Playwright; target running app (default `http://localhost:3000`)
- Test env uses **Podman** (`podman compose` / `podman-compose`), port 5433 for test DB

### Code Generation
- Templ generates `*_templ.go` files — these are **gitignored**, must be regenerated with `make templ` before `go build`
- Tailwind output (`ui/static/css/style.css`) is **gitignored** — must run `make tailwind-build` before any run/build

### Other
- Protobuf validation via `go-playground/validator/v10`
- Password hashing: `golang.org/x/crypto/bcrypt`
- JWT: `github.com/golang-jwt/jwt/v5` (access + refresh token pair)
- Token blacklist uses in-memory cache (`samber/hot`) with periodic DB-backed cleanup worker
- Event bus: in-process async, buffer size 100, used for side effects (e.g., ticket → warranty generation)
- This project uses **Podman**, not Docker. All `docker-compose.yml` files use `podman compose`
