//go:build integration

package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/handler"
	"github.com/denden-dr/openbench/apps/backend/internal/middleware"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type TicketIntegrationTestSuite struct {
	suite.Suite
	db               *sqlx.DB
	app              *fiber.App
	idempotencyStore *database.PostgresStorage
}

func (s *TicketIntegrationTestSuite) SetupSuite() {
	s.db = SetupTestDB()

	ticketRepo := repository.NewTicketRepository(s.db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	s.app = fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	s.idempotencyStore = database.NewPostgresStorage(s.db)
	s.app.Use(middleware.ScopeIdempotencyKey(s.idempotencyStore))
	s.app.Use(middleware.NewIdempotency(s.idempotencyStore))

	api := s.app.Group("/api/v1")
	tickets := api.Group("/tickets")
	tickets.Post("/", ticketHandler.Create)
	tickets.Get("/", ticketHandler.List)
	tickets.Get("/:id", ticketHandler.GetByID)
	tickets.Patch("/:id", ticketHandler.Update)
	tickets.Delete("/:id", ticketHandler.Delete)
}

func (s *TicketIntegrationTestSuite) SetupTest() {
	CleanTestDB(s.T(), s.db)
}

func (s *TicketIntegrationTestSuite) TearDownSuite() {
	if s.idempotencyStore != nil {
		_ = s.idempotencyStore.Close()
	}
}

func TestTicketIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TicketIntegrationTestSuite))
}

