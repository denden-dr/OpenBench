# Backend Part 4: Service & API Layer

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Note:** This plan depends on the database connection and repository established in [Backend Part 3: Data Layer](backend-part-3-data-layer.md).

**Goal:** Implement business logic in the service layer, HTTP handlers in the API layer, and wire everything together in `main.go`.

---

### Task 1: Service Layer Implementation

**Files to modify/create:**
- `apps/backend/internal/dto/ticket_dto.go`
- `apps/backend/internal/service/ticket_service.go`

- [ ] **Step 1: Define DTOs (Data Transfer Objects)**
  - Create `internal/dto/ticket_dto.go`.
  - Define a `CreateTicketRequest` struct to represent the incoming JSON payload. Include validation tags (e.g., using `go-playground/validator`) to enforce required fields.
  - Define a `TicketResponse` struct to dictate exactly what data should be exposed to the client, preventing accidental data leaks from the internal model.
  - **Important:** Ensure monetary fields like `DiagnosisFee` use a robust decimal library like `github.com/shopspring/decimal`, matching the internal model, to avoid floating-point precision issues.

- [ ] **Step 2: Implement Ticket Service**
  - Create `internal/service/ticket_service.go`.
  - Define a `TicketService` interface containing the business operations (e.g., `CreateTicket`, `GetTicket`).
  - Implement the interface. The service implementation should hold a dependency on the `TicketRepository` interface (Dependency Injection).
  - **Best Practice (Error Handling Boundaries):** Define domain-level errors (e.g., `ErrTicketNotFound`) in the service layer or a shared package. The handler must NEVER import the repository package to check errors; it should only check against service layer errors.
  - **Best Practice (Separation of Concerns):** The service layer should handle mapping incoming DTOs to the internal domain `model` structs, passing them to the repository, and then mapping the returned domain models back into response DTOs.
  - **Best Practice (Context Propagation):** Ensure `context.Context` is accepted as the first argument in all service methods and passed down to the repository layer.

---

### Task 2: API Handler & Wiring

**Files to modify/create:**
- `apps/backend/internal/handler/ticket_handler.go`
- `apps/backend/main.go`

- [ ] **Step 1: Implement Ticket Handler**
  - Create `internal/handler/ticket_handler.go`.
  - Define a `TicketHandler` struct that depends on the `TicketService` interface.
  - Implement Fiber handler methods for routing (e.g., `Create`, `Get`).
  - **Best Practice (Input Validation & Error Handling):** Parse the incoming request body into the DTO. If parsing or service layer execution fails, return proper HTTP status codes (e.g., 400 Bad Request, 404 Not Found, 500 Internal Server Error).
  - **Important:** When validation fails, parse `validator.ValidationErrors` to return user-friendly messages instead of exposing raw internal validator errors to the client.
  - **Important:** Validate path parameters (like UUIDs) before executing service logic or DB queries to avoid pushing validation down to the database layer.
  - **Best Practice (API Standardization):** Ensure all JSON responses, both successful and erroneous, follow the standard envelope format (`{"success": bool, "data": any, "error": string}`) as prescribed by the project's API patterns.

- [ ] **Step 2: Update main.go with Dependency Injection**
  - Modify `main.go` to construct the dependency graph manually.
  - Initialize the database connection.
  - Instantiate the `TicketRepository`, injecting the database connection.
  - Instantiate the `TicketService`, injecting the repository.
  - Instantiate the `TicketHandler`, injecting the service.
  - Register the handler methods to a versioned API router group (e.g., `/api/v1/tickets`).
  - Update the `/health` endpoint to perform a `db.PingContext` check rather than returning a static "ok" string, ensuring actual database connectivity is verified.

- [ ] **Step 3: Final Verification**
  - Ensure the database and migrations are up using the provided Make commands.
  - Run the application.
  - Perform integration testing using `curl` or a similar tool to create a ticket and retrieve it, confirming the end-to-end flow and JSON structures are correct.
