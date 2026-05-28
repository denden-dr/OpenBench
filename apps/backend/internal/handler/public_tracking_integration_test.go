//go:build integration

package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
	"github.com/denden-dr/openbench/apps/backend/internal/handler"
	"github.com/denden-dr/openbench/apps/backend/internal/middleware"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type PublicTrackingIntegrationTestSuite struct {
	suite.Suite
	db  *sqlx.DB
	app *fiber.App
}

func (s *PublicTrackingIntegrationTestSuite) SetupSuite() {
	s.db = SetupTestDB()

	ticketRepo := repository.NewTicketRepository(s.db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	warrantyClaimRepo := repository.NewWarrantyClaimRepository(s.db)
	warrantyClaimService := service.NewWarrantyClaimService(warrantyClaimRepo, ticketRepo)
	warrantyClaimHandler := handler.NewWarrantyClaimHandler(warrantyClaimService)

	healthHandler := handler.NewHealthHandler(s.db.DB)

	s.app = fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			Disable:   false,
			MaxPublic: 1000,
			MaxAdmin:  1000,
		},
	}

	handler.RegisterRoutes(s.app, cfg, ticketHandler, warrantyClaimHandler, healthHandler)
}

func (s *PublicTrackingIntegrationTestSuite) SetupTest() {
	CleanTestDB(s.T(), s.db)
}

func TestPublicTrackingIntegrationSuite(t *testing.T) {
	suite.Run(t, new(PublicTrackingIntegrationTestSuite))
}

func (s *PublicTrackingIntegrationTestSuite) TestGetPublicTicket_SuccessUUID() {
	// 1. Insert a ticket directly into the DB
	id := uuid.New().String()
	_, err := s.db.Exec(`
		INSERT INTO tickets (id, customer_name, customer_phone, customer_gender, brand, model, issue, status, payment_status, warranty_days, entry_date)
		VALUES ($1, 'Prabowo Subianto', '+62812-3456-7890', 'Male', 'Apple', 'iPhone 13', 'Baterai Kembung', 'service_in', 'unpaid', 30, CURRENT_TIMESTAMP)
	`, id)
	s.Require().NoError(err)

	// 2. Request GET /api/v1/public/tickets/:id
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tickets/"+id, nil)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var res map[string]interface{}
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&res))
	s.Equal(float64(200), res["code"])
	s.Equal("Success", res["message"])

	data := res["data"].(map[string]interface{})
	s.Equal(id, data["id"])
	s.Equal("P****** S*******", data["customer_name_masked"])
	s.Equal("0812******90", data["customer_phone_masked"])
	s.Equal("Apple", data["brand"])
	s.Equal("iPhone 13", data["model"])
	s.Equal("Baterai Kembung", data["issue"])
	s.Equal("service_in", data["status"])
	s.Nil(data["customer_name"])  // Ensure original fields are NOT leaked
	s.Nil(data["customer_phone"]) // Ensure original fields are NOT leaked
}

func (s *PublicTrackingIntegrationTestSuite) TestGetPublicTicket_SuccessShortID() {
	// 1. Insert a ticket directly into the DB
	id := "abcdef12-3456-7890-abcd-ef1234567890"
	_, err := s.db.Exec(`
		INSERT INTO tickets (id, customer_name, customer_phone, customer_gender, brand, model, issue, status, payment_status, warranty_days, entry_date)
		VALUES ($1, 'Joko Widodo', '081299998888', 'Male', 'Samsung', 'Galaxy S22', 'Layar Retak', 'fixed', 'paid', 14, CURRENT_TIMESTAMP)
	`, id)
	s.Require().NoError(err)

	// 2. Request GET /api/v1/public/tickets/abcdef12
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tickets/abcdef12", nil)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *PublicTrackingIntegrationTestSuite) TestGetPublicTicket_NotFoundAndInvalidFormat() {
	// 1. Request with invalid format (neither UUID nor 8-char hex)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tickets/invalidformat", nil)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)

	// 2. Request with valid UUID format but not found
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tickets/"+uuid.New().String(), nil)
	resp2, err := s.app.Test(req2)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusNotFound, resp2.StatusCode)

	// 3. Request with valid 8-char hex format (which is now invalid format since only full UUID is allowed)
	req3, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tickets/abcdefaa", nil)
	resp3, err := s.app.Test(req3)
	s.Require().NoError(err)
	defer resp3.Body.Close()
	s.Equal(http.StatusBadRequest, resp3.StatusCode)
}

