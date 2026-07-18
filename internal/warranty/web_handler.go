package warranty

import (
	"github.com/denden-dr/OpenBench/internal/utils"
	warranty_components "github.com/denden-dr/OpenBench/ui/views/components/warranties"
	admin_pages "github.com/denden-dr/OpenBench/ui/views/pages/admin"
	"github.com/denden-dr/OpenBench/ui/views/viewmodels"
	"github.com/gofiber/fiber/v3"
)

type WebHandler struct {
	service Service
}

func NewWebHandler(service Service) *WebHandler {
	return &WebHandler{
		service: service,
	}
}

// WarrantiesPage handles GET /warranties
func (h *WebHandler) WarrantiesPage(c fiber.Ctx) error {
	status := c.Query("status")
	search := c.Query("search")
	cursor := c.Query("cursor")

	claims, _, err := h.service.GetClaims(c.Context(), status, search, 50, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load claims")
	}

	vms := viewmodels.NewClaimVMs(claims)

	if c.Get("HX-Request") == "true" && search != "" {
		return utils.Render(c, admin_pages.WarrantiesRows(vms))
	}

	return utils.Render(c, admin_pages.WarrantiesPage(vms, nil))
}

// NewClaimPage handles GET /warranties/claims/new
func (h *WebHandler) NewClaimPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	drawerComponent := warranty_components.DrawerNewClaim()

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	claims, _, _ := h.service.GetClaims(c.Context(), "", "", 50, "")
	vms := viewmodels.NewClaimVMs(claims)

	return utils.Render(c, admin_pages.WarrantiesPage(vms, drawerComponent))
}

// ClaimDetailPage handles GET /warranties/claims/:id
func (h *WebHandler) ClaimDetailPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	claimID := c.Params("id")

	claimSummary, err := h.service.GetClaimSummaryByID(c.Context(), claimID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load claim")
	}

	vm := viewmodels.NewClaimVM(*claimSummary)
	drawerComponent := warranty_components.DrawerClaimDetail(vm)

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	claims, _, _ := h.service.GetClaims(c.Context(), "", "", 50, "")
	vms := viewmodels.NewClaimVMs(claims)

	return utils.Render(c, admin_pages.WarrantiesPage(vms, drawerComponent))
}

// VerifyClaim handles POST /warranties/claims/verify
func (h *WebHandler) VerifyClaim(c fiber.Ctx) error {
	ticketNumber := c.FormValue("ticket_number")
	if ticketNumber == "" {
		return utils.Render(c, warranty_components.ClaimVerificationForm("", "Ticket number is required"))
	}

	w, err := h.service.GetWarrantyByTicketNumber(c.Context(), ticketNumber)
	if err != nil {
		return utils.Render(c, warranty_components.ClaimVerificationForm(ticketNumber, "Warranty not found for this ticket"))
	}

	if w.Status != "ACTIVE" {
		return utils.Render(c, warranty_components.ClaimVerificationForm(ticketNumber, "Warranty is expired or void"))
	}

	return utils.Render(c, warranty_components.ClaimSubmissionForm(ticketNumber, ""))
}

// SubmitClaim handles POST /warranties/claims/submit
func (h *WebHandler) SubmitClaim(c fiber.Ctx) error {
	ticketNumber := c.FormValue("ticket_number")
	issueDesc := c.FormValue("issue_description")

	if ticketNumber == "" || issueDesc == "" {
		return utils.Render(c, warranty_components.ClaimSubmissionForm(ticketNumber, "Ticket Number and Issue Description are required"))
	}

	req := CreateClaimRequest{
		TicketNumber:     ticketNumber,
		IssueDescription: issueDesc,
	}

	_, err := h.service.CreateClaim(c.Context(), req)
	if err != nil {
		return utils.Render(c, warranty_components.ClaimSubmissionForm(ticketNumber, err.Error()))
	}

	c.Set("HX-Trigger", "claim-created")
	return c.SendStatus(fiber.StatusNoContent)
}
