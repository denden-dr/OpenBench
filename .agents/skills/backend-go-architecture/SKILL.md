---
name: backend-go-architecture
description: Build and refactor OpenBench Go/Fiber backend code. Use when adding domain models, repositories, services, handlers, SQL migrations, transactions, auth sessions, JWT/cookie behavior, public tracker endpoints, or database-backed tests in apps/backend.
---

# Backend Go Architecture

## Operating Rule

Keep backend changes in the existing layered shape: domain model, repository, service, handler, route registration, migration, and tests as needed. Do not skip layers for speed.

## Workflow

1. Inspect a nearby package under `apps/backend/internal` before adding new files.
2. Put DTO parsing and response formatting in handlers; put business rules and transaction orchestration in services; put SQL in repositories.
3. Use `response.JSON` and `response.Error` for all HTTP responses.
4. For DB mutations spanning multiple statements, start the transaction in the service, defer rollback, pass `*sqlx.Tx` only to repository methods, and commit from the service.
5. Add migrations for schema changes and keep domain structs free of ORM tags.
6. Verify with `gofmt` and targeted Go tests; use integration tests when SQL behavior changes.

## Load References

- Read `references/backend-patterns.md` before touching transactions, auth/session flow, public endpoints, migrations, config, or database tests.

## Hard Checks

- Do not expose internal sequential IDs from unauthenticated endpoints.
- Do not put transaction commits, rollbacks, or raw SQL in handlers.
- Do not use permissive production CORS.
- Do not clear cookies with flags that differ from the set-cookie flags.
