package auth

import (
	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/utils"
	auth_components "github.com/denden-dr/OpenBench/ui/views/components/auth"
	auth_pages "github.com/denden-dr/OpenBench/ui/views/pages/auth"
	"github.com/gofiber/fiber/v3"
)

type WebHandler struct {
	service Service
	cfg     *config.Config
}

func NewWebHandler(service Service, cfg *config.Config) *WebHandler {
	return &WebHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *WebHandler) LoginPage(c fiber.Ctx) error {
	return utils.Render(c, auth_pages.LoginPage())
}

func (h *WebHandler) LoginPost(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return utils.Render(c, auth_components.LoginError("Email and password are required."))
	}

	result, err := h.service.Login(c, email, password)
	if err != nil {
		// Render HTMX error snippet
		return utils.Render(c, auth_components.LoginError(err.Error()))
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

	// HTMX specific redirect
	c.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}
