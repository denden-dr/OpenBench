# User Local Database Design Plan

> Architectural blueprint for the local `users` entity schema, domain model, and repository layer. This design supports the Supabase Authentication integration by providing a local proxy for user records.

---

## A. Logical Requirements

### Problem Statement
The application relies on Supabase for identity and authentication but requires a local representation of users to associate with business domain entities. We need a local `users` table and a corresponding data access layer (Repository) that can:
1. Perform idempotent creation or updates (Just-in-Time provisioning) when a user authenticates via Supabase.
2. Retrieve user profiles efficiently for authorization and data association.
3. Strictly adhere to the project's Domain Service and Repository Rules (public interfaces, private implementations, constructor injection).
14. Provide a secure endpoint for users to retrieve their own profile information.

### Edge Cases
1. **Concurrent First-Logins**: Multiple concurrent requests with a newly minted Supabase token must not cause primary key constraint violations. Upserts must be atomic.
2. **Missing Users**: Fetching a non-existent user should return a clear, identifiable "not found" error, not a generic database failure.
3. **Data Synchronization**: Repeated logins must update local profile fields (like email or avatar changes) if they differ from the local copy.
4. **Unauthorized Access**: Ensure the `/me` endpoint only returns data if a valid Supabase session is present.

---

## B. Structural Strategy

### B.1 — File System Impact

| # | Path | Purpose |
|---|------|---------|
| 1 | `migrations/001_create_users_table.sql` | Up/Down SQL definitions for the `users` table. |
| 2 | `internal/domain/user.go` | Domain struct defining the `User` entity properties. |
| 3 | `internal/repository/user_repo.go` | The public `UserRepository` interface, private struct implementation, and constructor. |
| 4 | `internal/service/user_service.go` | The public `UserService` interface, private struct implementation, and constructor. |
| 5 | `internal/handlers/user_handler.go` | HTTP handler for user-related endpoints. |

### B.2 — Module Architecture

```
┌─────────────────────────┐
│ UserHandler             │
└────────────┬────────────┘
             │ (Injects Interface)
┌────────────▼────────────┐
│ UserService             │
│ (Public Interface)      │
└────────────┬────────────┘
             │ (Implemented By)
┌────────────▼────────────┐
│ userService             │
│ (Private Struct)        │
└────────────┬────────────┘
             │ (Injects Interface)
┌────────────▼────────────┐
│ UserRepository          │
│ (Public Interface)      │
└────────────┬────────────┘
             │ (Implemented By)
┌────────────▼────────────┐
│ userRepository          │
│ (Private Struct)        │
└────────────┬────────────┘
             │
┌────────────▼────────────┐
│ PostgreSQL (Local DB)   │
└─────────────────────────┘
```

### B.3 — Interface Specifications

#### Domain Model (`internal/domain/user.go`)
- `User`: Struct containing `ID` (UUID), `Email` (String), `FullName` (String pointer/nullable), `AvatarURL` (String pointer/nullable), and `UpdatedAt` (Timestamp).

#### Repository Layer (`internal/repository/user_repo.go`)
- **Public Interface**: `UserRepository`
  - `UpsertFromAuth(ctx, id, email, fullName, avatarURL) -> (User, error)`: Handles atomic insert-or-update based on the Supabase ID.
  - `FindByID(ctx, id) -> (User, error)`: Retrieves a single user by their ID.
- **Constructor**: `NewUserRepository(db) -> UserRepository`: Returns the interface.
- **Private Implementation**: `userRepository` holding the database connection dependency.

#### Service Layer (`internal/service/user_service.go`)
- **Public Interface**: `UserService`
  - `GetProfile(ctx, id) -> (User, error)`: Orchestrates user retrieval by ID.
- **Constructor**: `NewUserService(repo) -> UserService`: Injects the `UserRepository` interface.

#### Handler Layer (`internal/handlers/user_handler.go`)
- **Public Interface**: `UserHandler`
  - `GetMe(c) -> error`: Fiber handler for `GET /v1/users/me`.
- **Constructor**: `NewUserHandler(service) -> UserHandler`: Injects the `UserService` interface.

---

## C. Step-by-Step Logic

### Phase 1 — Database Schema Generation
1. **Migration Creation**: Write the SQL script to create the `users` table.
   - Set `id` as `UUID PRIMARY KEY`.
   - Set `email` as `VARCHAR UNIQUE NOT NULL`.
   - Set `full_name` and `avatar_url` as `TEXT` allowing `NULL`.
   - Add an `updated_at` timestamp defaulting to the current time to track the last sync from Supabase.
2. **Indexing**: Add an index on the `email` column if frequent email lookups are anticipated, though `id` will be the primary lookup vector.

### Phase 2 — Domain Modeling
1. **Entity Definition**: In `internal/domain/user.go`, create the `User` struct.
2. **Tagging**: Add appropriate database serialization tags (e.g., `db:"id"`) to map smoothly to SQL query results.

### Phase 3 — Repository Implementation
1. **Interface Definition**: Define the `UserRepository` interface with the required methods in `internal/repository/user_repo.go`.
2. **Private Struct**: Define `type userRepository struct { db *sqlx.DB }` (or standard `*sql.DB` depending on the database driver choice).
3. **Constructor Injection**: Create `func NewUserRepository(db *sqlx.DB) UserRepository` returning the private struct initialized with the DB connection.
4. **Implement `FindByID`**:
   - Execute a `SELECT * FROM users WHERE id = $1` query.
   - Handle "no rows" gracefully, returning a domain-specific "not found" error.
