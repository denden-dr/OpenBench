# Unit Testing Handlers with Mock Services

By using interfaces for both the repository and the service, we can mock them in tests without needing a database connection.

```go
package auth_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService implements the Service interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) RotateRefreshToken(ctx context.Context, rawToken string) (string, error) {
	args := m.Called(ctx, rawToken)
	return args.String(0), args.Error(1)
}

func TestRefreshHandler_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := new(MockAuthService)
	
	// Set up mock expectation
	mockSvc.On("RotateRefreshToken", mock.Anything, "valid-token").Return("new-token", nil)
	
	handler := auth.NewHandler(mockSvc, true)
	app.Post("/refresh", handler.Refresh)
	
	req := httptest.NewRequest("POST", "/refresh", nil)
	// Add mock cookie to request
	req.Header.Set("Cookie", "refresh_token=valid-token")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}
```
