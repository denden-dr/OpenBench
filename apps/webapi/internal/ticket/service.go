package ticket

import (
	"context"
	"errors"
	"fmt"

	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/events"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/utils"
	"log/slog"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrInvalidInput   = errors.New("invalid input data")
)

type WarrantyGenerator interface {
	CreateWarranty(ctx context.Context, ticketID string, warrantyDays int) (*models.Warranty, error)
}

type Service interface {
	CreateTicket(ctx context.Context, req CreateTicketRequest) (TicketResponse, error)
	GetTicketSummaries(ctx context.Context, status, search string, limit int, cursor string) ([]TicketSummaryResponse, string, error)
	SearchTicketSummaries(ctx context.Context, req TicketSearchRequest) ([]TicketSummaryResponse, string, error)
	GetTicketByID(ctx context.Context, id string) (TicketResponse, error)
	UpdateTicketStatus(ctx context.Context, id string, req ChangeStatusRequest) (TicketStatusResponse, error)
	UpdateTicketDetails(ctx context.Context, id string, req UpdateTicketRequest) (TicketResponse, error)
	EmergencyUpdateTicket(ctx context.Context, id string, req EmergencyUpdateTicketRequest) (TicketResponse, error)
}

type service struct {
	queryRepo     QueryRepository
	commandRepo   CommandRepository
	txManager     database.TxManager
	warrantyGen   WarrantyGenerator
	eventBus      events.EventBus
	encryptionKey string
}

func NewService(
	queryRepo QueryRepository,
	commandRepo CommandRepository,
	txManager database.TxManager,
	warrantyGen WarrantyGenerator,
	eventBus events.EventBus,
	encryptionKey string,
) Service {
	return &service{
		queryRepo:     queryRepo,
		commandRepo:   commandRepo,
		txManager:     txManager,
		warrantyGen:   warrantyGen,
		eventBus:      eventBus,
		encryptionKey: encryptionKey,
	}
}

