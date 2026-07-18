package pos

import (
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/utils"
	pos_components "github.com/denden-dr/OpenBench/ui/views/components/pos"
	admin_pages "github.com/denden-dr/OpenBench/ui/views/pages/admin"
	"github.com/denden-dr/OpenBench/ui/views/viewmodels"
	"github.com/gofiber/fiber/v3"
)

type WebHandler struct {
	service   Service
	invReader InventoryProductReader
}

func NewWebHandler(service Service, invReader InventoryProductReader) *WebHandler {
	return &WebHandler{
		service:   service,
		invReader: invReader,
	}
}

// CheckoutPage handles GET /pos
func (h *WebHandler) CheckoutPage(c fiber.Ctx) error {
	products, _, err := h.invReader.FindAll(c.Context(), "", 100, "")
	if err != nil {
		products = []models.Product{}
	}
	return utils.Render(c, admin_pages.POSCheckoutPage(viewmodels.NewProductVMs(products)))
}

// InventoryPage handles GET /pos/inventory
// When called by HTMX (search), it returns only the table rows partial.
func (h *WebHandler) InventoryPage(c fiber.Ctx) error {
	search := c.Query("search")
	products, _, err := h.invReader.FindAll(c.Context(), search, 100, "")
	if err != nil {
		products = []models.Product{}
	}

	vms := viewmodels.NewProductVMs(products)

	// HTMX search request: return only the table body rows (partial)
	if c.Get("HX-Request") == "true" && search != "" {
		return utils.Render(c, admin_pages.POSInventoryRows(vms))
	}
	return utils.Render(c, admin_pages.POSInventoryPage(vms, nil))
}

// NewProductPage handles GET /pos/inventory/new
func (h *WebHandler) NewProductPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	drawerComponent := pos_components.DrawerNewProduct()

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	products, _, err := h.invReader.FindAll(c.Context(), "", 100, "")
	if err != nil {
		products = []models.Product{}
	}
	return utils.Render(c, admin_pages.POSInventoryPage(viewmodels.NewProductVMs(products), drawerComponent))
}

// HistoryPage handles GET /pos/history
func (h *WebHandler) HistoryPage(c fiber.Ctx) error {
	transactions, _, err := h.service.GetTransactions(c.Context(), 100, "")
	if err != nil {
		transactions = []models.PosTransaction{}
	}
	return utils.Render(c, admin_pages.POSHistoryPage(viewmodels.NewTransactionVMs(transactions)))
}
