//go:build integration

package ticket_test

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/testutils"
	"github.com/denden-dr/OpenBench/apps/backend/internal/ticket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTicketRepository_Integration(t *testing.T) {
	ctx := context.Background()

	// Spin up PostgreSQL test container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	cmdRepo := ticket.NewCommandRepository(db)
	queryRepo := ticket.NewQueryRepository(db)

	t.Run("CRUD Operations", func(t *testing.T) {
		err := testutils.CleanTable(db, "service_tickets")
		require.NoError(t, err)

		repairAction := "Replaced Screen"
		notes := "Minor scratches on body"

		ticketObj := &models.ServiceTicket{
			ID:               uuid.New().String(),
			TicketNumber:     "TKT-20260713-0001",
			Status:           models.StatusReceived,
			CustomerName:     "Alice Smith",
			CustomerPhone:    "081234567890",
			DeviceBrand:      "Apple",
			DeviceModel:      "iPhone 13",
			DevicePasscode:   "123456",
			IssueDescription: "Cracked screen",
			RepairAction:     &repairAction,
			Cost:             1500000,
			WarrantyDays:     30,
			Notes:            &notes,
		}

		// 1. Create
		err = cmdRepo.Create(ctx, ticketObj)
		require.NoError(t, err)
		assert.NotEmpty(t, ticketObj.CreatedAt)
		assert.NotEmpty(t, ticketObj.UpdatedAt)

		// 2. FindByID
		found, err := queryRepo.FindByID(ctx, ticketObj.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, ticketObj.ID, found.ID)
		assert.Equal(t, ticketObj.TicketNumber, found.TicketNumber)
		assert.Equal(t, models.StatusReceived, found.Status)
		assert.Equal(t, "Alice Smith", found.CustomerName)
		assert.Equal(t, "081234567890", found.CustomerPhone)
		assert.Equal(t, "Apple", found.DeviceBrand)
		assert.Equal(t, "iPhone 13", found.DeviceModel)
		assert.Equal(t, "123456", found.DevicePasscode)
		assert.Equal(t, "Cracked screen", found.IssueDescription)
		require.NotNil(t, found.RepairAction)
		assert.Equal(t, "Replaced Screen", *found.RepairAction)
		assert.Equal(t, int64(1500000), found.Cost)
		assert.Equal(t, 30, found.WarrantyDays)
		require.NotNil(t, found.Notes)
		assert.Equal(t, "Minor scratches on body", *found.Notes)

		// 3. Update
		newRepairAction := "Replaced Screen and Battery"
		newNotes := "Battery was swollen"
		ticketObj.Status = models.StatusRepairing
		ticketObj.RepairAction = &newRepairAction
		ticketObj.Notes = &newNotes

		err = cmdRepo.Update(ctx, ticketObj)
		require.NoError(t, err)

		found, err = queryRepo.FindByID(ctx, ticketObj.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, models.StatusRepairing, found.Status)
		require.NotNil(t, found.RepairAction)
		assert.Equal(t, "Replaced Screen and Battery", *found.RepairAction)
		require.NotNil(t, found.Notes)
		assert.Equal(t, "Battery was swollen", *found.Notes)
	})

	t.Run("FindAll and Search", func(t *testing.T) {
		err := testutils.CleanTable(db, "service_tickets")
		require.NoError(t, err)

		ticketsToSeed := []*models.ServiceTicket{
			{
				ID:               uuid.New().String(),
				TicketNumber:     "TKT-001",
				Status:           models.StatusReceived,
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081111111111",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Battery drain",
			},
			{
				ID:               uuid.New().String(),
				TicketNumber:     "TKT-002",
				Status:           models.StatusRepairing,
				CustomerName:     "Citra Lestari",
				CustomerPhone:    "082222222222",
				DeviceBrand:      "Apple",
				DeviceModel:      "MacBook Pro",
				IssueDescription: "Keyboard issue",
			},
			{
				ID:               uuid.New().String(),
				TicketNumber:     "TKT-003",
				Status:           models.StatusCompleted,
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081111111111",
				DeviceBrand:      "Apple",
				DeviceModel:      "iPad Air",
				IssueDescription: "Charging port repair",
			},
		}

		for _, tk := range ticketsToSeed {
			err = cmdRepo.Create(ctx, tk)
			require.NoError(t, err)
		}

		// 1. FindAll - Filter by Status
		list, total, err := queryRepo.FindAll(ctx, string(models.StatusReceived), "", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		require.Len(t, list, 1)
		assert.Equal(t, "TKT-001", list[0].TicketNumber)

		// 2. FindAll - Search query (name/phone/number)
		list, total, err = queryRepo.FindAll(ctx, "", "Budi", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		require.Len(t, list, 2)

		// 3. Search - Active flag (IsActive true means NOT IN 'COMPLETED', 'RETURNED')
		isActive := true
		searchReq := ticket.TicketSearchRequest{
			IsActive: &isActive,
			Limit:    10,
			Offset:   0,
		}
		list, total, err = queryRepo.Search(ctx, searchReq)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		require.Len(t, list, 2)

		// 4. Search - Inactive flag (IsActive false means IN 'COMPLETED', 'RETURNED')
		isNotActive := false
		searchReq.IsActive = &isNotActive
		list, total, err = queryRepo.Search(ctx, searchReq)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		require.Len(t, list, 1)
		assert.Equal(t, "TKT-003", list[0].TicketNumber)

		// 5. Search - Text search (across ID, number, customer, brand, model)
		isActive = true
		searchReq.IsActive = &isActive
		searchReq.Search = "MacBook"
		list, total, err = queryRepo.Search(ctx, searchReq)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		require.Len(t, list, 1)
		assert.Equal(t, "TKT-002", list[0].TicketNumber)
	})
}
