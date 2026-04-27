# Authentication Implementation Plan

> Supabase Auth integration with local PostgreSQL user synchronization for the OpenBench API.
> **Note**: Database connection & environment infrastructure relies on [Local DB Setup](../database/setup-local-db/plan.md) and [Database Connection](../database/connection/plan.md).

---

## A. Logical Requirements

### Problem Statement
The OpenBench API needs to authenticate incoming requests using Supabase as its identity provider, while maintaining its own PostgreSQL database for business-model data. Supabase handles all identity concerns (sign-up, sign-in, password resets, OAuth). The API's responsibility is limited to **verifying Supabase-issued JWTs** and **synchronizing the authenticated user identity into a local `users` table** so that business entities can reference a local user record.

### Architectural Boundaries
- **Supabase owns**: User registration, login, password management, token issuance/refresh, OAuth providers.
- **OpenBench API owns**: JWT verification, local user profile storage, business-domain data, authorization rules.
- **The API never stores passwords** — it only stores a reference to the Supabase user UUID and supplementary profile fields.

### Edge Cases
1. **Expired JWT** — Must reject and return `401 Unauthorized` with a clear error message.
2. **Malformed or missing Authorization header** — Must reject with `401`.
3. **Valid JWT but user does not yet exist in local DB** — Must auto-provision a local user record on first authenticated request (Just-In-Time provisioning).
4. **JWKS endpoint unreachable at startup** — Application must fail fast with a fatal log.
5. **JWKS key rotation** — Cached key set must refresh periodically so rotated keys are picked up transparently.
6. **Concurrent first-requests for the same new user** — The user upsert must be idempotent (use `ON CONFLICT DO UPDATE` or equivalent) to avoid duplicate-key race conditions.
7. **Database unreachable during request** — Must return `503 Service Unavailable`, not panic.

---

## B. Structural Strategy

### B.1 — File System Impact

#### New Files to Create

| # | Path | Purpose |
|---|------|---------|
| 1 | `internal/domain/user.go` | Domain model struct for the local `User` entity (fields: `ID` as UUID primary key mirroring Supabase `sub`, `Email`, `FullName`, `AvatarURL`, `CreatedAt`, `UpdatedAt`). |
| 2 | `internal/repository/user_repo.go` | `UserRepository` interface and private implementation. Methods: `UpsertFromAuth` (idempotent create-or-update from JWT claims) and `FindByID`. |
| 3 | `internal/service/auth_service.go` | `AuthService` interface and private implementation. Orchestrates JWT parsing, claim extraction, and user upsert. Depends on the JWKS key set and `UserRepository`. |
| 4 | `internal/middleware/auth.go` | Fiber middleware that intercepts requests, delegates to `AuthService` for token verification, injects the authenticated user into the Fiber context locals, and short-circuits with `401` on failure. |
| 5 | `internal/handlers/auth.go` | HTTP handlers for auth-related endpoints: `GET /auth/me` (returns the current user profile from context). |
| 6 | `migrations/001_create_users_table.sql` | SQL migration file to create the `users` table in the local PostgreSQL database. |

#### Files to Modify

| # | Path | Change |
|---|------|--------|
| 1 | `pkg/config/config.go` | Add Supabase project URL and JWKS URL to the central configuration struct. |
| 2 | `.env.example` | Add Supabase required variables placeholders. |
| 3 | `cmd/api/main.go` | Add bootstrap steps: initialize JWKS key set, wire Auth repositories and services, register auth middleware on protected route groups, add auth routes. |
| 4 | `go.mod` | New dependencies will be added (see Dependencies section). |

*Note: Database definitions and general configuration loader implementations are defined in the associated DB plans.*

### B.2 — New Dependencies

| Library | Purpose |
|---------|---------|
| `github.com/lestrrat-go/jwx/v3` | JWKS fetching, caching, and JWT parsing/verification. Preferred over `golang-jwt` for its first-class JWKS support and automatic key-ID matching. |
| `github.com/google/uuid` | Standard UUID library for dealing with user IDs matching Supabase models. |

### B.3 — Module Architecture

```
Client
  │
  ▼
┌──────────────────────────────────────────────────────┐
│  Fiber App                                           │
│  ┌────────────────────┐   ┌───────────────────────┐  │
│  │ ZapLogger MW       │──▶│ AuthMiddleware        │  │
│  │ (existing)         │   │ (new)                 │  │
│  └────────────────────┘   └───────┬───────────────┘  │
│                                   │                  │
│                         ┌─────────▼──────────┐       │
│                         │ AuthService        │       │
│                         │ (JWT verify +      │       │
│                         │  user sync)        │       │
│                         └──┬──────────┬──────┘       │
│                            │          │              │
│              ┌─────────────▼┐   ┌─────▼───────────┐  │
│              │ JWKS KeySet  │   │ UserRepository  │  │
│              │ (lestrrat)   │   │ (sqlx + pgx)    │  │
│              └──────────────┘   └─────────────────┘  │
│                                        │             │
│                                  ┌─────▼──────┐      │
│                                  │ PostgreSQL │      │
│                                  └────────────┘      │
└──────────────────────────────────────────────────────┘
```

