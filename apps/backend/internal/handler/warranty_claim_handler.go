package handler

import (
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *WarrantyClaimHandler) List(c *fiber.Ctx) error {
	status := c.Query("status")
	if status != "" && status != "waiting_inspection" && status != "approved" && status != "void" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid warranty claim status")
	}
	res, err := h.service.ListClaims(c.UserContext(), status)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
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

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"claim":  res.Claim,
			"ticket": res.Ticket,
		},
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

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"claim":  res.Claim,
			"ticket": res.Ticket,
		},
	})
}
