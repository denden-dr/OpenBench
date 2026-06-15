---
name: backend-go-architecture
description: Use when building Go domain logic, handlers, repository layers, managing database transactions, securing JWT/cookie sessions, and building secure public trackers.
version: 1.0.0
---

# Backend Go Architecture

## Overview
Go backend packages follow a layered architecture (Domain -> Repository -> Service -> Handler -> Server). This skill covers server configuration, handling transactions, securing authentication sessions (RTR), and managing external identifiers.

## Layered Architecture
1. **Server Initialization**: Set up Fiber with read/write timeouts, `recover.New()`, strict CORS origins, and graceful shutdown listening for OS signals.
2. **Domain Models**: Pure structs without ORM tags.
3. **Repository Interface**: Data access signatures and SQL execution.
4. **Service Layer**: Business logic, orchestrating transactions, validations.
5. **Handlers**: Struct methods focusing on body parsing, validation via `validator`, and formatting responses via standard `response.JSON` / `response.Error`. DTOs belong here.

## Managing Transactions and Services
1. **Transaction Boundaries**: Handlers must never call DB commits/rollbacks. Service layer begins transaction, defers `tx.Rollback()`, and returns `tx.Commit()`.
2. **Refactoring Bloated Services**: Split large service methods into private helpers passing `tx *sqlx.Tx`. Do not leak `sqlx.Tx` outside the service.
3. **Database Locking**: Use `SELECT ... FOR UPDATE` (or pessimistic/optimistic locks) in the repository layer to prevent race conditions. Ensure search columns are indexed.

## Securing Auth Sessions (RTR & Cookies)
1. **Cookie Security**: Store Refresh Tokens in `HttpOnly` secure cookies.
   - `HttpOnly: true`, `Secure: !isDev`, `SameSite: "Lax"`, `Path: "/"`.
2. **Clear Cookies Cleanly**: Must match all original configuration parameters (`Path`, `HTTPOnly`, `Secure`, `SameSite`) to successfully delete.
3. **Refresh Token Rotation**:
   - Track Token Families (`family_id`).
   - Replay Attack Detection: If a used token is presented, revoke the entire family.
   - Grace Period: Allow < 5 seconds of latency to prevent false breach alarms on concurrent requests.
4. **Deliver Access Tokens**: Short-lived (5-15 min) delivered in JSON payload; refresh token long-lived in cookie.

## Securing Public Trackers
1. **Separate Identifiers**: Use a UUID v4 (Public ID) for unauthenticated lookup interfaces. Never accept internal sequential numbers (e.g. `OB-YYYYMM-XXXX`).
2. **Validate Input**: Check UUID format before querying the DB.
3. **Limit Exposed Data**: Do not expose PII, internal notes, or sequential internal IDs in the public tracker payload. Limit client-side history (`localStorage`) to essential data.

## Common Mistakes
- Leaking transaction logic to handlers.
- Using `SELECT` without `FOR UPDATE` in transactions.
- Permissive CORS (`AllowOrigins: "*"`) in production.
- Forgetting `HTTPOnly: true` or missing `Path: "/"` when clearing cookies.
- No grace period for token refresh leading to unexpected logouts.
