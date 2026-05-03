package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Ticket struct {
	ID               string          `db:"id"`
	DeviceType       string          `db:"device_type"`
	Brand            string          `db:"brand"`
	Model            string          `db:"model"`
	IssueDescription string          `db:"issue_description"`
	Status           string          `db:"status"`
	DiagnosisFee     decimal.Decimal `db:"diagnosis_fee"`
	TechnicianID     *string         `db:"technician_id"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
}
