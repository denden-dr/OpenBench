package auth

import (
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// RequireAuth is a middleware that validates JWT access token from HTTP cookies
func RequireAuth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Authentication required", nil)
		}

		claims, err := ParseAccessToken(accessToken, jwtSecret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired access token", err)
		}

		// Inject user info into Fiber locals context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}

// RequireRole is a middleware that restricts access based on user role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("user_role")
		role, ok := roleVal.(string)
		if !ok || role != requiredRole {
			return response.Error(c, fiber.StatusForbidden, "Access forbidden: insufficient permissions", nil)
		}

		return c.Next()
	}
}
