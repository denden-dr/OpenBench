# Design Spec: OpenBench Backend Setup

## 1. Overview
This document outlines the architecture and initial implementation plan for the OpenBench backend, a phone repair management system. The backend will be built using Go and the Fiber framework, prioritizing simplicity, performance, and testability.

## 2. Architecture
We will use a layered architecture with manual dependency injection to ensure clear boundaries and ease of testing.

### Layers:
1.  **Handlers (`internal/handler`)**: Responsible for HTTP request parsing, input validation, and calling the appropriate service.
2.  **Services (`internal/service`)**: Contains the core business logic, orchestrates repository calls, and manages transactions.
3.  **Repositories (`internal/repository`)**: Handles all database interactions using `sqlx` and `pgx`.
4.  **DTOs (`internal/dto`)**: Data Transfer Objects for request and response payloads, ensuring we don't leak database models to the API.
5.  **Models (`internal/model`)**: Database entities representing the schema.

### Dependency Injection:
-   **Approach**: Manual constructor injection in `main.go`.
-   **Example**: `ticketService := service.NewTicketService(ticketRepo)`.

## 3. Technology Stack
-   **Language**: Go 1.25+
-   **Web Framework**: [Fiber](https://gofiber.io/)
-   **Database**: PostgreSQL
-   **DB Library**: `sqlx` with `pgx/v5` driver
-   **Migrations**: `golang-migrate`
-   **Configuration**: `.env` files via `godotenv`
-   **Infrastructure**: Docker/Podman for local development
-   **Tooling**: `Makefile` for common tasks

## 4. Data Model (Initial Slice)
We will start with the core tables required for the "Ticket" booking flow:

### Tickets Table
- `id`: UUID (Primary Key)
- `customer_id`: UUID (Foreign Key)
- `device_type`: Enum (Android, Apple)
- `brand`, `model`: String
- `issue_description`: Text
- `status`: Enum (received, diagnosing, etc.)
- `diagnosis_fee`: Decimal
- `created_at`, `updated_at`: Timestamp

*(Full schema as per PRD.md will be implemented incrementally)*

## 5. API Endpoints (Initial Slice)
- `POST /api/v1/tickets`: Create a new repair ticket.
- `GET /api/v1/tickets/:id`: Retrieve ticket details by ID.

## 6. Infrastructure (Podman/Docker)
A `Makefile` will manage the containerized PostgreSQL instance:
- `make db-up`: Start the database container.
- `make db-down`: Stop and remove the container.
- `make migrate-up`: Apply migrations using `golang-migrate`.
- `make migrate-down`: Rollback the last migration.

## 7. Testing Strategy
- **Unit Tests**: For service layer logic using mocks or interfaces.
- **Integration Tests**: For handler and repository layers, using a real PostgreSQL container.
- **Note**: Auth (Supabase JWT) is intentionally deferred to simplify initial integration testing.

---
*Created: 2026-05-02*
