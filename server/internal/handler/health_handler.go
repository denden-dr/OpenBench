package handler

import (
	"openbench/server/internal/dto"
	"openbench/server/internal/service"
	"github.com/gofiber/fiber/v3"
)

type HealthHandler struct {
	healthService service.HealthService
}

func NewHealthHandler(healthService service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h *HealthHandler) GetHealth(c fiber.Ctx) error {
	data := h.healthService.CheckHealth()
	return c.Status(fiber.StatusOK).JSON(dto.NewAPIResponse("OK", fiber.StatusOK, data))
}
