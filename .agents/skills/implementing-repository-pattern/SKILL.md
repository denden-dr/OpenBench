---
name: implementing-repository-pattern
description: Use when adding new domain packages to the Go backend, refactoring handler-to-database coupling, or initializing the Fiber server. Do not use for frontend data stores.
version: 2.0.0
---

# Go Backend Architecture

## Overview
Go backend packages must follow a clean layered architecture: Domain Models → Repository → Service → Handler → Server. This skill covers the full lifecycle from Fiber server initialization to HTTP response formatting.

## When to Use
- Initializing a new Go Fiber web server with security middleware and graceful shutdown.
- Fiber routing handlers are growing too large (>100 lines) with SQL, validation, and serialization mixed.
- Adding or refactoring domain packages under `internal/`.
- Writing or modifying HTTP handlers that return standardized JSON responses.

### When NOT to Use
- Frontend data stores or non-Go backend codebases.
- Gin, Chi, Echo, or other non-Fiber web frameworks.

---

## Step-by-Step Instructions

### Phase 1: Server Initialization
1. **Initialize Fiber App**: Read `assets/server.go.template` to set up the Fiber server with read/write timeouts (slowloris protection).
2. **Register Security Middlewares**: Configure `recover.New()` first, then `cors.New(...)` with explicit domain origins. Set `AllowCredentials: true` for cookie-based auth.
3. **Implement Graceful Shutdown**: Use a buffered error channel + OS signal listener (`SIGINT`, `SIGTERM`) with `app.Shutdown()`.

### Phase 2: Domain & Data Layers
4. **Define Domain Models**: Read `assets/domain.go.template` to create pure business structs without ORM or database tags.
5. **Define Repository Interface**: Read `assets/repository.go.template` for data access signatures and SQL execution.
6. **Orchestrate Business Logic**: Read `assets/service.go.template` for the Service layer (transactions, validation, domain errors).

### Phase 3: HTTP Handler & Response Formatting
7. **Implement Handler Struct**: Read `assets/handler-response.go.template` to build handlers as struct methods with body parsing and structural validation via `validator.ValidateStruct`.
8. **Organize DTOs**: Create DTO schemas inside a local `dto.go` file (e.g. `internal/auth/dto.go`). Read `assets/dto.go.template` for validation tag structure.
9. **Standardize Responses**: Use `response.JSON` and `response.Error` helpers from the `response` package.
10. **Register Routes**: Read `references/routing-example.md` to map handler struct methods to Fiber routes.
11. **Ensure Search Columns are Indexed**: Define indexes for columns regularly queried or locked in the repository layer.

---

## Common Mistakes
- **Leaking Transactions to Handlers**: Transaction lifecycle belongs in the Service layer, not handlers.
- **Leaking SQL Errors to Client**: Map repository errors to domain sentinel errors (e.g. `ErrUserNotFound`), then serialize in handler.
- **Fat Handlers**: Move email validation, bcrypt, session logic to domain/service layers.
- **Injecting Structs Instead of Interfaces**: Pass `Service` interface (not `*authService`) to handlers for mockability.
- **Bloated Service Functions**: Extract private helper methods (e.g. `handleUsedToken`, `rotateToken`).
- **Closure Handlers**: Use Handler struct constructor + methods, not closure-returning functions.
- **Direct `fiber.Map` Responses**: Always use `response.JSON` / `response.Error` for format consistency.
- **No Recover Middleware**: Unhandled panics crash the entire server without `recover.New()`.
- **Permissive CORS**: Never use `AllowOrigins: "*"` in production.
- **No Read/Write Timeouts**: Leaves server vulnerable to slowloris attacks.
- **Abrupt Termination**: Call `app.Shutdown()` instead of `os.Exit` to finish in-flight requests.
