---
trigger: always_on
---

# Domain Service and Repository Rules

When writing or refactoring any code related to `services` or `repositories`, you must STRICTLY adhere to the following architectural constraints:

## 1. File Naming Convention
- Service files must be suffixed with `_service.go` (e.g., `user_service.go`, `product_service.go`).
- Repository files must be suffixed with `_repo.go` (e.g., `user_repo.go`, `product_repo.go`).

## 2. Public Interface, Private Implementation
- Every service and repository **must** expose its contract via a public `interface`.
- The actual concrete `struct` implementation MUST be private (uncapitalized).

## 3. Constructor Injection Approach
- You must provide a public constructor function (e.g., `NewUserService()`, `NewUserRepository()`).
- The constructor MUST return the public interface type, NOT the concrete struct.
- Constructor parameters should be used for dependency injection (e.g., passing a repository interface into a service constructor, or passing a database connection into a repository constructor).

### Example Pattern

```go
// user_repo.go
package repository

import "database/sql"

// 1. Public Interface
type UserRepository interface {
	FindByID(id string) (*User, error)
}

// 2. Private Implementation
type userRepository struct {
	db *sql.DB
}

// 3. Constructor returning the Interface
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindByID(id string) (*User, error) {
	// execution logic
	return nil, nil
}
```

## Why this is required
This standardizes dependency injection, makes code inherently mockable for unit tests, prevents tight coupling by hiding concrete internals, and complies with standard Go enterprise patterns.