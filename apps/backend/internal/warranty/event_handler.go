package warranty

import (
	"context"
	"fmt"
	"log/slog"

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

	slog.InfoContext(ctx, "TicketCompletedEvent received",
		slog.String("ticket_id", completedEvt.TicketID),
		slog.Int("warranty_days", completedEvt.WarrantyDays),
	)

	_, err := h.service.CreateWarranty(ctx, completedEvt.TicketID, completedEvt.WarrantyDays)
	if err != nil {
		slog.ErrorContext(ctx, "Error generating warranty", slog.String("ticket_id", completedEvt.TicketID), slog.Any("error", err))
		return err
	}

	slog.InfoContext(ctx, "Warranty successfully created", slog.String("ticket_id", completedEvt.TicketID))
	return nil
}
