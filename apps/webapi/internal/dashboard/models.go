package dashboard

import (
	"time"

	"github.com/denden-dr/OpenBench/internal/models"
)

type DashboardMetrics struct {
	ActiveTickets    int     `json:"active_tickets"`
	PendingDiagnoses int     `json:"pending_diagnoses"`
	SalesToday       float64 `json:"sales_today"`
	ActiveWarranties int     `json:"active_warranties"`
}

type RecentTicket struct {
	TicketID     string                     `db:"id" json:"ticket_id"`
	TicketNumber string                     `db:"ticket_number" json:"ticket_number"`
	Status       models.ServiceTicketStatus `db:"status" json:"status"`
	CustomerName string                     `db:"customer_name" json:"customer_name"`
	DeviceBrand  string                     `db:"device_brand" json:"device_brand"`
	DeviceModel  string                     `db:"device_model" json:"device_model"`
	Cost         int64                      `db:"cost" json:"cost"`
	CreatedAt    time.Time                  `db:"created_at" json:"created_at"`
}

type DashboardResponse struct {
	Metrics       DashboardMetrics `json:"metrics"`
	RecentTickets []RecentTicket   `json:"recent_tickets"`
}
