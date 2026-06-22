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

For list endpoints:

- Return an empty slice (`[]T{}`) instead of a nil slice when no rows exist.
- Ensure the response envelope serializes `data: []` for successful empty lists.
- Avoid `omitempty` on envelope `data` when clients rely on the key being present.
- Add a focused test when changing shared response helpers because envelope changes affect every endpoint.

## Transactions

- Start multi-step DB transactions in the service layer.
- Use `defer tx.Rollback()` immediately after a successful begin.
- Pass `*sqlx.Tx` into repository methods that participate in the transaction.
- Commit only in the service after all business rules and repository calls pass.
- Use `SELECT ... FOR UPDATE` or an equivalent locking strategy when concurrent updates can race.

## State Transitions

- Capture the previous state before applying update fields.
- Run side effects only on transitions, not because the final state has a value. Example: `oldStatus != "picked_up" && newStatus == "picked_up"`.
- Treat terminal-state timestamps and derived dates as immutable unless a product requirement explicitly says otherwise.
- Reject or ignore changes to fields that would recalculate immutable derived records after terminal state is reached.
- Keep source row and derived rows in the same transaction when a transition creates both.

## Sequential Identifiers

- Do not generate user-visible sequential identifiers with `MAX(value) + 1` unless the read is protected by a transactionally safe lock.
- Prefer a database sequence, counter table with `SELECT ... FOR UPDATE`, or advisory lock scoped to the sequence partition such as month prefix.
- Keep the uniqueness constraint, but do not rely on unique-constraint failures as normal concurrency control.

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
- Public tracker-style endpoints should map to a dedicated public DTO/schema instead of reusing admin detail DTOs.

## Migrations

- Put SQL migrations under `apps/backend/migrations`.
- Keep migrations forward-only and deterministic.
- Add indexes for lookup and locking columns touched by new repository queries.

## Config And Server Safety

- Reject unsafe production config such as empty DB passwords or disabled SSL outside local/test.
- Keep CORS origins explicit in production.
- Use context timeouts around connection retries and DB operations that can hang.
