package ticket

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var ticketNumberRegex = regexp.MustCompile(`(?i)^OB-\d{6}-\d{4,}-[A-Z0-9]{8,12}$`)

type Handler struct {
	adminService  AdminTicketService
	publicService TicketService
}

// NewHandler creates a new ticket handler
func NewHandler(adminService AdminTicketService, publicService TicketService) *Handler {
	return &Handler{
		adminService:  adminService,
		publicService: publicService,
	}
}

// CreateTicket handles the creation of a new ticket
func (h *Handler) CreateTicket(c *fiber.Ctx) error {
	var req api.TicketCreate
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate using go-playground/validator
	if err := validator.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	t, err := h.adminService.CreateTicket(c.UserContext(), &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create ticket", err)
	}

	return response.JSON(c, fiber.StatusCreated, "Ticket created successfully", ToTicketAPI(t))
}

// GetTicket handles retrieving a single ticket by ID
func (h *Handler) GetTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "Ticket ID is required", errors.New("missing ticket id"))
	}
	if _, err := uuid.Parse(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID format (must be UUID)", err)
	}

	t, err := h.adminService.GetTicket(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Ticket not found", err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve ticket", err)
	}

	return response.JSON(c, fiber.StatusOK, "Ticket retrieved successfully", ToTicketAPI(t))
}

// ListTickets handles listing all tickets
func (h *Handler) ListTickets(c *fiber.Ctx) error {
	tickets, err := h.adminService.ListTickets(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to list tickets", err)
	}

	return response.JSON(c, fiber.StatusOK, "Tickets retrieved successfully", ToTicketListAPI(tickets))
}

// UpdateTicket handles partial updates to a ticket
func (h *Handler) UpdateTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "Ticket ID is required", errors.New("missing ticket id"))
	}
	if _, err := uuid.Parse(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID format (must be UUID)", err)
	}

	var req api.TicketUpdate
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate using go-playground/validator
	if err := validator.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	t, err := h.adminService.UpdateTicket(c.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Ticket not found", err)
		}
		if errors.Is(err, ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update ticket", err)
	}

	return response.JSON(c, fiber.StatusOK, "Ticket updated successfully", ToTicketAPI(t))
}

// GetPublicTrackerTicket handles the public tracker route
func (h *Handler) GetPublicTrackerTicket(c *fiber.Ctx) error {
	ticketNumber := c.Params("ticket_number")
	if ticketNumber == "" {
		return response.Error(c, fiber.StatusBadRequest, "Ticket Number is required", errors.New("missing ticket number"))
	}

	if !ticketNumberRegex.MatchString(ticketNumber) {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket number format", errors.New("malformed ticket number"))
	}

	t, err := h.publicService.GetTicketByNumber(c.UserContext(), ticketNumber)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Ticket not found", err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve ticket", err)
	}

	return response.JSON(c, fiber.StatusOK, "Ticket retrieved successfully", ToPublicTrackerTicketAPI(t))
}

// ListWarranties handles listing all warranties for admin dashboard
func (h *Handler) ListWarranties(c *fiber.Ctx) error {
	warranties, err := h.adminService.ListWarranties(c.UserContext())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to list warranties", err)
	}

	return response.JSON(c, fiber.StatusOK, "Warranties retrieved successfully", ToWarrantyListAPI(warranties))
}

// EmergencyUpdateTicket handles emergency administrative updates to a ticket
func (h *Handler) EmergencyUpdateTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "Ticket ID is required", errors.New("missing ticket id"))
	}
	if _, err := uuid.Parse(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID format (must be UUID)", err)
	}

	var req api.TicketUpdate
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate using go-playground/validator
	if err := validator.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	t, err := h.adminService.EmergencyUpdateTicket(c.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Ticket not found", err)
		}
		if errors.Is(err, ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update ticket (emergency)", err)
	}

	return response.JSON(c, fiber.StatusOK, "Ticket updated successfully (emergency)", ToTicketAPI(t))
}
