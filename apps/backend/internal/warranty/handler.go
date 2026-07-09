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
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID tiket wajib diisi.")
	}

	w, err := h.service.GetWarrantyByTicketID(c.Context(), ticketID)
	if err != nil {
		if errors.Is(err, ErrWarrantyNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Data garansi tidak ditemukan.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengambil data garansi.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToWarrantyResponse(w),
	})
}

func (h *Handler) CreateClaim(c fiber.Ctx) error {
	var req CreateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	claim, err := h.service.CreateClaim(c.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		if errors.Is(err, ErrWarrantyNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Data garansi tidak ditemukan.")
		}
		if errors.Is(err, ErrWarrantyNotActive) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Masa garansi sudah habis atau tidak aktif.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal membuat klaim garansi.")
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
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengambil daftar klaim.")
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
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID klaim wajib diisi.")
	}

	claim, err := h.service.GetClaimByID(c.Context(), claimID)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Klaim garansi tidak ditemukan.")
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengambil detail klaim.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) UpdateClaimStatus(c fiber.Ctx) error {
	claimID := c.Params("claim_id")
	if claimID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID klaim wajib diisi.")
	}

	var req ChangeClaimStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	claim, err := h.service.UpdateClaimStatus(c.Context(), claimID, req.Status)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Klaim garansi tidak ditemukan.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal memperbarui status klaim.")
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
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID klaim wajib diisi.")
	}

	var req UpdateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	claim, err := h.service.UpdateClaim(c.Context(), claimID, req)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Klaim garansi tidak ditemukan.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal memperbarui data klaim.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) EvaluateClaim(c fiber.Ctx) error {
	claimID := c.Params("claim_id")
	if claimID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID klaim wajib diisi.")
	}

	var req EvaluateClaimRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	claim, err := h.service.EvaluateClaim(c.Context(), claimID, req)
	if err != nil {
		if errors.Is(err, ErrClaimNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Klaim garansi tidak ditemukan.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal mengevaluasi klaim.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToClaimResponse(claim),
	})
}

func (h *Handler) UpdateWarrantyStatus(c fiber.Ctx) error {
	warrantyID := c.Params("warranty_id")
	if warrantyID == "" {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "ID garansi wajib diisi.")
	}

	var req UpdateWarrantyStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", "Format JSON tidak valid.")
	}

	w, err := h.service.UpdateWarrantyStatus(c.Context(), warrantyID, req)
	if err != nil {
		if errors.Is(err, ErrWarrantyNotFound) {
			return utils.SendProblem(c, fiber.StatusNotFound, "https://openbench.local/errors/not-found", "Not Found", "Data garansi tidak ditemukan.")
		}
		if errors.Is(err, ErrInvalidInput) {
			return utils.SendProblem(c, fiber.StatusBadRequest, "https://openbench.local/errors/bad-request", "Bad Request", err.Error())
		}
		return utils.SendProblem(c, fiber.StatusInternalServerError, "https://openbench.local/errors/internal-server-error", "Internal Server Error", "Gagal memperbarui status garansi.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": MapToWarrantyResponse(w),
	})
}
