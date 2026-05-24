package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const DefaultWarrantyDays = 30

type TicketStatus string
type TicketPaymentStatus string

const (
	StatusServiceIn TicketStatus = "service_in"
	StatusOnProcess TicketStatus = "on_process"
	StatusFixed     TicketStatus = "fixed"
	StatusPickedUp  TicketStatus = "picked_up"

	PaymentUnpaid TicketPaymentStatus = "unpaid"
	PaymentPaid   TicketPaymentStatus = "paid"
)

type Ticket struct {
	ID                    string              `db:"id" json:"id"`
	CustomerName          string              `db:"customer_name" json:"customer_name"`
	CustomerGender        string              `db:"customer_gender" json:"customer_gender"`
	Brand                 string              `db:"brand" json:"brand"`
	Model                 string              `db:"model" json:"model"`
	Issue                 string              `db:"issue" json:"issue"`
	AdditionalDescription *string             `db:"additional_description" json:"additional_description"`
	Accessories           *string             `db:"accessories" json:"accessories"`
	Price                 decimal.Decimal     `db:"price" json:"price"`
	Status                TicketStatus        `db:"status" json:"status"`
	PaymentStatus         TicketPaymentStatus `db:"payment_status" json:"payment_status"`
	WarrantyDays          int                 `db:"warranty_days" json:"warranty_days"`
	EntryDate             time.Time           `db:"entry_date" json:"entry_date"`
	ExitDate              *time.Time          `db:"exit_date" json:"exit_date"`
}

type TicketUpdate struct {
	CustomerName          *string
	CustomerGender        *string
	Brand                 *string
	Model                 *string
	Issue                 *string
	AdditionalDescription *string
	Accessories           *string
	Price                 *decimal.Decimal
	Status                *string
	PaymentStatus         *string
	WarrantyDays          *int
	ExitDate              *time.Time
}

func ValidateTicketUpdate(update TicketUpdate) error {
	if update.Price != nil && update.Price.IsNegative() {
		return ErrNegativePrice
	}
	if update.WarrantyDays != nil && *update.WarrantyDays < 0 {
		return ErrNegativeWarranty
	}
	return nil
}

func (t *Ticket) PrepareForCreate() error {
	if err := t.validatePrice(); err != nil {
		return err
	}
	if err := t.validateWarrantyDays(); err != nil {
		return err
	}
	t.applyDefaultWarrantyDays()
	return nil
}

func (t *Ticket) ApplyUpdate(update TicketUpdate) error {
	t.applyBasicFields(update)

	if err := t.applyPricing(update.Price); err != nil {
		return err
	}
	if err := t.applyWarranty(update.WarrantyDays); err != nil {
		return err
	}

	t.applyStatusAndPaymentChanges(update.Status, update.PaymentStatus, update.ExitDate)

	return t.validateLifecycleInvariants()
}

func (t *Ticket) WarrantyExpiryDate() *time.Time {
	if t.ExitDate == nil {
		return nil
	}
	expiry := t.ExitDate.AddDate(0, 0, t.WarrantyDays)
	return &expiry
}

func (t *Ticket) applyBasicFields(update TicketUpdate) {
	if update.CustomerName != nil {
		t.CustomerName = *update.CustomerName
	}
	if update.CustomerGender != nil {
		t.CustomerGender = *update.CustomerGender
	}
	if update.Brand != nil {
		t.Brand = *update.Brand
	}
	if update.Model != nil {
		t.Model = *update.Model
	}
	if update.Issue != nil {
		t.Issue = *update.Issue
	}
	if update.AdditionalDescription != nil {
		t.AdditionalDescription = optionalText(update.AdditionalDescription)
	}
	if update.Accessories != nil {
		t.Accessories = optionalText(update.Accessories)
	}
}

func (t *Ticket) applyPricing(price *decimal.Decimal) error {
	if price == nil {
		return nil
	}
	if price.IsNegative() {
		return ErrNegativePrice
	}
	t.Price = *price
	return nil
}

func (t *Ticket) applyWarranty(warrantyDays *int) error {
	if warrantyDays == nil {
		return nil
	}
	if *warrantyDays < 0 {
		return ErrNegativeWarranty
	}
	t.WarrantyDays = *warrantyDays
	return nil
}
func (t *Ticket) applyStatusAndPaymentChanges(status *string, paymentStatus *string, exitDate *time.Time) {
	oldStatus := t.Status

	if status != nil {
		t.Status = TicketStatus(*status)
	}

	if paymentStatus != nil {
		t.PaymentStatus = TicketPaymentStatus(*paymentStatus)
	} else if status != nil && t.Status == StatusPickedUp && oldStatus != StatusPickedUp {
		t.PaymentStatus = PaymentPaid
	}

	if exitDate != nil {
		t.ExitDate = exitDate
	} else if status != nil && t.Status == StatusPickedUp && oldStatus != StatusPickedUp {
		now := time.Now().UTC()
		t.ExitDate = &now
	} else if status != nil && t.Status != StatusPickedUp && oldStatus == StatusPickedUp {
		t.ExitDate = nil
	}
}

func (t *Ticket) validateLifecycleInvariants() error {
	if t.Status == StatusPickedUp {
		if t.PaymentStatus != PaymentPaid {
			return ErrPickedUpRequiresPaid
		}
		if t.ExitDate == nil {
			return ErrPickedUpRequiresExitDate
		}
		return nil
	}

	if t.ExitDate != nil {
		return ErrNonPickedUpCannotHaveExitDate
	}
	return nil
}

func (t *Ticket) validatePrice() error {
	if t.Price.IsNegative() {
		return ErrNegativePrice
	}
	return nil
}

func (t *Ticket) validateWarrantyDays() error {
	if t.WarrantyDays < 0 {
		return ErrNegativeWarranty
	}
	return nil
}

func (t *Ticket) applyDefaultWarrantyDays() {
	if t.WarrantyDays <= 0 {
		t.WarrantyDays = DefaultWarrantyDays
	}
}

func optionalText(value *string) *string {
	if *value == "" {
		return nil
	}
	return value
}
