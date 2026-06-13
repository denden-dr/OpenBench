//go:build integration

package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/auth"
	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthIntegration(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	db, err := database.NewConnection(&cfg.DB)
	require.NoError(t, err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Clear tables before running tests
	_, err = db.DB.ExecContext(ctx, "DELETE FROM refresh_tokens")
	require.NoError(t, err)
	_, err = db.DB.ExecContext(ctx, "DELETE FROM users")
	require.NoError(t, err)

	// Seed test users
	userPasswordHash, err := auth.HashPassword("UserPassword123!")
	require.NoError(t, err)
	_, err = db.DB.ExecContext(ctx, "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)", "user@openbench.dev", userPasswordHash, "user")
	require.NoError(t, err)

	adminPasswordHash, err := auth.HashPassword("AdminPassword123!")
	require.NoError(t, err)
	_, err = db.DB.ExecContext(ctx, "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)", "admin_test@openbench.dev", adminPasswordHash, "admin")
	require.NoError(t, err)

	// Set up Fiber test application
	app := fiber.New()

	jwtSecret := cfg.JWTSecret
	accessExpiry := 5 * time.Minute
	refreshExpiry := 24 * time.Hour

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, db, jwtSecret)
	authHandler := auth.NewHandler(authService, accessExpiry, refreshExpiry, true)

	// Register test routes
	app.Post("/api/v1/auth/signin", authHandler.SignIn)
	app.Post("/api/v1/auth/refresh", authHandler.Refresh)
	app.Post("/api/v1/auth/signout", authHandler.SignOut)
	app.Get("/api/v1/auth/me", auth.RequireAuth(jwtSecret), authHandler.Me)

	// User-only route
	app.Get("/api/v1/user/profile", auth.RequireAuth(jwtSecret), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		role := c.Locals("user_role").(string)
		return c.JSON(fiber.Map{
			"user_id": userID,
			"role":    role,
		})
	})

	// Admin-only route
	app.Get("/api/v1/admin/dashboard", auth.RequireAuth(jwtSecret), auth.RequireRole("admin"), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		return c.JSON(fiber.Map{
			"message": "Welcome to the admin dashboard",
			"user_id": userID,
		})
	})

	// ==========================================
	// Scenario A: Sign-In Flow
	// ==========================================
	t.Run("SignIn - Invalid Credentials", func(t *testing.T) {
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "WrongPassword!",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	var userAccessTokenCookie, userRefreshTokenCookie string
	t.Run("SignIn - Successful User Login", func(t *testing.T) {
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "UserPassword123!",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, 2000)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp response.APIResponse[auth.SignInResponse]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.Equal(t, "user", apiResp.Data.Role)
		assert.Equal(t, "user@openbench.dev", apiResp.Data.Email)
		assert.NotEmpty(t, apiResp.Data.UserID)

		// Extract cookies
		cookies := resp.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "access_token=") {
				userAccessTokenCookie = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
			if strings.HasPrefix(c, "refresh_token=") {
				userRefreshTokenCookie = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
		}

		assert.NotEmpty(t, userAccessTokenCookie)
		assert.NotEmpty(t, userRefreshTokenCookie)
	})

	t.Run("SignIn - Validation Failure (Empty Password)", func(t *testing.T) {
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var apiResp response.APIResponse[any]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.Contains(t, apiResp.Error, "validation failed")
		assert.Contains(t, apiResp.Error, "Password")
	})

	// ==========================================
	// Scenario B: Middleware Access Controls
	// ==========================================
	t.Run("Middleware - RequireAuth Missing Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Middleware - RequireAuth Valid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
		req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", userAccessTokenCookie))

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "user", body["role"])
	})

	t.Run("Middleware - RequireRole Insufficient Permissions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/dashboard", nil)
		req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", userAccessTokenCookie))

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	var adminAccessTokenCookie string
	t.Run("Middleware - RequireRole Valid Permissions", func(t *testing.T) {
		// Log in as admin
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "admin_test@openbench.dev",
			Password: "AdminPassword123!",
		})
		reqSignIn := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		reqSignIn.Header.Set("Content-Type", "application/json")

		respSignIn, err := app.Test(reqSignIn)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respSignIn.StatusCode)

		cookies := respSignIn.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "access_token=") {
				adminAccessTokenCookie = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
		}

		// Access admin dashboard
		reqDashboard := httptest.NewRequest("GET", "/api/v1/admin/dashboard", nil)
		reqDashboard.Header.Set("Cookie", fmt.Sprintf("access_token=%s", adminAccessTokenCookie))

		respDashboard, err := app.Test(reqDashboard)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respDashboard.StatusCode)
	})

	// ==========================================
	// Scenario C: Token Rotation & Security (RTR)
	// ==========================================
	var nextRefreshTokenCookie string
	t.Run("RTR - Successful Token Rotation", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		req.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", userRefreshTokenCookie))

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cookies := resp.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "refresh_token=") {
				nextRefreshTokenCookie = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
		}
		assert.NotEmpty(t, nextRefreshTokenCookie)
		assert.NotEqual(t, userRefreshTokenCookie, nextRefreshTokenCookie)
	})

	t.Run("RTR - Reuse Old Token (Replay Attack & Revocation)", func(t *testing.T) {
		// Update used_at to be older than the 5-second grace period (e.g., 10 seconds ago)
		hashedOldToken := auth.HashSha256(userRefreshTokenCookie)
		_, err = db.DB.ExecContext(ctx, "UPDATE refresh_tokens SET used_at = $1 WHERE token_hash = $2", time.Now().Add(-10*time.Second), hashedOldToken)
		require.NoError(t, err)

		// Attempt to use the original userRefreshTokenCookie which has already been rotated (used)
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		req.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", userRefreshTokenCookie))

		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should fail
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var apiResp response.APIResponse[any]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.Contains(t, apiResp.Error, "compromise detected")

		// Confirm that the nextRefreshTokenCookie is now revoked in the database due to breach action
		hashedNextToken := auth.HashSha256(nextRefreshTokenCookie)
		var isRevoked bool
		err = db.DB.GetContext(ctx, &isRevoked, "SELECT is_revoked FROM refresh_tokens WHERE token_hash = $1", hashedNextToken)
		require.NoError(t, err)
		assert.True(t, isRevoked, "Token family must be revoked")
	})

	t.Run("RTR - Grace Period Validation", func(t *testing.T) {
		// 1. Log in again to get a fresh refresh token
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "UserPassword123!",
		})
		reqSignIn := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		reqSignIn.Header.Set("Content-Type", "application/json")
		respSignIn, err := app.Test(reqSignIn)
		require.NoError(t, err)

		var rawToken string
		cookies := respSignIn.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "refresh_token=") {
				rawToken = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
		}
		require.NotEmpty(t, rawToken)

		// 2. Refresh it to rotate (first call)
		reqRef1 := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		reqRef1.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", rawToken))
		respRef1, err := app.Test(reqRef1)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respRef1.StatusCode)

		// 3. Immediately refresh using the same old token again (within the 5s grace period)
		reqRef2 := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		reqRef2.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", rawToken))
		respRef2, err := app.Test(reqRef2)
		require.NoError(t, err)
		// Should succeed during the grace period!
		assert.Equal(t, http.StatusOK, respRef2.StatusCode)

		var apiResp response.APIResponse[any]
		err = json.NewDecoder(respRef2.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.Contains(t, apiResp.Message, "grace period")
	})

	t.Run("Session Info & SignOut Flow", func(t *testing.T) {
		// Get Me info without token
		reqMeFail := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
		respMeFail, err := app.Test(reqMeFail)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, respMeFail.StatusCode)

		// Get Me info with token
		reqMeSuccess := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
		reqMeSuccess.Header.Set("Cookie", fmt.Sprintf("access_token=%s", userAccessTokenCookie))
		respMeSuccess, err := app.Test(reqMeSuccess)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respMeSuccess.StatusCode)

		var apiResp response.APIResponse[auth.MeResponse]
		err = json.NewDecoder(respMeSuccess.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.Equal(t, "user", apiResp.Data.Role)
		assert.Equal(t, "user@openbench.dev", apiResp.Data.Email)
		assert.NotEmpty(t, apiResp.Data.UserID)

		// SignOut
		reqSignOut := httptest.NewRequest("POST", "/api/v1/auth/signout", nil)
		reqSignOut.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", userRefreshTokenCookie))
		respSignOut, err := app.Test(reqSignOut)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respSignOut.StatusCode)

		// Verify cookies were cleared in the header
		var hasAccessCleared, hasRefreshCleared bool
		cookies := respSignOut.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "access_token=") {
				hasAccessCleared = true
			}
			if strings.HasPrefix(c, "refresh_token=") {
				hasRefreshCleared = true
			}
		}
		assert.True(t, hasAccessCleared)
		assert.True(t, hasRefreshCleared)
	})
}
