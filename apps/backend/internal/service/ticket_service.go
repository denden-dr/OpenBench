package service

import (
    "context"

    "github.com/denden-dr/openbench/apps/backend/internal/dto"
    "github.com/denden-dr/openbench/apps/backend/internal/model"
    "github.com/denden-dr/openbench/apps/backend/internal/repository"
    "github.com/go-playground/validator/v10"
)

type TicketService interface {
    CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*dto.TicketResponse, error)
    GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error)
}

type ticketService struct {
    repo      repository.TicketRepository
    validate  *validator.Validate
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
        DiagnosisFee:     req.DiagnosisFee,
    }

    if err := s.repo.Create(ctx, ticket); err != nil {
        return nil, err
    }

    return s.mapToResponse(ticket), nil
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error) {
    ticket, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return s.mapToResponse(ticket), nil
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
