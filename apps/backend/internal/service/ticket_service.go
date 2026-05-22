package service

import (
	"context"
	"errors"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
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
	if ticket.WarrantyDays <= 0 {
		ticket.WarrantyDays = 30 // Default to 30 days
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

func (s *ticketService) UpdateTicket(ctx context.Context, id string, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	// Apply updates
	if req.CustomerName != nil {
		ticket.CustomerName = *req.CustomerName
	}
	if req.CustomerGender != nil {
		ticket.CustomerGender = *req.CustomerGender
	}
	if req.Brand != nil {
		ticket.Brand = *req.Brand
	}
	if req.Model != nil {
		ticket.Model = *req.Model
	}
	if req.Issue != nil {
		ticket.Issue = *req.Issue
	}
	if req.AdditionalDescription != nil {
		if *req.AdditionalDescription == "" {
			ticket.AdditionalDescription = nil
		} else {
			ticket.AdditionalDescription = req.AdditionalDescription
		}
	}
	if req.Accessories != nil {
		if *req.Accessories == "" {
			ticket.Accessories = nil
		} else {
			ticket.Accessories = req.Accessories
		}
	}
	if req.Price != nil {
		ticket.Price = *req.Price
	}
	if req.WarrantyDays != nil {
		ticket.WarrantyDays = *req.WarrantyDays
	}

	// Status change logic
	if req.Status != nil && *req.Status != ticket.Status {
		ticket.Status = *req.Status
		if *req.Status == "picked_up" {
			// When picked up, it is the payment moment
			ticket.PaymentStatus = "paid"
			now := time.Now()
			ticket.ExitDate = &now
			expiry := now.AddDate(0, 0, ticket.WarrantyDays)
			ticket.WarrantyExpiryDate = &expiry
		} else {
			// Transitioning to any other status clears exit and warranty dates
			ticket.ExitDate = nil
			ticket.WarrantyExpiryDate = nil
		}
	}
	if req.PaymentStatus != nil {
		ticket.PaymentStatus = *req.PaymentStatus
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	return s.mapToResponse(ticket), nil
}

func (s *ticketService) ListTickets(ctx context.Context) ([]dto.TicketResponse, error) {
	tickets, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
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
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTicketNotFound
		}
		return err
	}
	return nil
}

func (s *ticketService) mapToResponse(ticket *model.Ticket) *dto.TicketResponse {
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
		Status:                ticket.Status,
		PaymentStatus:         ticket.PaymentStatus,
		WarrantyDays:          ticket.WarrantyDays,
		EntryDate:             ticket.EntryDate,
		ExitDate:              ticket.ExitDate,
		WarrantyExpiryDate:    ticket.WarrantyExpiryDate,
	}
}
