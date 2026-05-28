//go:build integration

package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/config"
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

type WarrantyClaimIntegrationTestSuite struct {
	suite.Suite
	db               *sqlx.DB
	app              *fiber.App
	idempotencyStore *database.PostgresStorage
}

func (s *WarrantyClaimIntegrationTestSuite) SetupSuite() {
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

	s.idempotencyStore = database.NewPostgresStorage(s.db)
	s.app.Use(middleware.ScopeIdempotencyKey(s.idempotencyStore))
	s.app.Use(middleware.NewIdempotency(s.idempotencyStore))

	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			Disable:   false,
			MaxPublic: 1000,
			MaxAdmin:  1000,
		},
	}

	handler.RegisterRoutes(s.app, cfg, ticketHandler, warrantyClaimHandler, healthHandler)
}

func (s *WarrantyClaimIntegrationTestSuite) SetupTest() {
	CleanTestDB(s.T(), s.db)
}

func (s *WarrantyClaimIntegrationTestSuite) TearDownSuite() {
	if s.idempotencyStore != nil {
		_ = s.idempotencyStore.Close()
	}
}

func TestWarrantyClaimIntegrationSuite(t *testing.T) {
	suite.Run(t, new(WarrantyClaimIntegrationTestSuite))
}

func (s *WarrantyClaimIntegrationTestSuite) TestCreateWarrantyClaim_Success() {
	// 1. Create a ticket
	ticketID := s.createPickedUpTicket()

	// 2. Post a warranty claim for this ticket
	reqBody := map[string]interface{}{
		"ticket_id":              ticketID,
		"issue":                  "Layar flicker setelah 2 hari",
		"additional_description": "Flicker parah di bagian bawah layar",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/warranty-claims", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	s.Equal(float64(201), res["code"])

	data := res["data"].(map[string]interface{})
	s.NotEmpty(data["id"])
	s.Equal(ticketID, data["ticket_id"])
	s.Equal("waiting_inspection", data["status"])
	s.Equal("Layar flicker setelah 2 hari", data["issue"])
	s.Equal("Flicker parah di bagian bawah layar", data["additional_description"])
}

func (s *WarrantyClaimIntegrationTestSuite) TestCreateWarrantyClaim_NotPickedUp() {
	// 1. Create a ticket (status: service_in)
	ticketID := s.createTicket("service_in")

	// 2. Post a warranty claim for it
	reqBody := map[string]interface{}{
		"ticket_id": ticketID,
		"issue":     "Layar flicker",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/warranty-claims", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Expect 400 Bad Request because it has not been picked up
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	s.Equal(float64(400), res["status"])
	s.Equal("ticket has not been picked up by customer", res["detail"])
}

func (s *WarrantyClaimIntegrationTestSuite) TestCreateWarrantyClaim_Expired() {
	// 1. Create a ticket
	ticketID := s.createPickedUpTicket()

	// 2. Artificially make exit_date far in the past in DB
	pastDate := time.Now().UTC().AddDate(0, 0, -100)
	_, err := s.db.Exec("UPDATE tickets SET exit_date = $1 WHERE id = $2", pastDate, ticketID)
	s.Require().NoError(err)

	// 3. Post a warranty claim for it
	reqBody := map[string]interface{}{
		"ticket_id": ticketID,
		"issue":     "Baterai drop lagi",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/warranty-claims", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Expect 400 Bad Request because warranty period has expired
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	s.Equal(float64(400), res["status"])
	s.Equal("warranty period has expired", res["detail"])
}

func (s *WarrantyClaimIntegrationTestSuite) TestApproveClaim_Success() {
	// 1. Setup a valid claim waiting inspection
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Layar bergaris hijau")
	oldUpdatedAt := s.backdateClaimTimestamps(claimID)

	// 2. Approve the claim
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/approve", claimID), nil)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	s.Equal(float64(200), res["code"])

	// Check claim updates
	data := res["data"].(map[string]interface{})
	claimData := data["claim"].(map[string]interface{})
	s.Equal(claimID, claimData["id"])
	s.Equal("approved", claimData["status"])
	s.NotEmpty(claimData["claim_ticket_id"])
	s.NotEmpty(claimData["inspected_at"])
	responseUpdatedAt := s.parseJSONTime(claimData["updated_at"])
	s.True(responseUpdatedAt.After(oldUpdatedAt), "claim updated_at in response should reflect the DB update")

	// Check spawned ticket details
	ticketData := data["ticket"].(map[string]interface{})
	claimTicketID := claimData["claim_ticket_id"].(string)
	s.Equal(claimTicketID, ticketData["id"])
	s.Equal("on_process", ticketData["status"])
	s.Equal("paid", ticketData["payment_status"])
	s.Equal(true, ticketData["is_warranty"])
	s.Equal(ticketID, ticketData["parent_ticket_id"])
	s.Equal("[Klaim Garansi] Layar bergaris hijau", ticketData["issue"])

	// Double check that price is Rp 0 (serialized as string or float depending on config)
	priceVal := fmt.Sprintf("%v", ticketData["price"])
	s.Contains([]string{"0", "0.00", "0.0"}, priceVal)

	dbUpdatedAt := s.getClaimUpdatedAt(claimID)
	s.True(responseUpdatedAt.Equal(dbUpdatedAt), "claim updated_at in response should match the database")
}

func (s *WarrantyClaimIntegrationTestSuite) TestVoidClaim_Success() {
	// 1. Setup a valid claim waiting inspection
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Kamera retak")
	oldUpdatedAt := s.backdateClaimTimestamps(claimID)

	// 2. Void the claim
	reqBody := map[string]interface{}{
		"void_reason": "Kerusakan fisik akibat jatuh sendiri oleh user",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	s.Equal(float64(200), res["code"])

	// Check claim updates
	data := res["data"].(map[string]interface{})
	claimData := data["claim"].(map[string]interface{})
	s.Equal(claimID, claimData["id"])
	s.Equal("void", claimData["status"])
	s.Equal("Kerusakan fisik akibat jatuh sendiri oleh user", claimData["void_reason"])
	s.NotEmpty(claimData["claim_ticket_id"])
	responseUpdatedAt := s.parseJSONTime(claimData["updated_at"])
	s.True(responseUpdatedAt.After(oldUpdatedAt), "claim updated_at in response should reflect the DB update")

	// Check spawned ticket details (should be cancelled)
	ticketData := data["ticket"].(map[string]interface{})
	s.Equal("cancelled", ticketData["status"])
	s.Equal(true, ticketData["is_warranty"])
	s.Equal(ticketID, ticketData["parent_ticket_id"])
	s.Equal("[Klaim Ditolak] Kamera retak", ticketData["issue"])

	// Void/cancelled warranty tickets must have warranty_days = 0
	warrantyVal := fmt.Sprintf("%v", ticketData["warranty_days"])
	s.Equal("0", warrantyVal)

	dbUpdatedAt := s.getClaimUpdatedAt(claimID)
	s.True(responseUpdatedAt.Equal(dbUpdatedAt), "claim updated_at in response should match the database")
}

func (s *WarrantyClaimIntegrationTestSuite) TestCreateWarrantyClaim_Duplicate() {
	ticketID := s.createPickedUpTicket()

	// First claim succeeds
	reqBody := map[string]interface{}{
		"ticket_id": ticketID,
		"issue":     "Layar flicker",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/warranty-claims", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	// Second claim for the same ticket is rejected
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/warranty-claims", bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err2 := s.app.Test(req2)
	s.Require().NoError(err2)
	defer resp2.Body.Close()
	s.Require().Equal(http.StatusConflict, resp2.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&res)
	s.Equal(float64(409), res["status"])
	s.Equal("ticket already has an open warranty claim", res["detail"])
}

func (s *WarrantyClaimIntegrationTestSuite) TestApproveClaim_AlreadyDecided() {
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Layar bergaris")

	// First approve succeeds
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/approve", claimID), nil)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	// Second approve fails - claim already decided
	req2, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/approve", claimID), nil)
	resp2, err2 := s.app.Test(req2)
	s.Require().NoError(err2)
	s.Require().Equal(http.StatusBadRequest, resp2.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&res)
	s.Equal(float64(400), res["status"])
	s.Equal("warranty claim has already been approved or voided", res["detail"])
}

func (s *WarrantyClaimIntegrationTestSuite) TestVoidClaim_AlreadyDecided() {
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Kamera retak")

	// First void succeeds
	reqBody := map[string]interface{}{"void_reason": "Kerusakan fisik"}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	// Second void fails - claim already decided
	req2, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err2 := s.app.Test(req2)
	s.Require().NoError(err2)
	s.Require().Equal(http.StatusBadRequest, resp2.StatusCode)

	var res map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&res)
	s.Equal(float64(400), res["status"])
	s.Equal("warranty claim has already been approved or voided", res["detail"])
}

func (s *WarrantyClaimIntegrationTestSuite) TestConcurrentApprove_RaceCondition() {
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Layar bergaris hijau")

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	errMessages := make([]string, 0)

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/approve", claimID), nil)
			resp, err := s.app.Test(req)
			if err != nil {
				mu.Lock()
				errMessages = append(errMessages, err.Error())
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			var res map[string]interface{}
			_ = json.NewDecoder(resp.Body).Decode(&res)

			mu.Lock()
			if resp.StatusCode == http.StatusOK {
				successCount++
			} else {
				errMsg, _ := res["detail"].(string)
				errMessages = append(errMessages, errMsg)
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	// Only one approve should succeed; the rest must fail with already-decided error
	s.Equal(1, successCount, "exactly one approve should succeed")
	for _, msg := range errMessages {
		s.Equal("warranty claim has already been approved or voided", msg)
	}

	// Verify final state in DB
	var status string
	err := s.db.QueryRow("SELECT status FROM warranty_claims WHERE id = $1", claimID).Scan(&status)
	s.Require().NoError(err)
	s.Equal("approved", status)
}

func (s *WarrantyClaimIntegrationTestSuite) TestConcurrentVoid_RaceCondition() {
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Kamera retak")

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	errMessages := make([]string, 0)
	reqBody := map[string]interface{}{"void_reason": "Kerusakan fisik akibat jatuh"}
	bodyBytes, _ := json.Marshal(reqBody)

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp, err := s.app.Test(req)
			if err != nil {
				mu.Lock()
				errMessages = append(errMessages, err.Error())
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			var res map[string]interface{}
			_ = json.NewDecoder(resp.Body).Decode(&res)

			mu.Lock()
			if resp.StatusCode == http.StatusOK {
				successCount++
			} else {
				errMsg, _ := res["detail"].(string)
				errMessages = append(errMessages, errMsg)
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	s.Equal(1, successCount, "exactly one void should succeed")
	for _, msg := range errMessages {
		s.Equal("warranty claim has already been approved or voided", msg)
	}

	// Verify final state in DB
	var status string
	err := s.db.QueryRow("SELECT status FROM warranty_claims WHERE id = $1", claimID).Scan(&status)
	s.Require().NoError(err)
	s.Equal("void", status)
}

func (s *WarrantyClaimIntegrationTestSuite) TestConcurrentApproveAndVoid_RaceCondition() {
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Baterai kembung")

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	errMessages := make([]string, 0)
	voidBody, _ := json.Marshal(map[string]interface{}{"void_reason": "Bukan garansi"})

	for i := range 10 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			var req *http.Request
			if idx%2 == 0 {
				req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/approve", claimID), nil)
			} else {
				req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(voidBody))
				req.Header.Set("Content-Type", "application/json")
			}
			resp, err := s.app.Test(req)
			if err != nil {
				mu.Lock()
				errMessages = append(errMessages, err.Error())
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			var res map[string]interface{}
			_ = json.NewDecoder(resp.Body).Decode(&res)

			mu.Lock()
			if resp.StatusCode == http.StatusOK {
				successCount++
			} else {
				errMsg, _ := res["detail"].(string)
				errMessages = append(errMessages, errMsg)
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	// Exactly one operation (approve or void) should succeed
	s.Equal(1, successCount, "exactly one operation should succeed")
	for _, msg := range errMessages {
		s.Equal("warranty claim has already been approved or voided", msg)
	}

	// Verify final state is either approved or void (not waiting_inspection)
	var status string
	err := s.db.QueryRow("SELECT status FROM warranty_claims WHERE id = $1", claimID).Scan(&status)
	s.Require().NoError(err)
	s.NotEqual("waiting_inspection", status, "claim must not remain waiting_inspection after concurrent decision")
}

// Helpers
func (s *WarrantyClaimIntegrationTestSuite) createTicket(status string) string {
	var id string
	query := `
		INSERT INTO tickets (customer_name, customer_gender, brand, model, issue, status, price, warranty_days)
		VALUES ('Budi', 'Male', 'Apple', 'iPhone 13', 'Layar Rusak', $1, 1500000, 30)
		RETURNING id
	`
	err := s.db.QueryRow(query, status).Scan(&id)
	s.Require().NoError(err)
	return id
}

func (s *WarrantyClaimIntegrationTestSuite) createPickedUpTicket() string {
	var id string
	now := time.Now().UTC()
	query := `
		INSERT INTO tickets (customer_name, customer_gender, brand, model, issue, status, payment_status, price, warranty_days, exit_date)
		VALUES ('Andi', 'Male', 'Samsung', 'Galaxy S21', 'Baterai Kembung', 'picked_up', 'paid', 500000, 30, $1)
		RETURNING id
	`
	err := s.db.QueryRow(query, now).Scan(&id)
	s.Require().NoError(err)
	return id
}

func (s *WarrantyClaimIntegrationTestSuite) createClaim(ticketID string, issue string) string {
	var id string
	query := `
		INSERT INTO warranty_claims (ticket_id, issue, status)
		VALUES ($1, $2, 'waiting_inspection')
		RETURNING id
	`
	err := s.db.QueryRow(query, ticketID, issue).Scan(&id)
	s.Require().NoError(err)
	return id
}

func (s *WarrantyClaimIntegrationTestSuite) backdateClaimTimestamps(claimID string) time.Time {
	oldUpdatedAt := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	_, err := s.db.Exec(
		"UPDATE warranty_claims SET created_at = $1, updated_at = $1 WHERE id = $2",
		oldUpdatedAt,
		claimID,
	)
	s.Require().NoError(err)
	return oldUpdatedAt
}

func (s *WarrantyClaimIntegrationTestSuite) getClaimUpdatedAt(claimID string) time.Time {
	var updatedAt time.Time
	err := s.db.QueryRow("SELECT updated_at FROM warranty_claims WHERE id = $1", claimID).Scan(&updatedAt)
	s.Require().NoError(err)
	return updatedAt.UTC()
}

func (s *WarrantyClaimIntegrationTestSuite) parseJSONTime(value interface{}) time.Time {
	raw, ok := value.(string)
	s.Require().True(ok)

	parsed, err := time.Parse(time.RFC3339Nano, raw)
	s.Require().NoError(err)
	return parsed.UTC()
}

func (s *WarrantyClaimIntegrationTestSuite) TestIdempotency_CreateWarrantyClaim() {
	ticketID := s.createPickedUpTicket()
	key := uuid.New().String()

	reqBody := map[string]interface{}{
		"ticket_id":              ticketID,
		"issue":                  "Layar flicker setelah 2 hari",
		"additional_description": "Flicker parah di bagian bawah layar",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// First Request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/warranty-claims", bytes.NewReader(bodyBytes))
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
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/warranty-claims", bytes.NewReader(bodyBytes))
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

	// Verify only one claim was created in the DB
	var count int
	err = s.db.Get(&count, "SELECT COUNT(*) FROM warranty_claims WHERE ticket_id = $1", ticketID)
	s.Require().NoError(err)
	s.Equal(1, count)
}

func (s *WarrantyClaimIntegrationTestSuite) TestIdempotency_ApproveWarrantyClaim() {
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Layar flicker")
	key := uuid.New().String()

	// First Request
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/approve", claimID), nil)
	req.Header.Set("X-Idempotency-Key", key)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var res1 map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res1)
	s.Equal(float64(200), res1["code"])
	data1 := res1["data"].(map[string]interface{})
	claimData1 := data1["claim"].(map[string]interface{})
	s.Equal("approved", claimData1["status"])

	// Second Request (Duplicate)
	req2, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/approve", claimID), nil)
	req2.Header.Set("X-Idempotency-Key", key)
	resp2, err := s.app.Test(req2)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Require().Equal(http.StatusOK, resp2.StatusCode)

	var res2 map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&res2)
	s.Equal(float64(200), res2["code"])
	data2 := res2["data"].(map[string]interface{})
	claimData2 := data2["claim"].(map[string]interface{})
	s.Equal("approved", claimData2["status"])
}

func (s *WarrantyClaimIntegrationTestSuite) TestIdempotency_VoidWarrantyClaim() {
	ticketID := s.createPickedUpTicket()
	claimID := s.createClaim(ticketID, "Layar flicker")
	key := uuid.New().String()

	reqBody := map[string]interface{}{
		"void_reason": "Kerusakan fisik akibat jatuh sendiri oleh user",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// First Request
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Idempotency-Key", key)
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var res1 map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res1)
	s.Equal(float64(200), res1["code"])
	data1 := res1["data"].(map[string]interface{})
	claimData1 := data1["claim"].(map[string]interface{})
	s.Equal("void", claimData1["status"])
	s.Equal("Kerusakan fisik akibat jatuh sendiri oleh user", claimData1["void_reason"])

	// Second Request (Duplicate)
	req2, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Idempotency-Key", key)
	resp2, err := s.app.Test(req2)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Require().Equal(http.StatusOK, resp2.StatusCode)

	var res2 map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&res2)
	s.Equal(float64(200), res2["code"])
	data2 := res2["data"].(map[string]interface{})
	claimData2 := data2["claim"].(map[string]interface{})
	s.Equal("void", claimData2["status"])
	s.Equal("Kerusakan fisik akibat jatuh sendiri oleh user", claimData2["void_reason"])

	// Verify that sending a different body with the same key returns 409 Conflict
	differentBody := map[string]interface{}{
		"void_reason": "Berbeda alasan",
	}
	differentBodyBytes, _ := json.Marshal(differentBody)
	req3, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/warranty-claims/%s/void", claimID), bytes.NewReader(differentBodyBytes))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("X-Idempotency-Key", key)
	resp3, err := s.app.Test(req3)
	s.Require().NoError(err)
	defer resp3.Body.Close()
	s.Require().Equal(http.StatusConflict, resp3.StatusCode)
}
