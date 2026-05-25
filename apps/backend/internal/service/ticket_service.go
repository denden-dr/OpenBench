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
	UpdateTicket(ctx context.Context, id string, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error)
	ListTickets(ctx context.Context) ([]dto.TicketResponse, error)
	DeleteTicket(ctx context.Context, id string) error
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
		CustomerName:   req.CustomerName,
		CustomerGender: req.CustomerGender,
		Brand:          req.Brand,
		Model:          req.Model,
		Issue:          req.Issue,
		Price:          req.Price,
		WarrantyDays:   req.WarrantyDays,
	}
	if req.AdditionalDescription != "" {
		ticket.AdditionalDescription = &req.AdditionalDescription
	}
	if req.Accessories != "" {
		ticket.Accessories = &req.Accessories
	}

	if err := ticket.PrepareForCreate(); err != nil {
		return nil, MapModelError(err)
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, MapRepositoryError(err)
	}

	return s.mapToResponse(ticket), nil
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, MapRepositoryError(err)
	}
	return s.mapToResponse(ticket), nil
}

func (s *ticketService) UpdateTicket(ctx context.Context, id string, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	update := model.TicketUpdate{
		CustomerName:          req.CustomerName,
		CustomerGender:        req.CustomerGender,
		Brand:                 req.Brand,
		Model:                 req.Model,
		Issue:                 req.Issue,
		AdditionalDescription: req.AdditionalDescription,
		Accessories:           req.Accessories,
		Price:                 req.Price,
		Status:                req.Status,
		PaymentStatus:         req.PaymentStatus,
		WarrantyDays:          req.WarrantyDays,
		ExitDate:              req.ExitDate,
	}
	if err := model.ValidateTicketUpdate(update); err != nil {
		return nil, MapModelError(err)
	}

	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, MapRepositoryError(err)
	}

	if err := ticket.ApplyUpdate(update); err != nil {
		return nil, MapModelError(err)
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, MapRepositoryError(err)
	}

	return s.mapToResponse(ticket), nil
}

func (s *ticketService) ListTickets(ctx context.Context) ([]dto.TicketResponse, error) {
	tickets, err := s.repo.List(ctx)
	if err != nil {
		return nil, MapRepositoryError(err)
	}

	responses := make([]dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		responses[i] = *s.mapToResponse(&t)
	}
	return responses, nil
}

func (s *ticketService) DeleteTicket(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return MapRepositoryError(err)
	}
	return nil
}

func (s *ticketService) mapToResponse(ticket *model.Ticket) *dto.TicketResponse {
	return MapTicketToResponse(ticket)
}

func MapTicketToResponse(ticket *model.Ticket) *dto.TicketResponse {
	if ticket == nil {
		return nil
	}
	return &dto.TicketResponse{
		ID:                    ticket.ID,
		CustomerName:          ticket.CustomerName,
		CustomerGender:        ticket.CustomerGender,
		Brand:                 ticket.Brand,
		Model:                 ticket.Model,
		Issue:                 ticket.Issue,
		AdditionalDescription: ticket.AdditionalDescription,
		Accessories:           ticket.Accessories,
		Price:                 ticket.Price,
		Status:                string(ticket.Status),
		PaymentStatus:         string(ticket.PaymentStatus),
		WarrantyDays:          ticket.WarrantyDays,
		EntryDate:             ticket.EntryDate,
		ExitDate:              ticket.ExitDate,
		WarrantyExpiryDate:    ticket.WarrantyExpiryDate(),
		IsWarranty:            ticket.IsWarranty,
		ParentTicketID:        ticket.ParentTicketID,
	}
}
