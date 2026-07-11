package ticket

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/internal/events"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/google/uuid"
	"log/slog"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrInvalidInput   = errors.New("invalid input data")
)

type Service interface {
	CreateTicket(ctx context.Context, req CreateTicketRequest) (*TicketResponse, error)
	GetTickets(ctx context.Context, status, search string, limit, offset int) ([]TicketListResponse, int, error)
	SearchTickets(ctx context.Context, req TicketSearchRequest) ([]TicketListResponse, int, error)
	GetTicketByID(ctx context.Context, id string) (*TicketResponse, error)
	UpdateTicketStatus(ctx context.Context, id string, req ChangeStatusRequest) (*TicketStatusResponse, error)
	UpdateTicketDetails(ctx context.Context, id string, req UpdateTicketRequest) (*TicketResponse, error)
	EmergencyUpdateTicket(ctx context.Context, id string, req EmergencyUpdateTicketRequest) (*TicketResponse, error)
}

type service struct {
	repo          Repository
	eventBus      events.EventBus
	encryptionKey string
}

func NewService(repo Repository, eventBus events.EventBus, encryptionKey string) Service {
	return &service{repo: repo, eventBus: eventBus, encryptionKey: encryptionKey}
}

func (s *service) CreateTicket(ctx context.Context, req CreateTicketRequest) (*TicketResponse, error) {
	if req.CustomerName == "" || req.CustomerPhone == "" || req.DeviceBrand == "" || req.DeviceModel == "" || req.IssueDescription == "" {
		return nil, fmt.Errorf("%w: customer name, phone, device brand, model, and issue description are required", ErrInvalidInput)
	}

	ticketNum, err := s.generateTicketNumber()
	if err != nil {
		return nil, err
	}

	encryptedPasscode := ""
	if req.DevicePasscode != "" {
		enc, err := utils.Encrypt(req.DevicePasscode, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt passcode: %w", err)
		}
		encryptedPasscode = enc
	}

	ticket := &models.ServiceTicket{
		ID:               uuid.New().String(),
		TicketNumber:     ticketNum,
		Status:           models.StatusReceived,
		CustomerName:     req.CustomerName,
		CustomerPhone:    req.CustomerPhone,
		DeviceBrand:      req.DeviceBrand,
		DeviceModel:      req.DeviceModel,
		DevicePasscode:   encryptedPasscode,
		IssueDescription: req.IssueDescription,
		Cost:             req.Cost,
		WarrantyDays:     req.WarrantyDays,
	}

	if req.RepairAction != "" {
		ticket.RepairAction = &req.RepairAction
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Service ticket created",
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("customer", ticket.CustomerName),
		slog.String("brand", ticket.DeviceBrand),
		slog.String("model", ticket.DeviceModel),
	)

	// Decrypt for client response
	if ticket.DevicePasscode != "" {
		dec, err := utils.Decrypt(ticket.DevicePasscode, s.encryptionKey)
		if err == nil {
			ticket.DevicePasscode = dec
		}
	}

	res := MapToTicketResponse(ticket)
	return &res, nil
}

func (s *service) GetTickets(ctx context.Context, status, search string, limit, offset int) ([]TicketListResponse, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	tickets, total, err := s.repo.FindAll(ctx, status, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var res []TicketListResponse
	for _, t := range tickets {
		res = append(res, MapToTicketListResponse(t))
	}

	return res, total, nil
}

func (s *service) SearchTickets(ctx context.Context, req TicketSearchRequest) ([]TicketListResponse, int, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	tickets, total, err := s.repo.Search(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	var res []TicketListResponse
	for _, t := range tickets {
		res = append(res, MapToTicketListResponse(t))
	}

	return res, total, nil
}

func (s *service) GetTicketByID(ctx context.Context, id string) (*TicketResponse, error) {
	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrTicketNotFound
	}

	if ticket.DevicePasscode != "" {
		dec, err := utils.Decrypt(ticket.DevicePasscode, s.encryptionKey)
		if err == nil {
			ticket.DevicePasscode = dec
		} else {
			slog.ErrorContext(ctx, "Failed to decrypt device passcode", slog.String("error", err.Error()))
		}
	}

	res := MapToTicketResponse(ticket)
	return &res, nil
}

func (s *service) UpdateTicketStatus(ctx context.Context, id string, req ChangeStatusRequest) (*TicketStatusResponse, error) {
	if req.Status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	// Validate status enum
	switch req.Status {
	case models.StatusReceived, models.StatusRepairing, models.StatusPendingConfirmation,
		models.StatusFixed, models.StatusCompleted, models.StatusCancelled, models.StatusReturned:
		// Valid
	default:
		return nil, fmt.Errorf("%w: invalid ticket status", ErrInvalidInput)
	}

	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrTicketNotFound
	}

	ticket.Status = req.Status
	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Service ticket status updated",
		slog.String("ticket_id", ticket.ID),
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("status", string(ticket.Status)),
	)

	if ticket.Status == models.StatusCompleted && ticket.WarrantyDays > 0 {
		_ = s.eventBus.Publish(ctx, events.TicketCompletedEvent{
			TicketID:     ticket.ID,
			WarrantyDays: ticket.WarrantyDays,
			CompletedAt:  ticket.UpdatedAt,
		})
	}

	return &TicketStatusResponse{
		TicketID:  ticket.ID,
		Status:    ticket.Status,
		UpdatedAt: ticket.UpdatedAt,
	}, nil
}

