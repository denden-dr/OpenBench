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
