package middleware

import (
	"errors"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/gofiber/fiber/v3"
)

// RequireTech enforces that the authenticated user has an active Technician profile.
// It must be registered AFTER the standard authentication middleware.
func RequireTech(techRepo repository.TechnicianRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Extract the domain.User from the fiber context (set by previous auth middleware)
		user, ok := c.Locals("user").(*domain.User)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		// Query the technicians table via TechnicianRepository
		tech, err := techRepo.FindByUserID(c.Context(), user.ID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "technician access required"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error verifying role"})
		}

		// Check if technician profile is active
		if !tech.IsActive {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "technician profile is inactive"})
		}

		// Inject the technician profile into the context for downstream handlers
		c.Locals("tech", tech)
		return c.Next()
	}
}
