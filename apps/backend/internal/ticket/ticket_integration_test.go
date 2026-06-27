//go:build integration

package ticket_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/response"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/denden-dr/openbench/apps/backend/internal/ticket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type TicketHandlerTestSuite struct {
	testutil.IntegrationSuite
	app       *fiber.App
	jwtSecret string
}

func TestTicketHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TicketHandlerTestSuite))
}

func (s *TicketHandlerTestSuite) SetupTest() {
	s.IntegrationSuite.SetupTest() // Clean tables dynamically

	s.app = fiber.New()
	s.jwtSecret = "test_jwt_secret_key_123_456_789"

	ticketRepo := ticket.NewRepository(s.DB)
	publicTicketService := ticket.NewService(ticketRepo, s.DB)
	adminTicketService := ticket.NewAdminService(ticketRepo, s.DB)
	ticketHandler := ticket.NewHandler(adminTicketService, publicTicketService)

	// Auth routes/setup if we need authentication tests, but we can test the handler directly by registering the endpoints.
	// We can also protect these routes using standard middlewares or test them unprotected to focus specifically on ticket handler logic.
	// Since we want to test handler behavior, registering them directly is perfect and focused.
	s.app.Post("/api/v1/admin/tickets", ticketHandler.CreateTicket)
	s.app.Get("/api/v1/admin/tickets", ticketHandler.ListTickets)
	s.app.Get("/api/v1/admin/tickets/:id", ticketHandler.GetTicket)
	s.app.Patch("/api/v1/admin/tickets/:id", ticketHandler.UpdateTicket)
	s.app.Post("/api/v1/admin/tickets/:id/emergency", ticketHandler.EmergencyUpdateTicket)
	s.app.Get("/api/v1/admin/warranties", ticketHandler.ListWarranties)
	s.app.Get("/api/v1/tracker/:ticket_number", ticketHandler.GetPublicTrackerTicket)
}

func (s *TicketHandlerTestSuite) TestCreateTicket() {
	s.Run("Success", func() {
		serial := "SN-999"
		reqBody, _ := json.Marshal(api.TicketCreate{
			CustomerName:         "Budi Santoso",
			CustomerPhone:        "0812-3456-7890",
			BrandPhone:           "Samsung",
			ModelPhone:           "Galaxy S24",
			SerialNumber:         &serial,
			DamageDescription:    "LCD retak",
			WarrantyDurationDays: 30,
		})
		req := httptest.NewRequest("POST", "/api/v1/admin/tickets", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusCreated, resp.StatusCode)

		var apiResp response.APIResponse[api.Ticket]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Equal("Budi Santoso", apiResp.Data.CustomerName)
		s.Assert().Equal(30, apiResp.Data.WarrantyDurationDays)
		s.Assert().Contains(apiResp.Data.TicketNumber, "OB-")
	})

	s.Run("Validation Failure - Missing Name", func() {
		reqBody, _ := json.Marshal(api.TicketCreate{
			CustomerPhone:        "0812-3456-7890",
			BrandPhone:           "Samsung",
			ModelPhone:           "Galaxy S24",
			DamageDescription:    "LCD retak",
			WarrantyDurationDays: 30,
		})
		req := httptest.NewRequest("POST", "/api/v1/admin/tickets", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)

		var apiResp response.APIResponse[any]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Contains(apiResp.Error, "customer_name")
	})
}

func (s *TicketHandlerTestSuite) TestGetAndUpdateTicket() {
	// 1. Seed a ticket directly
	ctx := context.Background()
	ticketRepo := ticket.NewRepository(s.DB)
	tID := "00000000-0000-0000-0000-000000000001"
	tkt := &ticket.Ticket{
		ID:                   tID,
		TicketNumber:         "OB-202606-9999-A9X2B8Y3",
		CustomerName:         "Budi",
		CustomerPhone:        "0812",
		BrandPhone:           "Samsung",
		ModelPhone:           "S24",
		DamageDescription:    "Crack",
		Status:               "received",
		DevicePosition:       "warehouse",
		PaymentStatus:        "none",
		WarrantyDurationDays: 30,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	err := ticketRepo.Create(ctx, nil, tkt)
	s.Require().NoError(err)

	s.Run("Get Ticket Success", func() {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/admin/tickets/%s", tID), nil)
		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		var apiResp response.APIResponse[api.Ticket]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Equal("Budi", apiResp.Data.CustomerName)
	})

	s.Run("Get Public Tracker Success", func() {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/tracker/%s", tkt.TicketNumber), nil)
		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)
	})

	s.Run("Update to Picked Up", func() {
		statusVal := api.TicketUpdateStatus("completed")
		dpVal := api.TicketUpdateDevicePosition("picked_up")
		pmVal := api.TicketUpdatePaymentMethod("cash")
		reqBody, _ := json.Marshal(api.TicketUpdate{
			Status:         &statusVal,
			DevicePosition: &dpVal,
			PaymentMethod:  &pmVal,
		})
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/admin/tickets/%s", tID), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		var apiResp response.APIResponse[api.Ticket]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Equal(string(api.TicketStatusCompleted), string(apiResp.Data.Status))
		s.Assert().Equal(string(api.TicketDevicePositionPickedUp), string(apiResp.Data.DevicePosition))
		s.Assert().Equal(string(api.TicketPaymentStatusPaid), string(apiResp.Data.PaymentStatus))
		s.Assert().NotNil(apiResp.Data.PickedUpAt)
	})

	s.Run("List Warranties", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/warranties", nil)
		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		var apiResp response.APIResponse[[]api.Warranty]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Len(apiResp.Data, 1)
		s.Assert().Equal("Budi", apiResp.Data[0].CustomerName)
		s.Assert().Equal("Samsung S24", apiResp.Data[0].DeviceInfo)
	})

	s.Run("Update normal with status reversal fails", func() {
		dpVal := api.TicketUpdateDevicePosition("warehouse")
		reqBody, _ := json.Marshal(api.TicketUpdate{
			DevicePosition: &dpVal,
		})
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/admin/tickets/%s", tID), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.Run("Emergency Update with status reversal succeeds and deletes warranty", func() {
		dpVal := api.TicketUpdateDevicePosition("warehouse")
		newName := "Budi Santoso"
		reqBody, _ := json.Marshal(api.TicketUpdate{
			DevicePosition: &dpVal,
			CustomerName:   &newName,
		})
		req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/admin/tickets/%s/emergency", tID), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, resp.StatusCode)

		var apiResp response.APIResponse[api.Ticket]
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		s.Require().NoError(err)
		s.Assert().Equal(string(api.TicketDevicePositionWarehouse), string(apiResp.Data.DevicePosition))
		s.Assert().Equal("Budi Santoso", apiResp.Data.CustomerName)
		s.Assert().Nil(apiResp.Data.PickedUpAt)

		// List warranties should return 0 items
		reqList := httptest.NewRequest("GET", "/api/v1/admin/warranties", nil)
		respList, err := s.app.Test(reqList)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusOK, respList.StatusCode)

		var warrantiesResp response.APIResponse[[]api.Warranty]
		err = json.NewDecoder(respList.Body).Decode(&warrantiesResp)
		s.Require().NoError(err)
		s.Assert().Len(warrantiesResp.Data, 0)
	})

	s.Run("Invalid UUID format returns 400 Bad Request", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/tickets/invalid-uuid", nil)
		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
	})
}
