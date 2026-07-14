package inventory

import (
	"errors"

	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type AdjustStockRequest struct {
	QuantityChange int `json:"quantity_change"`
}

type ProductResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *Handler) CreateProduct(c fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	p, err := h.service.CreateProduct(c.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to create product.")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": p,
	})
}

func (h *Handler) UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Product ID is required.")
	}

	var req UpdateProductRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	p, err := h.service.UpdateProduct(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "/errors/not-found", "Not Found", "Product not found.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to update product.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": p,
	})
}

func (h *Handler) AdjustStock(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Product ID is required.")
	}

	var req AdjustStockRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	err := h.service.AdjustStock(c.Context(), id, req.QuantityChange)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "/errors/not-found", "Not Found", "Product not found.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to adjust stock.")
	}

	p, err := h.service.GetProductByID(c.Context(), id)
	if err != nil {
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to fetch updated product.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": p,
	})
}

func (h *Handler) GetProductByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Product ID is required.")
	}

	p, err := h.service.GetProductByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "/errors/not-found", "Not Found", "Product not found.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to retrieve product details.")
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
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to retrieve product list.")
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewCursorPaginatedResponse(products, limit, nextCursor))
}

func (h *Handler) DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "/errors/bad-request", "Bad Request", "Product ID is required.")
	}

	err := h.service.DeleteProduct(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "/errors/not-found", "Not Found", "Product not found.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "/errors/internal-server-error", "Internal Server Error", "Failed to delete product.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
