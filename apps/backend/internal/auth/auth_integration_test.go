//go:build integration

package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/auth"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	testutil.IntegrationSuite
	app       *fiber.App
	jwtSecret string
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (s *AuthHandlerTestSuite) SetupTest() {
	s.IntegrationSuite.SetupTest() // Clean tables dynamically

	// Seed basic users for testing
	ctx := context.Background()
	userPasswordHash, err := auth.HashPassword("UserPassword123!")
	s.Require().NoError(err)
	_, err = s.DB.DB.ExecContext(ctx, "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)", "user@openbench.dev", userPasswordHash, "user")
	s.Require().NoError(err)

	adminPasswordHash, err := auth.HashPassword("AdminPassword123!")
	s.Require().NoError(err)
	_, err = s.DB.DB.ExecContext(ctx, "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)", "admin_test@openbench.dev", adminPasswordHash, "admin")
	s.Require().NoError(err)

	// Initalize Fiber application
	s.app = fiber.New()
	s.jwtSecret = "test_jwt_secret_key_123_456_789"

	accessExpiry := 5 * time.Minute
	refreshExpiry := 24 * time.Hour

	authRepo := auth.NewRepository(s.DB)
	authService := auth.NewService(authRepo, s.DB, s.jwtSecret)
	authHandler := auth.NewHandler(authService, accessExpiry, refreshExpiry, true)

	// Register test routes
	s.app.Post("/api/v1/auth/signin", authHandler.SignIn)
	s.app.Post("/api/v1/auth/refresh", authHandler.Refresh)
	s.app.Post("/api/v1/auth/signout", authHandler.SignOut)
	s.app.Get("/api/v1/auth/me", auth.RequireAuth(s.jwtSecret), authHandler.Me)

	// User-only route
	s.app.Get("/api/v1/user/profile", auth.RequireAuth(s.jwtSecret), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		role := c.Locals("user_role").(string)
		return c.JSON(fiber.Map{
			"user_id": userID,
			"role":    role,
		})
	})

	// Admin-only route
	s.app.Get("/api/v1/admin/dashboard", auth.RequireAuth(s.jwtSecret), auth.RequireRole("admin"), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		return c.JSON(fiber.Map{
			"message": "Welcome to the admin dashboard",
			"user_id": userID,
		})
	})
}

