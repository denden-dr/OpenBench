//go:build integration

package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/testutils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandler_Integration(t *testing.T) {
	ctx := context.Background()

	// Spin up test db container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	// Setup Config
	cfg := &config.Config{
		App: config.AppConfig{
			Env: "test",
		},
		Auth: config.AuthConfig{
			AccessSecret:  "test_access_secret_key_which_is_long_enough",
			RefreshSecret: "test_refresh_secret_key_which_is_long_enough",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 24 * time.Hour,
		},
	}

	cmdRepo := auth.NewCommandRepository(db)
	queryRepo := auth.NewQueryRepository(db)
	authService := auth.NewService(queryRepo, cfg)
	authHandler := auth.NewHandler(authService, cfg)

	// Setup App
	app := fiber.New()
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.Refresh)
	authGroup.Post("/logout", authHandler.Logout)

	// Seed User
	password := "supersecret"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "handler@test.com",
		PasswordHash: string(passwordHash),
		FullName:     "Handler Test User",
		Role:         "ADMIN",
	}

	err = testutils.CleanTable(db, "users")
	require.NoError(t, err)

	err = cmdRepo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("Successful Login", func(t *testing.T) {
		body := map[string]string{
			"email":    user.Email,
			"password": password,
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check response body
		var loginResp struct {
			Data struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)
		assert.NotEmpty(t, loginResp.Data.AccessToken)

		// Check cookies
		var accessCookie, refreshCookie string
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "access_token" {
				accessCookie = cookie.Value
			}
			if cookie.Name == "refresh_token" {
				refreshCookie = cookie.Value
			}
		}
		assert.NotEmpty(t, accessCookie)
		assert.NotEmpty(t, refreshCookie)
	})

	t.Run("Login Failure - Wrong Password", func(t *testing.T) {
		body := map[string]string{
			"email":    user.Email,
			"password": "wrongpassword",
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Refresh Token Flow", func(t *testing.T) {
		// Log in first to get refresh token
		body := map[string]string{
			"email":    user.Email,
			"password": password,
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var refreshCookie *http.Cookie
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "refresh_token" {
				refreshCookie = cookie
			}
		}
		require.NotNil(t, refreshCookie)

		// Attempt to refresh
		refreshReq, err := http.NewRequest("POST", "/api/v1/auth/refresh", nil)
		require.NoError(t, err)
		refreshReq.AddCookie(refreshCookie)

		refreshResp, err := app.Test(refreshReq)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, refreshResp.StatusCode)

		var newAccessCookie string
		for _, cookie := range refreshResp.Cookies() {
			if cookie.Name == "access_token" {
				newAccessCookie = cookie.Value
			}
		}
		assert.NotEmpty(t, newAccessCookie)
	})

	t.Run("Logout Flow", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/v1/auth/logout", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var accessCookie, refreshCookie *http.Cookie
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "access_token" {
				accessCookie = cookie
			}
			if cookie.Name == "refresh_token" {
				refreshCookie = cookie
			}
		}
		assert.True(t, accessCookie == nil || accessCookie.Value == "" || accessCookie.Expires.Before(time.Now()))
		assert.True(t, refreshCookie == nil || refreshCookie.Value == "" || refreshCookie.Expires.Before(time.Now()))
	})
}

// httptest helpers for cookie inspection or request creation if needed.
func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if strings.EqualFold(c.Name, name) {
			return c
		}
	}
	return nil
}
