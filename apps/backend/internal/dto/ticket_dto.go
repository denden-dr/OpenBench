package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateTicketRequest struct {
	CustomerName          string          `json:"customer_name" validate:"required"`
	CustomerGender        string          `json:"customer_gender" validate:"required,oneof=Male Female Other"`
	Brand                 string          `json:"brand" validate:"required"`
	Model                 string          `json:"model" validate:"required"`
	Issue                 string          `json:"issue" validate:"required"`
	AdditionalDescription string          `json:"additional_description"`
	Accessories           string          `json:"accessories"`
	Price                 decimal.Decimal `json:"price"`
	WarrantyDays          int             `json:"warranty_days"`
}

type UpdateTicketRequest struct {
	CustomerName          *string          `json:"customer_name" validate:"omitempty"`
	CustomerGender        *string          `json:"customer_gender" validate:"omitempty,oneof=Male Female Other"`
	Brand                 *string          `json:"brand" validate:"omitempty"`
	Model                 *string          `json:"model" validate:"omitempty"`
	Issue                 *string          `json:"issue" validate:"omitempty"`
	AdditionalDescription *string          `json:"additional_description"`
	Accessories           *string          `json:"accessories"`
	Price                 *decimal.Decimal `json:"price" validate:"omitempty"`
	Status                *string          `json:"status" validate:"omitempty,oneof=service_in on_process fixed picked_up"`
	PaymentStatus         *string          `json:"payment_status" validate:"omitempty,oneof=unpaid paid"`
	WarrantyDays          *int             `json:"warranty_days" validate:"omitempty,min=0"`
}

type TicketResponse struct {
	ID                    string          `json:"id"`
	CustomerName          string          `json:"customer_name"`
	CustomerGender        string          `json:"customer_gender"`
	Brand                 string          `json:"brand"`
	Model                 string          `json:"model"`
	Issue                 string          `json:"issue"`
	AdditionalDescription *string         `json:"additional_description"`
	Accessories           *string         `json:"accessories"`
	Price                 decimal.Decimal `json:"price"`
	Status                string          `json:"status"`
	PaymentStatus         string          `json:"payment_status"`
	WarrantyDays          int             `json:"warranty_days"`
	EntryDate             time.Time       `json:"entry_date"`
	ExitDate              *time.Time      `json:"exit_date"`
	WarrantyExpiryDate    *time.Time      `json:"warranty_expiry_date"`
}
