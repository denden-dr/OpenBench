package warranty

import (
	"context"
	"fmt"
	"log"

	"github.com/denden-dr/OpenBench/apps/backend/internal/events"
)

type EventHandler struct {
	service Service
}

func NewEventHandler(service Service) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) HandleTicketCompleted(ctx context.Context, event events.Event) error {
	completedEvt, ok := event.(events.TicketCompletedEvent)
	if !ok {
		return fmt.Errorf("invalid event type, expected TicketCompletedEvent, got: %T", event)
	}

	log.Printf("[EventHandler] TicketCompletedEvent received for Ticket ID: %s with %d warranty days", completedEvt.TicketID, completedEvt.WarrantyDays)

	_, err := h.service.CreateWarranty(ctx, completedEvt.TicketID, completedEvt.WarrantyDays)
	if err != nil {
		log.Printf("[EventHandler] Error generating warranty: %v", err)
		return err
	}

	log.Printf("[EventHandler] Warranty successfully created for Ticket ID: %s", completedEvt.TicketID)
	return nil
}
