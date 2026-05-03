// apps/backend/internal/handler/middleware/role_middleware.go
package middleware

import "github.com/gofiber/fiber/v2"

const (
	RoleUser       = "user"
	RoleTechnician = "technician"
	LocalsRoleKey  = "role"
)

// RoleMiddleware extracts the X-Mock-Role header and stores it in Fiber locals.
// Defaults to "user" when the header is absent.
func RoleMiddleware(c *fiber.Ctx) error {
	role := c.Get("X-Mock-Role")
	if role != RoleTechnician {
		role = RoleUser
	}
	c.Locals(LocalsRoleKey, role)
	return c.Next()
}

// RequireTechnician rejects requests that don't have X-Mock-Role: technician.
func RequireTechnician(c *fiber.Ctx) error {
	if c.Locals(LocalsRoleKey) != RoleTechnician {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Access denied: technician role required",
		})
	}
	return c.Next()
}
