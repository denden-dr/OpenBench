# Plan: Setup Local PostgreSQL via Container

> **Goal**: Run a fully functional PostgreSQL instance inside a container (Docker or Podman) so that developers never need to install PostgreSQL on their host machine.

---

## A. Logical Requirements

### Problem Statement
The project currently references a PostgreSQL database at `localhost:5432` (see `.env.example`), but there is **no provisioning mechanism**. Every developer must manually install and configure PostgreSQL, which is error-prone and environment-dependent. We need a single-command, repeatable way to spin up the database.

### Edge Cases & Considerations
1. **Port collision** — Port 5432 may already be in use on the host (e.g., an existing PostgreSQL installed locally). The plan must allow overriding the host-side port via `DB_PORT`.
2. **Data persistence** — If the container is destroyed, all data is lost unless a named volume is used.
3. **First-run bootstrapping** — The database name `openbench` must exist before the application tries to connect.
4. **Migration ordering** — Migrations must only run after the database is healthy and accepting connections.
5. **Platform parity** — The setup must work identically on Linux and macOS (the two common development platforms).
6. **Credential mismatch** — The Compose credentials and the `DATABASE_URL` in `.env` must stay in sync; a mismatch silently causes connection errors.
7. **Container engine variance** — Developers may use Docker or Podman. The Makefile must auto-detect the available engine and compose command so that all targets work transparently regardless of which runtime is installed.

---

## B. Structural Strategy

### File System Impact

| Action   | File / Path                | Purpose |
|----------|----------------------------|---------|
| **Create** | `docker-compose.yml`       | Defines the PostgreSQL service, volume, health-check, and network. Compatible with both Docker Compose and Podman Compose. |
| **Modify** | `Makefile`                 | Add engine auto-detection (`docker` vs `podman`), lifecycle targets (`db-up`, `db-down`, `db-reset`, `db-logs`, `db-shell`), and upgrade `migrate-*` targets to use `golang-migrate` with a container health prerequisite. |
| **Modify** | `.env.example`             | Add clear comments anchoring the default values to the Compose defaults. |
| **Modify** | `.gitignore`               | Add `docker-compose.override.yml` so developers can have personal overrides without polluting version control. |

### Module Architecture
- **No Go code changes.** The existing `pkg/database/postgres.go` connects via a `DATABASE_URL` DSN string, which already supports any host/port/credentials. As long as the Compose service exposes a matching DSN, the Go application connects transparently.
- **Makefile is the integration surface.** All new commands are Makefile targets that shell-out to the detected compose command (Docker Compose or Podman Compose). The developer never interacts with the container engine directly for database work.

---

## C. Step-by-Step Logic

### Step 1 — Create `docker-compose.yml`

1. Define a single service named `postgres`.
2. Use the official `postgres:17-alpine` image (small footprint, latest stable LTS at time of writing).
3. Set environment variables inside the compose file that mirror the `.env.example` defaults:
   - `POSTGRES_USER` → `postgres`
   - `POSTGRES_PASSWORD` → `postgres`
   - `POSTGRES_DB` → `openbench`
4. Map the container's port `5432` to the host port `5432`. Use a variable (`${DB_PORT:-5432}`) so developers can override the host port without editing the file.
5. Attach a **named volume** called `openbench_pg_data` mounted at `/var/lib/postgresql/data` to persist data across container restarts and rebuilds.
6. Add a **health-check** that runs `pg_isready -U postgres` every 5 seconds, with 5 retries and a 5-second timeout. This allows dependent targets to gate on database readiness.
7. Set `restart: unless-stopped` so the database survives an accidental `docker stop` or host reboot (while Docker daemon is set to auto-start).

### Step 2 — Update `Makefile`

Add engine auto-detection and the following targets (descriptions are logic, not code):

**Engine Detection (mandatory first step)**:
- Define a `DOCKER_CMD` variable that resolves to `docker` if available, otherwise falls back to `podman`.
- Define a `COMPOSE_CMD` variable that resolves to `$(DOCKER_CMD) compose` if the subcommand is supported, otherwise falls back to `docker-compose` or `podman-compose` (whichever is found).
- All targets below must use `$(COMPOSE_CMD)` instead of hard-coding `docker compose`.

1. **`db-up`** — Runs `$(COMPOSE_CMD) up -d` to start the PostgreSQL container in detached mode. After launching, poll `$(COMPOSE_CMD) ps` for a `healthy` status **with a maximum timeout of 60 seconds**. If the database does not become healthy within the timeout, print an error message and exit non-zero to prevent the developer from waiting indefinitely. On success, print a confirmation message.
2. **`db-down`** — Runs `$(COMPOSE_CMD) down` to stop and remove the container (but keep the volume).
3. **`db-reset`** — Runs `$(COMPOSE_CMD) down -v` to stop and remove the container **and** the volume, then re-runs `db-up`. This gives a completely clean slate.
4. **`db-logs`** — Runs `$(COMPOSE_CMD) logs -f postgres` to tail the database container logs. Useful for debugging startup failures.
5. **`db-shell`** — Runs `$(COMPOSE_CMD) exec postgres psql -U postgres -d openbench` to open an interactive `psql` session inside the running container.
6. **Update `migrate-up`** — Switch from raw `psql` to `golang-migrate` (`migrate -path migrations -database ... up`). Add a prerequisite check that the container is healthy before executing.
7. **Update `migrate-down`** — Use `golang-migrate` (`migrate ... down -all`) with the same prerequisite check.
8. **Add `migrate-create`** — Add a target to scaffold new migration files (`*.up.sql` and `*.down.sql`) using `golang-migrate` to prevent manual file creation errors.
9. **Create `migrations/` directory** — Commit a `migrations/.gitkeep` file so the directory exists on fresh checkouts. Without this, `make migrate-up` fails with a "directory not found" error.

