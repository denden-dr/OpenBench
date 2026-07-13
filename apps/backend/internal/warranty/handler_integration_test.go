//go:build integration

package warranty_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/testutils"
	"github.com/denden-dr/OpenBench/apps/backend/internal/ticket"
	"github.com/denden-dr/OpenBench/apps/backend/internal/warranty"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarrantyHandler_Integration(t *testing.T) {
	ctx := context.Background()

	// Spin up PostgreSQL test container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	ticketCmdRepo := ticket.NewCommandRepository(db)
	cmdRepo := warranty.NewCommandRepository(db)
	queryRepo := warranty.NewQueryRepository(db)
	txManager := database.NewTxManager(db)

	err = testutils.CleanTable(db, "claims")
	require.NoError(t, err)
	err = testutils.CleanTable(db, "warranties")
	require.NoError(t, err)
	err = testutils.CleanTable(db, "service_tickets")
	require.NoError(t, err)

	// Seed base ticket
	ticketID := uuid.New().String()
	baseTicket := &models.ServiceTicket{
		ID:               ticketID,
		TicketNumber:     "TKT-WARR-999",
		Status:           models.StatusCompleted,
		CustomerName:     "Michael Scott",
		CustomerPhone:    "089988889999",
		DeviceBrand:      "Apple",
		DeviceModel:      "iPhone 15 Pro",
		IssueDescription: "Swollen Battery",
		Cost:             1200000,
		WarrantyDays:     90,
	}
	err = ticketCmdRepo.Create(ctx, baseTicket)
	require.NoError(t, err)

	// Setup service & handler
	svc := warranty.NewService(queryRepo, cmdRepo, txManager)
	h := warranty.NewHandler(svc)

	// Fiber router
	app := fiber.New()
	app.Get("/warranties/by-ticket/:ticket_id", h.GetWarrantyByTicketID)
	app.Patch("/warranties/:warranty_id/status", h.UpdateWarrantyStatus)
	app.Post("/claims", h.CreateClaim)
	app.Get("/claims", h.GetClaims)
	app.Get("/claims/:claim_id", h.GetClaimByID)
	app.Patch("/claims/:claim_id/status", h.UpdateClaimStatus)
	app.Put("/claims/:claim_id", h.UpdateClaim)
	app.Post("/claims/:claim_id/evaluate", h.EvaluateClaim)

	// 1. Trigger warranty creation by creating warranty manually via service (since we do not test ticket completion flow in this file, or we can just call svc.CreateWarranty)
	w, err := svc.CreateWarranty(ctx, ticketID, 90)
	require.NoError(t, err)
	assert.NotEmpty(t, w.ID)

	var createdClaimID string

	t.Run("Get Warranty By Ticket ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/warranties/by-ticket/%s", ticketID), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data warranty.WarrantyResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, w.ID, respData.Data.ID)
		assert.Equal(t, ticketID, respData.Data.TicketID)
		assert.Equal(t, models.WarrantyStatusActive, respData.Data.Status)
	})

	t.Run("Update Warranty Status", func(t *testing.T) {
		body := warranty.UpdateWarrantyStatusRequest{
			Status: models.WarrantyStatusVoid,
			Notes:  "Voided due to unauthorized repair",
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PATCH", fmt.Sprintf("/warranties/%s/status", w.ID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data warranty.WarrantyResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, models.WarrantyStatusVoid, respData.Data.Status)
		require.NotNil(t, respData.Data.Notes)
		assert.Equal(t, "Voided due to unauthorized repair", *respData.Data.Notes)

		// Activate warranty again for claims testing
		bodyActive := warranty.UpdateWarrantyStatusRequest{
			Status: models.WarrantyStatusActive,
			Notes:  "Re-activated after evaluation",
		}
		bodyBytesActive, _ := json.Marshal(bodyActive)
		reqActive, _ := http.NewRequest("PATCH", fmt.Sprintf("/warranties/%s/status", w.ID), bytes.NewBuffer(bodyBytesActive))
		reqActive.Header.Set("Content-Type", "application/json")
		respActive, err := app.Test(reqActive)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respActive.StatusCode)
	})

	t.Run("Create Claim", func(t *testing.T) {
		body := warranty.CreateClaimRequest{
			WarrantyID:       w.ID,
			IssueDescription: "Speaker still doesn't work after screen repair",
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/claims", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respData struct {
			Data warranty.ClaimResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.NotEmpty(t, respData.Data.ClaimID)
		assert.NotEmpty(t, respData.Data.ClaimNumber)
		assert.Equal(t, w.ID, respData.Data.WarrantyID)
		assert.Equal(t, models.StatusReceived, respData.Data.Status)
		assert.Equal(t, models.ClaimEvaluationPending, respData.Data.EvaluationStatus)

		createdClaimID = respData.Data.ClaimID
	})

	t.Run("Get Claim By ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/claims/%s", createdClaimID), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data warranty.ClaimResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, createdClaimID, respData.Data.ClaimID)
		assert.Equal(t, "Speaker still doesn't work after screen repair", respData.Data.IssueDescription)
	})

	t.Run("Update Claim Details", func(t *testing.T) {
		body := warranty.UpdateClaimRequest{
			IssueDescription: "Speaker is completely dead, headphone jack works",
			RepairAction:     "Replaced Speaker Module",
			Notes:            "Tested with multitone test",
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/claims/%s", createdClaimID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data warranty.ClaimResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, "Speaker is completely dead, headphone jack works", respData.Data.IssueDescription)
		require.NotNil(t, respData.Data.RepairAction)
		assert.Equal(t, "Replaced Speaker Module", *respData.Data.RepairAction)
	})

	t.Run("Update Claim Status", func(t *testing.T) {
		body := warranty.ChangeClaimStatusRequest{
			Status: models.StatusRepairing,
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PATCH", fmt.Sprintf("/claims/%s/status", createdClaimID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data warranty.ClaimStatusResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, models.StatusRepairing, respData.Data.Status)
	})

	t.Run("Evaluate Claim", func(t *testing.T) {
		body := warranty.EvaluateClaimRequest{
			Status: models.ClaimEvaluationAccepted,
			Notes:  "Approved for zero cost repair",
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", fmt.Sprintf("/claims/%s/evaluate", createdClaimID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data warranty.ClaimResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, models.ClaimEvaluationAccepted, respData.Data.EvaluationStatus)
		assert.Equal(t, models.StatusRepairing, respData.Data.Status)
		require.NotNil(t, respData.Data.EvaluationNotes)
		assert.Equal(t, "Approved for zero cost repair", *respData.Data.EvaluationNotes)
	})

	t.Run("Get Claims List", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/claims?status=REPAIRING&search=Speaker&limit=10&offset=0", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData warranty.ClaimListWrapper
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, 1, respData.Meta.TotalData)
		require.Len(t, respData.Data, 1)
		assert.Equal(t, createdClaimID, respData.Data[0].ClaimID)
	})
}
