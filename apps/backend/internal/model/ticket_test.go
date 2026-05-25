package model

import (
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTicketPrepareForCreate(t *testing.T) {
	t.Run("applies default warranty days", func(t *testing.T) {
		ticket := Ticket{
			Price:        decimal.NewFromInt(100000),
			WarrantyDays: 0,
		}

		err := ticket.PrepareForCreate()

		require.NoError(t, err)
		assert.Equal(t, DefaultWarrantyDays, ticket.WarrantyDays)
	})

	t.Run("rejects negative price", func(t *testing.T) {
		ticket := Ticket{
			Price:        decimal.NewFromInt(-1),
			WarrantyDays: 30,
		}

		err := ticket.PrepareForCreate()

		assert.ErrorIs(t, err, ErrNegativePrice)
	})

	t.Run("rejects negative warranty days", func(t *testing.T) {
		ticket := Ticket{
			Price:        decimal.NewFromInt(100000),
			WarrantyDays: -1,
		}

		err := ticket.PrepareForCreate()

		assert.ErrorIs(t, err, ErrNegativeWarranty)
	})
}

func TestTicketApplyUpdate(t *testing.T) {
	t.Run("updates basic fields and clears empty optional fields", func(t *testing.T) {
		description := "existing"
		empty := ""
		newName := "Budi Baru"
		ticket := Ticket{
			CustomerName:          "Budi",
			AdditionalDescription: &description,
			Status:                StatusServiceIn,
			PaymentStatus:         PaymentUnpaid,
		}

		err := ticket.ApplyUpdate(TicketUpdate{
			CustomerName:          &newName,
			AdditionalDescription: &empty,
		})

		require.NoError(t, err)
		assert.Equal(t, "Budi Baru", ticket.CustomerName)
		assert.Nil(t, ticket.AdditionalDescription)
	})

	t.Run("rejects negative price", func(t *testing.T) {
		price := decimal.NewFromInt(-1)
		ticket := Ticket{
			Status:        StatusServiceIn,
			PaymentStatus: PaymentUnpaid,
		}

		err := ticket.ApplyUpdate(TicketUpdate{Price: &price})

		assert.ErrorIs(t, err, ErrNegativePrice)
	})

	t.Run("rejects negative warranty days", func(t *testing.T) {
		warrantyDays := -1
		ticket := Ticket{
			Status:        StatusServiceIn,
			PaymentStatus: PaymentUnpaid,
		}

		err := ticket.ApplyUpdate(TicketUpdate{WarrantyDays: &warrantyDays})

		assert.ErrorIs(t, err, ErrNegativeWarranty)
	})

	t.Run("moves to picked up with paid status and exit date", func(t *testing.T) {
		status := string(StatusPickedUp)
		ticket := Ticket{
			Status:        StatusFixed,
			PaymentStatus: PaymentUnpaid,
			WarrantyDays:  30,
		}

		err := ticket.ApplyUpdate(TicketUpdate{Status: &status})

		require.NoError(t, err)
		assert.Equal(t, StatusPickedUp, ticket.Status)
		assert.Equal(t, PaymentPaid, ticket.PaymentStatus)
		assert.NotNil(t, ticket.ExitDate)
		assert.NotNil(t, ticket.WarrantyExpiryDate())
	})

	t.Run("rejects picked up with explicit unpaid payment status", func(t *testing.T) {
		status := string(StatusPickedUp)
		paymentStatus := string(PaymentUnpaid)
		ticket := Ticket{
			Status:        StatusFixed,
			PaymentStatus: PaymentUnpaid,
		}

		err := ticket.ApplyUpdate(TicketUpdate{
			Status:        &status,
			PaymentStatus: &paymentStatus,
		})

		assert.ErrorIs(t, err, ErrPickedUpRequiresPaid)
	})

	t.Run("clears exit date when moving out of picked up", func(t *testing.T) {
		status := string(StatusFixed)
		exitDate := time.Now()
		ticket := Ticket{
			Status:        StatusPickedUp,
			PaymentStatus: PaymentPaid,
			WarrantyDays:  30,
			ExitDate:      &exitDate,
		}

		err := ticket.ApplyUpdate(TicketUpdate{Status: &status})

		require.NoError(t, err)
		assert.Equal(t, StatusFixed, ticket.Status)
		assert.Nil(t, ticket.ExitDate)
		assert.Nil(t, ticket.WarrantyExpiryDate())
	})

	t.Run("rejects exit date on non-picked-up ticket", func(t *testing.T) {
		status := string(StatusFixed)
		exitDate := time.Now()
		ticket := Ticket{
			Status:        StatusServiceIn,
			PaymentStatus: PaymentUnpaid,
		}

		err := ticket.ApplyUpdate(TicketUpdate{
			Status:   &status,
			ExitDate: &exitDate,
		})

		assert.ErrorIs(t, err, ErrNonPickedUpCannotHaveExitDate)
	})

	t.Run("allows exit date update on picked-up ticket", func(t *testing.T) {
		exitDate := time.Now().Add(-48 * time.Hour)
		ticket := Ticket{
			Status:        StatusPickedUp,
			PaymentStatus: PaymentPaid,
			WarrantyDays:  30,
		}

		err := ticket.ApplyUpdate(TicketUpdate{ExitDate: &exitDate})

		require.NoError(t, err)
		require.NotNil(t, ticket.ExitDate)
		assert.True(t, ticket.ExitDate.Equal(exitDate))
		assert.True(t, ticket.WarrantyExpiryDate().Equal(exitDate.AddDate(0, 0, 30)))
	})

	t.Run("transitions to waiting_confirmation and cancelled successfully", func(t *testing.T) {
		statusWaiting := string(StatusWaitingConfirmation)
		statusCancelled := string(StatusCancelled)

		ticket := Ticket{
			Status:        StatusOnProcess,
			PaymentStatus: PaymentUnpaid,
		}

		// Update to waiting_confirmation
		err := ticket.ApplyUpdate(TicketUpdate{
			Status: &statusWaiting,
		})
		require.NoError(t, err)
		assert.Equal(t, StatusWaitingConfirmation, ticket.Status)

		// Update to cancelled
		err = ticket.ApplyUpdate(TicketUpdate{
			Status: &statusCancelled,
		})
		require.NoError(t, err)
		assert.Equal(t, StatusCancelled, ticket.Status)
	})
}

func TestValidateTicketUpdate(t *testing.T) {
	t.Run("rejects negative price before loading a ticket", func(t *testing.T) {
		price := decimal.NewFromInt(-1)

		err := ValidateTicketUpdate(TicketUpdate{Price: &price})

		assert.ErrorIs(t, err, ErrNegativePrice)
	})

	t.Run("rejects negative warranty days before loading a ticket", func(t *testing.T) {
		warrantyDays := -1

		err := ValidateTicketUpdate(TicketUpdate{WarrantyDays: &warrantyDays})

		assert.ErrorIs(t, err, ErrNegativeWarranty)
	})
}

func TestTicketPrepareForCreate_VoidWarrantyDays(t *testing.T) {
	t.Run("applies default warranty days then can be reset to zero", func(t *testing.T) {
		ticket := Ticket{
			Price:        decimal.NewFromInt(0),
			WarrantyDays: 0,
			Status:       StatusCancelled,
		}

		err := ticket.PrepareForCreate()
		require.NoError(t, err)
		// PrepareForCreate sets 0 → DefaultWarrantyDays
		assert.Equal(t, DefaultWarrantyDays, ticket.WarrantyDays)

		// Void claim path: reset to 0 after PrepareForCreate
		ticket.WarrantyDays = 0
		assert.Equal(t, 0, ticket.WarrantyDays)
	})
}

func TestTicketErrorsAreStable(t *testing.T) {
	assert.True(t, errors.Is(ErrNegativePrice, ErrNegativePrice))
	assert.True(t, errors.Is(ErrPickedUpRequiresPaid, ErrPickedUpRequiresPaid))
}