### Step 3 — Update `.env.example`

1. Add a comment block above `DATABASE_URL` that explicitly states: "These defaults match the docker-compose.yml configuration. Change them only if you have customised `docker-compose.override.yml`."
2. Keep the value unchanged (`postgres://postgres:postgres@localhost:5432/openbench?sslmode=disable`) — it already matches the Compose defaults we are setting.

### Step 4 — Update `.gitignore`

1. Add a "Docker" section.
2. Add `docker-compose.override.yml` to the ignore list — Docker Compose natively merges this override file when present, giving developers a private customisation point (e.g., changing the host port) without modifying the tracked `docker-compose.yml`.

### Step 5 — Developer Documentation (README Update)

1. Add a "Local Development" section to `README.md` with the following subsections:
   - **Prerequisites**: Docker & Docker Compose (v2+) **OR** Podman & Podman Compose, Go, Make, `golang-migrate` CLI, and Air (for hot-reloading).
   - **Quick Start**: Step-by-step instructions — copy `.env.example` to `.env`, run `make db-up`, run `make migrate-up`, run `make run`.
   - **Makefile Commands**: Comprehensive table of **all** Makefile targets (application lifecycle, database management, migrations) and their purpose.
   - **Port Overrides**: Instructions for using `DB_PORT=<port> make db-up` when port 5432 is already in use, and how to update `.env` accordingly.
   - **Podman Support**: Note that the Makefile auto-detects the container engine.

---

## D. Best Practice & Quality Guardrails

### Error Handling
- **`make migrate-up` without a running container** must fail fast with a human-readable message ("Database container is not running. Run `make db-up` first.") rather than a cryptic `psql: could not connect` error.
- **`make db-shell` on a stopped container** should similarly print a helpful error.

### Security
- The default credentials (`postgres`/`postgres`) are acceptable for local development only.
- The `docker-compose.yml` must include a comment warning that these are development-only credentials and must not be used in any deployed environment.
- The Compose file must **not** expose the port to `0.0.0.0` beyond the loopback interface. Bind to `127.0.0.1:${DB_PORT:-5432}:5432` so the database is only accessible from the host.

### Performance
- Use the `alpine` variant of the PostgreSQL image to minimise download time and disk usage.
- Using a named volume (instead of a bind-mount) ensures the PostgreSQL filesystem performance is optimal on macOS, where bind-mount I/O can be slow.

### Observability
- The `db-logs` Makefile target provides real-time log streaming from the container.
- The health-check status is visible via `$(COMPOSE_CMD) ps`, allowing developers to quickly verify if the database is healthy.

---

## E. Verification Plan

### Test Scenarios

#### Success Scenarios
| # | Scenario | Validation Steps |
|---|----------|-----------------|
| 1 | **Cold start** | Clone repo → `cp .env.example .env` → `make db-up` → compose ps shows `postgres` service "healthy" within 15 seconds. |
| 2 | **Application connects** | After scenario 1 → `make migrate-up` succeeds → `make run` → Hit `GET /health` → returns 200. |
| 3 | **Data persists** | Insert a row into the `users` table via `make db-shell` → `make db-down` → `make db-up` → Re-query the row; it should still exist (volume persisted). |
| 4 | **Clean reset** | `make db-reset` → `make db-shell` → `\dt` shows no tables (volume was destroyed and recreated). → `make migrate-up` re-creates them. |

#### Failure Scenarios
| # | Scenario | Expected Behaviour |
|---|----------|-------------------|
| 1 | **Port conflict** | If port 5432 is already bound, `make db-up` fails with a clear container runtime error about port allocation. Developer can override via `DB_PORT=5433 make db-up` and update their `.env` accordingly. |
| 2 | **Container engine not running** | `make db-up` fails with a daemon/service error. Developer sees a message indicating the container engine is not started. |
| 3 | **Migration without container** | `make migrate-up` before `make db-up` prints a controlled error message and exits non-zero. |
| 4 | **Neither Docker nor Podman installed** | `make db-up` fails because `COMPOSE_CMD` resolves to empty. The error output makes it clear a container engine is required. |

#### Edge Case Scenarios
| # | Scenario | Expected Behaviour |
|---|----------|-------------------|
| 1 | **Override file present** | Developer creates `docker-compose.override.yml` changing the host port to 5433. `make db-up` uses the merged config. The override file is not tracked by git. |
| 2 | **Repeated `db-up`** | Running `make db-up` when the container is already running is a no-op (Compose handles idempotency). |
| 3 | **Repeated `db-reset`** | Running `make db-reset` twice in a row works — the second run simply recreates a fresh volume. |
| 4 | **Podman instead of Docker** | Developer only has Podman installed. `make db-up` auto-detects `podman` and `podman-compose`, and all targets work identically. |
