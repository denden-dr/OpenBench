# Repository Guidelines

## Project Structure & Module Organization

OpenBench is a phone repair admin app split into two apps:

- `apps/backend`: Go Fiber API, PostgreSQL access through `sqlx`, migrations in `apps/backend/migrations`.
- `apps/backend/internal`: backend layers: `handler`, `service`, `repository`, `dto`, `model`, `config`.
- `apps/backend/internal/database`: shared database connection setup.
- `apps/backend/mocks`: generated Mockery mocks.
- `apps/frontend`: SvelteKit/Svelte 5 frontend. Source is in `apps/frontend/src`; static assets are in `apps/frontend/static`.
- `docs`: product specs and implementation plans.
- `compose.yaml`, `compose.test.yaml`: local and test container orchestration.

## Build, Test, and Development Commands

- `make compose-up`: start local containers with Docker or Podman.
- `make compose-down`: stop local containers.
- `make migrate-up`: apply backend SQL migrations to the local database.
- `make migrate-create NAME=add_field`: create a numbered migration.
- `make run-backend`: run the Go API from `apps/backend`.
- `make run-frontend`: run the SvelteKit dev server.
- `make up`: start database, backend, and frontend.
- `make test-backend-unit`: run `go test ./... -v`.
- `make test-backend-integration`: run integration-tagged Go tests.
- `cd apps/frontend && npm run check`: run Svelte type checking.
- `cd apps/frontend && npm run build`: build the frontend.

## Coding Style & Naming Conventions

Use `gofmt` for Go files and keep package names short and lowercase. Keep backend layers strict: handlers parse HTTP, services enforce business rules, repositories perform SQL. Name Go tests `*_test.go`; integration tests use the `integration` build tag.

Frontend code uses TypeScript, SvelteKit route files (`+page.svelte`, `+layout.svelte`), and Tailwind CSS. Keep route code colocated under `src`.

## Testing Guidelines

Backend tests use Go `testing`, `testify`, and Mockery mocks. Run `make mock-backend` after changing repository interfaces. Integration tests use Testcontainers/PostgreSQL; run `make test-backend-integration` before merging database changes, handler changes, or integration test changes.

Frontend changes should pass `npm run check`; use `npm run build` for production validation.

## Commit & Pull Request Guidelines

Recent history uses conventional prefixes such as `feat:`, `fix:`, `docs:`, and `chore:`. Keep commits focused, for example: `fix: enforce ticket pickup payment rules`.

Pull requests should include a summary, linked issue or plan, test results, and screenshots for UI changes. Call out migration, environment, or data-flow changes.

## Security & Configuration Tips

Do not commit secrets. Use `apps/backend/.env.example` for local configuration. Review down migrations with schema changes, and document container assumptions in Makefile targets or PR notes.

## Runner Notes

- When using `agent-browser` in this workspace, set `XDG_RUNTIME_DIR` to a writable path under `/tmp` before opening a browser session.
- When starting `compose-test-up` or other rootless Podman flows from the Codex sandbox, prefer running the command outside the sandbox with escalated permissions if the local runtime cannot write to `/run/user/*` or set up user namespaces.

## graphify

This project has a graphify knowledge graph at graphify-out/.

Rules:
- Before answering architecture or codebase questions, read graphify-out/GRAPH_REPORT.md for god nodes and community structure
- If graphify-out/wiki/index.md exists, navigate it instead of reading raw files
- For cross-module "how does X relate to Y" questions, prefer `graphify query "<question>"`, `graphify path "<A>" "<B>"`, or `graphify explain "<concept>"` over grep — these traverse the graph's EXTRACTED + INFERRED edges instead of scanning files
- After modifying code files in this session, run `graphify update .` to keep the graph current (AST-only, no API cost)
- Do not commit graphify artifacts automatically. Leave `graphify-out/` changes unstaged until all code review issues are resolved, then commit them only when explicitly instructed.

## Deployment & Schema Migrations

- Schema migrations that drop columns (such as `000005_remove_warranty_expiry_date`) require a downtime deployment window or an expand/contract rollout strategy to be rolling-compatible.
- Production deployments for OpenBench currently assume a single-instance deployment model with scheduled downtime for migrations.
