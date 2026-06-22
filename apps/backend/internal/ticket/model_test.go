package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTicket_NewTicket(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tkt, err := NewTicket(
			"OB-202606-0001",
			"John Doe",
			"08123456789",
			"Samsung",
			"S24",
			"SN123",
			"Broken Screen",
			"Clean dust, replace screen",
			150000,
			14,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, tkt.ID)
		assert.Equal(t, "OB-202606-0001", tkt.TicketNumber)
		assert.Equal(t, "John Doe", tkt.CustomerName)
		assert.Equal(t, "Clean dust, replace screen", tkt.RepairAction)
		assert.Equal(t, StatusReceived, tkt.Status)
		assert.Equal(t, PositionWarehouse, tkt.DevicePosition)
		assert.Equal(t, PaymentStatusNone, tkt.PaymentStatus)
		assert.Equal(t, 14, tkt.WarrantyDurationDays)
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		_, err := NewTicket(
			"OB-202606-0001",
			"", // missing CustomerName
			"08123456789",
			"Samsung",
			"S24",
			"SN123",
			"Broken Screen",
			"",
			150000,
			14,
		)
		assert.ErrorIs(t, err, ErrMissingRequiredFields)
	})

	t.Run("Negative Warranty Duration", func(t *testing.T) {
		_, err := NewTicket(
			"OB-202606-0001",
			"John Doe",
			"08123456789",
			"Samsung",
			"S24",
			"SN123",
			"Broken Screen",
			"",
			150000,
			-5,
		)
		assert.ErrorIs(t, err, ErrWarrantyDurationNegative)
	})
}

func TestTicket_ProcessPickup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tkt := &Ticket{
			ID:                   "tkt-1",
			TicketNumber:         "OB-202606-0001",
			CustomerName:         "John Doe",
			BrandPhone:           "Samsung",
			ModelPhone:           "S24",
			Status:               "completed",
			WarrantyDurationDays: 14,
		}

		now := time.Now()
		err := tkt.ProcessPickup(now)
		assert.NoError(t, err)
		assert.Equal(t, "completed", tkt.Status)
		assert.Equal(t, "picked_up", tkt.DevicePosition)
		assert.Equal(t, "paid", tkt.PaymentStatus)
		assert.Equal(t, now, *tkt.PickedUpAt)

		assert.NotNil(t, tkt.Warranty)
		assert.Equal(t, tkt.ID, tkt.Warranty.TicketID)
		assert.Equal(t, "active", tkt.Warranty.Status)
	})

	t.Run("Already Picked Up", func(t *testing.T) {
		tkt := &Ticket{
			DevicePosition: "picked_up",
		}
		err := tkt.ProcessPickup(time.Now())
		assert.ErrorIs(t, err, ErrTicketAlreadyPickedUp)
	})

	t.Run("Already set PickedUpAt is preserved", func(t *testing.T) {
		fixedTime := time.Now().Add(-24 * time.Hour)
		tkt := &Ticket{
			ID:                   "tkt-1",
			TicketNumber:         "OB-202606-0001",
			CustomerName:         "John Doe",
			BrandPhone:           "Samsung",
			ModelPhone:           "S24",
			Status:               "completed",
			WarrantyDurationDays: 14,
			PickedUpAt:           &fixedTime,
		}

		err := tkt.ProcessPickup(time.Now())
		assert.NoError(t, err)
		assert.Equal(t, fixedTime, *tkt.PickedUpAt)
	})
}

func TestTicket_UpdateWarrantyDuration(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tkt := &Ticket{
			Status:               "completed",
			WarrantyDurationDays: 14,
		}

		err := tkt.UpdateWarrantyDuration(30)
		assert.NoError(t, err)
		assert.Equal(t, 30, tkt.WarrantyDurationDays)
	})

	t.Run("Negative Duration", func(t *testing.T) {
		tkt := &Ticket{
			Status: "completed",
		}

		err := tkt.UpdateWarrantyDuration(-5)
		assert.ErrorIs(t, err, ErrWarrantyDurationNegative)
	})

	t.Run("Already Picked Up", func(t *testing.T) {
		tkt := &Ticket{
			DevicePosition: "picked_up",
		}

		err := tkt.UpdateWarrantyDuration(30)
		assert.ErrorIs(t, err, ErrWarrantyUpdateAfterPickup)
	})
}

func TestTicket_Validate(t *testing.T) {
	t.Run("Valid picked_up state", func(t *testing.T) {
		pm := "cash"
		tkt := &Ticket{
			Status:         "completed",
			DevicePosition: "picked_up",
			PaymentStatus:  "paid",
			PaymentMethod:  &pm,
		}
		err := tkt.Validate()
		assert.NoError(t, err)
	})

	t.Run("Non-picked_up state ignores payment validation", func(t *testing.T) {
		tkt := &Ticket{
			Status:         "received",
			DevicePosition: "warehouse",
			PaymentStatus:  "none",
			PaymentMethod:  nil,
		}
		err := tkt.Validate()
		assert.NoError(t, err)
	})

	t.Run("Missing paid status when picked_up", func(t *testing.T) {
		pm := "cash"
		tkt := &Ticket{
			Status:         "completed",
			DevicePosition: "picked_up",
			PaymentStatus:  "none",
			PaymentMethod:  &pm,
		}
		err := tkt.Validate()
		assert.ErrorIs(t, err, ErrPaymentRequiredForPickup)
	})

	t.Run("Missing payment method when picked_up", func(t *testing.T) {
		tkt := &Ticket{
			Status:         "completed",
			DevicePosition: "picked_up",
			PaymentStatus:  "paid",
			PaymentMethod:  nil,
		}
		err := tkt.Validate()
		assert.ErrorIs(t, err, ErrPaymentMethodRequired)
	})

	t.Run("Invalid payment method when picked_up", func(t *testing.T) {
		pm := "bank_transfer"
		tkt := &Ticket{
			Status:         "completed",
			DevicePosition: "picked_up",
			PaymentStatus:  "paid",
			PaymentMethod:  &pm,
		}
		err := tkt.Validate()
		assert.ErrorIs(t, err, ErrPaymentMethodRequired)
	})
}

func TestTicket_UpdateStatus(t *testing.T) {
	t.Run("Transition to other status", func(t *testing.T) {
		tkt := &Ticket{
			Status: "received",
		}
		err := tkt.UpdateStatus("in_repair", time.Now())
		assert.NoError(t, err)
		assert.Equal(t, "in_repair", tkt.Status)
	})
}
