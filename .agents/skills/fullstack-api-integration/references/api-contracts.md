# API Contracts

## Contract Source

The **OpenAPI 3.0 spec** at `docs/api/openapi.yml` is the single source of truth for all API contracts.

All generated types, DTOs, and service mocks must align with this spec.

## Files To Inspect First

- OpenAPI spec: `docs/api/openapi.yml`
- Generated Go types: `apps/backend/internal/pkg/api/openapi_types.gen.go`
- Generated TypeScript types: `apps/frontend/src/lib/api/openapi.gen.ts`
- Backend DTOs and handlers: `apps/backend/internal/**/dto.go`, `handler.go`
- Response helper: `apps/backend/internal/pkg/response/response.go`
- Frontend services: `apps/frontend/src/lib/services/*.ts`
- Mock services: `apps/frontend/src/lib/services/mocks`
- Seed data: `apps/backend/internal/database/seeder.go` and `apps/frontend/src/lib/services/mocks/seed.ts`

## Spec-First Workflow

1. **Edit** `docs/api/openapi.yml` when adding or changing endpoints, fields, or response shapes.
2. **Regenerate** code:
   - `make generate-api-go` — generates Go types from the spec
   - `make generate-api-ts` — generates TypeScript types from the spec
   - `make generate-api-types` — runs both
3. **Implement** backend handlers using the generated Go types.
4. **Update** frontend services to use the generated TypeScript types.
5. **Sync** mock service response shapes and seed data.

## Generated Type Gate

- After any OpenAPI edit, run `make generate-api-types` unless generation tooling is unavailable.
- If generated files do not change, inspect the generated Go and TypeScript types to confirm the edited fields already exist.
- Treat generated drift as a blocking contract bug, even if handwritten interfaces still compile.
- Prefer importing generated request/response types at service boundaries over duplicating shapes by hand.

## Generator Toolchain Compatibility

- Keep `openapi-typescript` compatible with the installed TypeScript version.
- Before changing TypeScript or OpenAPI generator versions, inspect peer dependency ranges in `apps/frontend/package-lock.json` or package metadata.
- Resolve generator peer conflicts by aligning versions in `apps/frontend/package.json` and regenerating `apps/frontend/package-lock.json`.
- Do not rely on `npm ci --legacy-peer-deps` for generated-type tooling; it can hide incompatible generator/runtime assumptions.
- After toolchain changes, run `make generate-api-types` and `cd apps/frontend && npm run check` when feasible.

## Response Envelope

All endpoints follow this envelope:

```json
{
  "code": 200,
  "message": "Message",
  "data": {}
}
```

Frontend consumers should read `resBody.data`, not top-level business fields.

List endpoints are a first-class contract:

- Successful list responses must include `data`, even when the list is empty.
- Empty list payloads must serialize as `[]`, not as a missing field.
- Backend response helpers must not use `omitempty` in a way that drops successful list `data`.
- Frontend list parsers should defensively normalize missing or null `data` to `[]` only as a compatibility guard, not as a substitute for fixing the backend contract.
- Mocks must exercise both populated and empty list responses.

## Naming Rules

- All JSON fields use **snake_case** (e.g., `user_id`, `ticket_number`, `payment_status`).
- Frontend services should map to camelCase at the service boundary if UI state prefers it.
- Do not let UI components consume raw envelope objects.

## Request Payload Discipline

- Create requests must send only fields defined by the OpenAPI create schema.
- Update requests must send only fields defined by the OpenAPI update schema.
- Do not use full response/domain types such as `Ticket` as create or update request types when `TicketCreate` or `TicketUpdate` exists.
- Backend may ignore unknown JSON today; frontend must still honor the documented schema.

## Public Endpoint Contracts

- Public unauthenticated endpoints must use a dedicated public schema when the admin schema includes PII, financial fields, internal workflow fields, or staff notes.
- Public mock responses must include every field the public UI reads and must omit fields intentionally hidden from customers.

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
- For derived side effects, compute mock results from the merged post-update record, not from stale pre-update values.
- When backend creates derived rows or dates, mock mode must produce the same rows/dates for the same request.

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
