package sales

import (
	"errors"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service SalesService
}

func NewHandler(service SalesService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListSales(c *fiber.Ctx) error {
	sales, err := h.service.ListSales(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to list sales", err)
	}
	return response.JSON(c, fiber.StatusOK, "Sales retrieved successfully", ToSaleListAPI(sales))
}

func (h *Handler) CreateSale(c *fiber.Ctx) error {
	var req api.CreateSaleJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	s, err := h.service.CreateSale(c.UserContext(), &req)
	if err != nil {
		if errors.Is(err, ErrInsufficientStock) || errors.Is(err, ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to record sale", err)
	}
	return response.JSON(c, fiber.StatusCreated, "Sale recorded successfully", ToSaleAPI(s))
}