func (s *service) UpdateTicketDetails(ctx context.Context, id string, req UpdateTicketRequest) (*TicketResponse, error) {
	if req.CustomerName == "" || req.CustomerPhone == "" || req.IssueDescription == "" {
		return nil, fmt.Errorf("%w: customer name, phone, and issue description are required", ErrInvalidInput)
	}

	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrTicketNotFound
	}

	ticket.CustomerName = req.CustomerName
	ticket.CustomerPhone = req.CustomerPhone
	ticket.IssueDescription = req.IssueDescription
	ticket.Cost = req.Cost
	ticket.WarrantyDays = req.WarrantyDays

	if req.RepairAction != "" {
		ticket.RepairAction = &req.RepairAction
	} else {
		ticket.RepairAction = nil
	}

	if req.Notes != "" {
		ticket.Notes = &req.Notes
	} else {
		ticket.Notes = nil
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Service ticket details updated",
		slog.String("ticket_id", ticket.ID),
		slog.String("ticket_number", ticket.TicketNumber),
	)

	if ticket.DevicePasscode != "" {
		dec, err := utils.Decrypt(ticket.DevicePasscode, s.encryptionKey)
		if err == nil {
			ticket.DevicePasscode = dec
		}
	}

	res := MapToTicketResponse(ticket)
	return &res, nil
}

func (s *service) EmergencyUpdateTicket(ctx context.Context, id string, req EmergencyUpdateTicketRequest) (*TicketResponse, error) {
	if req.CustomerName == "" || req.CustomerPhone == "" || req.DeviceBrand == "" || req.DeviceModel == "" || req.IssueDescription == "" || req.Status == "" {
		return nil, fmt.Errorf("%w: customer name, phone, device brand, model, issue description, and status are required", ErrInvalidInput)
	}

	// Validate status enum
	switch req.Status {
	case models.StatusReceived, models.StatusRepairing, models.StatusPendingConfirmation,
		models.StatusFixed, models.StatusCompleted, models.StatusCancelled, models.StatusReturned:
		// Valid
	default:
		return nil, fmt.Errorf("%w: invalid ticket status", ErrInvalidInput)
	}

	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrTicketNotFound
	}

	encryptedPasscode := ""
	if req.DevicePasscode != "" {
		enc, err := utils.Encrypt(req.DevicePasscode, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt passcode: %w", err)
		}
		encryptedPasscode = enc
	}

	ticket.CustomerName = req.CustomerName
	ticket.CustomerPhone = req.CustomerPhone
	ticket.DeviceBrand = req.DeviceBrand
	ticket.DeviceModel = req.DeviceModel
	ticket.DevicePasscode = encryptedPasscode
	ticket.Status = req.Status
	ticket.IssueDescription = req.IssueDescription
	ticket.Cost = req.Cost
	ticket.WarrantyDays = req.WarrantyDays

	if req.RepairAction != "" {
		ticket.RepairAction = &req.RepairAction
	} else {
		ticket.RepairAction = nil
	}

	if req.Notes != "" {
		ticket.Notes = &req.Notes
	} else {
		ticket.Notes = nil
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Emergency ticket update performed",
		slog.String("ticket_id", ticket.ID),
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("status", string(ticket.Status)),
	)

	if ticket.Status == models.StatusCompleted && ticket.WarrantyDays > 0 {
		_ = s.eventBus.Publish(ctx, events.TicketCompletedEvent{
			TicketID:     ticket.ID,
			WarrantyDays: ticket.WarrantyDays,
			CompletedAt:  ticket.UpdatedAt,
		})
	}

	if ticket.DevicePasscode != "" {
		dec, err := utils.Decrypt(ticket.DevicePasscode, s.encryptionKey)
		if err == nil {
			ticket.DevicePasscode = dec
		}
	}

	res := MapToTicketResponse(ticket)
	return &res, nil
}

func (s *service) generateTicketNumber() (string, error) {
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
