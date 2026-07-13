//go:build integration

package warranty_test

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/testutils"
	"github.com/denden-dr/OpenBench/apps/backend/internal/ticket"
	"github.com/denden-dr/OpenBench/apps/backend/internal/warranty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarrantyRepository_Integration(t *testing.T) {
	ctx := context.Background()

	// Spin up PostgreSQL test container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	ticketCmdRepo := ticket.NewCommandRepository(db)
	cmdRepo := warranty.NewCommandRepository(db)
	queryRepo := warranty.NewQueryRepository(db)

	// Clean tables
	err = testutils.CleanTable(db, "claims")
	require.NoError(t, err)
	err = testutils.CleanTable(db, "warranties")
	require.NoError(t, err)
	err = testutils.CleanTable(db, "service_tickets")
	require.NoError(t, err)

	// Seed base ticket to avoid foreign key violation
	ticketID := uuid.New().String()
	baseTicket := &models.ServiceTicket{
		ID:               ticketID,
		TicketNumber:     "TKT-WARR-001",
		Status:           models.StatusCompleted,
		CustomerName:     "Bob Vance",
		CustomerPhone:    "089988887777",
		DeviceBrand:      "Apple",
		DeviceModel:      "iPhone 14",
		IssueDescription: "Battery Health",
		Cost:             800000,
		WarrantyDays:     90,
	}
	err = ticketCmdRepo.Create(ctx, baseTicket)
	require.NoError(t, err)

	var warrID string

	t.Run("Create and Retrieve Warranty", func(t *testing.T) {
		warrID = uuid.New().String()
		notes := "90-day standard warranty"

		w := &models.Warranty{
			ID:        warrID,
			TicketID:  ticketID,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(90 * 24 * time.Hour),
			Status:    models.WarrantyStatusActive,
			Notes:     &notes,
		}

		err := cmdRepo.CreateWarranty(ctx, w)
		require.NoError(t, err)
		assert.NotEmpty(t, w.CreatedAt)
		assert.NotEmpty(t, w.UpdatedAt)

		// FindByID
		found, err := queryRepo.FindWarrantyByID(ctx, warrID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, warrID, found.ID)
		assert.Equal(t, ticketID, found.TicketID)
		assert.Equal(t, models.WarrantyStatusActive, found.Status)
		require.NotNil(t, found.Notes)
		assert.Equal(t, "90-day standard warranty", *found.Notes)

		// FindByTicketID
		foundByTicket, err := queryRepo.FindWarrantyByTicketID(ctx, ticketID)
		require.NoError(t, err)
		require.NotNil(t, foundByTicket)
		assert.Equal(t, warrID, foundByTicket.ID)
	})

	t.Run("Update Warranty Status", func(t *testing.T) {
		newNotes := "Voided due to customer damage"
		err := cmdRepo.UpdateWarrantyStatus(ctx, warrID, models.WarrantyStatusVoid, &newNotes)
		require.NoError(t, err)

		found, err := queryRepo.FindWarrantyByID(ctx, warrID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, models.WarrantyStatusVoid, found.Status)
		require.NotNil(t, found.Notes)
		assert.Equal(t, "Voided due to customer damage", *found.Notes)
	})

	t.Run("Claim Lifecycle Operations", func(t *testing.T) {
		claimID := uuid.New().String()
		notes := "Customer reported speaker is still muffled"

		c := &models.Claim{
			ID:               claimID,
			ClaimNumber:      "CLM-20260713-0001",
			WarrantyID:       warrID,
			Status:           models.StatusReceived,
			EvaluationStatus: models.ClaimEvaluationPending,
			IssueDescription: "Speaker issue persists",
			Notes:            &notes,
		}

		// 1. CreateClaim
		err := cmdRepo.CreateClaim(ctx, c)
		require.NoError(t, err)
		assert.NotEmpty(t, c.CreatedAt)
		assert.NotEmpty(t, c.UpdatedAt)

		// 2. FindClaimByID
		found, err := queryRepo.FindClaimByID(ctx, claimID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, claimID, found.ID)
		assert.Equal(t, "CLM-20260713-0001", found.ClaimNumber)
		assert.Equal(t, models.StatusReceived, found.Status)
		assert.Equal(t, models.ClaimEvaluationPending, found.EvaluationStatus)
		assert.Equal(t, "Speaker issue persists", found.IssueDescription)
		require.NotNil(t, found.Notes)
		assert.Equal(t, "Customer reported speaker is still muffled", *found.Notes)

		// 3. UpdateClaim
		repairAction := "Cleaned internal grill & replaced speaker unit"
		c.Status = models.StatusRepairing
		c.RepairAction = &repairAction

		err = cmdRepo.UpdateClaim(ctx, c)
		require.NoError(t, err)

		found, err = queryRepo.FindClaimByID(ctx, claimID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, models.StatusRepairing, found.Status)
		require.NotNil(t, found.RepairAction)
		assert.Equal(t, "Cleaned internal grill & replaced speaker unit", *found.RepairAction)

		// 4. UpdateClaimEvaluation
		evalNotes := "Authorized repair under warranty terms"
		err = cmdRepo.UpdateClaimEvaluation(ctx, claimID, models.StatusCompleted, models.ClaimEvaluationAccepted, &evalNotes)
		require.NoError(t, err)

		found, err = queryRepo.FindClaimByID(ctx, claimID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, models.StatusCompleted, found.Status)
		assert.Equal(t, models.ClaimEvaluationAccepted, found.EvaluationStatus)
		require.NotNil(t, found.EvaluationNotes)
		assert.Equal(t, "Authorized repair under warranty terms", *found.EvaluationNotes)
	})

	t.Run("FindAllClaims Filters", func(t *testing.T) {
		// Clean claims first
		err = testutils.CleanTable(db, "claims")
		require.NoError(t, err)

		claimsToSeed := []*models.Claim{
			{
				ID:               uuid.New().String(),
				ClaimNumber:      "CLM-001",
				WarrantyID:       warrID,
				Status:           models.StatusReceived,
				EvaluationStatus: models.ClaimEvaluationPending,
				IssueDescription: "Broken screen",
			},
			{
				ID:               uuid.New().String(),
				ClaimNumber:      "CLM-002",
				WarrantyID:       warrID,
				Status:           models.StatusRepairing,
				EvaluationStatus: models.ClaimEvaluationAccepted,
				IssueDescription: "Speaker crackle",
			},
		}

		for _, cl := range claimsToSeed {
			err = cmdRepo.CreateClaim(ctx, cl)
			require.NoError(t, err)
		}

		// 1. Filter by Status
		list, total, err := queryRepo.FindAllClaims(ctx, string(models.StatusReceived), "", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		require.Len(t, list, 1)
		assert.Equal(t, "CLM-001", list[0].ClaimNumber)

		// 2. Search query (number/description)
		list, total, err = queryRepo.FindAllClaims(ctx, "", "crackle", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		require.Len(t, list, 1)
		assert.Equal(t, "CLM-002", list[0].ClaimNumber)
	})
}
