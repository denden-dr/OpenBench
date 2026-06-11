---
name: requesting-full-code-review
description: Use when performing a comprehensive codebase-wide review across all directories (apps, config, environments, orchestration) to maintain architectural and pattern consistency.
---

# Requesting Full Code Review

## Overview
A codebase-wide review examines the entire repository structure, configurations, and cross-application interactions. This ensures that changes in one app (e.g., backend API contract) do not break or mismatch configurations in another (e.g., frontend client fetch configurations or orchestration setups).

## When to Use
- Before a release to production.
- When reorganizing directory architecture or adding new microservices.
- To audit the health of global variables, cross-origin rules, and Docker-Compose setups.

## Core Process

### 1. File & Directory Audit
Traverse all folders under the repository root to evaluate:
- **Project Structure**: Verify adherence to standard monorepo layouts (e.g., frontend and backend separated under `/apps`).
- **Configuration Alignment**: Ensure the `.env` configuration keys match the loaded variables in `apps/backend/internal/config` and the environment declarations in `docker-compose.yml`/`docker-compose-test.yml`.
- **Cross-Service Integrations**: Verify SvelteKit frontend `PUBLIC_API_URL` environment variables match the exposed ports of the backend service container.

### 2. Comprehensive Test Verification
Execute the full test suites of all services in the monorepo to ensure zero regressions:
```bash
# Verify backend
cd apps/backend && go test ./...
# Verify frontend
cd apps/frontend && npm run check
```

### 3. Tech Debt & Documentation Check
Ensure any compromise or unresolved design issue spotted during the review is documented in `docs/tech-debt.md`.

## Response Patterns
Present findings grouped by:
- Architecture & Monorepo Health
- Cross-Service Configuration Alignment
- Complete Test Verification Results
- Registered Technical Debt Items
