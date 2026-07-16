package pos

import (
	"github.com/denden-dr/OpenBench/internal/utils"
	pos_components "github.com/denden-dr/OpenBench/ui/views/components/pos"
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

// CheckoutPage handles GET /pos
func (h *WebHandler) CheckoutPage(c fiber.Ctx) error {
	return utils.Render(c, admin_pages.POSCheckoutPage())
}

// InventoryPage handles GET /pos/inventory
func (h *WebHandler) InventoryPage(c fiber.Ctx) error {
	return utils.Render(c, admin_pages.POSInventoryPage(nil))
}

// NewProductPage handles GET /pos/inventory/new
func (h *WebHandler) NewProductPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	drawerComponent := pos_components.DrawerNewProduct()

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	return utils.Render(c, admin_pages.POSInventoryPage(drawerComponent))
}

// HistoryPage handles GET /pos/history
func (h *WebHandler) HistoryPage(c fiber.Ctx) error {
	return utils.Render(c, admin_pages.POSHistoryPage())
}
