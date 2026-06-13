# OpenBench Skills Catalog

This catalog indexes all specialized skills configured in this repository. These skills serve as process guidelines, code standards, and patterns for AI agents and developers working on the OpenBench monorepo.

## Skills Directory

| Skill | Domain | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- | :--- |
| [`aligning-frontend-and-backend-contracts`](./aligning-frontend-and-backend-contracts/SKILL.md) | Fullstack | 1.0.0 | Use when implementing API contracts between Go backend and Svelte frontend, verifying database row locking for state transitions, and aligning seeder/mock user credentials. |
| [`configuring-postgres-compose`](./configuring-postgres-compose/SKILL.md) | Infrastructure | 1.0.0 | Use when setting up or debugging PostgreSQL containers, connection pooling, or docker-compose database services. |
| [`developing-ui-svelte-best-practices`](./developing-ui-svelte-best-practices/SKILL.md) | Frontend | 1.1.0 | Use when creating or modifying Svelte 5 UI components, route guards, or Vitest component tests. |
| [`implementing-repository-pattern`](./implementing-repository-pattern/SKILL.md) | Backend | 2.0.0 | Use when adding new domain packages to the Go backend, refactoring handler-to-database coupling, or initializing the Fiber server. |
| [`managing-multi-environment-config`](./managing-multi-environment-config/SKILL.md) | Infrastructure | 1.0.0 | Use when setting up environment-based configurations in a Go application, managing different environments (development, testing, production) using .env files, and loading variables cleanly. |
| [`practicing-atomic-commits`](./practicing-atomic-commits/SKILL.md) | Process | 1.0.0 | Use when staging and committing code changes to Git, or when git history needs to be clean, readable, and easy to roll back. |
| [`requesting-code-review`](./requesting-code-review/SKILL.md) | Process | 1.2.0 | Use when performing code reviews in the repository, either scoped to feature diffs or across the entire codebase to verify monorepo health, security, testing robustness, and config alignment. |
| [`secure-containerization`](./secure-containerization/SKILL.md) | Infrastructure | 1.0.0 | Use when containerizing applications (Go, SvelteKit, Node.js) with Docker or Podman, ensuring non-root execution, pinned image tags, multi-stage builds, and dockerignore configurations. |
| [`testing-frontend-e2e-with-playwright`](./testing-frontend-e2e-with-playwright/SKILL.md) | Testing | 1.1.0 | Use when SvelteKit browser tests are flaky, timing out on external resources (fonts/APIs), or failing due to client-side hydration race conditions. |
| [`testing-with-testify`](./testing-with-testify/SKILL.md) | Testing | 2.0.0 | Use when Go backend unit or database integration tests are failing, slow, resource-leaking, require mock dependencies, or need mock auto-generation with mockery. |

---

## Directory Conventions

Supporting resources for skills must be organized as follows:
*   `references/` - Extended guides, conceptual rules, and specific patterns that would otherwise bloat the main `SKILL.md`.
*   `assets/` - Reusable templates, boilerplate code, or configuration templates.

*Note: Skills without supporting resources should remain flat (only a `SKILL.md` file in their folder).*

## Merge History

| Version | Date | Change |
|---------|------|--------|
| 2.0.0 | 2026-06-13 | Merged `formatting-api-responses` and `initializing-go-fiber-api` → `implementing-repository-pattern` |
| 2.0.0 | 2026-06-13 | Merged `mocking-with-mockery` → `testing-with-testify` |
| 1.2.0 | 2026-06-13 | Merged `documenting-technical-debt` → `requesting-code-review` |