func (s *TicketIntegrationTestSuite) TestCreateAndListTicket() {
	// Create Ticket
	reqBody := map[string]interface{}{
		"customer_name":   "Budi",
		"customer_phone":  "081234567890",
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
	defer resp.Body.Close()
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
	defer respList.Body.Close()
	s.Require().Equal(http.StatusOK, respList.StatusCode)

	var listRes map[string]interface{}
	_ = json.NewDecoder(respList.Body).Decode(&listRes)
	s.True(listRes["success"].(bool))
	tickets := listRes["data"].([]interface{})
	s.Len(tickets, 1)
}

// TestDashboardStatusFlow is the contract test that exercises the exact status
// sequence the dashboard uses: service_in → on_process → fixed → picked_up.
// It verifies that each transition is accepted by the backend and that the
// picked_up side effects (payment_status=paid, exit_date, warranty_expiry_date)
// are correctly applied.
func (s *TicketIntegrationTestSuite) TestDashboardStatusFlow() {
	// 1. Create a ticket — initial status must be service_in
	reqBody := map[string]interface{}{
		"customer_name":   "Andi",
		"customer_phone":  "081234567890",
		"customer_gender": "Male",
		"brand":           "Samsung",
		"model":           "Galaxy S23",
		"issue":           "Layar Retak",
		"price":           850000,
		"warranty_days":   14,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	data := createRes["data"].(map[string]interface{})
	id := data["id"].(string)
	s.Equal("service_in", data["status"])
	s.Equal("unpaid", data["payment_status"])

	patchStatus := func(newStatus string) map[string]interface{} {
		body, _ := json.Marshal(map[string]interface{}{"status": newStatus})
		patchReq, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		patchReq.Header.Set("Content-Type", "application/json")
		patchResp, pErr := s.app.Test(patchReq)
		s.Require().NoError(pErr)
		defer patchResp.Body.Close()
		s.Require().Equal(http.StatusOK, patchResp.StatusCode,
			"PATCH to status=%q returned non-200", newStatus)
		var res map[string]interface{}
		_ = json.NewDecoder(patchResp.Body).Decode(&res)
		s.True(res["success"].(bool), "success=false for status=%q", newStatus)
		return res["data"].(map[string]interface{})
	}

	// 2. service_in → on_process (Quick Action: "Start Process")
	onProcessData := patchStatus("on_process")
	s.Equal("on_process", onProcessData["status"])
	s.Equal("unpaid", onProcessData["payment_status"])
	s.Nil(onProcessData["exit_date"])
	s.Nil(onProcessData["warranty_expiry_date"])

	// 3. on_process → fixed (Quick Action: "Mark Fixed")
	fixedData := patchStatus("fixed")
	s.Equal("fixed", fixedData["status"])
	s.Equal("unpaid", fixedData["payment_status"])
	s.Nil(fixedData["exit_date"])
	s.Nil(fixedData["warranty_expiry_date"])

	// 4. fixed → picked_up (Quick Action: "Pickup & Pay")
	//    Side effects: payment_status=paid, exit_date set, warranty_expiry_date set
	pickedUpData := patchStatus("picked_up")
	s.Equal("picked_up", pickedUpData["status"])
	s.Equal("paid", pickedUpData["payment_status"])
	s.NotEmpty(pickedUpData["exit_date"], "exit_date must be set on picked_up")
	s.NotEmpty(pickedUpData["warranty_expiry_date"], "warranty_expiry_date must be set on picked_up")
}

func (s *TicketIntegrationTestSuite) TestUpdateTicket() {
	// 1. Create a ticket first
	reqBody := map[string]interface{}{
		"customer_name":   "Budi",
		"customer_phone":  "081234567890",
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
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	data := createRes["data"].(map[string]interface{})
	id := data["id"].(string)

	// 2. Transition step-by-step to picked_up
	// A. service_in -> on_process
	{
		body, _ := json.Marshal(map[string]interface{}{"status": "on_process"})
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}
	// B. on_process -> fixed
	{
		body, _ := json.Marshal(map[string]interface{}{"status": "fixed"})
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}
	// C. fixed -> picked_up (final state — triggers side effects)
	{
		updateBody := map[string]interface{}{
			"status": "picked_up",
		}
		updateBytes, _ := json.Marshal(updateBody)
		reqUpdate, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(updateBytes))
		reqUpdate.Header.Set("Content-Type", "application/json")

		respUpdate, err := s.app.Test(reqUpdate)
		s.Require().NoError(err)
		defer respUpdate.Body.Close()
		s.Require().Equal(http.StatusOK, respUpdate.StatusCode)

		var updateRes map[string]interface{}
		_ = json.NewDecoder(respUpdate.Body).Decode(&updateRes)
		dataUpdate := updateRes["data"].(map[string]interface{})
		s.Equal("picked_up", dataUpdate["status"])
		s.Equal("paid", dataUpdate["payment_status"])
		s.NotEmpty(dataUpdate["exit_date"])
		s.NotEmpty(dataUpdate["warranty_expiry_date"])
	}

	// 3. Update status back to fixed (must fail because picked_up is terminal)
	updateBodyBack := map[string]interface{}{
		"status": "fixed",
	}
	updateBytesBack, _ := json.Marshal(updateBodyBack)
	reqUpdateBack, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(updateBytesBack))
	reqUpdateBack.Header.Set("Content-Type", "application/json")

	respUpdateBack, err := s.app.Test(reqUpdateBack)
	s.Require().NoError(err)
	defer respUpdateBack.Body.Close()
	s.Require().Equal(http.StatusBadRequest, respUpdateBack.StatusCode)

	// 4. Update price to 0 and check if it gets zeroed (instead of ignored)
	updatePriceBody := map[string]interface{}{
		"price": 0,
	}
	updatePriceBytes, _ := json.Marshal(updatePriceBody)
	reqUpdatePrice, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(updatePriceBytes))
	reqUpdatePrice.Header.Set("Content-Type", "application/json")

	respUpdatePrice, err := s.app.Test(reqUpdatePrice)
	s.Require().NoError(err)
	defer respUpdatePrice.Body.Close()
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
		"customer_phone":  "081234567890",
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
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	data := createRes["data"].(map[string]interface{})
	id := data["id"].(string)

	// 2. Delete the ticket
	reqDel, _ := http.NewRequest(http.MethodDelete, "/api/v1/tickets/"+id, nil)
	respDel, err := s.app.Test(reqDel)
	s.Require().NoError(err)
	defer respDel.Body.Close()
	s.Require().Equal(http.StatusOK, respDel.StatusCode)

	// 3. Try to get the ticket (should return 404)
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets/"+id, nil)
	respGet, err := s.app.Test(reqGet)
	s.Require().NoError(err)
	defer respGet.Body.Close()
	s.Require().Equal(http.StatusNotFound, respGet.StatusCode)
}

