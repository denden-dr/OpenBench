package ticket

import (
	"github.com/denden-dr/OpenBench/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateTicket(c fiber.Ctx) error {
	var req CreateTicketRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	res, err := h.service.CreateTicket(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) GetTickets(c fiber.Ctx) error {
	status := c.Query("status")
	search := c.Query("search")

	limit, cursor := utils.ParseCursorPagination(c)

	tickets, nextCursor, err := h.service.GetTickets(c.Context(), status, search, limit, cursor)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewCursorPaginatedResponse(tickets, limit, nextCursor))
}

func (h *Handler) SearchTickets(c fiber.Ctx) error {
	var req TicketSearchRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if req.Limit <= 0 {
		req.Limit = utils.DefaultLimit
	}
	if req.Limit > utils.MaxLimit {
		req.Limit = utils.MaxLimit
	}

	tickets, nextCursor, err := h.service.SearchTickets(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewCursorPaginatedResponse(tickets, req.Limit, nextCursor))
}

func (h *Handler) GetTicketByID(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Ticket ID is required.")
	}

	res, err := h.service.GetTicketByID(c.Context(), ticketID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) UpdateTicketStatus(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Ticket ID is required.")
	}

	var req ChangeStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	res, err := h.service.UpdateTicketStatus(c.Context(), ticketID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) UpdateTicketDetails(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Ticket ID is required.")
	}

	var req UpdateTicketRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	res, err := h.service.UpdateTicketDetails(c.Context(), ticketID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) EmergencyUpdateTicket(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Ticket ID is required.")
	}

	var req EmergencyUpdateTicketRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	res, err := h.service.EmergencyUpdateTicket(c.Context(), ticketID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
