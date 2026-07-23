package warranty

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/utils"
	"github.com/google/uuid"
	"log/slog"
)

var (
	ErrWarrantyNotFound  = errors.New("warranty not found")
	ErrWarrantyNotActive = errors.New("warranty is not active or has expired")
	ErrClaimNotFound     = errors.New("claim not found")
	ErrInvalidInput      = errors.New("invalid input data")
)

type TicketCreator interface {
	CreateWarrantyTicket(ctx context.Context, originalTicketID string, warrantyID string, issueDescription string) (string, error)
}

type Service interface {
	GetWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error)
	GetWarrantyByTicketNumber(ctx context.Context, ticketNumber string) (*models.Warranty, error)
	UpdateWarrantyStatus(ctx context.Context, id string, req UpdateWarrantyStatusRequest) (*models.Warranty, error)

	CreateClaim(ctx context.Context, req CreateClaimRequest) (*models.Claim, error)
	GetClaims(ctx context.Context, status, search string, limit int, cursor string) ([]models.ClaimSummary, string, error)
	GetClaimByID(ctx context.Context, id string) (*models.Claim, error)
	GetClaimSummaryByID(ctx context.Context, id string) (*models.ClaimSummary, error)
	UpdateClaim(ctx context.Context, id string, req UpdateClaimRequest) (*models.Claim, error)
	EvaluateClaim(ctx context.Context, claimID string, req EvaluateClaimRequest) (*models.Claim, error)
}

type service struct {
	queryRepo     QueryRepository
	commandRepo   CommandRepository
	txManager     database.TxManager
	ticketCreator TicketCreator
}

func NewService(
	queryRepo QueryRepository,
	commandRepo CommandRepository,
	txManager database.TxManager,
	ticketCreator TicketCreator,
) Service {
	return &service{
		queryRepo:     queryRepo,
		commandRepo:   commandRepo,
		txManager:     txManager,
		ticketCreator: ticketCreator,
	}
}

func (s *service) GetWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error) {
	w, err := s.queryRepo.FindWarrantyByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	s.autoExpireWarrantyIfNeeded(ctx, w)

	return w, nil
}

func (s *service) GetWarrantyByTicketNumber(ctx context.Context, ticketNumber string) (*models.Warranty, error) {
	w, err := s.queryRepo.FindWarrantyByTicketNumber(ctx, ticketNumber)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	s.autoExpireWarrantyIfNeeded(ctx, w)

	return w, nil
}

func (s *service) UpdateWarrantyStatus(ctx context.Context, id string, req UpdateWarrantyStatusRequest) (*models.Warranty, error) {
	trimmedNotes, err := validateUpdateWarrantyStatus(req)
	if err != nil {
		return nil, err
	}

	w, err := s.queryRepo.FindWarrantyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	var notesPtr *string
	if trimmedNotes != "" {
		notesPtr = &trimmedNotes
	} else if w.Notes != nil {
		notesPtr = w.Notes
	}

	err = s.commandRepo.UpdateWarrantyStatus(ctx, id, req.Status, notesPtr)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty status updated",
		slog.String("warranty_id", id),
		slog.String("status", string(req.Status)),
	)

	return s.queryRepo.FindWarrantyByID(ctx, id)
}

func (s *service) CreateClaim(ctx context.Context, req CreateClaimRequest) (*models.Claim, error) {
	if err := validateCreateClaim(req); err != nil {
		return nil, err
	}

	w, err := s.queryRepo.FindWarrantyByTicketNumber(ctx, req.TicketNumber)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	s.autoExpireWarrantyIfNeeded(ctx, w)

	if w.Status != models.WarrantyStatusActive {
		return nil, ErrWarrantyNotActive
	}

	claimID, err := s.generateClaimID()
	if err != nil {
		return nil, err
	}

	claimNum, err := s.generateClaimNumber()
	if err != nil {
		return nil, err
	}

	claim := &models.Claim{
		ID:               claimID,
		ClaimNumber:      claimNum,
		WarrantyID:       w.ID,
		EvaluationStatus: models.ClaimEvaluationPending,
		IssueDescription: req.IssueDescription,
	}

	if err := s.commandRepo.CreateClaim(ctx, claim); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Claim created successfully",
		slog.String("claim_id", claim.ID),
		slog.String("claim_number", claim.ClaimNumber),
		slog.String("warranty_id", claim.WarrantyID),
	)

	return claim, nil
}

func (s *service) GetClaims(ctx context.Context, status, search string, limit int, cursor string) ([]models.ClaimSummary, string, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > utils.MaxLimit {
		limit = utils.MaxLimit
	}
	return s.queryRepo.FindAllClaimSummaries(ctx, status, search, limit, cursor)
}

func (s *service) GetClaimByID(ctx context.Context, id string) (*models.Claim, error) {
	c, err := s.queryRepo.FindClaimByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}
	return c, nil
}

func (s *service) GetClaimSummaryByID(ctx context.Context, id string) (*models.ClaimSummary, error) {
	c, err := s.queryRepo.FindClaimSummaryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}
	return c, nil
}

