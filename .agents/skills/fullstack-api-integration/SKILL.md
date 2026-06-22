---
name: fullstack-api-integration
description: Keep OpenBench Go API contracts, Svelte services, mock API behavior, seed data, frontend payload types, and generated OpenAPI tooling aligned. Use when adding or changing endpoints, response envelopes, JSON tags, TypeScript interfaces, MOCK_API behavior, auth fetches, seeded credentials, OpenAPI generators, or generated API type commands.
---

# Fullstack API Integration & Mocking

## Operating Rule

Treat the backend response shape as the contract source and keep the frontend service, mock service, seed data, and tests in sync in the same change.

## Workflow

1. Start from the Go DTO or response struct and confirm its JSON tags.
2. Update the corresponding TypeScript interface and service parser in `apps/frontend/src/lib/services`.
3. Parse successful responses through `response.data`; parse errors through the response envelope message/error fields.
4. Preserve `credentials: 'include'` for auth-protected requests.
5. Mirror the endpoint in `src/lib/services/mocks` when mock mode must support the workflow.
6. Update seed data, mock tests, real service tests, and E2E assumptions together.

## Load References

- Read `references/api-contracts.md` before changing endpoints, mock behavior, seed credentials, API env variables, or payload naming.

## Hard Checks

- Do not add a frontend field name that lacks a matching Go JSON tag or mapping.
- Do not leave generated OpenAPI Go/TypeScript files stale after editing `docs/api/openapi.yml`.
- Do not send frontend create/update payloads that contain fields outside the OpenAPI request schema.
- Do not let successful list endpoints produce missing or nullable `data` when the frontend contract expects an array.
- Do not change OpenAPI generator or TypeScript versions without checking generator peer dependency compatibility.
- Do not return synchronous mock data for async workflows.
- Do not use one storage key for unrelated mock domains.
- Do not assume the skill text is the source of truth when current code and tests define an established contract; reconcile and update both if they disagree.