func (s *PublicTrackingIntegrationTestSuite) TestTrackPublic_Success() {
	// 1. Insert a ticket directly into the DB
	id := "12345678-3456-7890-abcd-ef1234567890"
	_, err := s.db.Exec(`
		INSERT INTO tickets (id, customer_name, customer_phone, customer_gender, brand, model, issue, status, payment_status, warranty_days, entry_date)
		VALUES ($1, 'Budi Gunawan', '+62 812-9876-5432', 'Male', 'Google', 'Pixel 7', 'Kamera Buram', 'service_in', 'unpaid', 30, CURRENT_TIMESTAMP)
	`, id)
	s.Require().NoError(err)

	// 2. Track using correct short ID and normalized Indonesian phone number
	// Input: "081298765432"
	reqBody := map[string]interface{}{
		"short_id": "12345678",
		"phone":    "081298765432",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/public/track", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var res map[string]interface{}
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&res))
	s.Equal(float64(200), res["code"])
	s.Equal("Success", res["message"])
	data := res["data"].(map[string]interface{})
	s.Equal(id, data["ticket_id"])

	// 3. Track using international Indonesian phone number format
	// Input: "+62 812-9876-5432"
	reqBody2 := map[string]interface{}{
		"short_id": "12345678",
		"phone":    "+62 812-9876-5432",
	}
	bodyBytes2, _ := json.Marshal(reqBody2)
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/public/track", bytes.NewReader(bodyBytes2))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := s.app.Test(req2)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)

	var res2 map[string]interface{}
	s.Require().NoError(json.NewDecoder(resp2.Body).Decode(&res2))
	s.Equal(float64(200), res2["code"])
	s.Equal("Success", res2["message"])
	data2 := res2["data"].(map[string]interface{})
	s.Equal(id, data2["ticket_id"])
}

func (s *PublicTrackingIntegrationTestSuite) TestTrackPublic_NotFound() {
	// 1. Insert a ticket directly into the DB
	id := "12345678-3456-7890-abcd-ef1234567890"
	_, err := s.db.Exec(`
		INSERT INTO tickets (id, customer_name, customer_phone, customer_gender, brand, model, issue, status, payment_status, warranty_days, entry_date)
		VALUES ($1, 'Budi Gunawan', '081298765432', 'Male', 'Google', 'Pixel 7', 'Kamera Buram', 'service_in', 'unpaid', 30, CURRENT_TIMESTAMP)
	`, id)
	s.Require().NoError(err)

	// 2. Try tracking with incorrect phone number
	reqBody := map[string]interface{}{
		"short_id": "12345678",
		"phone":    "081200000000",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/public/track", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *PublicTrackingIntegrationTestSuite) TestTrackPublic_RateLimiterTrigger() {
	// Test the rate limiter behavior on a custom app instance
	testApp := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Configure max of 2 requests per minute
	testLimiter := limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Minute,
	})
	testApp.Post("/api/v1/public/track", testLimiter, func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	reqBody := map[string]interface{}{
		"short_id": "12345678",
		"phone":    "081298765432",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Request 1
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/public/track", bytes.NewReader(bodyBytes))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := testApp.Test(req1)
	s.Require().NoError(err)
	resp1.Body.Close()
	s.Equal(http.StatusOK, resp1.StatusCode)

	// Request 2
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/public/track", bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := testApp.Test(req2)
	s.Require().NoError(err)
	resp2.Body.Close()
	s.Equal(http.StatusOK, resp2.StatusCode)

	// Request 3 (Limit Reached)
	req3, _ := http.NewRequest(http.MethodPost, "/api/v1/public/track", bytes.NewReader(bodyBytes))
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := testApp.Test(req3)
	s.Require().NoError(err)
	defer resp3.Body.Close()
	s.Equal(http.StatusTooManyRequests, resp3.StatusCode)
}

func (s *PublicTrackingIntegrationTestSuite) TestTrackPublic_NoPhoneNumber() {
	// 1. Insert a ticket directly into the DB with empty customer_phone
	id := "87654321-3456-7890-abcd-ef1234567890"
	_, err := s.db.Exec(`
		INSERT INTO tickets (id, customer_name, customer_phone, customer_gender, brand, model, issue, status, payment_status, warranty_days, entry_date)
		VALUES ($1, 'Budi Tanpa HP', '', 'Male', 'Google', 'Pixel 7', 'Kamera Buram', 'service_in', 'unpaid', 30, CURRENT_TIMESTAMP)
	`, id)
	s.Require().NoError(err)

	// 2. Try tracking it using the short ID and any phone number - should return 404 (Not Found)
	reqBody := map[string]interface{}{
		"short_id": "87654321",
		"phone":    "081298765432",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/public/track", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}
