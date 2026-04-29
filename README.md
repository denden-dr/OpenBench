# 🛠️ OpenBench

OpenBench is a robust, domain-driven web application designed to streamline phone repair services. It connects customers with skilled technicians, managing the entire lifecycle of a repair booking from initial request to diagnosis and completion.

---

## ✨ Key Features

- **🔐 Secure Authentication**: Integrated with Supabase Auth for robust identity management.
- **🔄 JIT User Provisioning**: Automatic synchronization of Supabase users into the local database upon first login.
- **📅 Booking Lifecycle**: Full management of repair requests, including creation, approval, and cancellation.
- **👷 Technician Workflow**: Specialized dashboard for technicians to claim bookings, provide diagnoses, and update status.
- **🛡️ RBAC Middleware**: Role-based access control ensuring only authorized technicians can access sensitive operations.
- **📈 Observability**: Structured logging using `uber-go/zap` and health monitoring.

---

## 🏗️ Architecture

OpenBench follows a **Clean Architecture** (Layered Architecture) pattern with a strong focus on **Domain-Driven Design (DDD)** principles:

- **Handlers**: HTTP transport layer using the Fiber framework.
- **Services**: Core business logic and orchestration.
- **Repositories**: Data access layer using the Repository pattern for database agnosticism.
- **Domain**: Pure business entities and interface definitions.

```mermaid
graph TD
    Client[HTTP Client] --> Handlers[Handlers Layer]
    Handlers --> Services[Service Layer]
    Services --> Repos[Repository Layer]
    Repos --> DB[(PostgreSQL)]
    Services --> Domain[Domain Entities]
    Repos --> Domain
```

---

## 🚀 Tech Stack

- **Language**: [Go (Golang)](https://go.dev/) 1.21+
- **Framework**: [Fiber v3](https://docs.gofiber.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **ORM/Query**: [sqlx](https://github.com/jmoiron/sqlx)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Logging**: [zap](https://github.com/uber-go/zap)
- **Auth**: [Supabase](https://supabase.com/)

---

## 📂 Project Structure

```text
.
├── cmd/api/            # Application entry point
├── internal/
│   ├── domain/         # Domain entities and error definitions
│   ├── dto/            # Data Transfer Objects
│   ├── handlers/       # HTTP handlers (Controller layer)
│   ├── middleware/     # Fiber middlewares (Auth, Logging, Role)
│   ├── repository/     # Data access implementation (Repository layer)
│   └── service/        # Business logic implementation (Service layer)
├── migrations/         # SQL migration files
├── pkg/                # Shared packages (Config, Database, Logger)
└── bin/                # Compiled binaries
```

---

## 🛠️ Local Development

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

## 📜 Makefile Commands

The project uses a `Makefile` to simplify common development tasks. It automatically detects if you are using `docker` or `podman`.

### Application Lifecycle
| Command | Portfolio | Description |
|:---|:---|:---|
| `make build` | Build | Compiles the Go application into `bin/api`. |
| `make run` | Dev | Starts the application with hot-reloading using `air`. |
| `make fmt` | Style | Formats all Go code in the project. |
| `make tidy` | Deps | Cleans up and synchronizes Go module dependencies. |

### Database Management (`db-*`)
| Command | Description |
|:---|:---|
| `make db-up` | Spins up the PostgreSQL container and waits for health. |
| `make db-down` | Stops and removes the PostgreSQL container. |
| `make db-reset` | **DANGER**: Wipes the database volume and restarts clean. |
| `make db-logs` | Tunnels into the container to stream PostgreSQL logs. |
| `make db-shell` | Opens an interactive `psql` session. |

### Migrations (`migrate-*`)
| Command | Description |
|:---|:---|
| `make migrate-up` | Applies all pending SQL migrations. |
| `make migrate-down` | **DANGER**: Rolls back all migrations (destructive). |
| `make migrate-create` | Prompts for a name and scaffolds new migration files. |

---

## ⚙️ Technical Configuration

### Port Overrides
By default, the database binds to `127.0.0.1:5432`. If this port is occupied, you can override it:

```bash
DB_PORT=5433 make db-up
```

After overriding, ensure your `DATABASE_URL` in `.env` reflects the change:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5433/openbench?sslmode=disable
```

### Podman Support
The `Makefile` is fully compatible with Podman. If `docker` is not found, it will automatically fall back to `podman` and `podman-compose`.

