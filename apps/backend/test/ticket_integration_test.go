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
	s.db = SetupTestDB()

	ticketRepo := repository.NewTicketRepository(s.db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	s.app = fiber.New()
	api := s.app.Group("/api/v1")
	tickets := api.Group("/tickets")
	tickets.Post("/", ticketHandler.Create)
	tickets.Get("/", ticketHandler.List)
	tickets.Get("/:id", ticketHandler.GetByID)
	tickets.Patch("/:id", ticketHandler.Update)
	tickets.Delete("/:id", ticketHandler.Delete)
}

func (s *TicketIntegrationTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *TicketIntegrationTestSuite) SetupTest() {
	CleanTestDB(s.db)
}

func TestTicketIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TicketIntegrationTestSuite))
}

func (s *TicketIntegrationTestSuite) TestCreateAndListTicket() {
	// Create Ticket
	reqBody := map[string]interface{}{
		"customer_name":   "Budi",
		"customer_gender": "Male",
		"brand":           "Apple",
		"model":           "iPhone 13 Pro",
		"issue":           "LCD Mati",
		"price":           1500000,
		"warranty_days":   30,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	s.True(createRes["success"].(bool))
	data := createRes["data"].(map[string]interface{})
	s.NotEmpty(data["id"])
	s.Equal("service_in", data["status"])

	// List Tickets
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	respList, err := s.app.Test(reqList)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respList.StatusCode)

	var listRes map[string]interface{}
	_ = json.NewDecoder(respList.Body).Decode(&listRes)
	s.True(listRes["success"].(bool))
	tickets := listRes["data"].([]interface{})
	s.Len(tickets, 1)
}

func (s *TicketIntegrationTestSuite) TestUpdateTicket() {
	// 1. Create a ticket first
	reqBody := map[string]interface{}{
		"customer_name":   "Budi",
		"customer_gender": "Male",
		"brand":           "Apple",
		"model":           "iPhone 13 Pro",
		"issue":           "LCD Mati",
		"price":           1500000,
		"warranty_days":   30,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	data := createRes["data"].(map[string]interface{})
	id := data["id"].(string)

	// 2. Update status to picked_up
	updateBody := map[string]interface{}{
		"status": "picked_up",
	}
	updateBytes, _ := json.Marshal(updateBody)
	reqUpdate, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(updateBytes))
	reqUpdate.Header.Set("Content-Type", "application/json")

	respUpdate, err := s.app.Test(reqUpdate)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respUpdate.StatusCode)

	var updateRes map[string]interface{}
	_ = json.NewDecoder(respUpdate.Body).Decode(&updateRes)
	dataUpdate := updateRes["data"].(map[string]interface{})
	s.Equal("picked_up", dataUpdate["status"])
	s.Equal("paid", dataUpdate["payment_status"])
	s.NotEmpty(dataUpdate["exit_date"])
	s.NotEmpty(dataUpdate["warranty_expiry_date"])

	// 3. Update status back to repaired (exit_date and warranty should clear)
	updateBodyBack := map[string]interface{}{
		"status": "repaired",
	}
	updateBytesBack, _ := json.Marshal(updateBodyBack)
	reqUpdateBack, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(updateBytesBack))
	reqUpdateBack.Header.Set("Content-Type", "application/json")

	respUpdateBack, err := s.app.Test(reqUpdateBack)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respUpdateBack.StatusCode)

	var updateResBack map[string]interface{}
	_ = json.NewDecoder(respUpdateBack.Body).Decode(&updateResBack)
	dataUpdateBack := updateResBack["data"].(map[string]interface{})
	s.Equal("repaired", dataUpdateBack["status"])
	s.Nil(dataUpdateBack["exit_date"])
	s.Nil(dataUpdateBack["warranty_expiry_date"])

	// 4. Update price to 0 and check if it gets zeroed (instead of ignored)
	updatePriceBody := map[string]interface{}{
		"price": 0,
	}
	updatePriceBytes, _ := json.Marshal(updatePriceBody)
	reqUpdatePrice, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(updatePriceBytes))
	reqUpdatePrice.Header.Set("Content-Type", "application/json")

	respUpdatePrice, err := s.app.Test(reqUpdatePrice)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respUpdatePrice.StatusCode)

	var updatePriceRes map[string]interface{}
	_ = json.NewDecoder(respUpdatePrice.Body).Decode(&updatePriceRes)
	dataUpdatePrice := updatePriceRes["data"].(map[string]interface{})
	priceVal, ok := dataUpdatePrice["price"].(string)
	if ok {
		s.Equal("0", priceVal)
	} else {
		s.Equal(0.0, dataUpdatePrice["price"].(float64))
	}
}


func (s *TicketIntegrationTestSuite) TestDeleteTicket() {
	// 1. Create a ticket
	reqBody := map[string]interface{}{
		"customer_name":   "Jane",
		"customer_gender": "Female",
		"brand":           "Samsung",
		"model":           "Galaxy S22",
		"issue":           "Baterai Hamil",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	data := createRes["data"].(map[string]interface{})
	id := data["id"].(string)

	// 2. Delete the ticket
	reqDel, _ := http.NewRequest(http.MethodDelete, "/api/v1/tickets/"+id, nil)
	respDel, err := s.app.Test(reqDel)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respDel.StatusCode)

	// 3. Try to get the ticket (should return 404)
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets/"+id, nil)
	respGet, err := s.app.Test(reqGet)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, respGet.StatusCode)
}

func (s *TicketIntegrationTestSuite) TestValidation() {
	// Try creating with invalid gender
	reqBody := map[string]interface{}{
		"customer_name":   "Budi",
		"customer_gender": "InvalidGender",
		"brand":           "Apple",
		"model":           "iPhone",
		"issue":           "Broken screen",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}
