# OpenBench Skills Catalog

This catalog indexes all specialized skills configured in this repository. These skills serve as process guidelines, code standards, and patterns for AI agents and developers working on the OpenBench monorepo.

## Skills Directory

### Frontend & UI
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`frontend-svelte5-architecture`](./frontend-svelte5-architecture/SKILL.md) | 1.0.0 | Use when creating/refactoring Svelte 5 components, building interactive forms, handling global state services, formatting inputs, slicing large pages, and handling client-side routing. |
| [`openbench-ui-design-system`](./openbench-ui-design-system/SKILL.md) | 1.0.0 | Use when styling components, building dashboard layouts, managing viewport constraints, or applying the Neubrutalism design language with Tailwind CSS v4. |

### Backend
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`backend-go-architecture`](./backend-go-architecture/SKILL.md) | 1.0.0 | Use when building Go domain logic, handlers, repository layers, managing database transactions, securing JWT/cookie sessions, and building secure public trackers. |

### Fullstack & Integration
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`fullstack-api-integration`](./fullstack-api-integration/SKILL.md) | 1.0.0 | Use when implementing API contracts between Go and Svelte, mocking endpoints in the frontend for development, or aligning seed data and payload structures. |

### Testing
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`openbench-testing-strategy`](./openbench-testing-strategy/SKILL.md) | 1.0.0 | Use when writing or debugging Playwright E2E tests for the frontend, or Go unit/integration tests with Testify and Testcontainers. |

### Workflow & Ops
| Skill | Version | Description (Trigger Conditions) |
| :--- | :--- | :--- |
| [`openbench-workflow-and-ops`](./openbench-workflow-and-ops/SKILL.md) | 1.0.0 | Use when setting up dev environments, Docker/Podman compose, securing container builds, committing code (Atomic Commits), or performing code reviews. |

---

## Directory Conventions

Supporting resources for skills must be organized as follows:
*   `references/` - Extended guides, conceptual rules, and specific patterns that would otherwise bloat the main `SKILL.md`.
*   `assets/` - Reusable templates, boilerplate code, or configuration templates.

*Note: Skills without supporting resources should remain flat (only a `SKILL.md` file in their folder).*

## Merge History

| Version | Date | Change |
|---------|------|--------|
| 3.0.0 | 2026-06-15 | **Massive Consolidation**: Merged 18 scattered skills into 6 core, comprehensive skills (`frontend-svelte5-architecture`, `openbench-ui-design-system`, `backend-go-architecture`, `fullstack-api-integration`, `openbench-testing-strategy`, `openbench-workflow-and-ops`) to streamline AI agent context and focus. |
| 2.0.0 | 2026-06-13 | Merged `formatting-api-responses` and `initializing-go-fiber-api` → `implementing-repository-pattern` |
| 2.0.0 | 2026-06-13 | Merged `mocking-with-mockery` → `testing-with-testify` |
| 1.2.0 | 2026-06-13 | Merged `documenting-technical-debt` → `requesting-code-review` |
| 1.0.0 | 2026-06-14 | Created `securing-public-trackers-with-uuids`, `mocking-fullstack-endpoints-in-frontend`, `managing-forms-and-types-in-svelte5` |
| 2.0.0 | 2026-06-14 | Merged `securing-cookie-auth-handlers` + `implementing-token-rotation-security` → `securing-auth-sessions` |
| 2.0.0 | 2026-06-14 | Merged `configuring-postgres-compose` + `managing-multi-environment-config` → `managing-dev-environment` |
| 2.0.0 | 2026-06-14 | Merged `managing-forms-and-types-in-svelte5` → `developing-ui-svelte-best-practices` v2.0 |
| 1.0.0 | 2026-06-14 | Created `slicing-svelte5-components` |
| 2.1.0 | 2026-06-14 | Updated `developing-ui-svelte-best-practices` with `$effect` and `svelte.ts` patterns |
| 1.0.0 | 2026-06-14 | Created `svelte5-global-state-services`, `handling-formatted-inputs-svelte`, `neubrutalism-ui-design-system`, `responsive-dashboard-layouts` |
