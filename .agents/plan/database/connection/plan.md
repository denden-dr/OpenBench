# Plan: PostgreSQL Connection & Health Integration

> **Goals**
> 1. Connect to the local PostgreSQL database with a fully-configured, best-practice connection pool.
> 2. Surface database reachability as a parameter in the `/health` endpoint.

---

## A. Logical Requirements

### Problem Statement

The project has a local PostgreSQL container (via `docker-compose.yml`) and a `DATABASE_URL` in `.env`, but **no Go code to connect to it**. The `pkg/config` and `pkg/database` packages do not exist yet. Additionally, the `/health` endpoint in `internal/handlers/health.go` returns a static `"status": "ok"` regardless of whether the database is actually reachable. Three problems must be solved:

1. **No configuration layer** — There is no `Config` struct or environment-variable loader. The application has no way to read `DATABASE_URL` or any pool-tuning parameters at runtime.
2. **No database connection** — There is no database package. The application cannot open, pool, or ping a PostgreSQL connection.
3. **Opaque health status** — The `/health` endpoint ignores its primary dependency. A health-check that always says "ok" is misleading to developers and orchestrators (Docker health-checks, load-balancers, Kubernetes probes).

### Edge Cases & Considerations

1. **Startup with database down** — `NewPostgresDB` must call `Ping()` and return an error at startup. This is correct for a hard dependency. However, health must degrade gracefully *after* startup if the database becomes temporarily unreachable (network blip, container restart).
2. **Context timeout on health ping** — A Ping during health-check must time-out quickly (≤ 2 seconds). A slow or hung Ping should not block the entire health endpoint.
3. **Connection pool exhaustion** — If `MaxOpenConns` is set too low for load, queries will block. Environment-specific tuning must be possible without code changes.
4. **Idle connection eviction** — If `ConnMaxIdleTime` is not set, idle connections can remain open indefinitely and eventually be killed server-side, resulting in `unexpected EOF` errors to the next query that reuses a stale connection.
5. **Graceful shutdown** — When the application exits, `db.Close()` must drain in-flight queries. The `main.go` must add `defer db.Close()` after opening the connection.
6. **Config validation** — Pool values must be sane: `MaxOpenConns ≥ MaxIdleConns`, `MaxIdleConns ≥ 1`, lifetimes > 0.

---

## B. Structural Strategy

### File System Impact

| Action     | File                                  | Purpose |
|------------|---------------------------------------|-----------------------------------------|
| **Create** | `pkg/config/config.go`               | Create the configuration package. Define a `Config` struct with `DatabaseURL` and connection-pool fields (`DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME_SECS`, `DB_CONN_MAX_IDLE_TIME_SECS`), a `LoadConfig` function using `envconfig`, and pool-parameter validation. |
| **Create** | `pkg/database/postgres.go`           | Create the database package. Implement `NewPostgresDB` that accepts the DSN and pool-tuning parameters, opens a `sqlx` connection, configures the pool, and pings on startup. |
| **Modify** | `cmd/api/main.go`                    | Add config loading via `config.LoadConfig`, database initialisation via `database.NewPostgresDB`, dependency injection into the health handler, and structured startup logging. |
| **Modify** | `internal/handlers/health.go`        | Refactor from a standalone function to a struct-based handler (or accept a dependency) so it can call `db.PingContext`. Return a structured response that includes database status. |
| **Modify** | `.env.example`                       | Document the new pool-related environment variables with their defaults. |
| **Modify** | `.env`                               | Add the new pool-related environment variables with sensible development defaults. |

### Module Architecture

- **`pkg/config`** owns the source-of-truth for every environment variable. The pool settings live here as first-class fields with documented defaults.
- **`pkg/database`** is the only module that touches `*sqlx.DB` lifecycle operations (open, configure, ping, close). No other package should directly configure pool settings.
- **`internal/handlers`** gains a dependency on the database handle purely for the health-check ping. It does **not** import `pkg/database` — it only needs the standard `*sqlx.DB` (or a narrow interface) to call `PingContext`.
- **`cmd/api/main.go`** is the composition root. It reads config, creates the DB, and injects it into the health handler. The wiring changes here, but no business logic does.

### Interface Specs

- **`NewPostgresDB`** revised signature: Accepts the DSN string **and** the four pool-tuning integers/durations. Returns `*sqlx.DB` and `error`, unchanged. This keeps the function simple and avoids coupling it to the `Config` struct, which lives in a different package.
- **Health response body**: The `/health` endpoint should return a JSON object with a top-level `status` field (`"healthy"` or `"degraded"`), a `checks` object containing named subsystem statuses. Each check has a `status` field (`"up"` or `"down"`) and an optional `message` field with details on failure. Example structure:
  - `status`: `"healthy"` when all checks pass, `"degraded"` when at least one check fails.
  - `checks.database.status`: `"up"` or `"down"`.
  - `checks.database.message`: Empty on success; contains the error string on failure.

---

## C. Step-by-Step Logic

### Step 0 — Install Go Dependencies