func (s *TicketIntegrationTestSuite) TestValidation() {
	// Try creating with invalid gender
	reqBody := map[string]interface{}{
		"customer_name":   "Budi",
		"customer_phone":  "081234567890",
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
	defer resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TicketIntegrationTestSuite) TestValidation_PhoneOptional() {
	reqBody := map[string]interface{}{
		"customer_name":   "Budi Tanpa HP",
		"customer_phone":  "",
		"customer_gender": "Male",
		"brand":           "Apple",
		"model":           "iPhone",
		"issue":           "Broken screen",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	s.True(res["success"].(bool))
	data := res["data"].(map[string]interface{})
	s.Equal("", data["customer_phone"])
}

func (s *TicketIntegrationTestSuite) TestInvalidStatusRejected() {
	// Create a ticket first
	reqBody := map[string]interface{}{
		"customer_name":   "Test",
		"customer_phone":  "081234567890",
		"customer_gender": "Male",
		"brand":           "Xiaomi",
		"model":           "Redmi Note 12",
		"issue":           "Speaker Mati",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	id := createRes["data"].(map[string]interface{})["id"].(string)

	// Try patching with a legacy/invalid status — must be rejected with 400
	for _, badStatus := range []string{"diagnostics", "in_progress", "waiting_parts", "repaired", "done", "cancel"} {
		body, _ := json.Marshal(map[string]interface{}{"status": badStatus})
		patchReq, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		patchReq.Header.Set("Content-Type", "application/json")
		patchResp, pErr := s.app.Test(patchReq)
		s.Require().NoError(pErr)
		defer patchResp.Body.Close()
		s.Equal(http.StatusBadRequest, patchResp.StatusCode,
			"expected 400 for invalid status %q", badStatus)
	}
}

func (s *TicketIntegrationTestSuite) TestBusinessValidationRules() {
	// 1. Create a ticket first
	reqBody := map[string]interface{}{
		"customer_name":   "Business Validation Tester",
		"customer_phone":  "081234567890",
		"customer_gender": "Female",
		"brand":           "Google",
		"model":           "Pixel 8",
		"issue":           "Camera glass cracked",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	id := createRes["data"].(map[string]interface{})["id"].(string)

	// 2. Try updating price to negative value - must return 400
	{
		body, _ := json.Marshal(map[string]interface{}{"price": -100})
		patchReq, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		patchReq.Header.Set("Content-Type", "application/json")
		patchResp, pErr := s.app.Test(patchReq)
		s.Require().NoError(pErr)
		defer patchResp.Body.Close()
		s.Equal(http.StatusBadRequest, patchResp.StatusCode)
	}

	// 3. Try updating warranty days to negative value - must return 400
	{
		body, _ := json.Marshal(map[string]interface{}{"warranty_days": -1})
		patchReq, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		patchReq.Header.Set("Content-Type", "application/json")
		patchResp, pErr := s.app.Test(patchReq)
		s.Require().NoError(pErr)
		defer patchResp.Body.Close()
		s.Equal(http.StatusBadRequest, patchResp.StatusCode)
	}

	// 4. Try updating status to picked_up with payment_status unpaid - must return 400
	{
		body, _ := json.Marshal(map[string]interface{}{
			"status":         "picked_up",
			"payment_status": "unpaid",
		})
		patchReq, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		patchReq.Header.Set("Content-Type", "application/json")
		patchResp, pErr := s.app.Test(patchReq)
		s.Require().NoError(pErr)
		defer patchResp.Body.Close()
		s.Equal(http.StatusBadRequest, patchResp.StatusCode)
	}

	// 5. Try creating a ticket with negative price - must return 400
	{
		badReqBody := map[string]interface{}{
			"customer_name":   "Negative Price Creator",
			"customer_phone":  "081234567890",
			"customer_gender": "Female",
			"brand":           "Google",
			"model":           "Pixel 8",
			"issue":           "Camera glass cracked",
			"price":           -50,
		}
		badBodyBytes, _ := json.Marshal(badReqBody)
		badReq, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(badBodyBytes))
		badReq.Header.Set("Content-Type", "application/json")
		badResp, cErr := s.app.Test(badReq)
		s.Require().NoError(cErr)
		defer badResp.Body.Close()
		s.Equal(http.StatusBadRequest, badResp.StatusCode)
	}

	// 6. Try creating a ticket with negative warranty_days - must return 400
	{
		badReqBody := map[string]interface{}{
			"customer_name":   "Negative Warranty Creator",
			"customer_phone":  "081234567890",
			"customer_gender": "Female",
			"brand":           "Google",
			"model":           "Pixel 8",
			"issue":           "Camera glass cracked",
			"warranty_days":   -10,
		}
		badBodyBytes, _ := json.Marshal(badReqBody)
		badReq, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(badBodyBytes))
		badReq.Header.Set("Content-Type", "application/json")
		badResp, cErr := s.app.Test(badReq)
		s.Require().NoError(cErr)
		defer badResp.Body.Close()
		s.Equal(http.StatusBadRequest, badResp.StatusCode)
	}

	// 7. Try patching with a non-picked_up ticket having an exit_date — must return 400
	{
		body, _ := json.Marshal(map[string]interface{}{
			"status":    "fixed",
			"exit_date": time.Now().Format(time.RFC3339),
		})
		patchReq, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		patchReq.Header.Set("Content-Type", "application/json")
		patchResp, pErr := s.app.Test(patchReq)
		s.Require().NoError(pErr)
		defer patchResp.Body.Close()
		s.Equal(http.StatusBadRequest, patchResp.StatusCode)
	}

	// 8. Try patching with a picked_up ticket and update warranty_days and exit_date — must succeed, persist exit_date, and return dynamically calculated warranty_expiry_date
	{
		// First transition service_in -> on_process
		{
			b, _ := json.Marshal(map[string]interface{}{"status": "on_process"})
			req, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			res, pErr := s.app.Test(req)
			s.Require().NoError(pErr)
			res.Body.Close()
			s.Equal(http.StatusOK, res.StatusCode)
		}
		// Then transition on_process -> fixed
		{
			b, _ := json.Marshal(map[string]interface{}{"status": "fixed"})
			req, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			res, pErr := s.app.Test(req)
			s.Require().NoError(pErr)
			res.Body.Close()
			s.Equal(http.StatusOK, res.StatusCode)
		}

		exitDate := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
		body, _ := json.Marshal(map[string]interface{}{
			"status":         "picked_up",
			"payment_status": "paid",
			"exit_date":      exitDate.Format(time.RFC3339),
			"warranty_days":  45,
		})
		patchReq, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(body))
		patchReq.Header.Set("Content-Type", "application/json")
		patchResp, pErr := s.app.Test(patchReq)
		s.Require().NoError(pErr)
		defer patchResp.Body.Close()
		s.Equal(http.StatusOK, patchResp.StatusCode)

		var updateRes map[string]interface{}
		_ = json.NewDecoder(patchResp.Body).Decode(&updateRes)
		data := updateRes["data"].(map[string]interface{})

		s.Equal("picked_up", data["status"])
		s.Equal("paid", data["payment_status"])
		s.Equal(float64(45), data["warranty_days"])

		resExitDateStr := data["exit_date"].(string)
		resExpiryDateStr := data["warranty_expiry_date"].(string)

		resExitDate, err := time.Parse(time.RFC3339, resExitDateStr)
		s.Require().NoError(err)
		resExpiryDate, err := time.Parse(time.RFC3339, resExpiryDateStr)
		s.Require().NoError(err)

		s.True(exitDate.Equal(resExitDate), "expected exit_date %v, got %v", exitDate, resExitDate)
		expectedExpiry := exitDate.AddDate(0, 0, 45)
		s.True(expectedExpiry.Equal(resExpiryDate), "expected warranty_expiry_date %v, got %v", expectedExpiry, resExpiryDate)
	}
}

func (s *TicketIntegrationTestSuite) TestIdempotency_DuplicatePOST() {
	key := uuid.New().String()
	reqBody := map[string]interface{}{
		"customer_name":   "Budi duplicate",
		"customer_phone":  "081234567890",
		"customer_gender": "Male",
		"brand":           "Apple",
		"model":           "iPhone 13 Pro",
		"issue":           "LCD Mati",
		"price":           1500000,
		"warranty_days":   30,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// First Request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Idempotency-Key", key)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var res1 map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res1)
	data1 := res1["data"].(map[string]interface{})
	id1 := data1["id"].(string)

	// Second Request (Duplicate)
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Idempotency-Key", key)
	resp2, err := s.app.Test(req2)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Require().Equal(http.StatusCreated, resp2.StatusCode)

	var res2 map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&res2)
	data2 := res2["data"].(map[string]interface{})
	id2 := data2["id"].(string)

	s.Equal(id1, id2)

	// Check DB has only 1 ticket
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	respList, err := s.app.Test(reqList)
	s.Require().NoError(err)
	defer respList.Body.Close()
	var listRes map[string]interface{}
	_ = json.NewDecoder(respList.Body).Decode(&listRes)
	tickets := listRes["data"].([]interface{})
	s.Len(tickets, 1)
}

func (s *TicketIntegrationTestSuite) TestIdempotency_DuplicatePATCH() {
	// Create ticket first
	reqBody := map[string]interface{}{
		"customer_name":   "PATCH test",
		"customer_phone":  "081234567890",
		"customer_gender": "Female",
		"brand":           "Xiaomi",
		"model":           "Redmi Note 10",
		"issue":           "Speaker Mati",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	id := createRes["data"].(map[string]interface{})["id"].(string)

	key := uuid.New().String()
	patchBody := map[string]interface{}{"status": "on_process"}
	patchBytes, _ := json.Marshal(patchBody)

	// First PATCH
	reqPatch, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(patchBytes))
	reqPatch.Header.Set("Content-Type", "application/json")
	reqPatch.Header.Set("X-Idempotency-Key", key)
	respPatch, err := s.app.Test(reqPatch)
	s.Require().NoError(err)
	defer respPatch.Body.Close()
	s.Equal(http.StatusOK, respPatch.StatusCode)

	var patchRes1 map[string]interface{}
	_ = json.NewDecoder(respPatch.Body).Decode(&patchRes1)
	s.Equal("on_process", patchRes1["data"].(map[string]interface{})["status"])

	// Second PATCH
	reqPatch2, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(patchBytes))
	reqPatch2.Header.Set("Content-Type", "application/json")
	reqPatch2.Header.Set("X-Idempotency-Key", key)
	respPatch2, err := s.app.Test(reqPatch2)
	s.Require().NoError(err)
	defer respPatch2.Body.Close()
	s.Equal(http.StatusOK, respPatch2.StatusCode)

	var patchRes2 map[string]interface{}
	_ = json.NewDecoder(respPatch2.Body).Decode(&patchRes2)
	s.Equal("on_process", patchRes2["data"].(map[string]interface{})["status"])
}

