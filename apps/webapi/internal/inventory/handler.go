package inventory

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

type AdjustStockRequest struct {
	QuantityChange int `json:"quantity_change" validate:"required"`
}

type ProductResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
func (h *Handler) getProductID(c fiber.Ctx) (string, error) {
	id := c.Params("id")
	if id == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Product ID is required.")
	}
	return id, nil
}

func (h *Handler) CreateProduct(c fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	p, err := h.service.CreateProduct(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": p,
	})
}

func (h *Handler) UpdateProduct(c fiber.Ctx) error {
	id, err := h.getProductID(c)
	if err != nil {
		return err
	}

	var req UpdateProductRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	p, err := h.service.UpdateProduct(c.Context(), id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": p,
	})
}

func (h *Handler) AdjustStock(c fiber.Ctx) error {
	id, err := h.getProductID(c)
	if err != nil {
		return err
	}

	var req AdjustStockRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	err = h.service.AdjustStock(c.Context(), id, req.QuantityChange)
	if err != nil {
		return err
	}

	p, err := h.service.GetProductByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": p,
	})
}

func (h *Handler) GetProductByID(c fiber.Ctx) error {
	id, err := h.getProductID(c)
	if err != nil {
		return err
	}

	p, err := h.service.GetProductByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": p,
	})
}

func (h *Handler) GetProducts(c fiber.Ctx) error {
	search := c.Query("search")

	limit, cursor := utils.ParseCursorPagination(c)

	products, nextCursor, err := h.service.GetProducts(c.Context(), search, limit, cursor)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewCursorPaginatedResponse(products, limit, nextCursor))
}

func (h *Handler) DeleteProduct(c fiber.Ctx) error {
	id, err := h.getProductID(c)
	if err != nil {
		return err
	}

	err = h.service.DeleteProduct(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
