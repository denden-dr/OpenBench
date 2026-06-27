package ticket_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/denden-dr/openbench/apps/backend/internal/ticket"
	"github.com/denden-dr/openbench/apps/backend/internal/ticket/mocks"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupServiceTest(t *testing.T) (*mocks.TicketRepository, ticket.AdminTicketService, sqlmock.Sqlmock) {
	mockDB, mockSQL, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "postgres")
	dbWrapper := &database.Database{DB: sqlxDB}

	repo := mocks.NewTicketRepository(t)
	service := ticket.NewAdminService(repo, dbWrapper)

	t.Cleanup(func() {
		mockDB.Close()
		assert.NoError(t, mockSQL.ExpectationsWereMet())
	})

	return repo, service, mockSQL
}

func TestService_CreateTicket(t *testing.T) {
	repo, service, mockSQL := setupServiceTest(t)
	ctx := context.Background()

	serial := "SN123456"
	validReq := &api.TicketCreate{
		CustomerName:         "John Doe",
		CustomerPhone:        "0812345678",
		BrandPhone:           "Apple",
		ModelPhone:           "iPhone 15",
		SerialNumber:         &serial,
		DamageDescription:    "Shattered Screen",
		WarrantyDurationDays: 30,
	}

	t.Run("Success", func(t *testing.T) {
		mockSQL.ExpectBegin()
		repo.On("GetMaxTicketNumberByPrefix", ctx, mock.Anything, mock.Anything).Return("", nil).Once()
		repo.On("Create", ctx, mock.Anything, mock.Anything).Return(nil).Once()
		mockSQL.ExpectCommit()

		tkt, err := service.CreateTicket(ctx, validReq)
		require.NoError(t, err)
		assert.NotNil(t, tkt)
		assert.NotEmpty(t, tkt.ID)
		assert.Contains(t, tkt.TicketNumber, "OB-")
		assert.Equal(t, "John Doe", tkt.CustomerName)
		assert.Equal(t, 30, tkt.WarrantyDurationDays)
		assert.Equal(t, "received", tkt.Status)
	})

	t.Run("Sequential Ticket Numbers", func(t *testing.T) {
		mockSQL.ExpectBegin()
		repo.On("GetMaxTicketNumberByPrefix", ctx, mock.Anything, mock.Anything).Return("OB-202606-0005-ABCD", nil).Once()
		repo.On("Create", ctx, mock.Anything, mock.MatchedBy(func(tk *ticket.Ticket) bool {
			parts := strings.Split(tk.TicketNumber, "-")
			return len(parts) == 4 && parts[2] == "0006" && len(parts[3]) == 4
		})).Return(nil).Once()
		mockSQL.ExpectCommit()

		tkt, err := service.CreateTicket(ctx, validReq)
		require.NoError(t, err)
		assert.NotNil(t, tkt)
		parts := strings.Split(tkt.TicketNumber, "-")
		assert.Equal(t, 4, len(parts))
		assert.Equal(t, "0006", parts[2])
		assert.Equal(t, 4, len(parts[3]))
	})

	t.Run("Validation - Missing Required Fields", func(t *testing.T) {
		invalidReq := &api.TicketCreate{
			CustomerPhone:        "0812345678",
			BrandPhone:           "Apple",
			ModelPhone:           "iPhone 15",
			DamageDescription:    "Shattered Screen",
			WarrantyDurationDays: 30,
		}

		tkt, err := service.CreateTicket(ctx, invalidReq)
		assert.ErrorIs(t, err, ticket.ErrInvalidInput)
		assert.Nil(t, tkt)
	})

	t.Run("Validation - Negative Warranty Days", func(t *testing.T) {
		invalidReq := &api.TicketCreate{
			CustomerName:         "John Doe",
			CustomerPhone:        "0812345678",
			BrandPhone:           "Apple",
			ModelPhone:           "iPhone 15",
			DamageDescription:    "Shattered Screen",
			WarrantyDurationDays: -5,
		}

		tkt, err := service.CreateTicket(ctx, invalidReq)
		assert.ErrorIs(t, err, ticket.ErrInvalidInput)
		assert.Nil(t, tkt)
	})
}

