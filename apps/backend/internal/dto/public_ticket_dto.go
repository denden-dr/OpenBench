package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type PublicTicketResponse struct {
	ID                    string           `json:"id"`
	CustomerNameMasked    string           `json:"customer_name_masked"`
	CustomerPhoneMasked   string           `json:"customer_phone_masked,omitempty"`
	Brand                 string           `json:"brand"`
	Model                 string           `json:"model"`
	Issue                 string           `json:"issue"`
	Status                string           `json:"status"`
	EntryDate             time.Time        `json:"entry_date"`
	ExitDate              *time.Time       `json:"exit_date,omitempty"`
	WarrantyDays          int              `json:"warranty_days,omitempty"`
	PaymentStatus         string           `json:"payment_status,omitempty"`
	Price                 *decimal.Decimal `json:"price,omitempty"`
	AdditionalDescription *string          `json:"additional_description,omitempty"`
	Accessories           *string          `json:"accessories,omitempty"`
	WarrantyExpiryDate    *time.Time       `json:"warranty_expiry_date,omitempty"`
}

type PublicTrackRequest struct {
	ShortID string `json:"short_id" validate:"required,len=8,hexadecimal"`
	Phone   string `json:"phone" validate:"required"`
}
