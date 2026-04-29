package handlers

import (
	"errors"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// BookingHandler handles customer-related repair booking requests.
type BookingHandler struct {
	svc service.BookingService
}

// NewBookingHandler returns a new BookingHandler.
func NewBookingHandler(svc service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

type CreateBookingRequest struct {
	DeviceName       string `json:"device_name"`
	IssueDescription string `json:"issue_description"`
}

// Create handles the creation of a new repair booking.
func (h *BookingHandler) Create() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(*domain.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		var req CreateBookingRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if req.DeviceName == "" || req.IssueDescription == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "device_name and issue_description are required"})
		}

		b, err := h.svc.CreateBooking(c.Context(), user.ID, req.DeviceName, req.IssueDescription)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create booking"})
		}

		return c.Status(fiber.StatusCreated).JSON(b)
	}
}

// GetByID handles retrieving details of a specific booking.
func (h *BookingHandler) GetByID() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(*domain.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id format"})
		}

		b, err := h.svc.GetBookingDetails(c.Context(), user.ID, id)
		if err != nil {
			if errors.Is(err, domain.ErrBookingNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "booking not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch booking"})
		}

		return c.JSON(b)
	}
}

// Approve handles the customer's approval of a diagnosis and cost estimate.
func (h *BookingHandler) Approve() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(*domain.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id format"})
		}

		if err := h.svc.ApproveRepair(c.Context(), user.ID, id); err != nil {
			if errors.Is(err, domain.ErrInvalidStateTransition) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "booking cannot be approved in its current state"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to approve repair"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// Cancel handles the customer's cancellation of a repair booking.
func (h *BookingHandler) Cancel() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(*domain.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id format"})
		}

		if err := h.svc.CancelRepair(c.Context(), user.ID, id); err != nil {
			if errors.Is(err, domain.ErrInvalidStateTransition) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "booking cannot be canceled in its current state"})
			}
			if errors.Is(err, domain.ErrBookingNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "booking not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to cancel repair"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
