package warranty

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

func (h *Handler) GetWarrantyByTicketID(c fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	if ticketID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Ticket ID is required.")
	}

	w, err := h.service.GetWarrantyByTicketID(c.Context(), ticketID)
	if err != nil {
		if errors.Is(err, ErrWarrantyNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Warranty data not found.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to retrieve warranty data.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToWarrantyResponse(w),
	})
}

func (h *Handler) CreateClaim(c fiber.Ctx) error {
	var req CreateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	claim, err := h.service.CreateClaim(c.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		if errors.Is(err, ErrWarrantyNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Warranty data not found.")
		}
		if errors.Is(err, ErrWarrantyNotActive) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Warranty period has expired or is inactive.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to create warranty claim.")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) GetClaims(c fiber.Ctx) error {
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

	claims, total, err := h.service.GetClaims(c.Context(), status, search, limit, offset)
	if err != nil {
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to retrieve claim list.")
	}

	var res []ClaimListResponse
	for _, cl := range claims {
		res = append(res, MapToClaimListResponse(cl))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 0
	}

	return c.Status(fiber.StatusOK).JSON(ClaimListWrapper{
		Data: res,
		Meta: ClaimMeta{
			TotalData:  total,
			Limit:      limit,
			Offset:     offset,
			TotalPages: totalPages,
		},
	})
}

func (h *Handler) GetClaimByID(c fiber.Ctx) error {
	claimID := c.Params("claim_id")
	if claimID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Claim ID is required.")
	}

	claim, err := h.service.GetClaimByID(c.Context(), claimID)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Warranty claim not found.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to retrieve claim details.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) UpdateClaimStatus(c fiber.Ctx) error {
	claimID := c.Params("claim_id")
	if claimID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Claim ID is required.")
	}

	var req ChangeClaimStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	claim, err := h.service.UpdateClaimStatus(c.Context(), claimID, req.Status)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Warranty claim not found.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to update claim status.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": ClaimStatusResponse{
			ClaimID:   claim.ID,
			Status:    claim.Status,
			UpdatedAt: claim.UpdatedAt,
		},
	})
}

func (h *Handler) UpdateClaim(c fiber.Ctx) error {
	claimID := c.Params("claim_id")
	if claimID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Claim ID is required.")
	}

	var req UpdateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	claim, err := h.service.UpdateClaim(c.Context(), claimID, req)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Warranty claim not found.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to update claim data.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) EvaluateClaim(c fiber.Ctx) error {
	claimID := c.Params("claim_id")
	if claimID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Claim ID is required.")
	}

	var req EvaluateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	claim, err := h.service.EvaluateClaim(c.Context(), claimID, req)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Warranty claim not found.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to evaluate claim.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) UpdateWarrantyStatus(c fiber.Ctx) error {
	warrantyID := c.Params("warranty_id")
	if warrantyID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Warranty ID is required.")
	}

	var req UpdateWarrantyStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Invalid JSON format.")
	}

	w, err := h.service.UpdateWarrantyStatus(c.Context(), warrantyID, req)
	if err != nil {
		if errors.Is(err, ErrWarrantyNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Warranty data not found.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Failed to update warranty status.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToWarrantyResponse(w),
	})
}
