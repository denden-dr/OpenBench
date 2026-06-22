# Testing Matrix

## Backend Commands

- Fast backend checks: `make test-unit`
- Full backend tests: `make test-backend`
- Targeted package test: `cd apps/backend && go test ./internal/auth -run TestName`
- Generate mocks after interface changes: `cd apps/backend && go generate ./...`

## Frontend Commands

- Strict dependency install: `cd apps/frontend && npm ci`
- Type and Svelte check: `cd apps/frontend && npm run check`
- Production frontend build: `cd apps/frontend && npm run build`
- Vitest tests: `make test-frontend`
- Targeted frontend tests: `cd apps/frontend && npm test -- path-or-pattern`
- Mock-mode Playwright: `make test-frontend-mock`

## Container And Compose Commands

- Frontend Docker dependency/build check: `make compose-test-build`

## Test Selection

- Pure formatting/helper change: focused unit test or existing helper test.
- Frontend service parser: service Vitest test with mocked `fetch`.
- Mock DB/auth behavior: colocated Vitest test and explicit storage cleanup.
- Svelte component state: component test when behavior is isolated; Playwright when the behavior depends on routing or browser workflow.
- Repository SQL: integration test with Testcontainers.
- Service transaction or auth behavior: unit test with mocks plus integration test when DB locking or persistence matters.
- API contract change: backend handler/service test, frontend service test, and mock update.
- List endpoint or empty-state change: test empty database/list responses, service parser normalization, and visible empty state after loading completes.
- OpenAPI schema change: regenerate types, run generated-type compile checks, and add at least one assertion using the new field on backend and frontend paths.
- Frontend dependency conflict: run strict install, inspect peer ranges, align manifest plus lockfile, then run frontend check/build and container build when Dockerfiles are affected.
- Public endpoint change: test that public payload includes required customer-visible fields and excludes PII/internal fields.
- State-transition side effect: test first transition, repeated terminal-state update, and same-request field change plus transition.

## Backend Patterns

- Use `require.NoError(t, err)` for setup that must stop the test.
- Use `assert.Equal(t, expected, actual)` for independent assertions after setup succeeds.
- Use `testutil.IntegrationSuite` for database-backed tests.
- Terminate Testcontainers once through `TestMain` or the existing suite helper.
- Keep Ryuk enabled unless the user explicitly directs otherwise.
- Do not hardcode container ports; use mapped ports.
- For sequential identifiers, add a concurrency test or document why database locking/sequence makes collision impossible.
- For derived records, assert the source row and derived row stay in sync.
- For list endpoints, assert successful empty results serialize as `data: []` when the frontend contract expects an array.

## Frontend Patterns

- Clear `localStorage` and `sessionStorage` in `beforeEach` when tests touch mock data or auth.
- Mock `fetch` at the service boundary for service tests.
- Assert credentials and envelope parsing for protected API calls.
- Keep tests colocated next to services when practical.
- Assert request bodies match OpenAPI create/update schemas and do not include response-only fields.
- Assert list service methods return `[]` for empty, missing, or null response data when defensive compatibility is required.
- Assert route components leave skeleton/loading state and show an empty state after a successful empty load.
- For mock parity, assert the same request produces the same derived dates/statuses as backend rules.
- For dependency fixes, prefer `npm ci` passing cleanly over Docker install bypass flags.

## Playwright Patterns

- Wait for hydration before interaction:

```typescript
await page.waitForSelector('main[data-hydrated="true"]');
```

- Keep `workers: 1` and `fullyParallel: false` for Vite-backed E2E stability unless the project proves otherwise.
- Abort or mock external network resources that can stall CI.
- Use mock mode for fast workflow tests; use composed backend/frontend only when validating real API integration.

## Reporting

In final responses, list commands run and whether they passed. If a command could not run because of Docker, network, or sandbox limitations, say that explicitly.
