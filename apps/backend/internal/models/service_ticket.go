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
	ID               string              `json:"id"`
	TicketNumber     string              `json:"ticket_number"`
	Status           ServiceTicketStatus `json:"status"`
	CustomerName     string              `json:"customer_name"`
	CustomerPhone    string              `json:"customer_phone"`
	DeviceBrand      string              `json:"device_brand"`
	DeviceModel      string              `json:"device_model"`
	DevicePasscode   string              `json:"device_passcode,omitempty"`
	IssueDescription string              `json:"issue_description"`
	RepairAction     *string             `json:"repair_action"`
	Cost             int64               `json:"cost"`
	WarrantyDays     int                 `json:"warranty_days"`
	Notes            *string             `json:"notes"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	DeletedAt        *time.Time          `json:"-"`
}
