# Naming Conventions & Best Practices

To ensure codebase consistency, maintainability, and alignment with idiomatic Go standards across the OpenBench project, all files, directories, and objects must strictly adhere to the following naming conventions.

## 1. General Guidelines
- **Clarity over Brevity**: Names should be descriptive enough to be understood without requiring comments. Avoid cryptic abbreviations.
- **Idiomatic Go Acronyms**: Acronyms must be consistently cased as a single unit (e.g., `ID`, `HTTP`, `URL`, `JSON`, `API`, `UUID`).
  - ✅ **Correct**: `userID`, `HTTPClient`, `DatabaseURL`
  - ❌ **Incorrect**: `userId`, `HttpClient`, `DatabaseUrl`

## 2. File & Directory Naming
- **Directories (Packages)**: Must be short, concise, and **all lowercase**. Avoid underscores, dashes, or camelCaps. Prefer single words.
  - ✅ **Correct**: `internal`, `handlers`, `database`, `config`
  - ❌ **Incorrect**: `Internal/`, `http_handlers`, `databaseManager`
- **Go Source Files**: Must use **`snake_case.go`**. The filename should concisely map to the primary types or functionality it contains.
  - ✅ **Correct**: `postgres_db.go`, `health.go`
  - ❌ **Incorrect**: `postgresDB.go`, `HealthHandler.go`
- **Test Files**: Must append `_test.go` to the file being tested (e.g., `health_test.go`).
- **Service & Repository Files**: Strictly follow the established domain rules:
  - Services must end in `_service.go` (e.g., `user_service.go`).
  - Repositories must end in `_repo.go` (e.g., `user_repo.go`).

## 3. Structural Naming (Structs & Interfaces)
- **Exported Types (Public)**: Use **`PascalCase`** (UpperCamelCase).
  - ✅ **Correct**: `type HealthHandler struct`, `type Config struct`
- **Unexported Types (Private)**: Use **`camelCase`**. This is especially required for hiding concrete implementations behind interfaces.
  - ✅ **Correct**: `type userRepository struct`
- **Interfaces**:
  - One-method interfaces should generally end in `-er` (e.g., `Reader`, `Writer`, `Validator`).
  - Domain contract interfaces should define their role without a suffix (e.g., `UserRepository`, `AuthService`).

## 4. Variable & Function Naming
- **Exported Functions/Variables**: Use **`PascalCase`**.
  - ✅ **Correct**: `func LoadConfig()`, `var ErrNotFound`
- **Unexported Functions/Variables**: Use **`camelCase`**.
  - ✅ **Correct**: `func validateConfig()`, `db *sqlx.DB`
- **Constructors**: Must start with `New`.
  - If returning the primary package type, use `New()` (e.g., `logger.New()`).
  - If returning a specific type, use `New<Type>()` (e.g., `NewPostgresDB()`, `NewHealthHandler()`).
  - **Constraint**: Domain service/repository constructors *must* return their public `interface` type, not the private concrete struct.

## 5. Constants & Enums
- **Idiomatic Go Constants**: Use **`PascalCase`** for exported constants and **`camelCase`** for private ones. Avoid C-style uppercase formatting unless mapping directly to system environment variables.
  - ✅ **Correct**: `const DefaultTimeoutSeconds = 30`
  - ❌ **Incorrect**: `const DEFAULT_TIMEOUT_SECONDS = 30`
- **Environment Variables**: When referring to literal `.env` keys, use standard uppercase snake case format.
  - ✅ **Correct**: `DB_MAX_OPEN_CONNS`

## 6. Package Naming Pitfalls
- Do not repeat the package name in the names of its members if it causes stutterting.
  - ✅ **Correct**: `config.Config` (acceptable, though `config.Settings` or just `config.Load()` is often better).
  - ❌ **Incorrect**: `database.DatabaseConnection`, `handlers.HealthHandlerFn` (Keep it clean: `handlers.HealthCheck`).
