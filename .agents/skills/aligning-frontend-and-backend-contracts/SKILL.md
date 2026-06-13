---
name: aligning-frontend-and-backend-contracts
description: Use when implementing API contracts between Go backend and Svelte frontend, verifying database row locking for state transitions, and aligning seeder/mock user credentials. Do not use for non-Svelte frontends or non-Go backend services.
version: 1.0.0
---

# Aligning Frontend and Backend Contracts

## Overview
This skill provides patterns and constraints for synchronizing Svelte/TypeScript frontends with Go backend services. It ensures consistent payload contracts, prevents database concurrency race conditions on sensitive state transitions, and aligns test/mock credentials.

## When to Use
Use when:
- Adding or modifying API endpoints and JSON response payloads.
- Implementing state transition checks (e.g., token rotation, payment status, resource checks).
- Updating mock data services or database seeders.

## Step-by-Step Instructions

1. **Verify Database Locking**: Read `references/db-locking.md` to ensure pessimistic or optimistic locks are implemented for sensitive state transitions.
2. **Align Payload Contracts**: Read `references/payload-mapping.md` to map structures between Svelte/TypeScript frontends and Go backends, confirming JSON struct tag naming conventions and response envelope formatting.
3. **Synchronize Seed & Mock Profiles**: Keep the local development database seeder credentials (e.g. `seeder.go`) and frontend mock credentials (e.g. `mockAuth.ts`) identical.
4. **Match Cookie Configuration Flags**: When setting or clearing authentication cookies in Go, ensure the flags (`Secure`, `HttpOnly`, `SameSite`, and `Path`) are exactly matched. Leaving flags mismatched when clearing cookies prevents browsers from discarding them.
5. **Handle API Response Envelopes**: Ensure Svelte TypeScript fetch consumers always parse payloads through the standard backend envelope structure (mapping response `.data` attributes), preventing layout guards from failing on `undefined` variables.

## Common Mistakes
- **Standard SELECT in Transactions**: Querying a state record using a plain `SELECT` inside a transaction, allowing concurrent sessions to read the same stale state and duplicate side-effects.
- **Mismatched JSON tags**: Returning fields in Go without the correct `json` struct tag or with mismatching case, causing the frontend fields to be initialized as `undefined` or `""`.
- **Diverging Seed Credentials**: Updating database seeder credentials but forgetting to update the frontend mock equivalents, confusing local developer verification.
