---
name: frontend-svelte5-architecture
description: Use when creating/refactoring Svelte 5 components, building interactive forms, handling global state services, formatting inputs, slicing large pages, and handling client-side routing.
version: 1.0.0
---

# Frontend Svelte 5 Architecture

## Overview
Guidelines and patterns for building scalable, accessible, type-safe user interfaces with Svelte 5 Runes. This covers UI components, global state services, formatted inputs, and component slicing.

## When to Use
- Creating or refactoring Svelte 5 UI components.
- Building interactive forms and handling formatted inputs (e.g. currency masks).
- Slicing large Svelte pages (+page.svelte) into smaller components.
- Managing global state using Svelte 5 Runes (`$state` in `.svelte.ts`).
- Implementing client-side navigation route guards.

## Step-by-Step Instructions

### Part A — Components & Routing
1. **Implement State with Svelte 5 Runes**: Use explicit runes (`$state`, `$props`, `$derived`, `$effect`) for reactivity. Always return a cleanup function inside `$effect` if listeners are registered.
2. **Establish Route Protection**: Implement layout load functions to verify sessions and redirect unauthenticated users. Treat API errors during session validation as unauthenticated.
3. **Clean Unused Imports**: Run `npm run check` or `svelte-check` to catch unused state variables and dead imports.
4. **Cross-Origin Fetches**: Configure client API requests with `credentials: 'include'` for cookie-protected routes.
5. **Race Condition Prevention in UI**: When building UIs that switch views/tabs instantly but fetch data asynchronously, use a `fetchId` counter. Increment it on each request, and before committing the fetched data to `$state`, verify `if (currentFetchId === fetchId)` to prevent old, slow requests from overwriting newer ones.

### Part B — Forms, Bindings & Type Safety
1. **Declare `$derived` States in Dependency Order**: TypeScript enforces block-scoping. Declare dependencies before deriving from them.
2. **Standard HTML Inputs for Numbers**: Use `<input type="number">` with direct `bind:value` instead of custom wrappers that coerce to string.
3. **Validate Form State Reactively**: Use `$derived` for real-time validation feedback.
4. **Custom Component Props**: Respect custom component `$props()` contracts. Pass all required properties.

### Part C — Handling Formatted Inputs
1. **Visual Masks (e.g., Currency)**: Use `type="text"` for visual masks.
2. **Separate Raw and Display State**: 
   ```svelte
   let rawAmount = $state(0);
   let displayAmount = $state(formatCurrency(0));
   ```
3. **Sync with `$effect`**:
   ```svelte
   $effect(() => { displayAmount = formatCurrency(rawAmount); });
   ```
4. **Use `oninput` Handler**: Do not use `bind:value` on formatted inputs. Parse and format in `oninput`.

### Part D — Global State Services
1. **Use `.svelte.ts`**: Replace Svelte 4 stores with native class instances using the `$state` rune in `.svelte.ts` files.
2. **Export a Singleton**:
   ```typescript
   class ToastService {
     messages = $state<string[]>([]);
     show(msg: string) { this.messages.push(msg); }
   }
   export const toastService = new ToastService();
   ```
3. **Consume directly**: Import and use `toastService.messages` without the `$` prefix.

### Part E — Slicing Components
1. **Collocation**: Put route-specific components in a `components/` sub-directory alongside the route. Put highly reusable ones in `src/lib/components/`.
2. **Use `$bindable` for Two-Way State**: When extracting filters/inputs, mark the prop with `$bindable()`.
3. **Replace Event Dispatchers**: Pass callback functions (`onEdit`, `onDelete`) as props instead of using `createEventDispatcher`.
4. **Extract Complex Forms**: Move inline creation forms out of tables and into dedicated child routes (e.g., `/new/+page.svelte`).

## Common Mistakes to Avoid
- Out-of-order `$derived` chains.
- Number-string coercion when using custom input wrappers.
- Over-slicing tiny trivial HTML snippets.
- Prop drilling instead of using Context API.
- Using `writable` from `svelte/store` instead of `.svelte.ts` classes.
