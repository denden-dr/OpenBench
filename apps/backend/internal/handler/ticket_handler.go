package handler

import (
	"strconv"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
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
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	res, err := h.service.CreateTicket(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ApiResponse{
		Code:    fiber.StatusCreated,
		Message: "Success",
		Data:    res,
	})
}

func (h *TicketHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ticket ID format")
	}

	res, err := h.service.GetTicket(c.UserContext(), id)
	if err != nil {
		return err
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

func (h *TicketHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ticket ID format")
	}

	var req dto.UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	res, err := h.service.UpdateTicket(c.UserContext(), id, &req)
	if err != nil {
		return err
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

func (h *TicketHandler) List(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	search := c.Query("search")
	status := c.Query("status")

	if len(search) > 200 {
		return fiber.NewError(fiber.StatusBadRequest, "Search query too long")
	}

	page := 1
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid page parameter")
		}
		if page < 1 {
			return fiber.NewError(fiber.StatusBadRequest, "Page parameter must be at least 1")
		}
		if page > 10000 {
			return fiber.NewError(fiber.StatusBadRequest, "Page parameter exceeds maximum limit")
		}
	}

	limit := 20
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid limit parameter")
		}
		if limit < 1 || limit > 100 {
			return fiber.NewError(fiber.StatusBadRequest, "Limit parameter must be between 1 and 100")
		}
	}

	res, err := h.service.ListTickets(c.UserContext(), page, limit, search, status)
	if err != nil {
		return err
	}

	return c.JSON(dto.PaginatedTicketsResponse{
		Code:         fiber.StatusOK,
		Message:      "Success",
		Data:         res.Data,
		Total:        res.Total,
		TotalPages:   res.TotalPages,
		Page:         res.Page,
		Limit:        res.Limit,
		StatusCounts: res.StatusCounts,
	})
}

func (h *TicketHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ticket ID format")
	}

	err := h.service.DeleteTicket(c.UserContext(), id)
	if err != nil {
		return err
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Ticket deleted successfully",
	})
}

func (h *TicketHandler) GetPublicByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ticket ID format. Only full UUID is supported.")
	}

	res, err := h.service.GetPublicTicket(c.UserContext(), id)
	if err != nil {
		return err
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

func (h *TicketHandler) TrackPublic(c *fiber.Ctx) error {
	var req dto.PublicTrackRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	uuidResult, err := h.service.TrackPublicTicket(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data: fiber.Map{
			"ticket_id": uuidResult,
		},
	})
}
