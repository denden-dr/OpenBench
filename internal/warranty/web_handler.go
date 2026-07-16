package warranty

import (
	"github.com/denden-dr/OpenBench/internal/utils"
	warranty_components "github.com/denden-dr/OpenBench/ui/views/components/warranties"
	admin_pages "github.com/denden-dr/OpenBench/ui/views/pages/admin"
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
	return utils.Render(c, admin_pages.WarrantiesPage(nil))
}

// NewClaimPage handles GET /warranties/claims/new
func (h *WebHandler) NewClaimPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	drawerComponent := warranty_components.DrawerNewClaim()

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	return utils.Render(c, admin_pages.WarrantiesPage(drawerComponent))
}

// ClaimDetailPage handles GET /warranties/claims/:id
func (h *WebHandler) ClaimDetailPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	claimID := c.Params("id")
	drawerComponent := warranty_components.DrawerClaimDetail(claimID)

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	return utils.Render(c, admin_pages.WarrantiesPage(drawerComponent))
}
