# Repository Guidelines

## Project Structure & Module Organization

This repository is a small monorepo. Backend code lives in `apps/backend`, with the Fiber entrypoint at `cmd/api/main.go` and internal packages under `internal/` (`config`, `database`). Tests sit beside packages, e.g. `internal/database/database_test.go`.

Frontend code lives in `apps/frontend` using SvelteKit. Routes are in `src/routes`, shared code/assets are in `src/lib`, and static files are in `static`. Root compose files and `Makefile` manage local services. Track active fixes in `issue.md` and non-urgent debt in `docs/tech-debt.md`.

## Build, Test, and Development Commands

- `make dev-db-up`: start the development PostgreSQL service via Podman Compose.
- `make dev-db-down`: stop the development database.
- `make run-backend`: run the Go API locally from `apps/backend`.
- `make test-backend`: run backend Go tests with `go test -v ./...`.
- `make test-env-build`: build the test compose images.
- `make test-env-up` / `make test-env-down`: manage the test environment.
- `cd apps/frontend && npm run dev`: start SvelteKit dev server.
- `cd apps/frontend && npm run check`: run SvelteKit and TypeScript checks.
- `cd apps/frontend && npm run build`: build the frontend.

## Coding Style & Naming Conventions

Use `gofmt` for Go files before committing. Keep Go package names lowercase and short, and keep implementation under `internal/` unless it must be public. Use existing `DB_` and `APP_` env prefixes.

For Svelte, keep route files in SvelteKit form such as `+page.svelte` and `+layout.svelte`. Prefer TypeScript for frontend logic.

## Testing Guidelines

Backend tests use Go's standard `testing` package. Name tests `TestXxx` and place them in `*_test.go` files beside the package under test. Current database tests require PostgreSQL, so rerun with `go test -count=1 ./...` for environment-sensitive behavior.

Frontend validation currently uses `npm run check`; add dedicated component or integration tests when frontend behavior becomes non-trivial.

## Commit & Pull Request Guidelines

History follows Conventional Commits with optional scopes, for example `feat(frontend): ...`, `fix(backend): ...`, `test(backend): ...`, and `refactor: ...`. Keep commits focused.

Pull requests should include a short summary, linked issue or debt item when applicable, commands run, and screenshots for visible frontend changes. Call out config, database, or security-impacting changes explicitly.

## Security & Configuration Tips

Do not commit real `.env` files or secrets. Use `apps/backend/.env.example` as a template, keep local credentials out of images, and prefer localhost-only service exposure during development.
