# Review And Ops Checklist

## Code Review

Lead with findings, ordered by severity. Include file and line references.

Check:

- Architecture: changed code follows existing package, route, service, and component boundaries.
- Security: no secrets, permissive CORS, unsafe cookie flags, public PII leaks, or unauthenticated internal identifiers.
- Data integrity: migrations, indexes, transaction boundaries, locking, and rollback behavior are correct.
- API contracts: Go JSON tags, response envelope, frontend interfaces, mocks, and seed data align.
- Testing: risk surface has appropriate unit, integration, or E2E coverage.
- Frontend: Svelte 5 runes are used correctly; async state cannot be overwritten by stale requests; loading/error states exist.
- UI: layout fits mobile and desktop; text does not overflow; active and hover states remain readable.
- Performance: database queries use indexed lookups and avoid avoidable N+1 behavior.

## Environment And Config

- Do not commit real `.env` files.
- Keep `.env.example` current when variables change.
- Reject unsafe production settings such as disabled DB SSL or empty DB passwords.
- Prefer Makefile targets over ad hoc commands when they exist.

## Docker And Compose

- Use health checks for PostgreSQL services.
- Keep build contexts clean with `.dockerignore`.
- Use multi-stage images for backend and frontend builds.
- Pin base image tags; do not use `latest`.
- Run final production containers as non-root.
- Do not bake local `node_modules`, `.git`, or secrets into images.

## Commit Readiness

- Inspect `git status` before staging.
- Stage only the intended logical change.
- Avoid `git add .` when unrelated work exists.
- Commit messages should be Conventional Commits: `type(scope): description`.
- Avoid commit descriptions joined by `and`; split mixed changes.
- Run relevant checks before committing or state why they could not be run.

## Technical Debt

- Use `docs/tech-debt.md` for deliberate, deferrable compromises when that file exists.
- Use `hotissue.md` for active blocking bugs when that convention exists.
- Do not add debt notes for routine follow-ups that can be completed inside the current change.
