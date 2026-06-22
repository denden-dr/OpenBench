# Review And Ops Checklist

## Code Review

Lead with findings, ordered by severity. Include file and line references.

Check:

- Architecture: changed code follows existing package, route, service, and component boundaries.
- Security: no secrets, permissive CORS, unsafe cookie flags, public PII leaks, or unauthenticated internal identifiers.
- Data integrity: migrations, indexes, transaction boundaries, locking, and rollback behavior are correct.
- API contracts: OpenAPI, generated Go/TypeScript types, Go JSON tags, response envelope, frontend services, mocks, and seed data align.
- Testing: risk surface has appropriate unit, integration, or E2E coverage.
- Frontend: Svelte 5 runes are used correctly; async state cannot be overwritten by stale requests; loading/error states exist.
- Empty data: list endpoints, service parsers, and UI empty states handle zero rows without missing `data`, nullable array state, or stuck skeletons.
- UI: layout fits mobile and desktop; text does not overflow; active and hover states remain readable.
- Performance: database queries use indexed lookups and avoid avoidable N+1 behavior.

Extra checks for API/domain changes:

- If `docs/api/openapi.yml` changed, generated Go and TypeScript files are refreshed or explicitly verified current.
- Public unauthenticated routes return narrow public DTOs and do not reuse admin/internal schemas.
- Create/update frontend payloads match the OpenAPI request schemas exactly.
- Successful list endpoints return stable array payloads, including `data: []` for empty results.
- State-transition side effects compare previous and next state, and terminal-state derived data does not move on later edits.
- Mock API side effects match backend side effects, including dates and derived rows.

## Environment And Config

- Do not commit real `.env` files.
- Keep `.env.example` current when variables change.
- Reject unsafe production settings such as disabled DB SSL or empty DB passwords.
- Prefer Makefile targets over ad hoc commands when they exist.

## Dependency And Toolchain

- Prefer resolving package-manager peer conflicts by aligning versions in manifests and lockfiles.
- Treat `npm ci` failures as dependency graph problems first, not Dockerfile problems.
- Do not add `npm ci --force` or `npm ci --legacy-peer-deps` to Dockerfiles unless it is an explicitly documented temporary workaround.
- When changing frontend toolchain dependencies, verify `apps/frontend/package.json` and `apps/frontend/package-lock.json` together.
- Check peer dependency requirements for TypeScript, Svelte, Vite, Vitest, and OpenAPI generation tools before upgrading major versions.

## Docker And Compose

- Use health checks for PostgreSQL services.
- Keep build contexts clean with `.dockerignore`.
- Use multi-stage images for backend and frontend builds.
- Pin base image tags; do not use `latest`.
- Run final production containers as non-root.
- Do not bake local `node_modules`, `.git`, or secrets into images.
- For Node images that run `npm ci`, verify lockfile reproducibility locally before changing Docker install flags.

## Commit Readiness

- Inspect `git status` before staging.
- Stage only the intended logical change.
- Avoid `git add .` when unrelated work exists.
- Commit messages should be Conventional Commits: `type(scope): description`.
- Avoid commit descriptions joined by `and`; split mixed changes.
- Run relevant checks before committing or state why they could not be run.

## Technical Debt

- Use root `issue.md` for active blocking bugs found during review.
- Use `docs/tech-debt.md` for deliberate, deferrable compromises when that file exists.
- If both are relevant, record blockers in `issue.md` and lower-risk follow-ups in `docs/tech-debt.md`.
- Do not add debt notes for routine follow-ups that can be completed inside the current change.
