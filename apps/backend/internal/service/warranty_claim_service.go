package service

import (
	"context"
	"errors"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type WarrantyClaimService interface {
	CreateClaim(ctx context.Context, req *dto.CreateWarrantyClaimRequest) (*dto.WarrantyClaimResponse, error)
	ListClaims(ctx context.Context, status string) ([]*dto.WarrantyClaimResponse, error)
	ApproveClaim(ctx context.Context, id string) (*dto.ClaimCreationResult, error)
	VoidClaim(ctx context.Context, id string, req *dto.VoidWarrantyClaimRequest) (*dto.ClaimCreationResult, error)
}

type warrantyClaimService struct {
	claimRepo  repository.WarrantyClaimRepository
	ticketRepo repository.TicketRepository
	validate   *validator.Validate
}

func NewWarrantyClaimService(claimRepo repository.WarrantyClaimRepository, ticketRepo repository.TicketRepository) WarrantyClaimService {
	return &warrantyClaimService{
		claimRepo:  claimRepo,
		ticketRepo: ticketRepo,
		validate:   validator.New(),
	}
}

func (s *warrantyClaimService) CreateClaim(ctx context.Context, req *dto.CreateWarrantyClaimRequest) (*dto.WarrantyClaimResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, MapRepositoryError(err)
	}

	if ticket.ExitDate == nil {
		return nil, ErrTicketNotPickedUp
	}
	expiry := ticket.ExitDate.AddDate(0, 0, ticket.WarrantyDays)
	if time.Now().UTC().After(expiry) {
		return nil, ErrWarrantyExpired
	}

	existing, err := s.claimRepo.GetOpenClaimByTicketID(ctx, req.TicketID)
	if err != nil {
		return nil, MapRepositoryError(err)
	}
	if existing != nil {
		return nil, ErrDuplicateWarrantyClaim
	}

	claim := &model.WarrantyClaim{
		TicketID: req.TicketID,
		Issue:    req.Issue,
		Status:   model.ClaimWaitingInspection,
	}
	if req.AdditionalDescription != "" {
		claim.AdditionalDescription = &req.AdditionalDescription
	}

	if err := s.claimRepo.Create(ctx, claim); err != nil {
		return nil, MapRepositoryError(err)
	}

	return s.mapToResponse(claim, ticket), nil
}

func (s *warrantyClaimService) ListClaims(ctx context.Context, status string) ([]*dto.WarrantyClaimResponse, error) {
	claims, err := s.claimRepo.List(ctx, status)
	if err != nil {
		return nil, MapRepositoryError(err)
	}

	responses := make([]*dto.WarrantyClaimResponse, len(claims))
	for i, c := range claims {
		ticket, err := s.ticketRepo.GetByID(ctx, c.TicketID)
		if err != nil {
			return nil, MapRepositoryError(err)
		}
		responses[i] = s.mapToResponse(c, ticket)
	}
	return responses, nil
}

func (s *warrantyClaimService) ApproveClaim(ctx context.Context, id string) (*dto.ClaimCreationResult, error) {
	tx, err := s.claimRepo.BeginTx(ctx)
	if err != nil {
		return nil, MapRepositoryError(err)
	}
	defer tx.Rollback()

	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrWarrantyClaimNotFound
		}
		return nil, MapRepositoryError(err)
	}

	if claim.Status != model.ClaimWaitingInspection {
		return nil, ErrClaimAlreadyDecided
	}

	ticket, err := s.ticketRepo.GetByID(ctx, claim.TicketID)
	if err != nil {
		return nil, MapRepositoryError(err)
	}

	warrantyTicket := &model.Ticket{
		CustomerName:   ticket.CustomerName,
		CustomerGender: ticket.CustomerGender,
		Brand:          ticket.Brand,
		Model:          ticket.Model,
		Issue:          "[Klaim Garansi] " + claim.Issue,
		Accessories:    ticket.Accessories,
		Price:          decimal.Zero,
		Status:         model.StatusOnProcess,
		PaymentStatus:  model.PaymentPaid,
		WarrantyDays:   model.DefaultWarrantyDays,
		IsWarranty:     true,
		ParentTicketID: &ticket.ID,
	}
	if claim.AdditionalDescription != nil {
		warrantyTicket.AdditionalDescription = claim.AdditionalDescription
	}

	if err := warrantyTicket.PrepareForCreate(); err != nil {
		return nil, MapModelError(err)
	}

	if err := s.ticketRepo.CreateTx(ctx, tx, warrantyTicket); err != nil {
		return nil, MapRepositoryError(err)
	}

	now := time.Now().UTC()
	claim.Status = model.ClaimApproved
	claim.ClaimTicketID = &warrantyTicket.ID
	claim.InspectedAt = &now

	if err := s.claimRepo.UpdateTx(ctx, tx, claim); err != nil {
		return nil, MapRepositoryError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, MapRepositoryError(err)
	}

	return &dto.ClaimCreationResult{
		Claim:  *s.mapToResponse(claim, ticket),
		Ticket: MapTicketToResponse(warrantyTicket),
	}, nil
}

