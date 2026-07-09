package ticket

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/internal/events"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrInvalidInput   = errors.New("invalid input data")
)

type Service interface {
	CreateTicket(ctx context.Context, req CreateTicketRequest) (*TicketResponse, error)
	GetTickets(ctx context.Context, status, search string, limit, offset int) ([]TicketListResponse, int, error)
	GetTicketByID(ctx context.Context, id string) (*TicketResponse, error)
	UpdateTicketStatus(ctx context.Context, id string, req ChangeStatusRequest) (*TicketStatusResponse, error)
	UpdateTicketDetails(ctx context.Context, id string, req UpdateTicketRequest) (*TicketResponse, error)
	EmergencyUpdateTicket(ctx context.Context, id string, req EmergencyUpdateTicketRequest) (*TicketResponse, error)
}

type service struct {
	repo     Repository
	eventBus events.EventBus
}

func NewService(repo Repository, eventBus events.EventBus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (s *service) CreateTicket(ctx context.Context, req CreateTicketRequest) (*TicketResponse, error) {
	if req.CustomerName == "" || req.CustomerPhone == "" || req.DeviceBrand == "" || req.DeviceModel == "" || req.IssueDescription == "" {
		return nil, fmt.Errorf("%w: customer name, phone, device brand, model, and issue description are required", ErrInvalidInput)
	}

	ticketNum, err := s.generateTicketNumber()
	if err != nil {
		return nil, err
	}

	ticket := &models.ServiceTicket{
		ID:               uuid.New().String(),
		TicketNumber:     ticketNum,
		Status:           models.StatusReceived,
		CustomerName:     req.CustomerName,
		CustomerPhone:    req.CustomerPhone,
		DeviceBrand:      req.DeviceBrand,
		DeviceModel:      req.DeviceModel,
		DevicePasscode:   req.DevicePasscode,
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

func (s *service) GetTicketByID(ctx context.Context, id string) (*TicketResponse, error) {
	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrTicketNotFound
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

	ticket.CustomerName = req.CustomerName
	ticket.CustomerPhone = req.CustomerPhone
	ticket.DeviceBrand = req.DeviceBrand
	ticket.DeviceModel = req.DeviceModel
	ticket.DevicePasscode = req.DevicePasscode
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

	if ticket.Status == models.StatusCompleted && ticket.WarrantyDays > 0 {
		_ = s.eventBus.Publish(ctx, events.TicketCompletedEvent{
			TicketID:     ticket.ID,
			WarrantyDays: ticket.WarrantyDays,
			CompletedAt:  ticket.UpdatedAt,
		})
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
