# Testing Svelte Applications (Vitest & E2E)

### 1. Mocking Runtime Modules
SvelteKit build modules (like `$app/navigation` or `$app/environment`) must be mocked in the test setup file:
```typescript
// tests/setup.ts
import { vi } from 'vitest';
vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
}));
```

### 2. Testing Library Render
Render components using `render(Component, { props })`.

### 3. Awaiting Reactivity
User interactions (such as clicks) that trigger UI changes must be resolved by awaiting helper actions from Svelte Testing Library (`tick()` or `userEvent`).

### 4. A11y & Role Selection
Select inputs and buttons using roles and labels (e.g. `getByRole('button', { name: /sign in/i })`).