func (s *TicketIntegrationTestSuite) TestIdempotency_ConcurrentPOST() {
	key := uuid.New().String()
	reqBody := map[string]interface{}{
		"customer_name":   "Budi concurrent",
		"customer_phone":  "081234567890",
		"customer_gender": "Male",
		"brand":           "Apple",
		"model":           "iPhone 13 Pro",
		"issue":           "LCD Mati",
		"price":           1500000,
		"warranty_days":   30,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	var wg sync.WaitGroup
	var resp1, resp2 *http.Response
	var err1, err2 error

	wg.Add(2)
	go func() {
		defer wg.Done()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Idempotency-Key", key)
		resp1, err1 = s.app.Test(req, 10000)
	}()

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Idempotency-Key", key)
		resp2, err2 = s.app.Test(req, 10000)
	}()

	wg.Wait()

	s.Require().NoError(err1)
	s.Require().NoError(err2)
	s.Require().NotNil(resp1)
	s.Require().NotNil(resp2)
	defer resp1.Body.Close()
	defer resp2.Body.Close()

	s.Require().Equal(http.StatusCreated, resp1.StatusCode)
	s.Require().Equal(http.StatusCreated, resp2.StatusCode)

	var res1, res2 map[string]interface{}
	_ = json.NewDecoder(resp1.Body).Decode(&res1)
	_ = json.NewDecoder(resp2.Body).Decode(&res2)

	id1 := res1["data"].(map[string]interface{})["id"].(string)
	id2 := res2["data"].(map[string]interface{})["id"].(string)

	s.Equal(id1, id2)

	// Check DB has only 1 ticket
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	respList, err := s.app.Test(reqList)
	s.Require().NoError(err)
	defer respList.Body.Close()
	var listRes map[string]interface{}
	_ = json.NewDecoder(respList.Body).Decode(&listRes)
	tickets := listRes["data"].([]interface{})
	s.Len(tickets, 1)
}

