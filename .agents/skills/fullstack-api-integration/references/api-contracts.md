# API Contracts

## Files To Inspect First

- Backend DTOs and handlers: `apps/backend/internal/**/dto.go`, `handler.go`
- Response helper: `apps/backend/internal/pkg/response/response.go`
- Frontend services: `apps/frontend/src/lib/services/*.ts`
- Mock services: `apps/frontend/src/lib/services/mocks`
- Seed data: `apps/backend/internal/database/seeder.go` and `apps/frontend/src/lib/services/mocks/seed.ts`

## Contract Source

The Go response DTO and JSON tags define the wire contract. TypeScript interfaces must either match those field names or explicitly map them.

Backend success responses use:

```json
{
  "code": 200,
  "message": "Message",
  "data": {}
}
```

Frontend consumers should read `resBody.data`, not top-level business fields.

## Naming Rules

- Go JSON tags should use the field names the frontend receives.
- Existing APIs commonly use snake_case payload fields such as `user_id`.
- If frontend state prefers camelCase, map explicitly at the service boundary.
- Do not let UI components consume raw envelope objects.

## Auth And Cookies

- Auth-protected frontend requests must include `credentials: 'include'`.
- Backend set-cookie and clear-cookie flags must match.
- Session checks should treat failed refresh as unauthenticated and clear local session state.

## Mock API

- Keep mock files under `apps/frontend/src/lib/services/mocks`.
- Split static seed data, DB-like data behavior, auth behavior, network interception, and shared mock types.
- Use distinct storage keys per domain.
- Use `sessionStorage` or explicit active-session storage for auth state, and `localStorage` for persisted mock data.
- Add artificial latency around 300-600 ms for async methods so loading states execute.
- Mirror real method signatures and important business rules.

## Mock Toggle Reconciliation

Inspect current `isMockEnabled()` before changing mock toggle priority. If behavior should change, update:

- implementation
- frontend tests
- `.env.example`
- package scripts
- Playwright assumptions

Do not rely on stale skill text when code and tests already define an intentional priority.

## Seed Parity

- Keep demo credentials aligned between backend seeder and frontend mock auth.
- Use realistic records that exercise status, payment, inventory, warranty, and error states.
- Avoid placeholder names like `foo` or `test` unless the test explicitly asserts malformed input handling.
