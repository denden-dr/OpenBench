package models

import "time"

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
