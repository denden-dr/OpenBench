package inventory

import (
	"errors"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service InventoryService
}

func NewHandler(service InventoryService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListInventory(c *fiber.Ctx) error {
	products, err := h.service.ListInventory(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to list inventory", err)
	}
	return response.JSON(c, fiber.StatusOK, "Inventory retrieved successfully", ToProductListAPI(products))
}

func (h *Handler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "Product ID is required", errors.New("missing product id"))
	}
	if _, err := uuid.Parse(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID format", err)
	}

	p, err := h.service.GetProduct(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Product not found", err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get product", err)
	}
	return response.JSON(c, fiber.StatusOK, "Product retrieved successfully", ToProductAPI(p))
}

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var req api.ProductCreate
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	p, err := h.service.CreateProduct(c.UserContext(), &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create product", err)
	}
	return response.JSON(c, fiber.StatusCreated, "Product created successfully", ToProductAPI(p))
}

func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "Product ID is required", errors.New("missing product id"))
	}
	if _, err := uuid.Parse(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID format", err)
	}

	var req api.ProductUpdate
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	p, err := h.service.UpdateProduct(c.UserContext(), id, &req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Product not found", err)
		}
		if errors.Is(err, ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update product", err)
	}
	return response.JSON(c, fiber.StatusOK, "Product updated successfully", ToProductAPI(p))
}

func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "Product ID is required", errors.New("missing product id"))
	}
	if _, err := uuid.Parse(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID format", err)
	}

	err := h.service.DeleteProduct(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Product not found", err)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete product", err)
	}
	return response.JSON[any](c, fiber.StatusOK, "Product deleted successfully", nil)
}
