package warranty

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/google/uuid"
)

type Generator interface {
	CreateWarranty(ctx context.Context, ticketID string, warrantyDays int) (*models.Warranty, error)
}

type generator struct {
	commandRepo CommandRepository
}

func NewGenerator(commandRepo CommandRepository) Generator {
	return &generator{
		commandRepo: commandRepo,
	}
}

func (g *generator) CreateWarranty(ctx context.Context, ticketID string, warrantyDays int) (*models.Warranty, error) {
	if ticketID == "" || warrantyDays <= 0 {
		return nil, fmt.Errorf("%w: ticketID and warrantyDays > 0 are required", ErrInvalidInput)
	}

	wID := uuid.New().String()

	now := time.Now()
	w := &models.Warranty{
		ID:        wID,
		TicketID:  ticketID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, warrantyDays),
		Status:    models.WarrantyStatusActive,
	}

	if err := g.commandRepo.CreateWarranty(ctx, w); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty created successfully",
		slog.String("warranty_id", w.ID),
		slog.String("ticket_id", w.TicketID),
		slog.Time("end_date", w.EndDate),
	)

	return w, nil
}
