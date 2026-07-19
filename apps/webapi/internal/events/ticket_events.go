package events

import "time"

const (
	TicketCompletedType EventType = "ticket.completed"
)

type TicketCompletedEvent struct {
	TicketID     string    `json:"ticket_id"`
	WarrantyDays int       `json:"warranty_days"`
	CompletedAt  time.Time `json:"completed_at"`
}

func (e TicketCompletedEvent) Type() EventType {
	return TicketCompletedType
}
