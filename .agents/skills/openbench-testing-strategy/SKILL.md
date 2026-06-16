---
name: openbench-testing-strategy
description: Test OpenBench backend and frontend changes. Use when adding or debugging Go unit tests, Testify mocks, Testcontainers PostgreSQL integration tests, Vitest service/component tests, Playwright E2E tests, mock-mode browser flows, or CI test reliability issues.
---

# OpenBench Testing Strategy

## Operating Rule

Choose the smallest test that proves the changed behavior, then add broader integration or E2E coverage only when contracts, database behavior, auth, or user workflows are affected.

## Workflow

1. Identify the risk surface: pure function, service parser, UI state, SQL/repository, auth flow, or browser workflow.
2. Add tests at the matching layer and keep fixtures realistic.
3. Reset mock state in every frontend test that touches `localStorage` or `sessionStorage`.
4. For DB behavior, use `testutil.IntegrationSuite` and Testcontainers instead of a shared local database.
5. For Playwright, wait for Svelte hydration before interaction and keep workers constrained for Vite reliability.
6. Report exactly which commands passed or could not be run.

## Load References

- Read `references/testing-matrix.md` before adding tests, changing test infrastructure, debugging Playwright flakes, or modifying integration test setup.

## Hard Checks

- Do not hardcode database ports for tests.
- Do not disable Ryuk globally unless the user explicitly requires it.
- Do not click Playwright forms before hydration.
- Do not leave browser tests dependent on external network resources.