func TestService_GetTicket(t *testing.T) {
	repo, service, _ := setupServiceTest(t)
	ctx := context.Background()
	testID := "tkt-1"

	t.Run("Success", func(t *testing.T) {
		expectedTicket := &ticket.Ticket{ID: testID, CustomerName: "Jane Doe"}
		repo.On("GetByID", ctx, mock.Anything, testID).Return(expectedTicket, nil).Once()

		res, err := service.GetTicket(ctx, testID)
		require.NoError(t, err)
		assert.Equal(t, expectedTicket, res)
	})

	t.Run("Not Found", func(t *testing.T) {
		repo.On("GetByID", ctx, mock.Anything, testID).Return(nil, sql.ErrNoRows).Once()

		res, err := service.GetTicket(ctx, testID)
		assert.ErrorIs(t, err, ticket.ErrTicketNotFound)
		assert.Nil(t, res)
	})
}

func TestService_GetTicketByNumber(t *testing.T) {
	repo, service, _ := setupServiceTest(t)
	ctx := context.Background()
	testNumber := "OB-202606-0001-A9X2"

	t.Run("Success", func(t *testing.T) {
		expectedTicket := &ticket.Ticket{TicketNumber: testNumber, CustomerName: "Jane Doe"}
		repo.On("GetByTicketNumber", ctx, mock.Anything, testNumber).Return(expectedTicket, nil).Once()

		res, err := service.GetTicketByNumber(ctx, testNumber)
		require.NoError(t, err)
		assert.Equal(t, expectedTicket, res)
	})

	t.Run("Not Found", func(t *testing.T) {
		repo.On("GetByTicketNumber", ctx, mock.Anything, testNumber).Return(nil, sql.ErrNoRows).Once()

		res, err := service.GetTicketByNumber(ctx, testNumber)
		assert.ErrorIs(t, err, ticket.ErrTicketNotFound)
		assert.Nil(t, res)
	})
}

func TestService_UpdateTicket_PickedUp(t *testing.T) {
	repo, service, mockSQL := setupServiceTest(t)
	ctx := context.Background()
	testID := "tkt-1"

	originalTicket := &ticket.Ticket{
		ID:                   testID,
		TicketNumber:         "OB-202606-0001",
		CustomerName:         "John",
		BrandPhone:           "Samsung",
		ModelPhone:           "S24",
		Status:               "completed",
		WarrantyDurationDays: 14,
	}

	t.Run("Status transitioning to completed with device pickup", func(t *testing.T) {
		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(originalTicket, nil).Once()
		repo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()
		mockSQL.ExpectCommit()

		posVal := api.TicketUpdateDevicePosition("picked_up")
		pmVal := api.TicketUpdatePaymentMethod("cash")
		updateReq := &api.TicketUpdate{
			DevicePosition: &posVal,
			PaymentMethod:  &pmVal,
		}

		res, err := service.UpdateTicket(ctx, testID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, "completed", res.Status)
		assert.Equal(t, "picked_up", res.DevicePosition)
		assert.Equal(t, "paid", res.PaymentStatus)
		assert.NotNil(t, res.PaymentMethod)
		assert.Equal(t, "cash", *res.PaymentMethod)
		assert.NotNil(t, res.PickedUpAt)
	})

	t.Run("Re-applying picked_up device position", func(t *testing.T) {
		fixedTime := time.Now().Add(-24 * time.Hour)
		pm := "cash"
		alreadyPickedUpTicket := &ticket.Ticket{
			ID:                   testID,
			TicketNumber:         "OB-202606-0001",
			CustomerName:         "John",
			BrandPhone:           "Samsung",
			ModelPhone:           "S24",
			Status:               "completed",
			DevicePosition:       "picked_up",
			PaymentStatus:        "paid",
			PaymentMethod:        &pm,
			PickedUpAt:           &fixedTime,
			WarrantyDurationDays: 14,
		}

		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(alreadyPickedUpTicket, nil).Once()
		repo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()
		mockSQL.ExpectCommit()

		posVal := api.TicketUpdateDevicePosition("picked_up")
		updateReq := &api.TicketUpdate{
			DevicePosition: &posVal,
		}

		res, err := service.UpdateTicket(ctx, testID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, fixedTime, *res.PickedUpAt)
	})

	t.Run("Updating warranty duration after pickup fails", func(t *testing.T) {
		fixedTime := time.Now().Add(-24 * time.Hour)
		pm := "cash"
		alreadyPickedUpTicket := &ticket.Ticket{
			ID:                   testID,
			TicketNumber:         "OB-202606-0001",
			CustomerName:         "John",
			BrandPhone:           "Samsung",
			ModelPhone:           "S24",
			Status:               "completed",
			DevicePosition:       "picked_up",
			PaymentStatus:        "paid",
			PaymentMethod:        &pm,
			PickedUpAt:           &fixedTime,
			WarrantyDurationDays: 14,
		}

		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(alreadyPickedUpTicket, nil).Once()
		repo.On("MockSQLRollback").Maybe() // Just referencing context
		mockSQL.ExpectRollback()

		warrantyDays := 30
		updateReq := &api.TicketUpdate{
			WarrantyDurationDays: &warrantyDays,
		}

		res, err := service.UpdateTicket(ctx, testID, updateReq)
		assert.ErrorIs(t, err, ticket.ErrInvalidInput)
		assert.Nil(t, res)
	})
}