func (s *TicketIntegrationTestSuite) TestIdempotency_KeyIsolation() {
	key := uuid.New().String()

	// 1. Send POST with key
	reqBody := map[string]interface{}{
		"customer_name":   "Budi isolation",
		"customer_phone":  "081234567890",
		"customer_gender": "Male",
		"brand":           "Apple",
		"model":           "iPhone 13 Pro",
		"issue":           "LCD Mati",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Idempotency-Key", key)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var createRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&createRes)
	id := createRes["data"].(map[string]interface{})["id"].(string)

	// 2. Reuse the same key for PATCH on the created ticket
	patchBody := map[string]interface{}{"status": "on_process"}
	patchBytes, _ := json.Marshal(patchBody)
	reqPatch, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id, bytes.NewReader(patchBytes))
	reqPatch.Header.Set("Content-Type", "application/json")
	reqPatch.Header.Set("X-Idempotency-Key", key) // Same key
	respPatch, err := s.app.Test(reqPatch)
	s.Require().NoError(err)
	defer respPatch.Body.Close()

	// Should succeed and return the update response, not the cached create response
	s.Equal(http.StatusOK, respPatch.StatusCode)
	var patchRes map[string]interface{}
	_ = json.NewDecoder(respPatch.Body).Decode(&patchRes)
	s.Equal("on_process", patchRes["data"].(map[string]interface{})["status"])
}

