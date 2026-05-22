package handler

import (
	"context"
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
				errs[f.Field()] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", f.Field(), f.Tag())
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

	// Validate UUID
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

func (h *TicketHandler) ClaimTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ticket ID format",
		})
	}

	// For mock auth: technician_id comes from the body.
	// In production this would come from the JWT.
	var body struct {
		TechnicianID string `json:"technician_id" validate:"required,uuid"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.service.ClaimTicket(c.Context(), id, body.TechnicianID); err != nil {
		if errors.Is(err, service.ErrClaimConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"error":   "Ticket is already claimed or not available",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Ticket claimed"})
}

func (h *TicketHandler) CompleteDiagnosis(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.CompleteDiagnosis)
}

func (h *TicketHandler) ApproveRepair(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.ApproveRepair)
}

func (h *TicketHandler) CancelRepair(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.CancelRepair)
}

func (h *TicketHandler) CompleteRepair(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.CompleteRepair)
}

func (h *TicketHandler) MarkPickedUp(c *fiber.Ctx) error {
	return h.handleTransition(c, h.service.MarkPickedUp)
}

// handleTransition is a DRY helper for all status transition endpoints.
func (h *TicketHandler) handleTransition(c *fiber.Ctx, action func(ctx context.Context, id string) error) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid ticket ID format",
		})
	}

	if err := action(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Ticket not found",
			})
		}
		if errors.Is(err, service.ErrInvalidTransition) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *TicketHandler) GetBoard(c *fiber.Ctx) error {
	board, err := h.service.ListForBoard(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    board,
	})
}
