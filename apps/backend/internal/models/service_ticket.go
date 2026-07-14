package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ServiceTicketStatus represents the current state of a service ticket
type ServiceTicketStatus string

const (
	StatusReceived            ServiceTicketStatus = "RECEIVED"
	StatusRepairing           ServiceTicketStatus = "REPAIRING"
	StatusPendingConfirmation ServiceTicketStatus = "PENDING_CONFIRMATION"
	StatusFixed               ServiceTicketStatus = "FIXED"
	StatusCompleted           ServiceTicketStatus = "COMPLETED"
	StatusCancelled           ServiceTicketStatus = "CANCELLED"
	StatusReturned            ServiceTicketStatus = "RETURNED"
)

var (
	ErrMissingCustomerName     = errors.New("customer name is required")
	ErrMissingCustomerPhone    = errors.New("customer phone is required")
	ErrMissingDeviceBrand      = errors.New("device brand is required")
	ErrMissingDeviceModel      = errors.New("device model is required")
	ErrMissingIssueDescription = errors.New("issue description is required")
	ErrNegativeCost            = errors.New("cost cannot be negative")
	ErrNegativeWarrantyDays    = errors.New("warranty days cannot be negative")
)

type ServiceTicket struct {
	ID               string              `json:"id" db:"id"`
	TicketNumber     string              `json:"ticket_number" db:"ticket_number"`
	Status           ServiceTicketStatus `json:"status" db:"status"`
	CustomerName     string              `json:"customer_name" db:"customer_name"`
	CustomerPhone    string              `json:"customer_phone" db:"customer_phone"`
	DeviceBrand      string              `json:"device_brand" db:"device_brand"`
	DeviceModel      string              `json:"device_model" db:"device_model"`
	DevicePasscode   string              `json:"device_passcode,omitempty" db:"device_passcode"`
	IssueDescription string              `json:"issue_description" db:"issue_description"`
	RepairAction     *string             `json:"repair_action" db:"repair_action"`
	Cost             int64               `json:"cost" db:"cost"`
	WarrantyDays     int                 `json:"warranty_days" db:"warranty_days"`
	Notes            *string             `json:"notes" db:"notes"`
	CreatedAt        time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time          `json:"-" db:"deleted_at"`
}

type CreateTicketParams struct {
	TicketNumber     string
	CustomerName     string
	CustomerPhone    string
	DeviceBrand      string
	DeviceModel      string
	DevicePasscode   string
	IssueDescription string
	RepairAction     *string
	Cost             int64
	WarrantyDays     int
}

// NewServiceTicket is a factory function that validates invariants and returns a valid ServiceTicket
func NewServiceTicket(params CreateTicketParams) (*ServiceTicket, error) {
	customerName := strings.TrimSpace(params.CustomerName)
	customerPhone := strings.TrimSpace(params.CustomerPhone)
	deviceBrand := strings.TrimSpace(params.DeviceBrand)
	deviceModel := strings.TrimSpace(params.DeviceModel)
	issueDescription := strings.TrimSpace(params.IssueDescription)

	if customerName == "" {
		return nil, ErrMissingCustomerName
	}
	if customerPhone == "" {
		return nil, ErrMissingCustomerPhone
	}
	if deviceBrand == "" {
		return nil, ErrMissingDeviceBrand
	}
	if deviceModel == "" {
		return nil, ErrMissingDeviceModel
	}
	if issueDescription == "" {
		return nil, ErrMissingIssueDescription
	}
	if params.Cost < 0 {
		return nil, ErrNegativeCost
	}
	if params.WarrantyDays < 0 {
		return nil, ErrNegativeWarrantyDays
	}

	return &ServiceTicket{
		ID:               uuid.New().String(),
		TicketNumber:     params.TicketNumber,
		Status:           StatusReceived,
		CustomerName:     customerName,
		CustomerPhone:    customerPhone,
		DeviceBrand:      deviceBrand,
		DeviceModel:      deviceModel,
		DevicePasscode:   params.DevicePasscode,
		IssueDescription: issueDescription,
		RepairAction:     params.RepairAction,
		Cost:             params.Cost,
		WarrantyDays:     params.WarrantyDays,
	}, nil
}
