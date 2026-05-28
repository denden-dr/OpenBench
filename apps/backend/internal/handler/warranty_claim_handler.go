package handler

import (
	"strconv"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WarrantyClaimHandler struct {
	service service.WarrantyClaimService
}

func NewWarrantyClaimHandler(service service.WarrantyClaimService) *WarrantyClaimHandler {
	return &WarrantyClaimHandler{service: service}
}

func (h *WarrantyClaimHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateWarrantyClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	res, err := h.service.CreateClaim(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ApiResponse{
		Code:    fiber.StatusCreated,
		Message: "Success",
		Data:    res,
	})
}

func (h *WarrantyClaimHandler) List(c *fiber.Ctx) error {
	status := c.Query("status")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if status != "" && status != "waiting_inspection" && status != "approved" && status != "void" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid warranty claim status")
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

	res, err := h.service.ListClaims(c.UserContext(), status, page, limit)
	if err != nil {
		return err
	}

	return c.JSON(dto.PaginatedWarrantyClaimsResponse{
		Code:       fiber.StatusOK,
		Message:    "Success",
		Data:       res.Data,
		Total:      res.Total,
		TotalPages: res.TotalPages,
		Page:       res.Page,
		Limit:      res.Limit,
	})
}

func (h *WarrantyClaimHandler) Approve(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid warranty claim ID format")
	}

	res, err := h.service.ApproveClaim(c.UserContext(), id)
	if err != nil {
		return err
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

func (h *WarrantyClaimHandler) Void(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid warranty claim ID format")
	}

	var req dto.VoidWarrantyClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	res, err := h.service.VoidClaim(c.UserContext(), id, &req)
	if err != nil {
		return err
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