func (s *warrantyClaimService) VoidClaim(ctx context.Context, id string, req *dto.VoidWarrantyClaimRequest) (*dto.ClaimCreationResult, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	tx, err := s.claimRepo.BeginTx(ctx)
	if err != nil {
		return nil, MapRepositoryError(err)
	}
	defer tx.Rollback()

	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrWarrantyClaimNotFound
		}
		return nil, MapRepositoryError(err)
	}

	if claim.Status != model.ClaimWaitingInspection {
		return nil, ErrClaimAlreadyDecided
	}

	ticket, err := s.ticketRepo.GetByID(ctx, claim.TicketID)
	if err != nil {
		return nil, MapRepositoryError(err)
	}

	voidTicket := &model.Ticket{
		CustomerName:   ticket.CustomerName,
		CustomerGender: ticket.CustomerGender,
		Brand:          ticket.Brand,
		Model:          ticket.Model,
		Issue:          "[Klaim Ditolak] " + claim.Issue,
		Accessories:    ticket.Accessories,
		Price:          decimal.Zero,
		Status:         model.StatusCancelled,
		PaymentStatus:  model.PaymentPaid,
		WarrantyDays:   0,
		IsWarranty:     true,
		ParentTicketID: &ticket.ID,
	}
	additionalDesc := "Klaim Garansi Ditolak. Alasan: " + req.VoidReason
	voidTicket.AdditionalDescription = &additionalDesc

	if err := voidTicket.PrepareForCreate(); err != nil {
		return nil, MapModelError(err)
	}
	voidTicket.WarrantyDays = 0

	if err := s.ticketRepo.CreateTx(ctx, tx, voidTicket); err != nil {
		return nil, MapRepositoryError(err)
	}

	now := time.Now().UTC()
	claim.Status = model.ClaimVoid
	claim.VoidReason = &req.VoidReason
	claim.ClaimTicketID = &voidTicket.ID
	claim.InspectedAt = &now

	if err := s.claimRepo.UpdateTx(ctx, tx, claim); err != nil {
		return nil, MapRepositoryError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, MapRepositoryError(err)
	}

	return &dto.ClaimCreationResult{
		Claim:  *s.mapToResponse(claim, ticket),
		Ticket: MapTicketToResponse(voidTicket),
	}, nil
}

func (s *warrantyClaimService) mapToResponse(claim *model.WarrantyClaim, ticket *model.Ticket) *dto.WarrantyClaimResponse {
	return &dto.WarrantyClaimResponse{
		ID:                    claim.ID,
		TicketID:              claim.TicketID,
		ClaimTicketID:         claim.ClaimTicketID,
		Issue:                 claim.Issue,
		AdditionalDescription: claim.AdditionalDescription,
		Status:                string(claim.Status),
		VoidReason:            claim.VoidReason,
		InspectedAt:           claim.InspectedAt,
		CreatedAt:             claim.CreatedAt,
		UpdatedAt:             claim.UpdatedAt,
		OriginalTicket:        MapTicketToResponse(ticket),
	}
}
