---
name: openbench-testing-strategy
description: Use when writing or debugging Playwright E2E tests for the frontend, or Go unit/integration tests with Testify and Testcontainers.
version: 1.0.0
---

# OpenBench Testing Strategy

## Overview
Testing covers frontend E2E via Playwright and backend isolated integration tests using `testify` and `testcontainers-go`. Tests must be reliable, hermetic, and fast.

## Backend: Testing with Testify & Testcontainers
1. **Isolated Integration Suites**: Use `testutil.IntegrationSuite` for database tests to spin up Postgres inside a Docker container using testcontainers-go.
2. **Mandatory Teardown**: Teardown the container once in `TestMain` via `tdb.Terminate()`.
3. **Keep Ryuk Enabled**: Ryuk safely cleans up orphaned containers on crash. Pass `TESTCONTAINERS_RYUK_DISABLED=false` unless explicitly overridden.
4. **Migration Pool**: Use a separate DB connection for migrations before testing to avoid connection pool contention.
5. **Mockery**: Use `//go:generate mockery` to auto-generate interface mocks. Run `go generate ./...` and use `mocks.NewRepository(t)`.

### Backend Assertions
- `require.NoError(t, err)`: Stops execution if setup fails.
- `assert.Equal(t, exp, act)`: Logs failure, continues test.

## Frontend: Unit Testing Services (Vitest)
1. **Mock Environment Cleanup**: When unit testing mock auth or database services (e.g., `mocks/auth.test.ts`), always explicitly clear `sessionStorage` and `localStorage` in the `beforeEach` hook to prevent test state pollution.
2. **Collocation**: Place unit tests next to the service they test (e.g., `mocks/auth.test.ts` next to `mocks/auth.ts`) to maintain a highly organized module structure.

## Frontend: E2E Testing with Playwright
1. **Safe Page Hydration**: Always wait for Svelte hydration before interacting to prevent native form submission.
   ```typescript
   await page.waitForSelector('main[data-hydrated="true"]');
   ```
2. **Network Interception**: Abort external requests (like Google Fonts) in `test.beforeEach` to prevent network timeouts in CI.
3. **Environment Alignment**: Validate mock mode variables. If testing against a mock, assert `process.env.PUBLIC_MOCK_API === 'true'`.
4. **Service Readiness**: Use Makefile checks (e.g., `curl /api/health`) or Playwright `webServer` config to ensure backend/frontend servers are fully active before running the test suite.

## Common Mistakes
- **Playwright Mismatches**: Clicking buttons before Svelte attaches events, causing URL query bloat instead of JS fetches.
- **Vite Dev Server Resource Starvation**: High parallel workers in Vite Dev Server for E2E. Use `workers: 1` and `fullyParallel: false`.
- **Backend Port Collisions**: Hardcoding test database ports instead of letting Testcontainers allocate dynamic ports.
- **Unmet Mock Expectations**: Failing to assert mock invocations.
- **Container Orphans**: Disabling Ryuk globally, leading to dangling test containers.
