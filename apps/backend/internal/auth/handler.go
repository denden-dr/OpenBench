package auth

import (
	"errors"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
	cfg     *config.Config
}

func NewHandler(service Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	if req.Email == "" || req.Password == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Email and Password are required.")
	}

	result, err := h.service.Login(c, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Invalid email or password.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "An internal server error occurred.")
	}

	// Set Access Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.App.Env == "production",
		SameSite: "Strict",
		MaxAge:   int(h.cfg.Auth.AccessExpiry.Seconds()),
	})

	// Set Refresh Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/api/v1/auth/refresh",
		HTTPOnly: true,
		Secure:   h.cfg.App.Env == "production",
		SameSite: "Strict",
		MaxAge:   int(h.cfg.Auth.RefreshExpiry.Seconds()),
	})

	return c.Status(fiber.StatusOK).JSON(SuccessResponse[*LoginResponse]{
		Data: result,
	})
}

func (h *Handler) Refresh(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Refresh token is missing.")
	}

	result, err := h.service.Refresh(c, refreshToken)
	if err != nil {
		return utils.SendProblem(c, fiber.StatusUnauthorized, "/errors/unauthorized", "Unauthorized Access", "Refresh token is invalid or expired.")
	}

	// Set new Access Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.App.Env == "production",
		SameSite: "Strict",
		MaxAge:   int(h.cfg.Auth.AccessExpiry.Seconds()),
	})

	return c.Status(fiber.StatusOK).JSON(SuccessResponse[*RefreshResponse]{
		Data: result,
	})
}

func (h *Handler) Logout(c fiber.Ctx) error {
	accessToken := c.Cookies("access_token")
	refreshToken := c.Cookies("refresh_token")

	_ = h.service.Logout(c.Context(), accessToken, refreshToken)

	// Clear Access Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.App.Env == "production",
		SameSite: "Strict",
		MaxAge:   -1,
	})

	// Clear Refresh Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		HTTPOnly: true,
		Secure:   h.cfg.App.Env == "production",
		SameSite: "Strict",
		MaxAge:   -1,
	})

	return c.Status(fiber.StatusOK).JSON(SuccessResponse[MessageResponse]{
		Data: MessageResponse{Message: "Logged out successfully"},
	})
}
