# Mock Generation with Mockery

Auto-generating mocks with Mockery eliminates manual boilerplate and keeps tests robust by dynamically synchronizing mock code with Go interface updates.

## Decision Matrix: Mocking Strategy

| Scenario | Strategy | Implementation |
|---|---|---|
| **Codebase Interface** | Mockery Auto-Generation | Add `//go:generate mockery` comment |
| **Simple Standard Library Interface** | Inline Manual Mock | Define inline stub/mock using `mock.Mock` |
| **External Package Interface** | Configured Mockery Packages | Set up `.mockery.yaml` package targeting |

## Steps

### 1. Declare Go Generate Directive
Add this directive inside the file containing the interface (e.g. `repository.go`):

```go
//go:generate mockery --name=Repository --output=mocks --outpkg=mocks --case=underscore
```

### 2. Auto-Generate the Mocks
From the workspace package root, run:
```bash
go generate ./...
```

### 3. Write Unit Test
Use the generated struct, which automatically registers `t.Cleanup` assertions on expectations:

```go
package auth_test

import (
	"testing"
	"github.com/denden-dr/openbench/apps/backend/internal/auth/mocks"
	"github.com/stretchr/testify/mock"
)

func TestSignIn(t *testing.T) {
	// Automatically registers t.Cleanup(func() { repo.AssertExpectations(t) })
	repo := mocks.NewRepository(t)
	repo.On("GetUserByID", mock.Anything, "user-123").Return(nil, nil)
}
```

## Quick Reference: Mockery CLI Flags

| Command / Flag | Purpose | Example |
|---|---|---|
| `--name=<name>` | Target interface to generate a mock for | `--name=Repository` |
| `--dir=<path>` | Directory containing interface definitions | `--dir=internal/auth` |
| `--output=<path>` | Output path for the generated code | `--output=internal/auth/mocks` |
| `--outpkg=<name>` | Package name of generated file | `--outpkg=mocks` |
| `--case=underscore` | Use snake_case filename casing | Output is `repository.go` |

## Mistakes to Avoid
*   **Forgetting Generation**: Modifying an interface definition and forgetting to run `go generate ./...`, leading to compilation errors.
*   **Import Cycle Pollution**: Outputting generated mocks into the package hosting the interface itself. Always output to a separate `mocks` package.
*   **Manual Assertion Redundancy**: Writing custom `mock.AssertExpectations(t)` calls at the end of tests instead of using mockery's built-in `mocks.NewRepository(t)` which registers cleanup hooks automatically.