func (s *AuthHandlerTestSuite) TestAuthFlow() {
	var userAccessTokenCookie, userRefreshTokenCookie string
	var adminAccessTokenCookie string
	var nextRefreshTokenCookie string

	// ==========================================
	// Scenario A: Sign-In Flow
	// ==========================================
	s.Run("SignIn - Invalid Credentials", func() {
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "WrongPassword!",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	s.Run("SignIn - Successful User Login", func() {
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "UserPassword123!",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req, 2000)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		var apiResp response.APIResponse[auth.SignInResponse]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Equal("user", apiResp.Data.Role)
		s.Assert().Equal("user@openbench.dev", apiResp.Data.Email)
		s.Assert().NotEmpty(apiResp.Data.UserID)

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

		s.Assert().NotEmpty(userAccessTokenCookie)
		s.Assert().NotEmpty(userRefreshTokenCookie)
	})

	s.Run("SignIn - Validation Failure (Empty Password)", func() {
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)

		var apiResp response.APIResponse[any]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Contains(apiResp.Error, "validation failed")
		s.Assert().Contains(apiResp.Error, "password")
	})

	s.Run("SignIn - Production Secure Cookies", func() {
		// Create a local app with isDev = false to test secure cookie generation
		prodApp := fiber.New()
		authRepo := auth.NewRepository(s.DB)
		authService := auth.NewService(authRepo, s.DB, s.jwtSecret)
		prodAuthHandler := auth.NewHandler(authService, 5*time.Minute, 24*time.Hour, false) // isDev = false
		prodApp.Post("/api/v1/auth/signin", prodAuthHandler.SignIn)

		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "UserPassword123!",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := prodApp.Test(req, 2000)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		cookies := resp.Header.Values("Set-Cookie")
		s.Require().NotEmpty(cookies, "Expected cookies to be set")

		var hasAccessSecure, hasRefreshSecure bool
		for _, c := range cookies {
			if strings.HasPrefix(c, "access_token=") && (strings.Contains(c, "Secure") || strings.Contains(c, "secure")) {
				hasAccessSecure = true
			}
			if strings.HasPrefix(c, "refresh_token=") && (strings.Contains(c, "Secure") || strings.Contains(c, "secure")) {
				hasRefreshSecure = true
			}
		}
		s.Assert().True(hasAccessSecure, "Access token cookie must have Secure flag in production: %v", cookies)
		s.Assert().True(hasRefreshSecure, "Refresh token cookie must have Secure flag in production: %v", cookies)
	})

	// ==========================================
	// Scenario B: Middleware Access Controls
	// ==========================================
	s.Run("Middleware - RequireAuth Missing Token", func() {
		req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	s.Run("Middleware - RequireAuth Valid Token", func() {
		req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
		req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", userAccessTokenCookie))

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		s.Require().NoError(err)
		s.Assert().Equal("user", body["role"])
	})

	s.Run("Middleware - RequireRole Insufficient Permissions", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/dashboard", nil)
		req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", userAccessTokenCookie))

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusForbidden, resp.StatusCode)
	})

	s.Run("Middleware - RequireRole Valid Permissions", func() {
		// Log in as admin
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "admin_test@openbench.dev",
			Password: "AdminPassword123!",
		})
		reqSignIn := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		reqSignIn.Header.Set("Content-Type", "application/json")

		respSignIn, err := s.app.Test(reqSignIn)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, respSignIn.StatusCode)

		cookies := respSignIn.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "access_token=") {
				adminAccessTokenCookie = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
		}

		// Access admin dashboard
		reqDashboard := httptest.NewRequest("GET", "/api/v1/admin/dashboard", nil)
		reqDashboard.Header.Set("Cookie", fmt.Sprintf("access_token=%s", adminAccessTokenCookie))

		respDashboard, err := s.app.Test(reqDashboard)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, respDashboard.StatusCode)
	})

	// ==========================================
	// Scenario C: Token Rotation & Security (RTR)
	// ==========================================
	s.Run("RTR - Successful Token Rotation", func() {
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		req.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", userRefreshTokenCookie))

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		cookies := resp.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "refresh_token=") {
				nextRefreshTokenCookie = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
		}
		s.Assert().NotEmpty(nextRefreshTokenCookie)
		s.Assert().NotEqual(userRefreshTokenCookie, nextRefreshTokenCookie)
	})

	s.Run("RTR - Reuse Old Token (Replay Attack & Revocation)", func() {
		ctx := context.Background()
		hashedOldToken := auth.HashSha256(userRefreshTokenCookie)
		_, err := s.DB.DB.ExecContext(ctx, "UPDATE refresh_tokens SET used_at = $1 WHERE token_hash = $2", time.Now().Add(-10*time.Second), hashedOldToken)
		s.Require().NoError(err)

		// Attempt to use the original userRefreshTokenCookie which has already been rotated (used)
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		req.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", userRefreshTokenCookie))

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		// Should fail
		s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)

		var apiResp response.APIResponse[any]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Contains(apiResp.Message, "Invalid or expired refresh token", "Should return generic message")
		s.Assert().Empty(apiResp.Error, "Should omit internal error details to prevent leakage")

		// Confirm that the nextRefreshTokenCookie is now revoked in the database due to breach action
		hashedNextToken := auth.HashSha256(nextRefreshTokenCookie)
		var isRevoked bool
		err = s.DB.DB.GetContext(ctx, &isRevoked, "SELECT is_revoked FROM refresh_tokens WHERE token_hash = $1", hashedNextToken)
		s.Require().NoError(err)
		s.Assert().True(isRevoked, "Token family must be revoked")
	})

	s.Run("RTR - Grace Period Validation", func() {
		// 1. Log in again to get a fresh refresh token
		reqBody, _ := json.Marshal(auth.SignInRequest{
			Email:    "user@openbench.dev",
			Password: "UserPassword123!",
		})
		reqSignIn := httptest.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(reqBody))
		reqSignIn.Header.Set("Content-Type", "application/json")
		respSignIn, err := s.app.Test(reqSignIn)
		s.Require().NoError(err)

		var rawToken string
		cookies := respSignIn.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.HasPrefix(c, "refresh_token=") {
				rawToken = strings.Split(strings.Split(c, ";")[0], "=")[1]
			}
		}
		s.Require().NotEmpty(rawToken)

		// 2. Refresh it to rotate (first call)
		reqRef1 := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		reqRef1.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", rawToken))
		respRef1, err := s.app.Test(reqRef1)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, respRef1.StatusCode)

		// 3. Immediately refresh using the same old token again (within the 5s grace period)
		reqRef2 := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		reqRef2.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", rawToken))
		respRef2, err := s.app.Test(reqRef2)
		s.Require().NoError(err)
		// Should succeed during the grace period!
		s.Assert().Equal(http.StatusOK, respRef2.StatusCode)

		var apiResp response.APIResponse[any]
		err = json.NewDecoder(respRef2.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Contains(apiResp.Message, "grace period")
	})

	s.Run("Session Info & SignOut Flow", func() {
		// Get Me info without token
		reqMeFail := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
		respMeFail, err := s.app.Test(reqMeFail)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusUnauthorized, respMeFail.StatusCode)

		// Get Me info with token
		reqMeSuccess := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
		reqMeSuccess.Header.Set("Cookie", fmt.Sprintf("access_token=%s", userAccessTokenCookie))
		respMeSuccess, err := s.app.Test(reqMeSuccess)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, respMeSuccess.StatusCode)

		var apiResp response.APIResponse[auth.MeResponse]
		err = json.NewDecoder(respMeSuccess.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Equal("user", apiResp.Data.Role)
		s.Assert().Equal("user@openbench.dev", apiResp.Data.Email)
		s.Assert().NotEmpty(apiResp.Data.UserID)

		// SignOut
		reqSignOut := httptest.NewRequest("POST", "/api/v1/auth/signout", nil)
		reqSignOut.Header.Set("Cookie", fmt.Sprintf("refresh_token=%s", userRefreshTokenCookie))
		respSignOut, err := s.app.Test(reqSignOut)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, respSignOut.StatusCode)

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
		s.Assert().True(hasAccessCleared)
		s.Assert().True(hasRefreshCleared)
	})
}

func TestMain(m *testing.M) {
	// Setup the global database container once for all integration tests in this package
	tdb, err := testutil.SetupTestDB()
	if err != nil {
		log.Fatalf("Failed to setup integration test database: %v", err)
	}

	code := m.Run()

	// Terminate the container
	tdb.Terminate()

	os.Exit(code)
}
