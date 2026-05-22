package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/denden-dr/openbench/apps/backend/internal/handler"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type TicketIntegrationTestSuite struct {
	suite.Suite
	db  *sqlx.DB
	app *fiber.App
}

func (s *TicketIntegrationTestSuite) SetupSuite() {
	// Connect to test database
	s.db = SetupTestDB()

	// Initialize application layers
	ticketRepo := repository.NewTicketRepository(s.db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	s.app = fiber.New()
	
	// Register routes to match main.go
	api := s.app.Group("/api/v1")
	tickets := api.Group("/tickets")
	tickets.Post("/", ticketHandler.Create)
	tickets.Get("/:id", ticketHandler.GetByID)
}

func (s *TicketIntegrationTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *TicketIntegrationTestSuite) SetupTest() {
	// Clean DB before each test to ensure isolation
	CleanTestDB(s.db)
}

// Launcher function
func TestTicketIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TicketIntegrationTestSuite))
}

func (s *TicketIntegrationTestSuite) TestCreateTicket() {
	tests := []struct {
		name         string
		reqBody      map[string]interface{}
		expectedCode int
		expectErr    bool
	}{
		{
			name: "Success - valid ticket",
			reqBody: map[string]interface{}{
				"device_type":       "Smartphone",
				"brand":             "Apple",
				"model":             "iPhone 13",
				"issue_description": "Screen broken",
			},
			expectedCode: http.StatusCreated,
			expectErr:    false,
		},
		{
			name: "Error - missing required fields",
			reqBody: map[string]interface{}{
				"device_type": "Smartphone",
			},
			expectedCode: http.StatusBadRequest,
			expectErr:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			bodyBytes, _ := json.Marshal(tc.reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := s.app.Test(req)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedCode, resp.StatusCode)

			var resBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&resBody)
			s.Require().NoError(err)

			if !tc.expectErr {
				s.Equal(true, resBody["success"])
				data := resBody["data"].(map[string]interface{})
				s.NotEmpty(data["id"])
				s.Equal(tc.reqBody["device_type"], data["device_type"])
				s.Equal("received", data["status"]) // default status from migration
			} else {
				s.Equal(false, resBody["success"])
				s.NotEmpty(resBody["error"])
			}
		})
	}
}

func (s *TicketIntegrationTestSuite) TestGetTicket() {
	// 1. Setup: Create a ticket first directly via API
	createReqBody := map[string]interface{}{
		"device_type":       "Smartphone",
		"brand":             "Samsung",
		"model":             "Galaxy S22",
		"issue_description": "Battery replacement",
	}
	createBytes, _ := json.Marshal(createReqBody)
	
	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/", bytes.NewReader(createBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := s.app.Test(createReq)
	
	var createResBody map[string]interface{}
	_ = json.NewDecoder(createResp.Body).Decode(&createResBody)
	data := createResBody["data"].(map[string]interface{})
	validTicketID := data["id"].(string)

	tests := []struct {
		name         string
		ticketID     string
		expectedCode int
		expectErr    bool
	}{
		{
			name:         "Success - ticket found",
			ticketID:     validTicketID,
			expectedCode: http.StatusOK,
			expectErr:    false,
		},
		{
			name:         "Error - ticket not found",
			ticketID:     "00000000-0000-0000-0000-000000000000",
			expectedCode: http.StatusNotFound,
			expectErr:    true,
		},
		{
			name:         "Error - invalid UUID format",
			ticketID:     "invalid-uuid",
			expectedCode: http.StatusBadRequest,
			expectErr:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets/"+tc.ticketID, nil)
			resp, err := s.app.Test(req)

			s.Require().NoError(err)
			s.Require().Equal(tc.expectedCode, resp.StatusCode)

			var resBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&resBody)
			s.Require().NoError(err)

			if !tc.expectErr {
				s.Equal(true, resBody["success"])
				fetchedData := resBody["data"].(map[string]interface{})
				s.Equal(tc.ticketID, fetchedData["id"])
				s.Equal("Samsung", fetchedData["brand"])
				s.Equal("Galaxy S22", fetchedData["model"])
			} else {
				s.Equal(false, resBody["success"])
				s.NotEmpty(resBody["error"])
			}
		})
	}
}