**Relationship descriptions:**
- The **AuthMiddleware** depends on `AuthService` (injected via constructor). It extracts the Bearer token from the request header and passes it to the service for verification.
- The **AuthService** depends on a cached JWKS `jwk.Set` (for signature verification) and `UserRepository` (for local user upsert). It parses and validates the JWT, then calls the repository to ensure the user exists locally.
- The **UserRepository** depends on `*sqlx.DB` (injected via constructor). It performs idempotent upsert and lookup operations on the `users` table.

### B.4 — Interface Specifications

#### UserRepository (interface)
- `UpsertFromAuth(ctx, supabaseUserID uuid, email string, fullName string, avatarURL string)` → returns the local User model or an error. Performs an `INSERT ... ON CONFLICT (id) DO UPDATE` to handle both first-time and returning users idempotently.
- `FindByID(ctx, id uuid)` → returns the User model or an error.

#### AuthService (interface)
- `VerifyAndSync(ctx, rawToken string)` → returns the authenticated User model or an error. Internally: parse JWT, validate signature via JWKS, validate claims (expiry, issuer), extract `sub`/`email`/`user_metadata`, call `UserRepository.UpsertFromAuth`, return User.

#### AuthMiddleware (function — not an interface)
- A function that accepts `AuthService` and returns a `fiber.Handler`. On each request it reads the `Authorization: Bearer <token>` header, calls `AuthService.VerifyAndSync`, and on success stores the User in Fiber's `c.Locals("user", user)`. On failure, responds with `401`.

---

## C. Step-by-Step Logic

### Phase 1 — Configuration Additions & Migrations

> **Note:** The underlying database connection pool, Make targets for migrations, and environment loading architectures are defined in [Local DB Setup](../database/setup-local-db/plan.md) and [Database Connection](../database/connection/plan.md). The steps below only detail the Auth-specific additions to that foundation.

#### Step 1: Environment & Config Updates
1. Add `SupabaseURL` setting to the `.env` file and `Config` struct (defined in `pkg/config/config.go`).
2. Derive `SupabaseJWKSURL` inside `LoadConfig()` by concatenating `SupabaseURL + "/auth/v1/.well-known/jwks.json"`.

#### Step 2: Database Migration
1. Create `migrations/001_create_users_table.sql` with the following schema intent:
   - Table `users` with columns: `id` (UUID, primary key — this mirrors the Supabase `auth.users.id`), `email` (text, unique, not null), `full_name` (text, nullable), `avatar_url` (text, nullable), `created_at` (timestamptz, default now), `updated_at` (timestamptz, default now).
   - An index on `email` for lookup performance.

### Phase 2 — JWT Verification

#### Step 3: JWKS Key Set Initialization
1. At application startup (in `main.go`), use `jwk.Fetch(ctx, supabaseJWKSURL)` from `lestrrat-go/jwx` to retrieve the key set.
2. Wrap the fetch in a context with a 10-second timeout. If it fails, log fatal and exit — the app cannot verify tokens without keys.
3. Consider using `jwk.NewCache()` with `jwk.WithRefreshInterval(1 hour)` for automatic background refresh of keys to handle key rotation seamlessly.
4. Pass the resulting `jwk.Set` into the `AuthService` constructor.

#### Step 4: AuthService Implementation
1. The private struct holds a reference to the JWKS `jwk.Set` and the `UserRepository` interface.
2. `VerifyAndSync` logic:
   - **Validation**: Call `jwt.Parse(rawToken, jwt.WithKeySet(keySet))` which automatically matches the `kid` header to the correct key and verifies the signature. If parsing fails, return an "invalid token" error.
   - **Claim Extraction**: Read the `sub` claim (Supabase user UUID), `email` claim, and optionally `user_metadata.full_name` and `user_metadata.avatar_url` from the parsed token.
   - **Issuer Check**: Validate that the `iss` claim matches the expected Supabase URL pattern. Reject if it does not.
   - **Expiry Check**: The `jwt.Parse` function in `lestrrat-go/jwx` automatically validates `exp` — no manual check needed.
   - **User Sync**: Call `UserRepository.UpsertFromAuth(ctx, sub, email, fullName, avatarURL)`. This ensures the local database always has a fresh copy of the user's basic profile.
   - **Resolution**: Return the resulting User model to the caller.

### Phase 3 — Middleware & Route Integration

