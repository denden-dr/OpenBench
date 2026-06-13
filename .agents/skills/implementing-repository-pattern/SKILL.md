---
name: implementing-repository-pattern
description: Use when structuring Go backend codebases, decoupling database operations from business logic and routing handlers, and implementing domain models. Do not use for frontend data stores or non-Go backend codebases.
version: 1.1.0
---

# Implementing Repository Pattern

## Overview
Decouple HTTP/routing handlers from database logic by introducing clean architectural layers: Domain Models (pure business structs), Repository Layer (data access), and Service Layer (business logic and transaction orchestration).

## When to Use
- Fiber routing handlers are growing too large (>100 lines) with SQL queries, validation, and serialization mixed together.
- Database queries are directly embedded in routing handlers, making unit testing difficult without a running database.
- Multiple handlers need to perform the same database operations or reuse identical business rules.
- You need to introduce clean mock layers for testing business logic.

## Step-by-Step Instructions

1. **Define Domain Models**: Read `assets/domain.go.template` and create pure business structs without ORM or database tags.
2. **Define Repository Interface**: Read `assets/repository.go.template` and specify data access signatures. Implement database SQL execution in the repository implementation.
3. **Orchestrate Business Use Cases**: Read `assets/service.go.template` to implement the Service layer, which manages database transactions, domain validation, and custom domain errors.
4. **Delegate HTTP Routing**: Read `assets/handler.go.template` to implement the HTTP routing handler, which is strictly responsible for request parsing, service invocation, and JSON/cookie response mapping.
5. **Verify with Mock Testing**: Read `references/mock-testing.md` to write unit tests for your handlers using mock implementations of the service interfaces.
6. **Ensure Search Columns are Indexed**: In database migration files, always define a unique constraint or create an index (`CREATE INDEX`) for any columns that are regularly queried or locked in the repository layer (such. e.g. `token_hash` in `WHERE token_hash = $1 FOR UPDATE`), preventing costly table scans.

## Common Mistakes
- **Leaking Transactions to Handlers**: Handlers opening and committing transactions (`BeginTxx`/`Commit`). Transaction lifecycle belongs to the Service layer using a context or helper.
- **Leaking SQL errors to HTTP client**: Returning raw PG errors like `sql: no rows in result set` or `violates foreign key constraint` in JSON responses. Map repository database errors to custom Domain/Service sentinel errors (e.g. `ErrUserNotFound`) in the service, then serialize cleanly in the handler.
- **Fat Handlers**: Doing email validation, bcrypt hashing, or session storage logic inside routing handlers. Move these to domain/service layers respectively.
- **Injecting Structs Instead of Interfaces**: Passing `*authService` instead of `Service` to the handler, which breaks mocking and prevents unit testing the HTTP handler in isolation.
- **Bloated Service Functions**: Allowing service functions to grow excessively large (>100 lines) with multiple validation steps and data mutations under a single transaction. Instead, extract private helper methods on the service implementation struct (e.g., `handleUsedToken` or `rotateToken`) to maintain readability and clean logical separation while preserving transaction boundaries.
