package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var (
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrClaimConflict     = errors.New("ticket is already claimed or not in received status")
)

type TicketService interface {
	CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*dto.TicketResponse, error)
	GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error)
	ClaimTicket(ctx context.Context, ticketID string, technicianID string) error
	CompleteDiagnosis(ctx context.Context, ticketID string) error
	ApproveRepair(ctx context.Context, ticketID string) error
	CancelRepair(ctx context.Context, ticketID string) error
	CompleteRepair(ctx context.Context, ticketID string) error
	MarkPickedUp(ctx context.Context, ticketID string) error
	ListForBoard(ctx context.Context) ([]dto.TicketBoardDTO, error)
}

type ticketService struct {
	repo     repository.TicketRepository
	validate *validator.Validate
}

func NewTicketService(repo repository.TicketRepository) TicketService {
	return &ticketService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *ticketService) CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*dto.TicketResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	ticket := &model.Ticket{
		DeviceType:       req.DeviceType,
		Brand:            req.Brand,
		Model:            req.Model,
		IssueDescription: req.IssueDescription,
		DiagnosisFee:     decimal.NewFromInt(50), // Standard diagnosis fee
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	return s.mapToResponse(ticket), nil
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return s.mapToResponse(ticket), nil
}

func (s *ticketService) ClaimTicket(ctx context.Context, ticketID string, technicianID string) error {
	if err := s.repo.ClaimTicket(ctx, ticketID, technicianID); err != nil {
		if errors.Is(err, repository.ErrClaimConflict) {
			return ErrClaimConflict
		}
		return err
	}
	return nil
}

// transitionStatus is a private helper that validates and executes a status transition.
func (s *ticketService) transitionStatus(ctx context.Context, ticketID string, expectedCurrent string, newStatus string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTicketNotFound
		}
		return err
	}
	if ticket.Status != expectedCurrent {
		return fmt.Errorf("%w: cannot move from '%s' to '%s'", ErrInvalidTransition, ticket.Status, newStatus)
	}
	return s.repo.UpdateStatus(ctx, ticketID, newStatus)
}

func (s *ticketService) CompleteDiagnosis(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusDiagnosing, model.StatusPendingApproval)
}

func (s *ticketService) ApproveRepair(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusPendingApproval, model.StatusRepairing)
}

func (s *ticketService) CancelRepair(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusPendingApproval, model.StatusCancelled)
}

func (s *ticketService) CompleteRepair(ctx context.Context, ticketID string) error {
	return s.transitionStatus(ctx, ticketID, model.StatusRepairing, model.StatusReady)
}

func (s *ticketService) MarkPickedUp(ctx context.Context, ticketID string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTicketNotFound
		}
		return err
	}
	if ticket.Status != model.StatusReady && ticket.Status != model.StatusCancelled {
		return fmt.Errorf("%w: cannot pick up from '%s'", ErrInvalidTransition, ticket.Status)
	}
	return s.repo.UpdateStatus(ctx, ticketID, model.StatusPickedUp)
}

func (s *ticketService) ListForBoard(ctx context.Context) ([]dto.TicketBoardDTO, error) {
	tickets, err := s.repo.ListForBoard(ctx)
	if err != nil {
		return nil, err
	}

	board := make([]dto.TicketBoardDTO, len(tickets))
	for i, t := range tickets {
		board[i] = dto.TicketBoardDTO{
			ID:         t.ID,
			DeviceType: t.DeviceType,
			Brand:      t.Brand,
			Model:      t.Model,
			Status:     t.Status,
			CreatedAt:  t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return board, nil
}

func (s *ticketService) mapToResponse(ticket *model.Ticket) *dto.TicketResponse {
	return &dto.TicketResponse{
		ID:               ticket.ID,
		DeviceType:       ticket.DeviceType,
		Brand:            ticket.Brand,
		Model:            ticket.Model,
		IssueDescription: ticket.IssueDescription,
		Status:           ticket.Status,
		DiagnosisFee:     ticket.DiagnosisFee,
		CreatedAt:        ticket.CreatedAt,
		UpdatedAt:        ticket.UpdatedAt,
	}
}