func (s *service) CreateTicket(ctx context.Context, req CreateTicketRequest) (TicketResponse, error) {
	ticketNum, err := GenerateTicketNumber()
	if err != nil {
		return TicketResponse{}, err
	}

	encryptedPasscode, err := s.encryptPasscode(req.DevicePasscode)
	if err != nil {
		return TicketResponse{}, err
	}

	var repairAction *string
	if req.RepairAction != "" {
		repairAction = &req.RepairAction
	}

	ticket, err := models.NewServiceTicket(models.CreateTicketParams{
		TicketNumber:     ticketNum,
		CustomerName:     req.CustomerName,
		CustomerPhone:    req.CustomerPhone,
		DeviceBrand:      req.DeviceBrand,
		DeviceModel:      req.DeviceModel,
		DevicePasscode:   encryptedPasscode,
		IssueDescription: req.IssueDescription,
		RepairAction:     repairAction,
		Cost:             req.Cost,
		WarrantyDays:     req.WarrantyDays,
	})
	if err != nil {
		return TicketResponse{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	if err := s.commandRepo.Create(ctx, ticket); err != nil {
		return TicketResponse{}, err
	}

	slog.InfoContext(ctx, "Service ticket created",
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("customer", ticket.CustomerName),
		slog.String("brand", ticket.DeviceBrand),
		slog.String("model", ticket.DeviceModel),
	)

	ticket.DevicePasscode = s.decryptPasscode(ctx, ticket.DevicePasscode)

	res := MapToTicketResponse(ticket)
	return res, nil
}

func (s *service) GetTicketSummaries(ctx context.Context, status, search string, limit int, cursor string) ([]TicketSummaryResponse, string, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > utils.MaxLimit {
		limit = utils.MaxLimit
	}

	tickets, nextCursor, err := s.queryRepo.FindAll(ctx, status, search, limit, cursor)
	if err != nil {
		return nil, "", err
	}

	var res []TicketSummaryResponse
	for _, t := range tickets {
		res = append(res, MapToTicketSummaryResponse(t))
	}

	return res, nextCursor, nil
}

func (s *service) SearchTicketSummaries(ctx context.Context, req TicketSearchRequest) ([]TicketSummaryResponse, string, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > utils.MaxLimit {
		req.Limit = utils.MaxLimit
	}

	tickets, nextCursor, err := s.queryRepo.Search(ctx, req)
	if err != nil {
		return nil, "", err
	}

	var res []TicketSummaryResponse
	for _, t := range tickets {
		res = append(res, MapToTicketSummaryResponse(t))
	}

	return res, nextCursor, nil
}

func (s *service) GetTicketByID(ctx context.Context, id string) (TicketResponse, error) {
	ticket, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return TicketResponse{}, err
	}
	if ticket == nil {
		return TicketResponse{}, ErrTicketNotFound
	}

	ticket.DevicePasscode = s.decryptPasscode(ctx, ticket.DevicePasscode)

	res := MapToTicketResponse(ticket)
	return res, nil
}

func (s *service) UpdateTicketStatus(ctx context.Context, id string, req ChangeStatusRequest) (TicketStatusResponse, error) {
	if req.Status == "" {
		return TicketStatusResponse{}, fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	if err := validateStatus(req.Status); err != nil {
		return TicketStatusResponse{}, err
	}

	ticket, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return TicketStatusResponse{}, err
	}
	if ticket == nil {
		return TicketStatusResponse{}, ErrTicketNotFound
	}

	if ticket.Status == req.Status {
		return TicketStatusResponse{}, fmt.Errorf("%w: ticket is already in %s status", ErrInvalidInput, req.Status)
	}

	ticket.Status = req.Status

	err = s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.commandRepo.Update(txCtx, ticket); err != nil {
			return err
		}

		if ticket.Status == models.StatusCompleted {
			if err := s.handleTicketCompletion(txCtx, ticket); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return TicketStatusResponse{}, err
	}

	slog.InfoContext(ctx, "Service ticket status updated",
		slog.String("ticket_id", ticket.ID),
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("status", string(ticket.Status)),
	)

	return TicketStatusResponse{
		TicketID:  ticket.ID,
		Status:    ticket.Status,
		UpdatedAt: ticket.UpdatedAt,
	}, nil
}

func (s *service) UpdateTicketDetails(ctx context.Context, id string, req UpdateTicketRequest) (TicketResponse, error) {
	if req.CustomerName == "" || req.CustomerPhone == "" || req.IssueDescription == "" {
		return TicketResponse{}, fmt.Errorf("%w: customer name, phone, and issue description are required", ErrInvalidInput)
	}

	ticket, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return TicketResponse{}, err
	}
	if ticket == nil {
		return TicketResponse{}, ErrTicketNotFound
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

	if err := s.commandRepo.Update(ctx, ticket); err != nil {
		return TicketResponse{}, err
	}

	slog.InfoContext(ctx, "Service ticket details updated",
		slog.String("ticket_id", ticket.ID),
		slog.String("ticket_number", ticket.TicketNumber),
	)

	ticket.DevicePasscode = s.decryptPasscode(ctx, ticket.DevicePasscode)

	res := MapToTicketResponse(ticket)
	return res, nil
}

func (s *service) EmergencyUpdateTicket(ctx context.Context, id string, req EmergencyUpdateTicketRequest) (TicketResponse, error) {
	if req.CustomerName == "" || req.CustomerPhone == "" || req.DeviceBrand == "" || req.DeviceModel == "" || req.IssueDescription == "" || req.Status == "" {
		return TicketResponse{}, fmt.Errorf("%w: customer name, phone, device brand, model, issue description, and status are required", ErrInvalidInput)
	}

	if err := validateStatus(req.Status); err != nil {
		return TicketResponse{}, err
	}

	ticket, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return TicketResponse{}, err
	}
	if ticket == nil {
		return TicketResponse{}, ErrTicketNotFound
	}

	encryptedPasscode, err := s.encryptPasscode(req.DevicePasscode)
	if err != nil {
		return TicketResponse{}, err
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

	err = s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.commandRepo.Update(txCtx, ticket); err != nil {
			return err
		}

		if ticket.Status == models.StatusCompleted {
			if err := s.handleTicketCompletion(txCtx, ticket); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return TicketResponse{}, err
	}

	slog.InfoContext(ctx, "Emergency ticket update performed",
		slog.String("ticket_id", ticket.ID),
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("status", string(ticket.Status)),
	)

	ticket.DevicePasscode = s.decryptPasscode(ctx, ticket.DevicePasscode)

	res := MapToTicketResponse(ticket)
	return res, nil
}

func (s *service) handleTicketCompletion(ctx context.Context, ticket *models.ServiceTicket) error {
	if ticket.WarrantyDays > 0 {
		_, err := s.warrantyGen.CreateWarranty(ctx, ticket.ID, ticket.WarrantyDays)
		if err != nil {
			return fmt.Errorf("failed to create warranty within transaction: %w", err)
		}
	}
	return nil
}

func (s *service) encryptPasscode(passcode string) (string, error) {
	if passcode == "" {
		return "", nil
	}
	enc, err := utils.Encrypt(passcode, s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt passcode: %w", err)
	}
	return enc, nil
}

func (s *service) decryptPasscode(ctx context.Context, passcode string) string {
	if passcode == "" {
		return ""
	}
	dec, err := utils.Decrypt(passcode, s.encryptionKey)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to decrypt device passcode", slog.String("error", err.Error()))
		return passcode
	}
	return dec
}

func validateStatus(status models.ServiceTicketStatus) error {
	switch status {
	case models.StatusReceived, models.StatusRepairing, models.StatusPendingConfirmation,
		models.StatusFixed, models.StatusCompleted, models.StatusCancelled, models.StatusReturned:
		return nil
	default:
		return fmt.Errorf("%w: invalid ticket status", ErrInvalidInput)
	}
}