#### Step 5: Auth Middleware
1. Define a function `RequireAuth(authService AuthService) fiber.Handler`.
2. Inside the returned handler:
   - Read the `Authorization` header from the request.
   - Check that it starts with `Bearer `. If not, respond with `401` and a JSON body `{"error": "missing or malformed authorization header"}`.
   - Extract the raw token string (everything after `Bearer `).
   - Call `authService.VerifyAndSync(ctx, rawToken)`.
   - If an error is returned, respond with `401` and a JSON body `{"error": "invalid or expired token"}`. Log the actual error at `warn` level for observability.
   - If successful, store the returned User in `c.Locals("user", user)` so downstream handlers can access it.
   - Call `c.Next()` to proceed to the route handler.

#### Step 6: Auth Handlers
1. Create a `GET /auth/me` handler that:
   - Retrieves the user from `c.Locals("user")`.
   - Type-asserts it to the domain User model.
   - Returns a `200` JSON response with the user's profile fields (id, email, full_name, avatar_url).

#### Step 7: Route Group Wiring (main.go updates)
1. Fetch/cache the JWKS key set.
2. Initialize `UserRepository` with the DB connection (from connection plan).
3. Initialize `AuthService` with the key set and user repository.
4. Create a Fiber route group `/api` (or `/api/v1`) and apply `RequireAuth(authService)` as group-level middleware.
5. Register `GET /auth/me` under the protected group.
6. Keep `GET /health` outside the protected group (public endpoint).

---

## D. Best Practice & Quality Guardrails

### Error Handling
- **JWKS fetch failure at startup** → `log.Fatal`. App cannot verify tokens.
- **JWKS refresh failure at runtime** → Log at `error` level but continue using the cached key set. Only fatal on initial fetch.
- **JWT parse/verify failure** → Return `401` to client. Log token error at `warn` level (do NOT log the raw token — security risk).
- **User upsert DB error** → Return `503` to client. Log full error at `error` level.
- **All errors** must be wrapped with context using `fmt.Errorf("operation: %w", err)` for traceability.

### Security
- **Never log raw JWT tokens** — they are bearer credentials.
- **Validate the `iss` (issuer) claim** against the known Supabase project URL to prevent token injection from other Supabase projects.
- **Use JWKS (asymmetric verification)** instead of the legacy JWT secret (symmetric). This is the Supabase-recommended approach.

### Performance
- **JWKS caching** — Fetch once at startup, auto-refresh every hour in the background. Never fetch per-request.
- **User upsert is a single round-trip** — The `INSERT ... ON CONFLICT` pattern combines existence check and insert/update into one SQL statement.
- **No external HTTP calls in the hot path** — JWT verification is entirely local (signature check against cached keys).

### Observability
- **JWKS fetch**: Log success/failure of initial and refresh fetches.
- **Auth middleware**: Log every authentication failure at `warn` level with the request method, path, and error reason (not the token).
- **User upsert**: Log first-time user provisioning at `info` level (`"new user synced from Supabase"` with the user UUID).

---

## E. Verification Plan

### Test Scenarios

#### Success Scenarios
| # | Scenario | Expected Result |
|---|----------|-----------------|
| S1 | Valid Supabase JWT sent to `GET /api/auth/me` | `200 OK` with user profile JSON. User exists in local `users` table. |
| S2 | First-time user with valid JWT hits a protected endpoint | `200 OK`. New row created in `users` table with matching UUID and email. |
| S3 | Returning user with valid JWT hits a protected endpoint | `200 OK`. Existing `users` row `updated_at` timestamp is refreshed. |
| S4 | `GET /health` without any Authorization header | `200 OK` — health check is unprotected. |

#### Failure Scenarios
| # | Scenario | Expected Result |
|---|----------|-----------------|
| F1 | No `Authorization` header on protected route | `401 Unauthorized` with `{"error": "missing or malformed authorization header"}` |
| F2 | `Authorization: Bearer <expired_token>` | `401 Unauthorized` with `{"error": "invalid or expired token"}` |
| F3 | `Authorization: Bearer <garbage_string>` | `401 Unauthorized` with `{"error": "invalid or expired token"}` |
| F4 | `Authorization: Basic <credentials>` (wrong scheme) | `401 Unauthorized` with `{"error": "missing or malformed authorization header"}` |
| F5 | Valid JWT from a **different** Supabase project (wrong issuer) | `401 Unauthorized` — issuer claim mismatch |

#### Edge Scenarios
| # | Scenario | Expected Result |
|---|----------|-----------------|
| E1 | Two concurrent requests for the same brand-new user | Both succeed with `200`. Only one row exists in `users` table (idempotent upsert). |
| E2 | Database goes down after startup | Protected endpoints return `503`. Health check should reflect unhealthy. |
| E3 | JWKS keys rotated by Supabase | After the cached key set auto-refreshes (within 1 hour), new tokens signed with the new key verify successfully. |
