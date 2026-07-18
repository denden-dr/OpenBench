package viewmodels

import (
	"time"

	"github.com/denden-dr/OpenBench/internal/models"
)

// TicketListVM represents a single row in the tickets table.
type TicketListVM struct {
	TicketID      string
	TicketNumber  string
	Status        models.ServiceTicketStatus
	CustomerName  string
	CustomerPhone string
	DeviceBrand   string
	DeviceModel   string
	CreatedAt     time.Time
}

// TicketDetailVM represents the full ticket details for the drawer.
type TicketDetailVM struct {
	TicketID         string
	TicketNumber     string
	Status           models.ServiceTicketStatus
	CustomerName     string
	CustomerPhone    string
	DeviceBrand      string
	DeviceModel      string
	DevicePasscode   string
	IssueDescription string
	RepairAction     string
	Cost             int64
	WarrantyDays     int
	Notes            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