1. Install the required Go packages that do not yet exist in `go.mod`:
   - `github.com/jmoiron/sqlx` — SQL extensions for Go's `database/sql`.
   - `github.com/jackc/pgx/v5/stdlib` — PostgreSQL driver registered as `pgx` for `database/sql`.
   - `github.com/kelseyhightower/envconfig` — Struct-based environment variable loader.
2. Run `go get` for each package, then `go mod tidy` to clean up.

### Step 1 — Create the Config Package

1. Create the directory `pkg/config/` and the file `pkg/config/config.go`.
2. Define a `Config` struct with the following fields:
   - `DatabaseURL` — string, env var `DATABASE_URL`, required.
   - `DBMaxOpenConns` — integer, env var `DB_MAX_OPEN_CONNS`, default `25`.
   - `DBMaxIdleConns` — integer, env var `DB_MAX_IDLE_CONNS`, default `5`.
   - `DBConnMaxLifetimeSecs` — integer (seconds), env var `DB_CONN_MAX_LIFETIME_SECS`, default `300` (5 minutes).
   - `DBConnMaxIdleTimeSecs` — integer (seconds), env var `DB_CONN_MAX_IDLE_TIME_SECS`, default `60` (1 minute).
3. Implement a `LoadConfig() (*Config, error)` function that:
   - Loads `.env` using `godotenv.Load()` (best-effort, no error if file is missing — supports container deployments without `.env`).
   - Calls `envconfig.Process("", &cfg)` to populate the struct from environment variables.
   - Validates pool parameters after loading:
     - Verify `DBMaxOpenConns >= 1`. If not, return an error — a value of 0 means "unlimited" in Go's `sql.DB`, which is unsafe for production.
     - Verify `DBMaxOpenConns >= DBMaxIdleConns`. If not, return an error with a descriptive message.
     - Verify `DBMaxIdleConns >= 1`. If not, return an error.
     - Verify both lifetime values are > 0. If not, return an error.
4. The `godotenv` package (`github.com/joho/godotenv`) must also be installed in Step 0.

### Step 2 — Create the Database Package

1. Create the directory `pkg/database/` and the file `pkg/database/postgres.go`.
2. Implement `NewPostgresDB` with the signature: `func NewPostgresDB(dsn string, maxOpenConns, maxIdleConns, connMaxLifetimeSecs, connMaxIdleTimeSecs int) (*sqlx.DB, error)`.
3. Inside the function:
   - Call `sqlx.Connect("pgx", dsn)` to open the connection using the `pgx` driver. Note: `sqlx.Connect` internally calls `Open` + `Ping`, so an explicit `db.Ping()` call is **not** needed — it would be a redundant second ping that adds unnecessary startup latency.
   - Configure the connection pool with the passed-in values:
     - `SetMaxOpenConns` ← `maxOpenConns` parameter
     - `SetMaxIdleConns` ← `maxIdleConns` parameter
     - `SetConnMaxLifetime` ← convert `connMaxLifetimeSecs` to `time.Duration`
     - `SetConnMaxIdleTime` ← convert `connMaxIdleTimeSecs` to `time.Duration`
   - Return `*sqlx.DB` and `nil` error on success.
4. Keep the function pure — let `main.go` handle logging the applied pool settings.

### Step 3 — Update the Health Handler

1. Open `internal/handlers/health.go`.
2. Create a `HealthHandler` struct that holds a `*sqlx.DB` field. This follows the pattern of injecting dependencies into handlers while keeping the handler testable.
3. Add a constructor function `NewHealthHandler` that accepts `*sqlx.DB` and returns the struct.
4. Rewrite the `HealthCheck` method on the struct:
   - Create a child context with a 2-second timeout from the request context.
   - Call `db.PingContext(ctx)` using the timeout context.
   - Build the response object:
     - If the ping succeeds: top-level `status` is `"healthy"`, `checks.database.status` is `"up"`. Return **HTTP 200**.
     - If the ping fails: top-level `status` is `"degraded"`, `checks.database.status` is `"down"`, `checks.database.message` contains the error string. Return **HTTP 503 Service Unavailable**.
   - The database is a hard dependency. When it is down, the service cannot fulfil its core purpose, so a `503` is the correct signal to load balancers and orchestrators (Kubernetes liveness/readiness probes, ECS health checks) to stop routing traffic to this instance.

### Step 4 — Update `main.go` Wiring

1. Open `cmd/api/main.go`.
2. Add imports for `pkg/config` and `pkg/database`.
3. After logger initialisation, add a call to `config.LoadConfig()`. Fatal if it returns an error.
4. Add a call to `database.NewPostgresDB` passing `cfg.DatabaseURL`, `cfg.DBMaxOpenConns`, `cfg.DBMaxIdleConns`, `cfg.DBConnMaxLifetimeSecs`, `cfg.DBConnMaxIdleTimeSecs`. Fatal if it returns an error.
5. Add `defer db.Close()` immediately after the database is successfully opened.
6. Log the applied pool configuration values using `zap.Int` fields so developers can confirm what was applied at startup.
7. Create a `HealthHandler` instance: call `handlers.NewHealthHandler(db)`.
8. Update the route registration from `handlers.HealthCheck` (standalone function) to the method on the `HealthHandler` instance (e.g., `healthHandler.HealthCheck`).

