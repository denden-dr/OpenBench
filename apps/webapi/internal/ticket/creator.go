package ticket

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/denden-dr/OpenBench/internal/models"
)

type Creator interface {
	CreateWarrantyTicket(ctx context.Context, originalTicketID string, warrantyID string, issueDescription string) (string, error)
}

type creator struct {
	queryRepo   QueryRepository
	commandRepo CommandRepository
}

func NewCreator(queryRepo QueryRepository, commandRepo CommandRepository) Creator {
	return &creator{
		queryRepo:   queryRepo,
		commandRepo: commandRepo,
	}
}

func GenerateTicketNumber() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	dateStr := time.Now().Format("20060102")
	return fmt.Sprintf("TKT-%s-%s", dateStr, string(b)), nil
}

func (c *creator) CreateWarrantyTicket(ctx context.Context, originalTicketID string, warrantyID string, issueDescription string) (string, error) {
	orig, err := c.queryRepo.FindByID(ctx, originalTicketID)
	if err != nil {
		return "", err
	}
	if orig == nil {
		return "", ErrTicketNotFound
	}

	ticketNum, err := GenerateTicketNumber()
	if err != nil {
		return "", err
	}

	fullIssue := fmt.Sprintf("[Garansi %s] %s", warrantyID, issueDescription)

	t, err := models.NewServiceTicket(models.CreateTicketParams{
		TicketNumber:     ticketNum,
		CustomerName:     orig.CustomerName,
		CustomerPhone:    orig.CustomerPhone,
		DeviceBrand:      orig.DeviceBrand,
		DeviceModel:      orig.DeviceModel,
		DevicePasscode:   orig.DevicePasscode,
		IssueDescription: fullIssue,
		Cost:             0,
		WarrantyDays:     0,
	})
	if err != nil {
		return "", err
	}

	if err := c.commandRepo.Create(ctx, t); err != nil {
		return "", err
	}

	slog.InfoContext(ctx, "Warranty repair ticket created",
		slog.String("ticket_id", t.ID),
		slog.String("ticket_number", t.TicketNumber),
		slog.String("warranty_id", warrantyID),
	)

	return t.ID, nil
}
