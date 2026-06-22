# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## Changelog Standard

When writing entries to this changelog, adhere to the following rules:

1. **Target Audience**: Write entries for human developers and stakeholders, not automated machines or git commit summaries.
2. **Version Headers**: Use the format `## [Version] - YYYY-MM-DD` (e.g. `## [0.1.0] - 2026-06-11`).
3. **Chronological Order**: The latest version or `## [Unreleased]` changes must always be at the top of the file.
4. **Change Classification**: Group changes under these explicit bulleted sections:
   - `### Added` for new features or modules.
   - `### Changed` for changes in existing behavior or style.
   - `### Deprecated` for soon-to-be-removed features.
   - `### Removed` for deleted features.
   - `### Fixed` for any bug fixes.
   - `### Security` in case of vulnerabilities fixed or authentication updates.

---

## [Unreleased]

### Added
- Dedicated `apps/e2e` Playwright package with its own package metadata, TypeScript config, browser tests, and E2E scripts.
- E2E Makefile targets: `test-e2e-mock`, `test-e2e-dev`, and `test-e2e-env`, with backward-compatible aliases for the old frontend E2E target names.
- Project skill reference guides for frontend architecture, UI patterns, backend architecture, API contracts, testing strategy, and workflow/review operations.
- Authentication module with sign-in, token rotation (RTR), sign-out, and session info
- Database migrations for users and refresh_tokens tables
- Dev seeder creating default admin account
- Health probe module (liveness, readiness with DB pool stats)
- Standardized API response helpers and struct validation
- Admin sign-in page with neubrutalism-styled form
- Admin dashboard with metrics cards, ticket list, system info
- Reusable UI components: Card, Button, Input
- Real and mock authentication services with session caching
- Landing page with hero section, feature cards, auth-aware CTA
- Neubrutalism design system: fonts, custom Tailwind theme, shadows
- Test infrastructure: Vitest, jsdom, testing-library/svelte
- Component and service unit tests (Button, auth, mockAuth)
- Vite proxy configuration for `/api` and `/health` routes
- Makefile targets: `dev`, `migrate-up`/`migrate-down`, `run-frontend`, `test-frontend`
- Frontend `.env.example` and `.env` files

### Changed
- Playwright E2E ownership moved out of `apps/frontend` into `apps/e2e`, leaving the frontend package focused on SvelteKit, Vitest, and UI development.
- `AGENTS.md` updated to document the new `apps/e2e` test location and E2E commands.
- Agent skills were refactored into short runbooks with selectively loaded reference files for lower context overhead and clearer implementation guidance.
- Home page rewritten from SvelteKit placeholder to feature landing page
- Layout CSS expanded from bare Tailwind import to full neubrutalism theme
- `vite.config.ts` switched to `vitest/config`, added proxy and test config
- `app.html` now loads Google Fonts (Plus Jakarta Sans, Space Grotesk, Space Mono)
- `index.ts` now exports Card, Button, Input as barrel module
- Backend config normalizes `"dev"` → `"development"` and validates JWT/production settings
- CORS fallback deduplicated — uses config default only
- `docker-compose-test.yml` added `CORS_ALLOWED_ORIGINS`, `JWT_SECRET`, `PUBLIC_MOCK_API` env vars
- Refactored Playwright E2E tests to dynamically manage tickets and capture generated UUIDs, allowing E2E suites to run successfully on both Mock and real environments (with empty databases).
- Skipped E2E tests for non-integrated Inventory and Sales features.

### Removed
- Playwright scripts, configuration, and dependency ownership from `apps/frontend`.
- Legacy Playwright test files from `apps/frontend/tests` after moving them to `apps/e2e/tests`.

### Fixed
- Bounded containerized E2E startup and readiness waits with clear timeout errors to prevent `make test-e2e-env` from hanging indefinitely.
- Ensured `test-e2e-env` tears down the compose test environment through a shell trap on failure or interruption.
- Fixed Makefile `test-e2e-env` exit trap to run compose down relative to `$(CURDIR)` preventing command failure when exiting from `apps/e2e` subdirectory.
- Fixed strict mode violation on "received" status locator in `tracker.test.ts` using exact text matching.

### Security
- JWT access tokens with HMAC-SHA256 signing and algorithm-switching protection
- Bcrypt password hashing (cost factor 12)
- HTTPOnly cookies with Secure flag and SameSite Strict
- Refresh token rotation with replay-attack detection and family revocation
- Role-based access control middleware
- Production-only validation: rejects empty DB passwords, disabled SSL, default JWT secret

---
