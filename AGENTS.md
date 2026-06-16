# Repository Guidelines

## Project Structure & Module Organization

OpenBench is a monorepo with a Go/Fiber API in `apps/backend`, a SvelteKit frontend in `apps/frontend`, and Playwright E2E browser tests in `apps/e2e`. Backend entrypoint code lives in `apps/backend/cmd/api/main.go`; domain packages live under `apps/backend/internal/`, including `auth`, `config`, `database`, `health`, and shared `pkg` helpers. SQL migrations are in `apps/backend/migrations`.

Frontend routes are in `apps/frontend/src/routes`, shared components and services are in `apps/frontend/src/lib`, and static assets are in `apps/frontend/static`. Browser E2E specs are in `apps/e2e/tests`. Root-level `docker-compose*.yml` files and `Makefile` targets coordinate local and test environments.

## Build, Test, and Development Commands

- `make dev`: start the dev database, run migrations, then launch backend and frontend.
- `make dev-db-up` / `make dev-db-down`: start or stop the local PostgreSQL service.
- `make run-backend`: run the API from `apps/backend`.
- `make run-frontend`: run SvelteKit against the API.
- `make run-frontend-mock`: run SvelteKit with `MOCK_API=true`.
- `make test-backend`: run Go unit and integration test targets.
- `make test-frontend`: run frontend Vitest tests.
- `make test-e2e-mock`: run Playwright E2E tests in mock mode.
- `make test-e2e-dev` / `make test-e2e-env`: run E2E tests against dev/test environments.
- `cd apps/frontend && npm run check`: run SvelteKit and TypeScript checks.

## Coding Style & Naming Conventions

Format Go with `gofmt`; `make test-unit` enforces formatting for backend source. Keep Go package names short, lowercase, and scoped under `internal/` unless public reuse is required. Use `*_test.go` for Go tests and `TestXxx` test names.

Use SvelteKit route conventions such as `+page.svelte` and `+layout.svelte`. Prefer TypeScript for frontend logic, colocate component tests near components as `*.test.ts`, and use existing service/component patterns before adding new abstractions.

## Testing Guidelines

Backend tests use Go `testing`, `testify`, SQL mocks, and Testcontainers for PostgreSQL integration coverage. Run `make test-unit` for fast checks and `make test-integration` when database behavior changes.

Frontend tests use Vitest with Testing Library for components. E2E browser tests use Playwright in `apps/e2e` for E2E flows. Run `npm run test:e2e:ui` from `apps/e2e` when debugging browser tests interactively.

## Commit & Pull Request Guidelines

Git history follows Conventional Commits with optional scopes, for example `feat(frontend): ...`, `chore(backend): ...`, and `docs: ...`. Keep commits focused on one logical change.

Pull requests should include a concise summary, linked issue when relevant, commands run, and screenshots for visible UI changes. Call out migrations, environment changes, or security-sensitive behavior explicitly.

## Security & Configuration Tips

Do not commit real `.env` files or secrets. Start from `apps/backend/.env.example` and `apps/frontend/.env.example`, and keep local credentials restricted to development or test environments.

## Agent Skills & Guidelines

This repository relies on several curated Agent Skills to enforce consistent architecture and styling. AI Agents MUST consult these skills before implementing changes:

- **Frontend & UI**: `frontend-svelte5-architecture` and `openbench-ui-design-system`
- **Backend & Integration**: `backend-go-architecture` and `fullstack-api-integration`
- **Quality & Workflow**: `openbench-testing-strategy` and `openbench-workflow-and-ops`

The full catalog and detailed instructions are located in `.agents/skills/CATALOG.md`.
