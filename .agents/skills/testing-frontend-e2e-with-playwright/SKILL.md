---
name: testing-frontend-e2e-with-playwright
description: Use when SvelteKit browser tests are flaky, timing out on external resources (fonts/APIs), or failing due to client-side hydration race conditions.
version: 1.1.0
---

# Frontend E2E Testing with Playwright

## Overview
End-to-End (E2E) tests validate user journeys across multiple pages by controlling browser sessions. We use Playwright to test SvelteKit client components against a mock or live API server.

## When to Use
- Simulating logins, form entries, page navigations, and redirects.
- Verifying UI feedback, toasts, error messages, and reactive states.
- Running parallel or sequential multi-browser validation.

### When NOT to Use
- Testing Go backend logic or SQL queries directly (use backend testify instead).

## Core Pattern: Safe Page Hydration & Network Interception

### Before (Anti-Pattern: Dynamic Races & External Bloat)
```typescript
test('should login', async ({ page }) => {
	// ❌ BAD: page.goto waits for google fonts / maps which timeout in CI
	await page.goto('/auth/signin');
	// ❌ BAD: clicks immediately, causing native form submission if Svelte is not hydrated
	await page.click('button[type="submit"]'); 
});
```

### After (Best Practice: Hydration Checking & Font Aborts)
```typescript
test.beforeEach(async ({ page }) => {
	// 1. Abort external Google Fonts to prevent network timeouts in isolated CI
	await page.route('**/*', route => {
		const url = route.request().url();
		if (url.includes('fonts.googleapis.com') || url.includes('fonts.gstatic.com')) {
			route.abort();
		} else {
			route.continue();
		}
	});

	await page.goto('/auth/signin');
	// 2. Wait for Svelte hydration to complete before interacting
	await page.waitForSelector('main[data-hydrated="true"]');
});

test('should login', async ({ page }) => {
	await page.fill('#email', 'admin@openbench.dev');
	await page.fill('#password', 'SecureAdminPassword123!');
	await page.click('button[type="submit"]');
	await expect(page).toHaveURL('/admin');
});
```

## Quick Reference: Operations

| Task | Playwright Call | Description |
|---|---|---|
| **Form Input** | `await page.fill('#id', 'value')` | Safely inputs values |
| **Action** | `await page.click('button')` | Clicks elements |
| **URL Assertion** | `await expect(page).toHaveURL('/path')` | Asserts exact path or pattern |
| **Visibility** | `await expect(locator).toBeVisible()` | Checks element visibility |
| **Hydration Hook** | `onMount(() => hydrated = true)` | Svelte hook bound to root markup |

## Common Mistakes
- **Vite Dev Server Resource Starvation**: Running Vite dev servers with high parallel workers (`fullyParallel: true`) can cause compilation timeouts. Set `workers: 1` and `fullyParallel: false` for dev-run E2E.
- **Hydration Race Condition**: Triggers raw form submission (e.g. appends `?email=...` to browser URL) because Svelte's client JS hasn't attached event handlers yet.
- **Mismatched Mock Environment Variables**: Exporting `MOCK_API=true` in the npm script but reading `PUBLIC_MOCK_API` in SvelteKit code. SvelteKit client-side code requires the `PUBLIC_` prefix to access env vars via `$env/static/public`. Always verify that the variable name exported in the script matches exactly what the code imports.
- **No Service Readiness Wait**: Running Playwright immediately after `docker-compose up` without verifying that backend and frontend are accepting requests, causing `ECONNREFUSED` or early test failures.

---

## Environment Variable Alignment for Mock Mode

Every environment variable that controls mock mode MUST use the prefix appropriate to the framework:
- **SvelteKit**: `PUBLIC_` prefix for client-side variables accessed via `$env/static/public`
- **Vite (non-SvelteKit)**: `VITE_` prefix for variables accessed via `import.meta.env`

### Checklist before writing/modifying mock scripts:
1. Identify the variable **read** by the frontend code (grep `$env/static/public` or `import.meta.env`).
2. Ensure the npm script exports **the exact same name**.
3. Add an assertion at the start of the test suite to verify mock mode is active:
```typescript
test.beforeAll(() => {
  const isMock = process.env.PUBLIC_MOCK_API === 'true';
  if (!isMock) throw new Error('Tests must run in mock mode. Check env var name alignment.');
});
```

---

## Service Readiness Before Test Execution

Makefile or CI scripts that run Playwright MUST wait for all services (backend API, frontend, database) to be ready before starting the test runner.

### Pattern: Wait-for-service in Makefile
```makefile
.PHONY: wait-for-services
wait-for-services:
	@echo "Waiting for backend..."
	@until curl -sf http://localhost:3000/api/health > /dev/null 2>&1; do sleep 1; done
	@echo "Waiting for frontend..."
	@until curl -sf http://localhost:5173 > /dev/null 2>&1; do sleep 1; done
	@echo "All services ready."

test-frontend-e2e: wait-for-services
	cd apps/frontend && npx playwright test
```

### Alternative: Playwright webServer config
```typescript
// playwright.config.ts
export default defineConfig({
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: true,
    timeout: 30_000,  // fail fast if service doesn't come up
  },
});
```
