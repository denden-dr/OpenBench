package auth

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestToken(secret string, method jwt.SigningMethod, userID string, role string, expiry time.Time) (string, error) {
	token := jwt.NewWithClaims(method, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  expiry.Unix(),
	})
	return token.SignedString([]byte(secret))
}

type mockQueryRepo struct{}

func (m *mockQueryRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *mockQueryRepo) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	return false, nil
}

func TestRequireAuth(t *testing.T) {
	secret := "my_test_jwt_access_secret_longer_than_32_bytes_12345"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AccessSecret: secret,
		},
	}

	app := fiber.New()
	app.Use(RequireAuth(cfg, &mockQueryRepo{}))
	app.Get("/test", func(c fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		userRole := c.Locals("userRole").(string)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"userID":   userID,
			"userRole": userRole,
		})
	})

	tests := []struct {
		name           string
		cookieValue    string
		expectedStatus int
	}{
		{
			name: "valid token",
			cookieValue: func() string {
				tok, _ := generateTestToken(secret, jwt.SigningMethodHS256, "u123", "ADMIN", time.Now().Add(time.Hour))
				return tok
			}(),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "missing token",
			cookieValue:    "",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name: "expired token",
			cookieValue: func() string {
				tok, _ := generateTestToken(secret, jwt.SigningMethodHS256, "u123", "ADMIN", time.Now().Add(-time.Hour))
				return tok
			}(),
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name: "wrong signing key",
			cookieValue: func() string {
				tok, _ := generateTestToken("wrong_secret_key_that_is_long_enough", jwt.SigningMethodHS256, "u123", "ADMIN", time.Now().Add(time.Hour))
				return tok
			}(),
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name: "wrong signing method",
			cookieValue: func() string {
				// JWT None signing method or similar (we can test with an unsupported method, but generating it using HMAC with a different header might be easier)
				// Or simply use a token signed by a different method like RSA (we'll just use a mock signing method)
				token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
					"sub":  "u123",
					"role": "ADMIN",
					"exp":  time.Now().Add(time.Hour).Unix(),
				})
				tok, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				return tok
			}(),
			expectedStatus: fiber.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: tt.cookieValue,
				})
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestRequireRole(t *testing.T) {
	app := fiber.New()

	// Middleware mock/setup to populate Locals
	app.Use(func(c fiber.Ctx) error {
		roleHeader := c.Get("X-Test-Role")
		if roleHeader != "" {
			c.Locals("userRole", roleHeader)
		}
		return c.Next()
	})

	app.Use(RequireRole("ADMIN", "MANAGER"))

	app.Get("/test", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	tests := []struct {
		name           string
		roleHeader     string
		expectedStatus int
	}{
		{
			name:           "matching role ADMIN",
			roleHeader:     "ADMIN",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "matching role MANAGER",
			roleHeader:     "MANAGER",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "non-matching role USER",
			roleHeader:     "USER",
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:           "missing role",
			roleHeader:     "",
			expectedStatus: fiber.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			if tt.roleHeader != "" {
				req.Header.Set("X-Test-Role", tt.roleHeader)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
