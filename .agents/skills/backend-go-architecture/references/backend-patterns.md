# Backend Patterns

## Files To Inspect First

- Package examples: `apps/backend/internal/auth`, `health`, `database`
- Response envelope: `apps/backend/internal/pkg/response/response.go`
- Validation helper: `apps/backend/internal/pkg/validator`
- Test database helper: `apps/backend/internal/pkg/testutil/db.go`
- Migrations: `apps/backend/migrations`

## Layer Responsibilities

- Domain models: pure Go structs and domain constants. Avoid ORM tags.
- Repository interfaces: database operations and SQL signatures.
- Repository implementations: SQL, scanning, locking, and persistence concerns.
- Services: business rules, transaction boundaries, orchestration, and domain errors.
- Handlers: request parsing, validation, auth context extraction, response mapping, and DTOs.

## Response Envelope

Use the shared helpers:

```go
return response.JSON(c, fiber.StatusOK, "Message", payload)
return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
```

Frontend services expect the successful payload under `data`.

## Transactions

- Start multi-step DB transactions in the service layer.
- Use `defer tx.Rollback()` immediately after a successful begin.
- Pass `*sqlx.Tx` into repository methods that participate in the transaction.
- Commit only in the service after all business rules and repository calls pass.
- Use `SELECT ... FOR UPDATE` or an equivalent locking strategy when concurrent updates can race.

## Auth Sessions

- Store refresh tokens in `HttpOnly` cookies.
- Use `Secure: !isDev`, `SameSite: Lax`, and `Path: "/"`.
- Clear cookies with the same path, secure, HTTPOnly, and SameSite settings used when setting them.
- Keep access tokens short lived and return them in JSON; keep refresh tokens long lived and cookie-backed.
- Refresh token rotation should track token families, revoke families on replay outside grace, and allow the existing short grace period for concurrent refreshes.

## Public Endpoints

- Use UUID public identifiers for unauthenticated lookup flows.
- Validate UUID input before querying.
- Do not expose internal notes, PII beyond the workflow requirement, or sequential internal IDs.
- Keep public payloads intentionally narrow.

## Migrations

- Put SQL migrations under `apps/backend/migrations`.
- Keep migrations forward-only and deterministic.
- Add indexes for lookup and locking columns touched by new repository queries.

## Config And Server Safety

- Reject unsafe production config such as empty DB passwords or disabled SSL outside local/test.
- Keep CORS origins explicit in production.
- Use context timeouts around connection retries and DB operations that can hang.
