# Testing Matrix

## Backend Commands

- Fast backend checks: `make test-unit`
- Full backend tests: `make test-backend`
- Targeted package test: `cd apps/backend && go test ./internal/auth -run TestName`
- Generate mocks after interface changes: `cd apps/backend && go generate ./...`

## Frontend Commands

- Type and Svelte check: `cd apps/frontend && npm run check`
- Vitest tests: `make test-frontend`
- Targeted frontend tests: `cd apps/frontend && npm test -- path-or-pattern`
- Mock-mode Playwright: `make test-frontend-mock`

## Test Selection

- Pure formatting/helper change: focused unit test or existing helper test.
- Frontend service parser: service Vitest test with mocked `fetch`.
- Mock DB/auth behavior: colocated Vitest test and explicit storage cleanup.
- Svelte component state: component test when behavior is isolated; Playwright when the behavior depends on routing or browser workflow.
- Repository SQL: integration test with Testcontainers.
- Service transaction or auth behavior: unit test with mocks plus integration test when DB locking or persistence matters.
- API contract change: backend handler/service test, frontend service test, and mock update.

## Backend Patterns

- Use `require.NoError(t, err)` for setup that must stop the test.
- Use `assert.Equal(t, expected, actual)` for independent assertions after setup succeeds.
- Use `testutil.IntegrationSuite` for database-backed tests.
- Terminate Testcontainers once through `TestMain` or the existing suite helper.
- Keep Ryuk enabled unless the user explicitly directs otherwise.
- Do not hardcode container ports; use mapped ports.

## Frontend Patterns

- Clear `localStorage` and `sessionStorage` in `beforeEach` when tests touch mock data or auth.
- Mock `fetch` at the service boundary for service tests.
- Assert credentials and envelope parsing for protected API calls.
- Keep tests colocated next to services when practical.

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
