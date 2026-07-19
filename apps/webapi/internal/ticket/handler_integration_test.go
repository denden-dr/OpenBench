//go:build integration

package ticket_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/events"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/testutils"
	"github.com/denden-dr/OpenBench/internal/ticket"
	"github.com/denden-dr/OpenBench/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWarrantyGenerator struct{}

func (f *fakeWarrantyGenerator) CreateWarranty(ctx context.Context, ticketID string, warrantyDays int) (*models.Warranty, error) {
	return &models.Warranty{
		ID:       "warr-123",
		TicketID: ticketID,
		Status:   models.WarrantyStatusActive,
		Notes:    nil,
	}, nil
}

func TestTicketHandler_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup db container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	err = testutils.CleanTable(db, "service_tickets")
	require.NoError(t, err)

	// Setup dependencies
	cmdRepo := ticket.NewCommandRepository(db)
	queryRepo := ticket.NewQueryRepository(db)
	txManager := database.NewTxManager(db)
	bus := events.NewAsyncEventBus(10)
	defer bus.Close()
	wgen := &fakeWarrantyGenerator{}

	service := ticket.NewService(queryRepo, cmdRepo, txManager, wgen, bus, "this_is_a_secret_key_32_chars_ok")
	handler := ticket.NewHandler(service)

	// Setup Fiber App
	app := fiber.New()
	app.Post("/tickets", handler.CreateTicket)
	app.Get("/tickets", handler.GetTickets)
	app.Add([]string{"QUERY"}, "/tickets/search", handler.SearchTickets)
	app.Get("/tickets/:ticket_id", handler.GetTicketByID)
	app.Patch("/tickets/:ticket_id/status", handler.UpdateTicketStatus)
	app.Put("/tickets/:ticket_id", handler.UpdateTicketDetails)
	app.Put("/tickets/:ticket_id/emergency", handler.EmergencyUpdateTicket)

	var createdTicketID string

	t.Run("Create Ticket", func(t *testing.T) {
		body := ticket.CreateTicketRequest{
			CustomerName:     "John Doe",
			CustomerPhone:    "081234567891",
			DeviceBrand:      "Samsung",
			DeviceModel:      "Galaxy S22",
			DevicePasscode:   "password123",
			IssueDescription: "USB port issue",
			Cost:             250000,
			WarrantyDays:     15,
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/tickets", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respData struct {
			Data ticket.TicketResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.NotEmpty(t, respData.Data.TicketID)
		assert.Equal(t, "John Doe", respData.Data.CustomerName)
		assert.Equal(t, models.StatusReceived, respData.Data.Status)
		assert.Equal(t, "password123", respData.Data.DevicePasscode) // Decrypted passcode returned

		createdTicketID = respData.Data.TicketID
	})

	t.Run("Get Ticket By ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/tickets/%s", createdTicketID), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data ticket.TicketResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, createdTicketID, respData.Data.TicketID)
		assert.Equal(t, "John Doe", respData.Data.CustomerName)
		assert.Equal(t, "password123", respData.Data.DevicePasscode)
	})

	t.Run("Update Ticket Details", func(t *testing.T) {
		body := ticket.UpdateTicketRequest{
			CustomerName:     "John Doe Update",
			CustomerPhone:    "081234567892",
			IssueDescription: "USB port & Headphone jack issue",
			RepairAction:     "Cleaned USB port, replaced jack",
			Cost:             350000,
			WarrantyDays:     30,
			Notes:            "Tested and working",
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/tickets/%s", createdTicketID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data ticket.TicketResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, "John Doe Update", respData.Data.CustomerName)
		assert.Equal(t, int64(350000), respData.Data.Cost)
		assert.Equal(t, 30, respData.Data.WarrantyDays)
		require.NotNil(t, respData.Data.RepairAction)
		assert.Equal(t, "Cleaned USB port, replaced jack", *respData.Data.RepairAction)
	})

	t.Run("Update Ticket Status", func(t *testing.T) {
		body := ticket.ChangeStatusRequest{
			Status: models.StatusRepairing,
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PATCH", fmt.Sprintf("/tickets/%s/status", createdTicketID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data ticket.TicketStatusResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, createdTicketID, respData.Data.TicketID)
		assert.Equal(t, models.StatusRepairing, respData.Data.Status)
	})

	t.Run("Emergency Update Ticket", func(t *testing.T) {
		body := ticket.EmergencyUpdateTicketRequest{
			CustomerName:     "John Doe Emergency",
			CustomerPhone:    "081234567899",
			DeviceBrand:      "Google",
			DeviceModel:      "Pixel 7",
			DevicePasscode:   "9999",
			Status:           models.StatusCompleted, // This triggers handleTicketCompletion -> wgen.CreateWarranty
			IssueDescription: "Total breakdown",
			RepairAction:     "Board swap",
			Cost:             2000000,
			WarrantyDays:     90,
			Notes:            "Swapped board",
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/tickets/%s/emergency", createdTicketID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data ticket.TicketResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, "John Doe Emergency", respData.Data.CustomerName)
		assert.Equal(t, models.StatusCompleted, respData.Data.Status)
		assert.Equal(t, "9999", respData.Data.DevicePasscode)
	})

	t.Run("Search Tickets QUERY Endpoint", func(t *testing.T) {
		isActive := false
		searchReq := ticket.TicketSearchRequest{
			Search:   "Google",
			IsActive: &isActive,
			Limit:    10,
			Cursor:   "",
		}
		bodyBytes, _ := json.Marshal(searchReq)

		req, err := http.NewRequest("QUERY", "/tickets/search", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData utils.CursorPaginatedResponse[ticket.TicketListResponse]
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, 10, respData.Meta.Limit)
		require.Len(t, respData.Data, 1)
		assert.Equal(t, createdTicketID, respData.Data[0].TicketID)
		assert.Equal(t, "John Doe Emergency", respData.Data[0].CustomerName)
	})

	t.Run("Get Tickets List", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/tickets?status=COMPLETED&search=John&limit=10&cursor=", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData utils.CursorPaginatedResponse[ticket.TicketListResponse]
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, 10, respData.Meta.Limit)
		require.Len(t, respData.Data, 1)
		assert.Equal(t, createdTicketID, respData.Data[0].TicketID)
	})
}