### Step 5 — Update Environment Files

1. Open `.env.example` and `.env`.
2. Add a "Connection Pool" sub-section under "Database Configuration" with the four variables and their defaults:
   - `DB_MAX_OPEN_CONNS=25`
   - `DB_MAX_IDLE_CONNS=5`
   - `DB_CONN_MAX_LIFETIME_SECS=300`
   - `DB_CONN_MAX_IDLE_TIME_SECS=60`
3. Add inline comments explaining the purpose of each setting in plain language.

---

## D. Best Practice & Quality Guardrails

### Error Handling
- **Config validation failure**: `LoadConfig` returns a descriptive error if pool parameters are invalid (e.g., "`DB_MAX_IDLE_CONNS (30) cannot exceed DB_MAX_OPEN_CONNS (25)`"). The application fatals at startup with this message — there is no sensible default the app can silently fall back to.
- **Startup Ping failure**: Remains a `log.Fatal`. A hard dependency that is unreachable at boot time should block deployment rather than start an app that will fail every request.
- **Health-check Ping failure**: Logged as a `zap.Warn` (not Fatal/Error). The application continues serving; only the health response body signals degradation. This prevents a transient database restart from crashing the entire process.

### Security
- No new secrets are introduced. The DSN is already in `.env` and loaded via `envconfig`.
- Pool tuning parameters are integers/durations with no security implications.

### Performance
- `ConnMaxIdleTime` of 60 seconds prevents stale-connection errors by proactively closing idle connections before the database server's own `idle_in_transaction_session_timeout` kills them.
- `ConnMaxLifetime` of 5 minutes forces periodic connection recycling, ensuring the pool re-resolves DNS and picks up database failovers (relevant in cloud-managed PostgreSQL).
- The health-check Ping has a 2-second hard timeout so it never blocks the HTTP response beyond a reasonable window.

### Observability
- At startup, log the four applied pool settings as structured `zap` fields. This makes misconfiguration immediately visible in log output.
- On health-check Ping failure, log a `Warn`-level message with the error so monitoring can alert on repeated degradation even if no external probe inspects the `/health` body.

---

## E. Verification Plan

### Test Scenarios

#### Success Scenarios

| # | Scenario | Validation Steps |
|---|----------|-----------------|
| 1 | **Default pool config** | Start the application without setting any `DB_*` pool env vars. Check startup logs — they should report `MaxOpenConns=25`, `MaxIdleConns=5`, `ConnMaxLifetime=5m`, `ConnMaxIdleTime=1m`. |
| 2 | **Custom pool config** | Set `DB_MAX_OPEN_CONNS=10` and `DB_MAX_IDLE_CONNS=3` in `.env`. Restart the app. Startup logs should reflect the overridden values. |
| 3 | **Healthy health-check** | With the database running, `GET /health` returns HTTP 200 with `"status": "healthy"` and `"checks": { "database": { "status": "up" } }`. |
| 4 | **Degraded health-check** | Stop the PostgreSQL container (`make db-down`). Hit `GET /health`. Response is **HTTP 503** with `"status": "degraded"` and `"checks": { "database": { "status": "down", "message": "..." } }`. The app remains running and serving the endpoint, but signals to upstream load balancers that it cannot process requests. |
| 5 | **Recovery after degradation** | After scenario 4, restart the database (`make db-up`). Hit `GET /health` again. Response flips back to `"healthy"` with `database: "up"`. |

#### Failure Scenarios

| # | Scenario | Expected Behaviour |
|---|----------|-------------------|
| 1 | **Invalid pool config** | Set `DB_MAX_IDLE_CONNS=50` with `DB_MAX_OPEN_CONNS=10`. Application should refuse to start with a descriptive fatal log message explaining the constraint violation. |
| 2 | **Database unreachable at startup** | Stop the database before starting the app. Application fatals with `"Failed to connect to database"` — it does not start in a degraded state. |
| 3 | **Missing DATABASE_URL** | Remove `DATABASE_URL` from `.env`. Application fatals with envconfig's required-field error. |

#### Edge Case Scenarios

| # | Scenario | Expected Behaviour |
|---|----------|-------------------|
| 1 | **Health-check during brief network blip** | Simulate a ~1 second network interruption. The health-check Ping times out within 2 seconds and returns `"degraded"`. Subsequent requests succeed normally once connectivity is restored. |
| 2 | **Concurrent health-checks under load** | Multiple simultaneous `GET /health` requests each get their own `PingContext`; none block the others beyond the 2-second timeout. The connection pool services these pings without exhausting the pool (Ping uses minimal resources). |
| 3 | **Zero-value environment override** | Set `DB_CONN_MAX_LIFETIME_SECS=0`. Validation in `LoadConfig` rejects this and fatals with a descriptive message — zero means "no limit", which is architecturally unsafe with managed databases that enforce idle timeouts. |
