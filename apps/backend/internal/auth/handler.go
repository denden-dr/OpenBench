package auth

import (
	"errors"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

// Handler handles auth HTTP requests
type Handler struct {
	service       Service
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	isDev         bool
}

// NewHandler creates a new auth Handler instance
func NewHandler(service Service, accessExpiry, refreshExpiry time.Duration, isDev bool) *Handler {
	return &Handler{
		service:       service,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		isDev:         isDev,
	}
}

// SignIn handles user login authentication
func (h *Handler) SignIn(c *fiber.Ctx) error {
	var req SignInRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	result, err := h.service.SignIn(c.UserContext(), req.Email, req.Password, h.accessExpiry, h.refreshExpiry)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid email or password", err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Internal server error", err)
	}

	secure := !h.isDev
	sameSite := "Strict"
	if h.isDev {
		sameSite = "Lax"
	}

	// Set HTTPOnly cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Expires:  time.Now().Add(h.accessExpiry),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RawRefreshToken,
		Expires:  result.ExpiresAt,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
	})

	return response.JSON(c, fiber.StatusOK, "Successfully signed in", SignInResponse{
		Role:   result.User.Role,
		UserID: result.User.ID,
		Email:  result.User.Email,
	})
}

// Refresh handles Token Rotation (RTR)
func (h *Handler) Refresh(c *fiber.Ctx) error {
	rawToken := c.Cookies("refresh_token")
	if rawToken == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Refresh token required", errors.New("refresh token required"))
	}

	result, err := h.service.Refresh(c.UserContext(), rawToken, h.accessExpiry, h.refreshExpiry)
	if err != nil {
		if errors.Is(err, ErrTokenCompromised) {
			secure := !h.isDev
			sameSite := "Strict"
			if h.isDev {
				sameSite = "Lax"
			}
			c.Cookie(&fiber.Cookie{
				Name:     "access_token",
				Value:    "",
				Expires:  time.Now().Add(-24 * time.Hour),
				HTTPOnly: true,
				Secure:   secure,
				SameSite: sameSite,
				Path:     "/",
			})
			c.Cookie(&fiber.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Expires:  time.Now().Add(-24 * time.Hour),
				HTTPOnly: true,
				Secure:   secure,
				SameSite: sameSite,
				Path:     "/",
			})
			return response.Error(c, fiber.StatusUnauthorized, "Token compromise detected, session revoked", err)
		}
		if errors.Is(err, ErrInvalidRefreshToken) {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid refresh token", err)
		}
		if errors.Is(err, ErrSessionRevoked) {
			return response.Error(c, fiber.StatusUnauthorized, "Session is revoked", err)
		}
		if errors.Is(err, ErrTokenExpired) {
			return response.Error(c, fiber.StatusUnauthorized, "Refresh token expired", err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Internal server error", err)
	}

	secure := !h.isDev
	sameSite := "Strict"
	if h.isDev {
		sameSite = "Lax"
	}

	if result.GraceRefreshed {
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    result.AccessToken,
			Expires:  time.Now().Add(h.accessExpiry),
			HTTPOnly: true,
			Secure:   secure,
			SameSite: sameSite,
			Path:     "/",
		})

		return response.JSON[any](c, fiber.StatusOK, "Token refreshed during grace period", nil)
	}

	// Set new cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Expires:  time.Now().Add(h.accessExpiry),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RawRefreshToken,
		Expires:  result.ExpiresAt,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
	})

	return response.JSON[any](c, fiber.StatusOK, "Tokens rotated successfully", nil)
}

// SignOut clears authentication cookies and revokes the active refresh token
func (h *Handler) SignOut(c *fiber.Ctx) error {
	rawToken := c.Cookies("refresh_token")
	if rawToken != "" {
		_ = h.service.SignOut(c.UserContext(), rawToken)
	}

	secure := !h.isDev
	sameSite := "Strict"
	if h.isDev {
		sameSite = "Lax"
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
	})

	return response.JSON[any](c, fiber.StatusOK, "Successfully signed out", nil)
}

// Me returns the current authenticated user's ID, role and email
func (h *Handler) Me(c *fiber.Ctx) error {
	userID, okID := c.Locals("user_id").(string)
	_, okRole := c.Locals("user_role").(string)
	if !okID || !okRole {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", errors.New("unauthorized"))
	}

	user, err := h.service.GetUserByID(c.UserContext(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve user information", err)
	}

	return response.JSON(c, fiber.StatusOK, "Successfully retrieved profile info", MeResponse{
		UserID: user.ID,
		Role:   user.Role,
		Email:  user.Email,
	})
}
