# OpenBench
Web APP project to provide repair phone services

## Local Development

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose **OR** [Podman](https://podman.io/) & Podman Compose
- [Go](https://go.dev/doc/install) (1.21+)
- [Make](https://www.gnu.org/software/make/)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI
- [Air](https://github.com/cosmtrek/air) (for live reloading)

### Quick Start
1. **Initialize Environment**:
   ```bash
   cp .env.example .env
   ```
2. **Start Infrastructure**:
   ```bash
   make db-up
   ```
3. **Run Migrations**:
   ```bash
   make migrate-up
   ```
4. **Launch Application**:
   ```bash
   make run
   ```

---

## Makefile Commands

The project uses a `Makefile` to simplify common development tasks. It automatically detects if you are using `docker` or `podman`.

### Application Lifecycle
| Command | Portfolio | Description |
|---------|-----------|-------------|
| `make build` | Build | Compiles the Go application into `bin/api`. |
| `make run` | Dev | Starts the application with hot-reloading using `air`. |
| `make fmt` | Style | Formats all Go code in the project. |
| `make tidy` | Deps | Cleans up and synchronizes Go module dependencies. |

### Database Management (`db-*`)
| Command | Description |
|---------|-------------|
| `make db-up` | Spins up the PostgreSQL container in detached mode and waits for health. |
| `make db-down` | Stops and removes the PostgreSQL container, preserving data. |
| `make db-reset` | **DANGER**: Wipes the database volume and restarts with a clean slate. |
| `make db-logs` | Tunnels into the container to stream PostgreSQL logs. |
| `make db-shell` | Opens an interactive `psql` session inside the running container. |

### Migrations (`migrate-*`)
| Command | Description |
|---------|-------------|
| `make migrate-up` | Applies all pending SQL migrations to the database. |
| `make migrate-down` | **DANGER**: Rolls back all migrations (destructive). |
| `make migrate-create` | Prompts for a name and scaffolds new `.up.sql` and `.down.sql` files. |

---

## Technical Configuration

### Port Overrides
By default, the database binds to `127.0.0.1:5432`. If this port is occupied on your host, you can override it using the `DB_PORT` environment variable:

```bash
DB_PORT=5433 make db-up
```

After overriding the port, ensure your `DATABASE_URL` in `.env` reflects the change:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5433/openbench?sslmode=disable
```

### Podman Support
The `Makefile` is compatible with Podman. If `docker` is not found, it will automatically attempt to use `podman` and `podman-compose`.
