package ticket

import (
	"strconv"
	"strings"

	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/utils"
	ticket_components "github.com/denden-dr/OpenBench/ui/views/components/tickets"
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

// TicketsPage handles GET /tickets
func (h *WebHandler) TicketsPage(c fiber.Ctx) error {
	status := c.Query("status")
	search := c.Query("search")

	list, _, err := h.service.GetTickets(c.Context(), status, search, 100, "") // Simple unpaginated list for now
	if err != nil {
		return err // Handle appropriately in a real app (e.g., flash message or error page)
	}

	vms := make([]viewmodels.TicketListVM, len(list))
	for i, t := range list {
		vms[i] = viewmodels.TicketListVM{
			TicketID:     t.TicketID,
			TicketNumber: t.TicketNumber,
			Status:       t.Status,
			CustomerName: t.CustomerName,
			DeviceBrand:  t.DeviceBrand,
			DeviceModel:  t.DeviceModel,
			CreatedAt:    t.CreatedAt,
		}
	}

	return utils.Render(c, admin_pages.TicketsPage("tickets", vms, nil))
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

	status := c.Query("status")
	search := c.Query("search")
	list, _, _ := h.service.GetTickets(c.Context(), status, search, 100, "")
	vms := make([]viewmodels.TicketListVM, len(list))
	for i, t := range list {
		vms[i] = viewmodels.TicketListVM{
			TicketID:     t.TicketID,
			TicketNumber: t.TicketNumber,
			Status:       t.Status,
			CustomerName: t.CustomerName,
			DeviceBrand:  t.DeviceBrand,
			DeviceModel:  t.DeviceModel,
			CreatedAt:    t.CreatedAt,
		}
	}

	return utils.Render(c, admin_pages.TicketsPage("tickets", vms, drawerComponent))
}

// TicketDetailPage handles GET /tickets/:id
func (h *WebHandler) TicketDetailPage(c fiber.Ctx) error {
	isHTMX := c.Get("HX-Request") == "true"
	ticketID := c.Params("id")

	t, err := h.service.GetTicketByID(c.Context(), ticketID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Ticket not found")
	}

	repairAction := ""
	if t.RepairAction != nil {
		repairAction = *t.RepairAction
	}
	notes := ""
	if t.Notes != nil {
		notes = *t.Notes
	}

	vm := viewmodels.TicketDetailVM{
		TicketID:         t.TicketID,
		TicketNumber:     t.TicketNumber,
		Status:           t.Status,
		CustomerName:     t.CustomerName,
		CustomerPhone:    t.CustomerPhone,
		DeviceBrand:      t.DeviceBrand,
		DeviceModel:      t.DeviceModel,
		DevicePasscode:   t.DevicePasscode,
		IssueDescription: t.IssueDescription,
		RepairAction:     repairAction,
		Cost:             t.Cost,
		WarrantyDays:     t.WarrantyDays,
		Notes:            notes,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}

	drawerComponent := ticket_components.DrawerTicketDetail(vm)

	if isHTMX {
		return utils.Render(c, drawerComponent)
	}

	list, _, _ := h.service.GetTickets(c.Context(), "", "", 100, "")
	vms := make([]viewmodels.TicketListVM, len(list))
	for i, l := range list {
		vms[i] = viewmodels.TicketListVM{
			TicketID:     l.TicketID,
			TicketNumber: l.TicketNumber,
			Status:       l.Status,
			CustomerName: l.CustomerName,
			DeviceBrand:  l.DeviceBrand,
			DeviceModel:  l.DeviceModel,
			CreatedAt:    l.CreatedAt,
		}
	}

	return utils.Render(c, admin_pages.TicketsPage("tickets", vms, drawerComponent))
}

// CreateTicketWeb handles POST /tickets
func (h *WebHandler) CreateTicketWeb(c fiber.Ctx) error {
	var req CreateTicketRequest

	req.CustomerName = c.FormValue("customer_name")
	req.CustomerPhone = c.FormValue("customer_phone")
	req.DeviceBrand = c.FormValue("device_brand")
	req.DeviceModel = c.FormValue("device_model")
	req.DevicePasscode = c.FormValue("device_passcode")
	req.IssueDescription = c.FormValue("issue_description")
	req.RepairAction = c.FormValue("repair_action")

	cost, _ := strconv.ParseInt(c.FormValue("cost"), 10, 64)
	req.Cost = cost

	warranty, _ := strconv.Atoi(c.FormValue("warranty_days"))
	req.WarrantyDays = warranty

	_, err := h.service.CreateTicket(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateTicketWeb handles PUT /tickets/:id
func (h *WebHandler) UpdateTicketWeb(c fiber.Ctx) error {
	ticketID := c.Params("id")

	// We use EmergencyUpdateTicketRequest here because the form can update status, cost, etc.
	var req EmergencyUpdateTicketRequest

	req.CustomerName = c.FormValue("customer_name")
	req.CustomerPhone = c.FormValue("customer_phone")
	req.DeviceBrand = c.FormValue("device_brand")
	req.DeviceModel = c.FormValue("device_model")
	req.DevicePasscode = c.FormValue("device_passcode")
	req.IssueDescription = c.FormValue("issue_description")
	req.RepairAction = c.FormValue("repair_action")
	req.Notes = c.FormValue("notes")

	// Parse status
	statusStr := c.FormValue("status")
	// The option value could have description in it if not set carefully, but we fixed it in templ to have value="STATUS"
	// However, just to be safe, split by space and take first word
	statusStr = strings.Split(statusStr, " ")[0]
	// Handle parsing models.ServiceTicketStatus (which is a string alias)
	req.Status = models.ServiceTicketStatus(statusStr)

	cost, _ := strconv.ParseInt(c.FormValue("cost"), 10, 64)
	req.Cost = cost

	warranty, _ := strconv.Atoi(c.FormValue("warranty_days"))
	req.WarrantyDays = warranty

	_, err := h.service.EmergencyUpdateTicket(c.Context(), ticketID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
