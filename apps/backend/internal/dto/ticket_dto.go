package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateTicketRequest struct {
	DeviceType       string          `json:"device_type" validate:"required"`
	Brand            string          `json:"brand" validate:"required"`
	Model            string          `json:"model" validate:"required"`
	IssueDescription string          `json:"issue_description" validate:"required"`
	DiagnosisFee     decimal.Decimal `json:"diagnosis_fee" validate:"required"`
}

type TicketResponse struct {
	ID               string          `json:"id"`
	DeviceType       string          `json:"device_type"`
	Brand            string          `json:"brand"`
	Model            string          `json:"model"`
	IssueDescription string          `json:"issue_description"`
	Status           string          `json:"status"`
	DiagnosisFee     decimal.Decimal `json:"diagnosis_fee"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