func (s *TicketIntegrationTestSuite) TestIdempotency_TicketSpecificPATCHKeyIsolation() {
	key := uuid.New().String()

	// Create Ticket 1
	body1 := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "Ticket 1", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b1, _ := json.Marshal(body1)
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b1))
	req1.Header.Set("Content-Type", "application/json")
	r1, _ := s.app.Test(req1)
	defer r1.Body.Close()
	var res1 map[string]interface{}
	_ = json.NewDecoder(r1.Body).Decode(&res1)
	id1 := res1["data"].(map[string]interface{})["id"].(string)

	// Create Ticket 2
	body2 := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "Ticket 2", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b2, _ := json.Marshal(body2)
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	r2, _ := s.app.Test(req2)
	defer r2.Body.Close()
	var res2 map[string]interface{}
	_ = json.NewDecoder(r2.Body).Decode(&res2)
	id2 := res2["data"].(map[string]interface{})["id"].(string)

	// Send PATCH on Ticket 1 with key
	patchBody := map[string]interface{}{"status": "on_process"}
	pb, _ := json.Marshal(patchBody)
	reqPatch1, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id1, bytes.NewReader(pb))
	reqPatch1.Header.Set("Content-Type", "application/json")
	reqPatch1.Header.Set("X-Idempotency-Key", key)
	rPatch1, _ := s.app.Test(reqPatch1)
	defer rPatch1.Body.Close()
	s.Equal(http.StatusOK, rPatch1.StatusCode)

	// Send PATCH on Ticket 2 with same key and same body
	reqPatch2, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id2, bytes.NewReader(pb))
	reqPatch2.Header.Set("Content-Type", "application/json")
	reqPatch2.Header.Set("X-Idempotency-Key", key)
	rPatch2, _ := s.app.Test(reqPatch2)
	defer rPatch2.Body.Close()

	// Should succeed (200), updating Ticket 2, and not returning Ticket 1's cached response
	s.Equal(http.StatusOK, rPatch2.StatusCode)
	var patchRes2 map[string]interface{}
	_ = json.NewDecoder(rPatch2.Body).Decode(&patchRes2)
	s.Equal(id2, patchRes2["data"].(map[string]interface{})["id"].(string))
}

func (s *TicketIntegrationTestSuite) TestIdempotency_BodyFingerprintConflict() {
	key := uuid.New().String()

	body1 := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "Budi conflict 1", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b1, _ := json.Marshal(body1)
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-Idempotency-Key", key)
	r1, _ := s.app.Test(req1)
	defer r1.Body.Close()
	s.Equal(http.StatusCreated, r1.StatusCode)

	// Send different body with same key
	body2 := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "Budi conflict 2", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b2, _ := json.Marshal(body2)
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Idempotency-Key", key)
	r2, _ := s.app.Test(req2)
	defer r2.Body.Close()

	s.Equal(http.StatusConflict, r2.StatusCode)
}

