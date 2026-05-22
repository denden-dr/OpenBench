package model

import (
	"time"
	"github.com/shopspring/decimal"
)

type Ticket struct {
	ID                    string          `db:"id" json:"id"`
	CustomerName          string          `db:"customer_name" json:"customer_name"`
	CustomerGender        string          `db:"customer_gender" json:"customer_gender"`
	Brand                 string          `db:"brand" json:"brand"`
	Model                 string          `db:"model" json:"model"`
	Issue                 string          `db:"issue" json:"issue"`
	AdditionalDescription *string         `db:"additional_description" json:"additional_description"`
	Accessories           *string         `db:"accessories" json:"accessories"`
	Price                 decimal.Decimal `db:"price" json:"price"`
	Status                string          `db:"status" json:"status"`
	PaymentStatus         string          `db:"payment_status" json:"payment_status"`
	WarrantyDays          int             `db:"warranty_days" json:"warranty_days"`
	EntryDate             time.Time       `db:"entry_date" json:"entry_date"`
	ExitDate              *time.Time      `db:"exit_date" json:"exit_date"`
	WarrantyExpiryDate    *time.Time      `db:"warranty_expiry_date" json:"warranty_expiry_date"`
}
