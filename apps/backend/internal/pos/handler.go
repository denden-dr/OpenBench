package pos

import (
	"errors"
	"math"
	"strconv"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Checkout(c fiber.Ctx) error {
	var req models.CheckoutRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	tx, err := h.service.Checkout(c.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", err.Error())
		}
		if errors.Is(err, ErrInsufficientStock) {
			return utils.SendProblem(c, fiber.StatusConflict, "/errors/conflict", "Conflict - Insufficient Stock", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to process checkout.")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": tx,
	})
}

func (h *Handler) GetTransactionByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Transaction ID is required.")
	}

	tx, err := h.service.GetTransactionByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "/errors/not-found", "Not Found", "Transaction not found.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to retrieve transaction details.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tx,
	})
}

func (h *Handler) GetTransactions(c fiber.Ctx) error {
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

	transactions, total, err := h.service.GetTransactions(c.Context(), limit, offset)
	if err != nil {
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to retrieve transaction list.")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": transactions,
		"meta": fiber.Map{
			"total_data":  total,
			"limit":       limit,
			"offset":      offset,
			"total_pages": totalPages,
		},
	})
}
