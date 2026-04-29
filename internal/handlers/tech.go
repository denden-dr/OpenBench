package handlers

import (
	"errors"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// TechHandler handles technician-related HTTP requests.
type TechHandler struct {
	techBookingSvc service.TechBookingService
}

// NewTechHandler returns a new TechHandler.
func NewTechHandler(techBookingSvc service.TechBookingService) *TechHandler {
	return &TechHandler{
		techBookingSvc: techBookingSvc,
	}
}

// GetAvailableBookings returns bookings pending diagnosis without an assigned technician.
func (h *TechHandler) GetAvailableBookings() fiber.Handler {
	return func(c fiber.Ctx) error {
		bookings, err := h.techBookingSvc.GetAvailableBookings(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch available bookings"})
		}
		return c.JSON(bookings)
	}
}

// GetMyBookings returns bookings assigned to the current technician.
func (h *TechHandler) GetMyBookings() fiber.Handler {
	return func(c fiber.Ctx) error {
		tech, ok := c.Locals("tech").(*domain.Technician)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve technician profile from context"})
		}

		bookings, err := h.techBookingSvc.GetMyBookings(c.Context(), tech.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch my bookings"})
		}
		return c.JSON(bookings)
	}
}

// AssignBooking claims a booking for the current technician.
func (h *TechHandler) AssignBooking() fiber.Handler {
	return func(c fiber.Ctx) error {
		tech, ok := c.Locals("tech").(*domain.Technician)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve technician profile from context"})
		}

		bookingID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id format"})
		}

		if err := h.techBookingSvc.AssignBooking(c.Context(), tech.UserID, bookingID); err != nil {
			if errors.Is(err, domain.ErrConflict) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "booking already assigned or not in pending state"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to assign booking"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

type DiagnoseRequest struct {
	Diagnosis  string  `json:"diagnosis"`
	Cost       float64 `json:"cost"`
	RepairTime string  `json:"repair_time"`
}

// DiagnoseBooking submits diagnosis details for a booking.
func (h *TechHandler) DiagnoseBooking() fiber.Handler {
	return func(c fiber.Ctx) error {
		tech, ok := c.Locals("tech").(*domain.Technician)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve technician profile from context"})
		}

		bookingID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id format"})
		}

		var req DiagnoseRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		input := service.DiagnoseInput{
			Diagnosis:  req.Diagnosis,
			Cost:       req.Cost,
			RepairTime: req.RepairTime,
		}

		if err := h.techBookingSvc.DiagnoseBooking(c.Context(), tech.UserID, bookingID, input); err != nil {
			if errors.Is(err, domain.ErrInvalidStateTransition) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "booking cannot be diagnosed or does not belong to you"})
			}
			if errors.Is(err, domain.ErrInvalidInput) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid diagnosis data provided"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to submit diagnosis"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

type UpdateStatusRequest struct {
	Status domain.BookingStatus `json:"status"`
}

// UpdateBookingStatus progresses the state of a booking.
func (h *TechHandler) UpdateBookingStatus() fiber.Handler {
	return func(c fiber.Ctx) error {
		tech, ok := c.Locals("tech").(*domain.Technician)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve technician profile from context"})
		}

		bookingID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id format"})
		}

		var req UpdateStatusRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := h.techBookingSvc.UpdateBookingStatus(c.Context(), tech.UserID, bookingID, req.Status); err != nil {
			if errors.Is(err, domain.ErrInvalidStateTransition) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid state transition or booking not assigned to you"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update status"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