5. **Implement `UpsertFromAuth`**:
   - Construct an `INSERT INTO users (id, email, full_name, avatar_url) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET email = EXCLUDED.email, full_name = EXCLUDED.full_name, avatar_url = EXCLUDED.avatar_url, updated_at = NOW() RETURNING *` query.
   - This ensures atomic synchronization of user data without race conditions.
   - Return the resulting `User` struct.

### Phase 4 — Service Implementation
1. **Interface & Struct**: Define `UserService` and `userService` in `internal/service/user_service.go`.
2. **Constructor**: Implement `NewUserService` taking `UserRepository`.
3. **GetProfile**: Implement the method to call `repo.FindByID`. Wrap errors contextually.

### Phase 5 — Handler & Routing
1. **Interface & Struct**: Define `UserHandler` and `userHandler` in `internal/handlers/user_handler.go`.
2. **Constructor**: Implement `NewUserHandler` taking `UserService`.
3. **GetMe Implementation**:
   - Extract the user ID from the authentication context (set by middleware).
   - Call `service.GetProfile`.
   - Return a JSON response or appropriate error status (404 for missing, 500 for failures).
4. **Wiring**: Register the handler in the main router using the `AuthMiddleware` to protect the route.

### Phase 6 — Unit Testing (Service Layer)
1. **Mock Generation**: Create a mock implementation of `UserRepository`.
2. **Table-Driven Tests**:
   - Implement tests using the table-driven approach to minimize boilerplate.
   - Define a test case struct including: `name`, `userID`, `mockRepoSetup`, `expectedUser`, and `expectedError`.
   - Test Cases:
     - `Success`: Verify the service returns a user when the repository find succeeds.
     - `User Not Found`: Verify the service returns the correct error when the user is missing.
     - `Repository Failure`: Verify the service wraps repository errors correctly.
3. **Assertions**: Use `testify/assert` and `testify/mock` for expectations.

### Phase 7 — Unit Testing (Repository Layer)
1. **Mocking Database**: Use `DATA-DOG/go-sqlmock` to simulate database interactions without a live PostgreSQL instance.
2. **Table-Driven Tests**:
   - `TestUserRepository_FindByID`: Verify correct scanning of rows and handling of `sql.ErrNoRows`.
   - `TestUserRepository_UpsertFromAuth`: Verify the `INSERT ... ON CONFLICT` query and the `RETURNING` clause mapping.
3. **Assertions**: Use `testify/assert` to verify the returned domain models match the expected database rows.

---

## D. Best Practice & Quality Guardrails

### Architectural Compliance
- **Strict Adherence to User Rules**: The repository layer must strictly use the interface/private-struct/constructor-injection pattern. No public structs or concrete return types from constructors are permitted.

### Error Handling
- **Wrapping**: All database errors must be contextually wrapped (e.g., `fmt.Errorf("upserting user: %w", err)`).
- **Isolation**: Do not expose raw database drivers or SQL errors directly to the service layer. Transform errors like `sql.ErrNoRows` into a domain-recognized error type if necessary.

### Concurrency & Performance
- **Idempotency**: Using `ON CONFLICT` guarantees safety against race conditions when two simultaneous requests are made by a user logging in for the very first time.
- **Connection Management**: Ensure the repository uses the context `ctx` parameter in all database calls to respect request timeouts and cancellations.

---

## E. Verification Plan

### Test Scenarios

#### Scenario 1: First-time User Upsert
- **Action**: Call `UpsertFromAuth` with a new UUID and email.
- **Verification**: Ensure the user is inserted into the database and the returned `User` struct contains the correct fields and generated timestamp.

#### Scenario 2: Existing User Sync
- **Action**: Call `UpsertFromAuth` with an existing UUID but a changed `email` or `fullName`.
- **Verification**: Ensure the database record is updated without throwing duplicate key errors, and the `UpdatedAt` timestamp advances.

#### Scenario 3: Fetching Users
- **Action**: Call `FindByID` with a known UUID.
- **Verification**: Assert the correct user profile is returned.
- **Action**: Call `FindByID` with a non-existent UUID.
- **Verification**: Assert a clean "not found" error is returned, completely distinct from connection or syntax errors.

#### Scenario 4: Profile Endpoint
- **Action**: Call `GET /v1/users/me` with a valid `user_id` set in context.
- **Verification**: Ensure the response is `200 OK` with the correct JSON body.
- **Action**: Call `GET /v1/users/me` without a `user_id` in context.
- **Verification**: Ensure the response is `401 Unauthorized`.

#### Scenario 5: Service Unit Tests
- **Action**: Run `go test ./internal/service/...`.
- **Verification**: Ensure all unit tests for `UserService` pass, demonstrating correct business logic and error handling without requiring a live database.

#### Scenario 6: Repository Unit Tests
- **Action**: Run `go test ./internal/repository/...`.
- **Verification**: Ensure all unit tests for `UserRepository` pass, demonstrating correct SQL execution and row scanning logic using `sqlmock`.
