package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceTicket_Invariants(t *testing.T) {
	tests := []struct {
		name        string
		params      CreateTicketParams
		expectedErr error
	}{
		{
			name: "Success - valid parameters",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
				Cost:             1500000,
				WarrantyDays:     30,
			},
			expectedErr: nil,
		},
		{
			name: "Failure - empty customer name",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
			},
			expectedErr: ErrMissingCustomerName,
		},
		{
			name: "Failure - customer name only whitespace",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "   ",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
			},
			expectedErr: ErrMissingCustomerName,
		},
		{
			name: "Failure - empty customer phone",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "Budi Santoso",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
			},
			expectedErr: ErrMissingCustomerPhone,
		},
		{
			name: "Failure - empty device brand",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
			},
			expectedErr: ErrMissingDeviceBrand,
		},
		{
			name: "Failure - empty device model",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				IssueDescription: "Layar pecah",
			},
			expectedErr: ErrMissingDeviceModel,
		},
		{
			name: "Failure - empty issue description",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
			},
			expectedErr: ErrMissingIssueDescription,
		},
		{
			name: "Failure - negative cost",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
				Cost:             -1,
			},
			expectedErr: ErrNegativeCost,
		},
		{
			name: "Failure - negative warranty days",
			params: CreateTicketParams{
				TicketNumber:     "TKT-1234",
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
				WarrantyDays:     -5,
			},
			expectedErr: ErrNegativeWarrantyDays,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket, err := NewServiceTicket(tt.params)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Nil(t, ticket)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ticket)
				assert.NotEmpty(t, ticket.ID)
				assert.Equal(t, StatusReceived, ticket.Status)
				assert.Equal(t, tt.params.CustomerName, ticket.CustomerName)
				assert.Equal(t, tt.params.CustomerPhone, ticket.CustomerPhone)
				assert.Equal(t, tt.params.DeviceBrand, ticket.DeviceBrand)
				assert.Equal(t, tt.params.DeviceModel, ticket.DeviceModel)
				assert.Equal(t, tt.params.IssueDescription, ticket.IssueDescription)
				assert.Equal(t, tt.params.Cost, ticket.Cost)
				assert.Equal(t, tt.params.WarrantyDays, ticket.WarrantyDays)
			}
		})
	}
}
