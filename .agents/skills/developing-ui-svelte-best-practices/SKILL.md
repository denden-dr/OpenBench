---
name: developing-ui-svelte-best-practices
description: Use when creating or modifying Svelte 5 UI components, route guards, or Vitest component tests. Do not use for React, Vue, or Angular.
version: 1.1.0
---

# Svelte UI Development Best Practices

## Overview
Guidelines and patterns for building accessible user interfaces with Svelte (Svelte 5 Runes mode), SvelteKit routing, Neubrutalism aesthetics (Tailwind CSS v4 CSS overrides), and robust Vitest/E2E testing setups.

## When to Use
- Creating or refactoring Svelte 5 UI components (e.g., Buttons, Inputs, Cards, Layouts).
- Implementing component states, native form bindings, or client-side navigation route guards.
- Overriding Tailwind CSS v4 themes and typography styles for Neubrutalism.
- Writing Vitest component tests or Playwright E2E browser tests for Svelte applications.

## Step-by-Step Instructions

1. **Implement State with Svelte 5 Runes**: Use explicit runes (`$state`, `$props`, `$derived`, `$effect`) for component-specific reactivity. Always return a cleanup function inside `$effect` if listeners or intervals are registered.
2. **Apply Neubrutalism styling**: Read `assets/layout.css.template` and apply theme overrides within the CSS entry point using the `@theme` directive.
3. **Establish Route Protection**: Read `assets/layout-load.ts.template` and implement layout load functions or client navigation hooks to verify sessions and redirect unauthenticated users.
4. **Implement Component Testing**: Read `references/svelte-testing-rules.md` and write tests using Vitest, rendering components via testing library helpers and mocking SvelteKit runtime modules.
5. **Clean Unused Imports and Dead Code**: Run compiler checks (such as `npm run check` or `svelte-check`) to catch unused state variables (declared via `$state()`) and dead imports. Safely remove them before finalizing changes.
6. **Maintain Environment Templates**: Ensure `.env` and `.env.example` templates exist in the frontend root to document all required environment properties (e.g., `PUBLIC_MOCK_API=true/false`).
7. **Verify Response Properties**: Map reactive state fields (such as login email or session roles) strictly from the fields returned in the API response data, rather than local function parameter values, to keep client state synced.
8. **Handle Cross-Origin and Production Fetches**: Configure client-side API requests to send cookies by setting `credentials: 'include'` when fetching from cookie-protected routes. Avoid hardcoding relative URLs unless a production reverse proxy routes them. Avoid overriding Vite env loader configurations (e.g. `define` in `vite.config.ts` overriding `VITE_` variables) so that standard `.env` variables can be loaded.
9. **Design Secure Session Checks**: Ensure route guards treat API errors (like server 500s or network failures) during session validation (e.g., `/auth/me`) as unauthenticated, refusing to fall back to cached sessionStorage data that might grant stale UI access.

## Common Mistakes to Avoid
- **Unused State Runes / Dead Imports**: Leaving state variables (declared via `$state()`) or imported functions unused in the component template, which triggers linter warnings. Always run `npm run check` to verify.
- **Tailwind JS Configuration**: Creating `tailwind.config.js` files. Tailwind v4 uses CSS `@theme` properties exclusively.
- **Hanging Reactivity / State Pollution**: Modifying states inside `$effect` without returning cleanups, or failing to clear `sessionStorage` in test suites, causing state leakage between tests.
- **Not Awaiting tick()**: Checking DOM state assertions immediately after state mutations without awaiting updates.
- **Direct DOM Querying**: Altering elements directly via the `document` namespace instead of binding to Svelte reactive variables.
