# Multi-Dimensional Review Checklist & Subagents

This document contains detailed verification checklists and subagent prompts to support code reviews.

---

## 1. Architecture & Boundaries
*   **Monorepo Separation**: Frontend (`apps/frontend`) and backend (`apps/backend`) should communicate exclusively via defined API contracts. No shared code imports between the two apps.
*   **Domain Isolation**: Go code under `internal/` must be organized by domain packages (e.g. `auth`, `database`).
*   **Dependency Cycle Check**: Ensure package interactions do not introduce circular dependencies. Utilize tests under separate `*_test` packages to break import cycles.

## 2. Security
*   **No Hardcoded Secrets**: Real keys, passphrases, and `.env` files must never be committed. Ensure new configs default to placeholders.
*   **Auth & Session Tokens**: Verify HttpOnly, Secure, and SameSite (Lax/Strict) cookie flags are present on JWT tokens.
*   **Refresh Token Rotation (RTR)**: Verify that old refresh tokens are invalidated upon reuse and family revocation triggers on breaches.
*   **Input Validation**: Ensure all handler requests utilize validator rules (e.g. `validator.ValidateStruct`).

## 3. Testing
*   **Hermetic & Isolated**: Test containers must choose random host ports to prevent port collision.
*   **Resource Leak Prevention**: Integration tests using Testcontainers must declare `TestMain` and invoke `tdb.Terminate()` to cleanly shut down containers.
*   **Mock Expectation Verification**: Verify all SQL or service mock expectations are verified. Go mock tests must assert expectations (e.g. `mockSQL.ExpectationsWereMet()` or `mock.AssertExpectations(t)`) via `t.Cleanup()`.
*   **Hydration Race Mitigation**: Playwright tests must wait for `main[data-hydrated="true"]` before page interactions.
*   **Network Isolation**: Playwright tests must abort external assets (`fonts.googleapis.com`) to prevent timeout delays.

## 4. Best-Practices
*   **Error Handling**: Go errors must be wrapped with `%w` and logged using structured fields (`slog`).
*   **State Management**: Svelte frontend components must utilize Svelte 5 runes syntax (`$state`, `$derived`, `$effect`) instead of older Svelte 4 store patterns.
*   **Dead Code Elimination**: Ensure there are no unused imports or unused reactive state variables (such as Svelte runes). Verify by running `npm run check` on the frontend before finalizing.
*   **Conventional Commits**: Commits must match conventional scopes (`feat(backend)`, `fix(frontend)`).

## 5. Performance
*   **Database Connections**: DB pool sizes (MaxOpenConns, MaxIdleConns) must be configured with timeouts and idle durations.
*   **Vite Dev Compilation**: Lucide icons and other massive dynamic ES modules must be pre-bundled via `vite.config.ts` `optimizeDeps` to speed up dev compilation times.
*   **E2E Worker Concurrency**: Playwright E2E tests targeting local dev instances should limit concurrent workers (`workers: 1`) to prevent server resource starvation.

---

## Subagent Specializations & Prompts

For large changesets or full-codebase audits, run parallel subagents with these specific prompt templates:

1. **Architecture & Boundaries Specialist**
   *   *Focus*: Monorepo boundaries, package cycles, domain isolation.
   *   *Prompt*: `Audit the architectural layout of the following diff/files: [file paths]. Check for correct domain isolation under 'apps/backend/internal/', verify SvelteKit/Go boundary separation, and ensure no circular package imports are created.`
2. **Security & Cryptography Auditor**
   *   *Focus*: Input injection, credential leaks, token rotation safety, secure cookie headers.
   *   *Prompt*: `Audit the security stance of the following diff/files: [file paths]. Look for hardcoded credentials, unvalidated struct inputs, unsafe JWT signature validation, missing authorization checks in Go handlers, and incorrect cookie attributes (HttpOnly, Secure, SameSite).`
3. **Testing & Automation Auditor**
   *   *Focus*: Container cleanup, hydration races, test isolation.
   *   *Prompt*: `Verify the reliability of the test scripts in the following diff/files: [file paths]. Ensure integration suites use 'testutil.SetupTestDB()' with 'TestMain' and 'tdb.Terminate()', check for dynamic port bindings, and confirm Playwright E2E files handle hydration safely and block Google Fonts.`
4. **Performance & Optimization Specialist**
   *   *Focus*: Connection pooling, goroutine lifecycle, Vite pre-bundling.
   *   *Prompt*: `Audit the performance patterns in the following diff/files: [file paths]. Check for context cancellation (preventing goroutine leaks), database connection pooling configuration, N+1 queries, and check if third-party modules are pre-bundled under 'vite.config.ts' to accelerate compilation.`
