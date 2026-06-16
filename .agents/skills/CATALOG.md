# OpenBench Skills Catalog

This catalog indexes all specialized skills configured in this repository. These skills serve as short agent runbooks plus selectively loaded references for OpenBench development.

## Skills Directory

### Frontend & UI
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`frontend-svelte5-architecture`](./frontend-svelte5-architecture/SKILL.md) | 1.1.0 | Use when creating/refactoring Svelte 5 components, building interactive forms, handling rune-based global state services, formatting inputs, slicing route pages, async views, and client-side routing. |
| [`openbench-ui-design-system`](./openbench-ui-design-system/SKILL.md) | 1.1.0 | Use when styling Svelte components, building admin dashboards, public tracker pages, responsive sidebars, loading states, or Tailwind CSS v4 Neubrutalism UI. |

### Backend
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`backend-go-architecture`](./backend-go-architecture/SKILL.md) | 1.1.0 | Use when adding Go/Fiber domain models, repositories, services, handlers, SQL migrations, transactions, auth sessions, JWT/cookie behavior, public tracker endpoints, or database-backed tests. |

### Fullstack & Integration
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`fullstack-api-integration`](./fullstack-api-integration/SKILL.md) | 1.1.0 | Use when keeping Go API contracts, Svelte services, mock API behavior, seed data, response envelopes, JSON tags, TypeScript interfaces, auth fetches, and payload structures aligned. |

### Testing
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`openbench-testing-strategy`](./openbench-testing-strategy/SKILL.md) | 1.1.0 | Use when adding or debugging Go unit tests, Testify mocks, Testcontainers PostgreSQL tests, Vitest service/component tests, Playwright E2E tests, mock-mode browser flows, or CI reliability issues. |

### Workflow & Ops
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`openbench-workflow-and-ops`](./openbench-workflow-and-ops/SKILL.md) | 1.1.0 | Use when setting up local dev, Docker/Podman compose, environment config, container builds, Makefile workflows, commits, pull-request readiness, code reviews, security checks, or technical-debt tracking. |

---

## Directory Conventions

Supporting resources for skills must be organized as follows:
*   `references/` - Extended guides, concrete local patterns, command matrices, and checklists that would otherwise bloat the main `SKILL.md`.
*   `scripts/` - Deterministic helpers for repetitive validation or generation.
*   `assets/` - Reusable templates, boilerplate code, or configuration templates.

`SKILL.md` files should stay short and act as routers: operating rule, workflow, reference-loading rule, and hard checks. Detailed examples belong in directly linked reference files one level deep.

## Merge History

| Version | Date | Change |
|---------|------|--------|
| 3.1.0 | 2026-06-16 | **Operational Refactor**: Converted the 6 core skills into short runbooks, removed non-trigger frontmatter fields, and added domain references for frontend patterns, UI patterns, backend patterns, API contracts, testing matrix, and review/ops checklist. |
| 3.0.0 | 2026-06-15 | **Massive Consolidation**: Merged 18 scattered skills into 6 core, comprehensive skills (`frontend-svelte5-architecture`, `openbench-ui-design-system`, `backend-go-architecture`, `fullstack-api-integration`, `openbench-testing-strategy`, `openbench-workflow-and-ops`) to streamline AI agent context and focus. |
| 2.0.0 | 2026-06-13 | Merged `formatting-api-responses` and `initializing-go-fiber-api` -> `implementing-repository-pattern` |
| 2.0.0 | 2026-06-13 | Merged `mocking-with-mockery` -> `testing-with-testify` |
| 1.2.0 | 2026-06-13 | Merged `documenting-technical-debt` -> `requesting-code-review` |
| 1.0.0 | 2026-06-14 | Created `securing-public-trackers-with-uuids`, `mocking-fullstack-endpoints-in-frontend`, `managing-forms-and-types-in-svelte5` |
| 2.0.0 | 2026-06-14 | Merged `securing-cookie-auth-handlers` + `implementing-token-rotation-security` -> `securing-auth-sessions` |
| 2.0.0 | 2026-06-14 | Merged `configuring-postgres-compose` + `managing-multi-environment-config` -> `managing-dev-environment` |
| 2.0.0 | 2026-06-14 | Merged `managing-forms-and-types-in-svelte5` -> `developing-ui-svelte-best-practices` v2.0 |
| 1.0.0 | 2026-06-14 | Created `slicing-svelte5-components` |
| 2.1.0 | 2026-06-14 | Updated `developing-ui-svelte-best-practices` with `$effect` and `svelte.ts` patterns |
| 1.0.0 | 2026-06-14 | Created `svelte5-global-state-services`, `handling-formatted-inputs-svelte`, `neubrutalism-ui-design-system`, `responsive-dashboard-layouts` |
