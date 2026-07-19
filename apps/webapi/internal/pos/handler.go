package pos

import (
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/utils"
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
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	tx, err := h.service.Checkout(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": tx,
	})
}

func (h *Handler) GetTransactionByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Transaction ID is required.")
	}

	tx, err := h.service.GetTransactionByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tx,
	})
}

func (h *Handler) GetTransactions(c fiber.Ctx) error {
	limit, cursor := utils.ParseCursorPagination(c)

	transactions, nextCursor, err := h.service.GetTransactions(c.Context(), limit, cursor)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewCursorPaginatedResponse(transactions, limit, nextCursor))
}
