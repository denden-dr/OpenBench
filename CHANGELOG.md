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
- Home page rewritten from SvelteKit placeholder to feature landing page
- Layout CSS expanded from bare Tailwind import to full neubrutalism theme
- `vite.config.ts` switched to `vitest/config`, added proxy and test config
- `app.html` now loads Google Fonts (Plus Jakarta Sans, Space Grotesk, Space Mono)
- `index.ts` now exports Card, Button, Input as barrel module
- Backend config normalizes `"dev"` → `"development"` and validates JWT/production settings
- CORS fallback deduplicated — uses config default only
- `docker-compose-test.yml` added `CORS_ALLOWED_ORIGINS`, `JWT_SECRET`, `PUBLIC_MOCK_API` env vars

### Security
- JWT access tokens with HMAC-SHA256 signing and algorithm-switching protection
- Bcrypt password hashing (cost factor 12)
- HTTPOnly cookies with Secure flag and SameSite Strict
- Refresh token rotation with replay-attack detection and family revocation
- Role-based access control middleware
- Production-only validation: rejects empty DB passwords, disabled SSL, default JWT secret

---