# Frontend Patterns

## Files To Inspect First

- Route page or layout: `apps/frontend/src/routes/**/+page.svelte` or `+layout.svelte`
- Route-local components: `apps/frontend/src/routes/**/components/*.svelte`
- Shared primitives: `apps/frontend/src/lib/components`
- Services and tests: `apps/frontend/src/lib/services/*.ts` and `*.test.ts`
- Global services: `apps/frontend/src/lib/services/*.svelte.ts`

## Svelte 5 State

- Use `$state` for mutable local state.
- Use `$derived` for values computed from state; declare dependencies before the derived value.
- Use `$effect` for side effects only. If the effect registers listeners, timers, or subscriptions, return cleanup.
- Use `.svelte.ts` classes for app-level services. Export singleton instances and consume properties directly, without the Svelte 4 `$store` prefix.

```typescript
class ToastService {
  messages = $state<Toast[]>([]);
  show(message: string) {
    this.messages = [...this.messages, { id: crypto.randomUUID(), message }];
  }
}

export const toastService = new ToastService();
```

## Component Slicing

- Keep route-specific components beside the route in `components/`.
- Extract when the parent page is doing multiple jobs, when the child has a clear API, or when testability improves.
- Pass callbacks such as `onEdit` and `onDelete` instead of using `createEventDispatcher`.
- Use `$bindable()` only for intentional two-way state such as filters, pagination controls, or form drafts.
- Avoid slicing one-line markup or splitting components so far that data flow becomes harder to read.

## Forms

- Use normal `<input type="number">` with direct `bind:value` for numeric state.
- Use `$derived` for validation and disabled states.
- Respect shared component prop contracts in `src/lib/components`; pass required `id`, `label`, `error`, and value props consistently.

## Formatted Inputs

Use text inputs for visual masks such as IDR currency:

```svelte
let rawAmount = $state(0);
let displayAmount = $state(formatCurrency(0));

$effect(() => {
  displayAmount = formatCurrency(rawAmount);
});

function handleAmountInput(event: Event) {
  const input = event.currentTarget as HTMLInputElement;
  rawAmount = parseCurrency(input.value);
}
```

- Do not use `bind:value` on masked inputs; parse in `oninput`.
- Keep raw numeric state separate from display text.

## Async Views

When tabs, filters, or route-local controls trigger fetches, guard against stale responses:

```typescript
let fetchId = 0;

async function loadCurrentView() {
  const currentFetchId = ++fetchId;
  const data = await service.list(activeFilter);
  if (currentFetchId === fetchId) {
    items = data;
  }
}
```

## Route Guards And Auth

- Treat session validation errors as unauthenticated.
- Preserve `credentials: 'include'` for cookie-protected API calls.
- Keep local session state consistent with `authService.checkSession()` and sign-out cleanup.
