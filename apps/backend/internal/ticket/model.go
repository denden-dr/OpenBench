package ticket

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Status Constants
const (
	StatusReceived  = "received"
	StatusInRepair  = "in_repair"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

// DevicePosition Constants
const (
	PositionWarehouse      = "warehouse"
	PositionServiceCounter = "service_counter"
	PositionPickedUp       = "picked_up"
)

// PaymentStatus Constants
const (
	PaymentStatusNone = "none"
	PaymentStatusPaid = "paid"
)

// Sentinel Domain Errors
var (
	ErrTicketAlreadyPickedUp     = errors.New("ticket already picked up")
	ErrWarrantyDurationNegative  = errors.New("warranty duration days cannot be negative")
	ErrWarrantyUpdateAfterPickup = errors.New("warranty duration cannot be updated after pickup")
	ErrPaymentRequiredForPickup  = errors.New("payment status must be paid when device is picked up")
	ErrPaymentMethodRequired     = errors.New("payment method must be set to cash or qris when device is picked up")
	ErrMissingRequiredFields     = errors.New("missing required fields")
	ErrStatusReversalNotAllowed  = errors.New("location reversal not allowed under normal updates")
)

// Ticket represents a repair ticket in the domain layer
type Ticket struct {
	ID                   string     `db:"id"`
	TicketNumber         string     `db:"ticket_number"`
	CustomerName         string     `db:"customer_name"`
	CustomerPhone        string     `db:"customer_phone"`
	BrandPhone           string     `db:"brand_phone"`
	ModelPhone           string     `db:"model_phone"`
	SerialNumber         string     `db:"serial_number"`
	DamageDescription    string     `db:"damage_description"`
	RepairAction         string     `db:"repair_action"`
	Cost                 float64    `db:"cost"`
	Status               string     `db:"status"`
	DevicePosition       string     `db:"device_position"`
	PaymentStatus        string     `db:"payment_status"`
	PaymentMethod        *string    `db:"payment_method"`
	WarrantyDurationDays int        `db:"warranty_duration_days"`
	PickedUpAt           *time.Time `db:"picked_up_at"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`

	// Warranty loaded as part of Ticket aggregate
	Warranty *Warranty `db:"-"`
}

// NewTicket constructs a new Ticket instance with validated fields and default status values
func NewTicket(
	ticketNumber string,
	customerName string,
	customerPhone string,
	brandPhone string,
	modelPhone string,
	serialNumber string,
	damageDescription string,
	repairAction string,
	cost float64,
	warrantyDurationDays int,
) (*Ticket, error) {
	if customerName == "" || customerPhone == "" || brandPhone == "" || modelPhone == "" || damageDescription == "" {
		return nil, ErrMissingRequiredFields
	}
	if warrantyDurationDays < 0 {
		return nil, ErrWarrantyDurationNegative
	}

	return &Ticket{
		ID:                   uuid.New().String(),
		TicketNumber:         ticketNumber,
		CustomerName:         customerName,
		CustomerPhone:        customerPhone,
		BrandPhone:           brandPhone,
		ModelPhone:           modelPhone,
		SerialNumber:         serialNumber,
		DamageDescription:    damageDescription,
		RepairAction:         repairAction,
		Cost:                 cost,
		Status:               StatusReceived,
		DevicePosition:       PositionWarehouse,
		PaymentStatus:        PaymentStatusNone,
		PaymentMethod:        nil,
		WarrantyDurationDays: warrantyDurationDays,
		PickedUpAt:           nil,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}

// ProcessPickup transitions the device position to picked_up and generates the warranty
func (t *Ticket) ProcessPickup(now time.Time) error {
	if t.DevicePosition == PositionPickedUp {
		return ErrTicketAlreadyPickedUp
	}

	t.DevicePosition = PositionPickedUp

	// Default to completed if not already cancelled or completed
	if t.Status != StatusCompleted && t.Status != StatusCancelled {
		t.Status = StatusCompleted
	}

	// If it is completed, ensure payment status is paid (or set it if it was none)
	if t.Status == StatusCompleted && t.PaymentStatus != PaymentStatusPaid {
		t.PaymentStatus = PaymentStatusPaid
	}

	// Invariant: Gunakan PickedUpAt yang sudah ada jika tersedia, jika tidak set ke 'now'
	if t.PickedUpAt == nil {
		t.PickedUpAt = &now
	}

	expiryDate := t.PickedUpAt.AddDate(0, 0, t.WarrantyDurationDays)

	if t.Status == StatusCompleted && t.Warranty == nil {
		t.Warranty = &Warranty{
			ID:           uuid.New().String(),
			TicketID:     t.ID,
			TicketNumber: t.TicketNumber,
			CustomerName: t.CustomerName,
			DeviceInfo:   fmt.Sprintf("%s %s", t.BrandPhone, t.ModelPhone),
			StartDate:    *t.PickedUpAt,
			EndDate:      expiryDate,
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}
	return nil
}

// UpdateWarrantyDuration updates the ticket's warranty duration after validating rules
func (t *Ticket) UpdateWarrantyDuration(days int) error {
	if t.DevicePosition == PositionPickedUp {
		return ErrWarrantyUpdateAfterPickup
	}
	if days < 0 {
		return ErrWarrantyDurationNegative
	}
	t.WarrantyDurationDays = days
	return nil
}

// Validate checks business rules for the ticket's current state
func (t *Ticket) Validate() error {
	if t.DevicePosition == PositionPickedUp {
		if t.Status != StatusCompleted && t.Status != StatusCancelled {
			return fmt.Errorf("service status must be completed or cancelled when device is picked up")
		}
		if t.Status == StatusCompleted {
			if t.PaymentStatus != PaymentStatusPaid {
				return ErrPaymentRequiredForPickup
			}
			if t.PaymentMethod == nil || (*t.PaymentMethod != "cash" && *t.PaymentMethod != "qris") {
				return ErrPaymentMethodRequired
			}
		}
	}
	return nil
}

// UpdateStatus handles transition to a new status and runs validations
func (t *Ticket) UpdateStatus(newStatus string, pickupTime time.Time) error {
	t.Status = newStatus
	return t.Validate()
}

// ReversePickupLocation resets the device position to warehouse and voids the associated warranty info
func (t *Ticket) ReversePickupLocation() {
	t.DevicePosition = PositionWarehouse
	t.PickedUpAt = nil
	t.Warranty = nil
}

// EmergencyUpdateWarrantyDuration updates warranty duration under emergency administrative updates
func (t *Ticket) EmergencyUpdateWarrantyDuration(days int) error {
	if days < 0 {
		return ErrWarrantyDurationNegative
	}
	t.WarrantyDurationDays = days
	if t.DevicePosition == PositionPickedUp && t.Warranty != nil && t.PickedUpAt != nil {
		t.Warranty.EndDate = t.PickedUpAt.AddDate(0, 0, days)
	}
	return nil
}

// Warranty represents a repair warranty in the domain layer
type Warranty struct {
	ID           string    `db:"id"`
	TicketID     string    `db:"ticket_id"`
	TicketNumber string    `db:"ticket_number"`
	CustomerName string    `db:"customer_name"`
	DeviceInfo   string    `db:"device_info"`
	StartDate    time.Time `db:"start_date"`
	EndDate      time.Time `db:"end_date"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
