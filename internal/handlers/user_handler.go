package handlers

import (
	"github.com/denden-dr/OpenBench/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests for users.
type UserHandler interface {
	GetMe(c fiber.Ctx) error
}

type userHandler struct {
	service service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) GetMe(c fiber.Ctx) error {
	// Assuming the user ID is stored in the context by AuthMiddleware
	// For now, we'll try to get it from Locals
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	user, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
