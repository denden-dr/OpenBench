package ticket

import (
	"errors"
	"math"
	"strconv"

	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
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
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	res, err := h.service.CreateTicket(c.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal membuat tiket servis.")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) GetTickets(c fiber.Ctx) error {
	status := c.Query("status")
	search := c.Query("search")

	limitStr := c.Query("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offsetStr := c.Query("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	tickets, total, err := h.service.GetTickets(c.Context(), status, search, limit, offset)
	if err != nil {
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengambil daftar tiket.")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 0
	}

	return c.Status(fiber.StatusOK).JSON(TicketListWrapper{
		Data: tickets,
		Meta: TicketMeta{
			TotalData:  total,
			Limit:      limit,
			Offset:     offset,
			TotalPages: totalPages,
		},
	})
}

func (h *Handler) GetTicketByID(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID tiket wajib diisi.")
	}

	res, err := h.service.GetTicketByID(c.Context(), ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Tiket servis tidak ditemukan.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengambil detail tiket.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) UpdateTicketStatus(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID tiket wajib diisi.")
	}

	var req ChangeStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	res, err := h.service.UpdateTicketStatus(c.Context(), ticketID, req)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Tiket servis tidak ditemukan.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengupdate status tiket.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) UpdateTicketDetails(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID tiket wajib diisi.")
	}

	var req UpdateTicketRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	res, err := h.service.UpdateTicketDetails(c.Context(), ticketID, req)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Tiket servis tidak ditemukan.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengupdate detail tiket.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *Handler) EmergencyUpdateTicket(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID tiket wajib diisi.")
	}

	var req EmergencyUpdateTicketRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	res, err := h.service.EmergencyUpdateTicket(c.Context(), ticketID, req)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Tiket servis tidak ditemukan.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal melakukan emergency update tiket.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
