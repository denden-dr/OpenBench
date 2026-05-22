package handler

import (
	"errors"
	"fmt"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TicketHandler struct {
	service service.TicketService
}

func NewTicketHandler(service service.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (h *TicketHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	res, err := h.service.CreateTicket(c.Context(), &req)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errs := make(map[string]string)
			for _, f := range ve {
				errs[f.Field()] = fmt.Sprintf("Field validation for '%s' failed", f.Field())
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Validation failed",
				"details": errs,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *TicketHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ticket ID format",
		})
	}

	res, err := h.service.GetTicket(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Ticket not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *TicketHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ticket ID format",
		})
	}

	var req dto.UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	res, err := h.service.UpdateTicket(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Ticket not found",
			})
		}
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errs := make(map[string]string)
			for _, f := range ve {
				errs[f.Field()] = fmt.Sprintf("Field validation for '%s' failed", f.Field())
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Validation failed",
				"details": errs,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *TicketHandler) List(c *fiber.Ctx) error {
	res, err := h.service.ListTickets(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *TicketHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ticket ID format",
		})
	}

	err := h.service.DeleteTicket(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Ticket not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ticket deleted successfully",
	})
}
