---
name: fullstack-api-integration
description: Use when implementing API contracts between Go and Svelte, mocking endpoints in the frontend for development, or aligning seed data and payload structures.
version: 1.0.0
---

# Fullstack API Integration & Mocking

## Overview
Patterns for synchronizing Svelte frontends with Go backends, defining consistent payload contracts, and building a robust client-side mock layer (`MOCK_API`) for parallel development.

## Step-by-Step Instructions

### Part A — Aligning Contracts
1. **Align Payload Contracts**: Ensure JSON struct tags in Go exactly match TypeScript interfaces in Svelte.
2. **Handle API Response Envelopes**: Svelte fetch consumers must always parse payloads through the standard backend envelope structure (`response.data`).
3. **Synchronize Seed & Mock Profiles**: Keep the backend seed credentials (e.g., `seeder.go`) exactly aligned with frontend mock credentials (`mockAuth.ts`).
4. **Match Cookie Flags**: When Go sets/clears auth cookies, ensure the frontend handles cross-origin fetch correctly (`credentials: 'include'`).

### Part B — Mocking Endpoints in Frontend
1. **Dynamic Mock Toggle**: Wrap the environment toggle in an `isMockEnabled()` check. Priority: `localStorage.getItem('MOCK_API')` -> `env.PUBLIC_MOCK_API` -> `import.meta.env.VITE_MOCK_API`.
2. **Structure the Mock Environment**:
   - Centralize all mock assets inside a dedicated `src/lib/services/mocks/` directory to prevent root-level clutter.
   - Split responsibilities into discrete files: `types.ts` (interfaces), `seed.ts` (static arrays), `db.ts` (data logic), `auth.ts` (auth logic), and `network.ts` (fetch interceptor).
   - Use distinct `localStorage` keys per entity (e.g., `openbench_mock_tickets`, `openbench_mock_inventory`) to avoid JSON serialization bottlenecks and accidental cross-domain wipes.
3. **Simulate Network Latency**: Every mock method MUST include an artificial delay (300–600ms) to trigger loading states and expose async bugs.
4. **Mirror Real Signatures**: Mock service methods must match real API signatures so transitioning requires swapping only the implementation body.
5. **Separate Auth from Data**: Keep auth simulation (`mockAuth.ts`) and `sessionStorage` separate from data CRUD (`mockDb.ts`) and `localStorage`.
6. **Simulate Business Logic & Errors**: Replicate critical business rules (e.g., stock depletion, status transitions) and include error scenarios to test UI feedback.

## Common Mistakes to Avoid
- Mismatched JSON tags causing frontend fields to be initialized as `undefined`.
- Exporting `MOCK_API=true` in npm but reading `PUBLIC_MOCK_API` in code. Variable names must align perfectly.
- Synchronous mock returns hiding loading states.
- Using placeholder seed data ("test", "foo") instead of realistic data.
- Shared storage keys between Auth and Data mocks, causing accidental wipes.