func TestService_UpdateTicket_CustomerFields(t *testing.T) {
	repo, service, mockSQL := setupServiceTest(t)
	ctx := context.Background()
	testID := "tkt-1"

	t.Run("Update customer and device fields with no active warranty", func(t *testing.T) {
		originalTicket := &ticket.Ticket{
			ID:            testID,
			TicketNumber:  "OB-202606-0001",
			CustomerName:  "Original Name",
			CustomerPhone: "0812000000",
			BrandPhone:    "Original Brand",
			ModelPhone:    "Original Model",
			SerialNumber:  "Original SN",
			Status:        "received",
		}

		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(originalTicket, nil).Once()
		repo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()
		mockSQL.ExpectCommit()

		newName := "New Name"
		newPhone := "0812999999"
		newBrand := "New Brand"
		newModel := "New Model"
		newSN := "New SN"

		updateReq := &api.TicketUpdate{
			CustomerName:  &newName,
			CustomerPhone: &newPhone,
			BrandPhone:    &newBrand,
			ModelPhone:    &newModel,
			SerialNumber:  &newSN,
		}

		res, err := service.UpdateTicket(ctx, testID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, newName, res.CustomerName)
		assert.Equal(t, newPhone, res.CustomerPhone)
		assert.Equal(t, newBrand, res.BrandPhone)
		assert.Equal(t, newModel, res.ModelPhone)
		assert.Equal(t, newSN, res.SerialNumber)
		assert.Nil(t, res.Warranty)
	})

	t.Run("Update customer and device fields with active warranty sync", func(t *testing.T) {
		pm := "cash"
		tktWarranty := &ticket.Warranty{
			ID:           "w-1",
			TicketID:     testID,
			CustomerName: "Original Name",
			DeviceInfo:   "Original Brand Original Model",
		}
		originalTicket := &ticket.Ticket{
			ID:             testID,
			TicketNumber:   "OB-202606-0001",
			CustomerName:   "Original Name",
			CustomerPhone:  "0812000000",
			BrandPhone:     "Original Brand",
			ModelPhone:     "Original Model",
			SerialNumber:   "Original SN",
			Status:         "completed",
			DevicePosition: "picked_up",
			PaymentStatus:  "paid",
			PaymentMethod:  &pm,
			Warranty:       tktWarranty,
		}

		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(originalTicket, nil).Once()
		repo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()
		mockSQL.ExpectCommit()

		newName := "New Name"
		newBrand := "New Brand"
		newModel := "New Model"

		updateReq := &api.TicketUpdate{
			CustomerName: &newName,
			BrandPhone:   &newBrand,
			ModelPhone:   &newModel,
		}

		res, err := service.UpdateTicket(ctx, testID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, newName, res.CustomerName)
		assert.Equal(t, newBrand, res.BrandPhone)
		assert.Equal(t, newModel, res.ModelPhone)
		assert.NotNil(t, res.Warranty)
		assert.Equal(t, newName, res.Warranty.CustomerName)
		assert.Equal(t, "New Brand New Model", res.Warranty.DeviceInfo)
	})

	t.Run("UpdateTicket prevents status reversal when already picked up", func(t *testing.T) {
		pm := "cash"
		originalTicket := &ticket.Ticket{
			ID:             testID,
			TicketNumber:   "OB-202606-0001",
			Status:         "completed",
			DevicePosition: "picked_up",
			PaymentStatus:  "paid",
			PaymentMethod:  &pm,
		}

		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(originalTicket, nil).Once()
		mockSQL.ExpectRollback()

		newPos := api.TicketUpdateDevicePositionWarehouse
		updateReq := &api.TicketUpdate{
			DevicePosition: &newPos,
		}

		res, err := service.UpdateTicket(ctx, testID, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "location reversal not allowed under normal updates")
		assert.Nil(t, res)
	})

	t.Run("EmergencyUpdateTicket allows status reversal and deletes warranty", func(t *testing.T) {
		pm := "cash"
		tktWarranty := &ticket.Warranty{
			ID:           "w-1",
			TicketID:     testID,
			CustomerName: "John Doe",
		}
		originalTicket := &ticket.Ticket{
			ID:             testID,
			TicketNumber:   "OB-202606-0001",
			CustomerName:   "John Doe",
			Status:         "completed",
			DevicePosition: "picked_up",
			PaymentStatus:  "paid",
			PaymentMethod:  &pm,
			Warranty:       tktWarranty,
		}

		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(originalTicket, nil).Once()
		repo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()
		repo.On("DeleteWarrantyByTicketID", ctx, mock.Anything, testID).Return(nil).Once()
		mockSQL.ExpectCommit()

		newPos := api.TicketUpdateDevicePositionWarehouse
		updateReq := &api.TicketUpdate{
			DevicePosition: &newPos,
		}

		res, err := service.EmergencyUpdateTicket(ctx, testID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, "warehouse", res.DevicePosition)
		assert.Nil(t, res.PickedUpAt)
		assert.Nil(t, res.Warranty)
	})

	t.Run("EmergencyUpdateTicket updates fundamental fields and recalculates warranty duration", func(t *testing.T) {
		pm := "cash"
		now := time.Now()
		tktWarranty := &ticket.Warranty{
			ID:           "w-1",
			TicketID:     testID,
			CustomerName: "John Doe",
			StartDate:    now,
			EndDate:      now.AddDate(0, 0, 30),
		}
		originalTicket := &ticket.Ticket{
			ID:                   testID,
			TicketNumber:         "OB-202606-0001",
			CustomerName:         "John Doe",
			Status:               "completed",
			DevicePosition:       "picked_up",
			PaymentStatus:        "paid",
			PaymentMethod:        &pm,
			PickedUpAt:           &now,
			WarrantyDurationDays: 30,
			Warranty:             tktWarranty,
		}

		mockSQL.ExpectBegin()
		repo.On("GetByIDWithLock", ctx, mock.Anything, testID).Return(originalTicket, nil).Once()
		repo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()
		mockSQL.ExpectCommit()

		newName := "Jane Doe"
		newDays := 60
		updateReq := &api.TicketUpdate{
			CustomerName:         &newName,
			WarrantyDurationDays: &newDays,
		}

		res, err := service.EmergencyUpdateTicket(ctx, testID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, "Jane Doe", res.CustomerName)
		assert.Equal(t, 60, res.WarrantyDurationDays)
		assert.NotNil(t, res.Warranty)
		assert.Equal(t, "Jane Doe", res.Warranty.CustomerName)
		// Check that the end date was updated to 60 days from now
		expectedEnd := now.AddDate(0, 0, 60)
		assert.True(t, res.Warranty.EndDate.Equal(expectedEnd))
	})
}