func (s *TicketIntegrationTestSuite) TestIdempotency_ValidationErrorRequiresFreshKeyForCorrectedPayload() {
	key := uuid.New().String()

	invalidBody := map[string]interface{}{
		"customer_name":   "Validation retry",
		"customer_phone":  "081234567890",
		"customer_gender": "Male",
	}
	invalidBytes, _ := json.Marshal(invalidBody)
	invalidReq, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(invalidBytes))
	invalidReq.Header.Set("Content-Type", "application/json")
	invalidReq.Header.Set("X-Idempotency-Key", key)
	invalidResp, err := s.app.Test(invalidReq)
	s.Require().NoError(err)
	defer invalidResp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, invalidResp.StatusCode)

	correctedBody := map[string]interface{}{
		"customer_name":   "Validation retry",
		"customer_phone":  "081234567890",
		"customer_gender": "Male",
		"brand":           "Apple",
		"model":           "iPhone 13",
		"issue":           "LCD Mati",
	}
	correctedBytes, _ := json.Marshal(correctedBody)

	reusedKeyReq, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(correctedBytes))
	reusedKeyReq.Header.Set("Content-Type", "application/json")
	reusedKeyReq.Header.Set("X-Idempotency-Key", key)
	reusedKeyResp, err := s.app.Test(reusedKeyReq)
	s.Require().NoError(err)
	defer reusedKeyResp.Body.Close()
	s.Require().Equal(http.StatusConflict, reusedKeyResp.StatusCode)

	freshKeyReq, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(correctedBytes))
	freshKeyReq.Header.Set("Content-Type", "application/json")
	freshKeyReq.Header.Set("X-Idempotency-Key", uuid.New().String())
	freshKeyResp, err := s.app.Test(freshKeyReq)
	s.Require().NoError(err)
	defer freshKeyResp.Body.Close()
	s.Require().Equal(http.StatusCreated, freshKeyResp.StatusCode)
}

func (s *TicketIntegrationTestSuite) TestIdempotency_InvalidKey() {
	reqBody := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "Invalid key test", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Idempotency-Key", "not-a-valid-uuid")
	resp, _ := s.app.Test(req)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TicketIntegrationTestSuite) TestIdempotency_NoKeyCompatibility() {
	reqBody := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "No key test 1", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := s.app.Test(req)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	reqBody2 := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "No key test 2", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b2, _ := json.Marshal(reqBody2)
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := s.app.Test(req2)
	defer resp2.Body.Close()
	s.Equal(http.StatusCreated, resp2.StatusCode)

	// Verify both were created
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	respList, _ := s.app.Test(reqList)
	defer respList.Body.Close()
	var listRes map[string]interface{}
	_ = json.NewDecoder(respList.Body).Decode(&listRes)
	tickets := listRes["data"].([]interface{})
	s.Len(tickets, 2)
}

func (s *TicketIntegrationTestSuite) TestIdempotency_ScopedHeaderSpoofing() {
	reqBody := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "Spoof test", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Scoped-Idempotency-Key", "POST:/api/v1/tickets:some-uuid")
	resp, _ := s.app.Test(req)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	id := res["data"].(map[string]interface{})["id"].(string)

	// Re-send with same spoofed header, should create another ticket instead of replaying
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Scoped-Idempotency-Key", "POST:/api/v1/tickets:some-uuid")
	resp2, _ := s.app.Test(req2)
	defer resp2.Body.Close()
	s.Equal(http.StatusCreated, resp2.StatusCode)

	var res2 map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&res2)
	id2 := res2["data"].(map[string]interface{})["id"].(string)

	s.NotEqual(id, id2)
}

func (s *TicketIntegrationTestSuite) TestIdempotency_ScopeExclusionDelete() {
	// Create ticket
	reqBody := map[string]interface{}{"customer_phone": "081234567890", "customer_name": "Delete test", "customer_gender": "Male", "brand": "A", "model": "B", "issue": "C"}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := s.app.Test(req)
	defer resp.Body.Close()
	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	id := res["data"].(map[string]interface{})["id"].(string)

	key := uuid.New().String()

	// First DELETE
	reqDel1, _ := http.NewRequest(http.MethodDelete, "/api/v1/tickets/"+id, nil)
	reqDel1.Header.Set("X-Idempotency-Key", key)
	respDel1, _ := s.app.Test(reqDel1)
	defer respDel1.Body.Close()
	s.Equal(http.StatusOK, respDel1.StatusCode)

	// Second DELETE
	reqDel2, _ := http.NewRequest(http.MethodDelete, "/api/v1/tickets/"+id, nil)
	reqDel2.Header.Set("X-Idempotency-Key", key)
	respDel2, _ := s.app.Test(reqDel2)
	defer respDel2.Body.Close()

	// Should not be cached, so it returns 404 since it was already deleted
	s.Equal(http.StatusNotFound, respDel2.StatusCode)
}

