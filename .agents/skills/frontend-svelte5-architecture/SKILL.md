---
name: frontend-svelte5-architecture
description: Build and refactor OpenBench Svelte 5 frontend code. Use when creating or slicing Svelte components, implementing route guards, managing rune-based state services, wiring forms, formatted inputs, async views, client-side navigation, or frontend toolchain dependencies in apps/frontend.
---

# Frontend Svelte 5 Architecture

## Operating Rule

Follow existing route and service patterns before adding abstractions. Inspect the nearest `+page.svelte`, route `components/`, and matching service tests before editing.

## Workflow

1. Locate the relevant route under `apps/frontend/src/routes` and shared services under `apps/frontend/src/lib/services`.
2. Use Svelte 5 runes explicitly: `$state`, `$props`, `$derived`, `$effect`, and `$bindable` where appropriate.
3. Keep route-specific UI in a local `components/` folder. Move only genuinely reusable primitives to `src/lib/components`.
4. Preserve API response envelope handling and `credentials: 'include'` for cookie-protected requests.
5. Add or update focused Vitest tests for service logic and component behavior when behavior changes.
6. For dependency changes, inspect peer requirements and update `package.json` plus `package-lock.json` together.
7. Verify with `cd apps/frontend && npm run check`; run targeted Vitest tests when service or component behavior changes.

## Load References

- Read `references/frontend-patterns.md` when touching forms, formatted inputs, route guards, global services, async tabs/views, component slicing, or frontend package/toolchain dependencies.

## Hard Checks

- Do not introduce Svelte 4 `writable` stores for new app state.
- Do not use `bind:value` on display-masked inputs.
- Do not type create/update service payloads as full response models when generated request types exist.
- Do not allow older async requests to overwrite newer tab/view state.
- Do not assign nullable or missing API list data into array state; normalize to an array at the service or loader boundary.
- Do not call array methods such as `.filter()` or `.map()` on state that can be `undefined` or `null`.
- Do not add route components to `src/lib/components` unless another route will reuse them now.
- Do not upgrade TypeScript, Svelte, Vite, Vitest, or OpenAPI generator packages across major peer boundaries without checking compatibility.
