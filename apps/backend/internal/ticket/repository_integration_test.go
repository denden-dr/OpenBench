//go:build integration

package ticket_test

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/denden-dr/openbench/apps/backend/internal/ticket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TicketRepositoryTestSuite struct {
	testutil.IntegrationSuite
	repo ticket.TicketRepository
}

func TestTicketRepositorySuite(t *testing.T) {
	suite.Run(t, new(TicketRepositoryTestSuite))
}

func (s *TicketRepositoryTestSuite) SetupTest() {
	s.IntegrationSuite.SetupTest() // Clean tables dynamically
	s.repo = ticket.NewRepository(s.DB)
}

func (s *TicketRepositoryTestSuite) TestTicketQueries() {
	ctx := context.Background()

	tID := uuid.New().String()
	ticketNumber := "OB-202606-0001"
	tkt := &ticket.Ticket{
		ID:                   tID,
		TicketNumber:         ticketNumber,
		CustomerName:         "Alice Budi",
		CustomerPhone:        "08111222333",
		BrandPhone:           "Xiaomi",
		ModelPhone:           "Redmi Note 13",
		SerialNumber:         "SN-XIAOMI-1",
		DamageDescription:    "Charging port broken",
		RepairAction:         "",
		Cost:                 0.0,
		Status:               "received",
		DevicePosition:       "warehouse",
		PaymentStatus:        "none",
		PaymentMethod:        nil,
		WarrantyDurationDays: 30,
		PickedUpAt:           nil,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// 1. Create Ticket
	s.Run("Create Ticket", func() {
		err := s.repo.Create(ctx, nil, tkt)
		s.Require().NoError(err)
	})

	// 2. GetByID
	s.Run("GetByID - Success", func() {
		res, err := s.repo.GetByID(ctx, nil, tID)
		s.Require().NoError(err)
		s.Assert().Equal(tkt.CustomerName, res.CustomerName)
		s.Assert().Equal(tkt.TicketNumber, res.TicketNumber)
		s.Assert().Equal(tkt.Status, res.Status)
		s.Assert().Equal(30, res.WarrantyDurationDays)
	})

	// 3. GetMaxTicketNumberByPrefix
	s.Run("GetMaxTicketNumberByPrefix", func() {
		maxNum, err := s.repo.GetMaxTicketNumberByPrefix(ctx, nil, "OB-202606-")
		s.Require().NoError(err)
		s.Assert().Equal("OB-202606-0001", maxNum)
	})

	// 4. Update Ticket
	s.Run("Update Ticket", func() {
		tkt.Status = "completed"
		tkt.DevicePosition = "picked_up"
		tkt.PaymentStatus = "paid"
		payMethod := "cash"
		tkt.PaymentMethod = &payMethod
		now := time.Now().Truncate(time.Second)
		tkt.PickedUpAt = &now

		err := s.repo.Update(ctx, nil, tkt)
		s.Require().NoError(err)

		res, err := s.repo.GetByID(ctx, nil, tID)
		s.Require().NoError(err)
		s.Assert().Equal("completed", res.Status)
		s.Assert().Equal("picked_up", res.DevicePosition)
		s.Assert().Equal("paid", res.PaymentStatus)
		s.Assert().Equal(now.Unix(), res.PickedUpAt.Unix())
	})

	// 4.5. Test Aggregate Warranty Cascading Upsert and Loading
	s.Run("Aggregate Warranty Cascading Upsert and Loading", func() {
		// Reset ticket state to allow ProcessPickup
		tkt.Status = "completed"
		tkt.DevicePosition = "warehouse"
		tkt.PaymentStatus = "paid"
		payMethod := "cash"
		tkt.PaymentMethod = &payMethod
		pickupTime := time.Now().Truncate(time.Second)
		err := tkt.ProcessPickup(pickupTime)
		s.Require().NoError(err)

		s.Require().NotNil(tkt.Warranty)
		s.Assert().Equal("active", tkt.Warranty.Status)

		// Save the entire aggregate (Ticket + Warranty)
		err = s.repo.Update(ctx, nil, tkt)
		s.Require().NoError(err)

		// Fetch the ticket from DB and verify the Warranty is loaded
		res, err := s.repo.GetByID(ctx, nil, tID)
		s.Require().NoError(err)
		s.Require().NotNil(res.Warranty)
		s.Assert().Equal(tkt.Warranty.ID, res.Warranty.ID)
		s.Assert().Equal("active", res.Warranty.Status)
		s.Assert().Equal(tkt.CustomerName, res.Warranty.CustomerName)
	})

	// 5. List Tickets
	s.Run("List Tickets", func() {
		list, err := s.repo.List(ctx, nil)
		s.Require().NoError(err)
		s.Assert().Len(list, 1)
		s.Assert().Equal(tID, list[0].ID)
	})
}

func (s *TicketRepositoryTestSuite) TestWarrantyQueries() {
	ctx := context.Background()

	// Seed ticket first
	tID := uuid.New().String()
	payMethod := "qris"
	tkt := &ticket.Ticket{
		ID:                   tID,
		TicketNumber:         "OB-202606-0002",
		CustomerName:         "Charlie",
		CustomerPhone:        "08112345",
		BrandPhone:           "Apple",
		ModelPhone:           "iPad Pro",
		SerialNumber:         "SN-IPAD-1",
		DamageDescription:    "Battery replacement",
		Status:               "completed",
		DevicePosition:       "picked_up",
		PaymentStatus:        "paid",
		PaymentMethod:        &payMethod,
		WarrantyDurationDays: 14,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	err := s.repo.Create(ctx, nil, tkt)
	s.Require().NoError(err)

	wID := uuid.New().String()
	startDate := time.Now().Truncate(time.Second)
	endDate := startDate.AddDate(0, 0, 14)
	warr := &ticket.Warranty{
		ID:           wID,
		TicketID:     tID,
		TicketNumber: tkt.TicketNumber,
		CustomerName: tkt.CustomerName,
		DeviceInfo:   "Apple iPad Pro",
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       "active",
	}

	// 1. CreateWarranty
	s.Run("CreateWarranty", func() {
		err := s.repo.CreateWarranty(ctx, nil, warr)
		s.Require().NoError(err)
	})

	// 2. GetWarrantyByTicketID
	s.Run("GetWarrantyByTicketID", func() {
		res, err := s.repo.GetWarrantyByTicketID(ctx, nil, tID)
		s.Require().NoError(err)
		s.Assert().Equal(wID, res.ID)
		s.Assert().Equal("Apple iPad Pro", res.DeviceInfo)
		s.Assert().Equal("active", res.Status)
	})

	// 3. ListWarranties
	s.Run("ListWarranties", func() {
		list, err := s.repo.ListWarranties(ctx, nil)
		s.Require().NoError(err)
		s.Assert().Len(list, 1)
		s.Assert().Equal(wID, list[0].ID)
	})
}