func (s *TicketIntegrationTestSuite) TestPostgresStorage_GetAndSet() {
	store := database.NewPostgresStorage(s.db)
	defer store.Close()

	// empty key and empty value are ignored
	err := store.Set("", nil, 0)
	s.Require().NoError(err)
	val, err := store.Get("")
	s.Require().NoError(err)
	s.Nil(val)

	// Get() returns nil, nil for missing or expired keys
	val, err = store.Get("missing-key")
	s.Require().NoError(err)
	s.Nil(val)

	// Get() returns a copied byte slice
	key := "test-key"
	expectedVal := []byte("cached-response")
	err = store.Set(key, expectedVal, 10*time.Minute)
	s.Require().NoError(err)

	val, err = store.Get(key)
	s.Require().NoError(err)
	s.Equal(expectedVal, val)

	// Modify the returned slice and verify stored one is not mutated
	val[0] = 'X'
	val2, err := store.Get(key)
	s.Require().NoError(err)
	s.Equal(expectedVal, val2)
}

func (s *TicketIntegrationTestSuite) TestPostgresStorage_ReserveRequest() {
	store := database.NewPostgresStorage(s.db)
	defer store.Close()
	key := "reserve-key"
	hash := "some-hash"

	// ReserveRequest() accepts the same key and same hash
	err := store.ReserveRequest(key, hash, 10*time.Minute)
	s.Require().NoError(err)

	err = store.ReserveRequest(key, hash, 10*time.Minute)
	s.Require().NoError(err)

	// ReserveRequest() returns a conflict for the same key and different hash
	err = store.ReserveRequest(key, "different-hash", 10*time.Minute)
	s.ErrorIs(err, database.ErrIdempotencyConflict)
}

func (s *TicketIntegrationTestSuite) TestPostgresStorage_DeleteExpired() {
	store := database.NewPostgresStorage(s.db)
	defer store.Close()
	key := "expired-key"

	// Set with negative duration (already expired)
	err := store.Set(key, []byte("expired-val"), -10*time.Minute)
	s.Require().NoError(err)

	// Get only returns active keys
	val, err := store.Get(key)
	s.Require().NoError(err)
	s.Nil(val)

	// DeleteExpired removes expired rows
	err = store.DeleteExpired()
	s.Require().NoError(err)

	// Verify it's gone from database entirely
	var count int
	err = s.db.Get(&count, "SELECT COUNT(*) FROM idempotency_keys WHERE key = $1", key)
	s.Require().NoError(err)
	s.Equal(0, count)
}

func (s *TicketIntegrationTestSuite) TestPostgresStorage_RequestPathDoesNotSynchronouslyDeleteExpiredKeys() {
	store := database.NewPostgresStorage(s.db, database.WithCleanupInterval(0))
	defer store.Close()

	_, err := s.db.Exec(`
		INSERT INTO idempotency_keys (key, request_hash, value, expires_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP - INTERVAL '1 minute')
	`, "expired-request-path-key", "expired-hash", []byte("expired-value"))
	s.Require().NoError(err)

	err = store.ReserveRequest("active-reserve-key", "active-hash", 10*time.Minute)
	s.Require().NoError(err)

	err = store.Set("active-cache-key", []byte("active-value"), 10*time.Minute)
	s.Require().NoError(err)

	var count int
	err = s.db.Get(&count, "SELECT COUNT(*) FROM idempotency_keys WHERE key = $1", "expired-request-path-key")
	s.Require().NoError(err)
	s.Equal(1, count)
}

func (s *TicketIntegrationTestSuite) TestPostgresStorage_CacheWriteFailure() {
	// Close a temporary database connection to simulate DB failure
	tempDB, err := sqlx.Open("postgres", "postgres://invalid:invalid@localhost:5432/invalid?sslmode=disable")
	s.Require().NoError(err)
	tempStore := database.NewPostgresStorage(tempDB)
	defer tempStore.Close()
	tempDB.Close() // closed db connection

	// Set on closed DB should return an error
	err = tempStore.Set("any-key", []byte("val"), 10*time.Minute)
	s.Error(err)
}
