package warranty

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

func (h *Handler) GetWarrantyByTicketID(c fiber.Ctx) error {
	ticketID, err := getTicketID(c)
	if err != nil {
		return err
	}

	w, err := h.service.GetWarrantyByTicketID(c.Context(), ticketID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToWarrantyResponse(w),
	})
}

func (h *Handler) CreateClaim(c fiber.Ctx) error {
	var req CreateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	claim, err := h.service.CreateClaim(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) GetClaims(c fiber.Ctx) error {
	status := c.Query("status")
	search := c.Query("search")

	limit, cursor := utils.ParseCursorPagination(c)

	claims, nextCursor, err := h.service.GetClaims(c.Context(), status, search, limit, cursor)
	if err != nil {
		return err
	}

	var res []ClaimListResponse
	for _, cl := range claims {
		res = append(res, MapToClaimListResponse(cl))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewCursorPaginatedResponse(res, limit, nextCursor))
}

func (h *Handler) GetClaimByID(c fiber.Ctx) error {
	claimID, err := getClaimID(c)
	if err != nil {
		return err
	}

	claim, err := h.service.GetClaimByID(c.Context(), claimID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) UpdateClaim(c fiber.Ctx) error {
	claimID, err := getClaimID(c)
	if err != nil {
		return err
	}

	var req UpdateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	claim, err := h.service.UpdateClaim(c.Context(), claimID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) EvaluateClaim(c fiber.Ctx) error {
	claimID, err := getClaimID(c)
	if err != nil {
		return err
	}

	var req EvaluateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	claim, err := h.service.EvaluateClaim(c.Context(), claimID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) GetWarrantyByTicketNumber(c fiber.Ctx) error {
	ticketNumber, err := getTicketNumber(c)
	if err != nil {
		return err
	}

	w, err := h.service.GetWarrantyByTicketNumber(c.Context(), ticketNumber)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToWarrantyResponse(w),
	})
}

func (h *Handler) UpdateWarrantyStatus(c fiber.Ctx) error {
	warrantyID, err := getWarrantyID(c)
	if err != nil {
		return err
	}

	var req UpdateWarrantyStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format: "+err.Error())
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	w, err := h.service.UpdateWarrantyStatus(c.Context(), warrantyID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToWarrantyResponse(w),
	})
}

func getTicketID(c fiber.Ctx) (string, error) {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Ticket ID is required.")
	}
	return ticketID, nil
}

func getClaimID(c fiber.Ctx) (string, error) {
	claimID := c.Params("claim_id")
	if claimID == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Claim ID is required.")
	}
	return claimID, nil
}

func getTicketNumber(c fiber.Ctx) (string, error) {
	ticketNumber := c.Params("ticket_number")
	if ticketNumber == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Ticket number is required.")
	}
	return ticketNumber, nil
}

func getWarrantyID(c fiber.Ctx) (string, error) {
	warrantyID := c.Params("warranty_id")
	if warrantyID == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Warranty ID is required.")
	}
	return warrantyID, nil
}
