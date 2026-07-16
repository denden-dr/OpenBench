package ticket

import (
	"github.com/denden-dr/OpenBench/internal/utils"
	ticket_components "github.com/denden-dr/OpenBench/ui/views/components/tickets"
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

// TicketsPage handles GET /tickets
func (h *WebHandler) TicketsPage(c fiber.Ctx) error {
	return utils.Render(c, admin_pages.TicketsPage(nil))
}

// NewTicketPage handles GET /tickets/new
func (h *WebHandler) NewTicketPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	drawerComponent := ticket_components.DrawerNewTicket()

	if isHTMX {
		// Just return the drawer snippet
		return utils.Render(c, drawerComponent)
	}

	// Full page refresh: return layout with drawer open
	return utils.Render(c, admin_pages.TicketsPage(drawerComponent))
}

// TicketDetailPage handles GET /tickets/:id
func (h *WebHandler) TicketDetailPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	ticketID := c.Params("id")

	drawerComponent := ticket_components.DrawerTicketDetail(ticketID)

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	return utils.Render(c, admin_pages.TicketsPage(drawerComponent))
}