func (s *service) UpdateClaim(ctx context.Context, id string, req UpdateClaimRequest) (*models.Claim, error) {
	if err := validateUpdateClaim(req); err != nil {
		return nil, err
	}

	c, err := s.queryRepo.FindClaimByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}

	c.IssueDescription = req.IssueDescription
	if req.Notes != "" {
		c.Notes = &req.Notes
	} else {
		c.Notes = nil
	}

	if err := s.commandRepo.UpdateClaim(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *service) EvaluateClaim(ctx context.Context, claimID string, req EvaluateClaimRequest) (*models.Claim, error) {
	trimmedNotes, err := validateEvaluateClaim(req)
	if err != nil {
		return nil, err
	}

	c, err := s.queryRepo.FindClaimByID(ctx, claimID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}

	var notesPtr *string
	if trimmedNotes != "" {
		notesPtr = &trimmedNotes
	}

	isVoidWarranty := req.Status == models.ClaimEvaluationVoid
	var warrantyNotes *string
	if isVoidWarranty {
		wNotes := fmt.Sprintf("Batal dari Klaim %s: %s", c.ClaimNumber, trimmedNotes)
		warrantyNotes = &wNotes
	}

	err = s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		var ticketRefIDPtr *string
		if req.Status == models.ClaimEvaluationAccepted {
			if s.ticketCreator == nil {
				return fmt.Errorf("ticket creator is not configured")
			}
			w, err := s.queryRepo.FindWarrantyByID(txCtx, c.WarrantyID)
			if err != nil {
				return err
			}
			if w == nil {
				return ErrWarrantyNotFound
			}
			tID, err := s.ticketCreator.CreateWarrantyTicket(txCtx, w.TicketID, c.WarrantyID, c.IssueDescription)
			if err != nil {
				return fmt.Errorf("failed to create warranty ticket: %w", err)
			}
			ticketRefIDPtr = &tID
		}

		// 1. Update claim evaluation
		err = s.commandRepo.UpdateClaimEvaluation(txCtx, c.ID, req.Status, notesPtr, ticketRefIDPtr)
		if err != nil {
			return err
		}

		// 2. Void warranty if requested
		if isVoidWarranty {
			err = s.commandRepo.UpdateWarrantyStatus(txCtx, c.WarrantyID, models.WarrantyStatusVoid, warrantyNotes)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Claim evaluated successfully",
		slog.String("claim_id", claimID),
		slog.String("evaluation_status", string(req.Status)),
	)

	return s.queryRepo.FindClaimByID(ctx, claimID)
}

// Helpers
func (s *service) generateClaimID() (string, error) {
	return uuid.New().String(), nil
}

func (s *service) generateClaimNumber() (string, error) {
	randNum, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	dateStr := time.Now().Format("20060102")
	return fmt.Sprintf("CLM-%s-%04d", dateStr, randNum.Int64()), nil
}

func (s *service) autoExpireWarrantyIfNeeded(ctx context.Context, w *models.Warranty) {
	if w.Status == models.WarrantyStatusActive && time.Now().After(w.EndDate) {
		w.Status = models.WarrantyStatusExpired
		if err := s.commandRepo.UpdateWarrantyStatus(ctx, w.ID, models.WarrantyStatusExpired, nil); err != nil {
			slog.WarnContext(ctx, "Failed to auto-expire warranty",
				slog.String("warranty_id", w.ID),
				slog.Any("error", err),
			)
		} else {
			slog.InfoContext(ctx, "Warranty auto-expired",
				slog.String("warranty_id", w.ID),
				slog.Time("end_date", w.EndDate),
			)
		}
	}
}

func validateUpdateWarrantyStatus(req UpdateWarrantyStatusRequest) (string, error) {
	if req.Status == "" {
		return "", fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	switch req.Status {
	case models.WarrantyStatusActive, models.WarrantyStatusExpired, models.WarrantyStatusVoid:
		// Valid
	default:
		return "", fmt.Errorf("%w: invalid warranty status", ErrInvalidInput)
	}

	trimmedNotes := strings.TrimSpace(req.Notes)
	if req.Status == models.WarrantyStatusVoid && trimmedNotes == "" {
		return "", fmt.Errorf("%w: notes/reason is required when voiding a warranty", ErrInvalidInput)
	}

	return trimmedNotes, nil
}

func validateCreateClaim(req CreateClaimRequest) error {
	if req.TicketNumber == "" || req.IssueDescription == "" {
		return fmt.Errorf("%w: ticket_number and issue_description are required", ErrInvalidInput)
	}
	return nil
}

func validateUpdateClaim(req UpdateClaimRequest) error {
	if req.IssueDescription == "" {
		return fmt.Errorf("%w: issue_description is required", ErrInvalidInput)
	}
	return nil
}

func validateEvaluateClaim(req EvaluateClaimRequest) (string, error) {
	if req.Status == "" {
		return "", fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	switch req.Status {
	case models.ClaimEvaluationAccepted, models.ClaimEvaluationRejected, models.ClaimEvaluationVoid:
		// Valid
	default:
		return "", fmt.Errorf("%w: invalid claim evaluation status", ErrInvalidInput)
	}

	trimmedNotes := strings.TrimSpace(req.Notes)
	if (req.Status == models.ClaimEvaluationRejected || req.Status == models.ClaimEvaluationVoid) && trimmedNotes == "" {
		return "", fmt.Errorf("%w: notes/reason is required when status is REJECTED or VOID", ErrInvalidInput)
	}

	return trimmedNotes, nil
}
