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

type Service interface {
	CreateWarranty(ctx context.Context, ticketID string, warrantyDays int) (*models.Warranty, error)
	GetWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error)
	UpdateWarrantyStatus(ctx context.Context, id string, req UpdateWarrantyStatusRequest) (*models.Warranty, error)

	CreateClaim(ctx context.Context, req CreateClaimRequest) (*models.Claim, error)
	GetClaims(ctx context.Context, status, search string, limit int, cursor string) ([]models.Claim, string, error)
	GetClaimByID(ctx context.Context, id string) (*models.Claim, error)
	UpdateClaimStatus(ctx context.Context, id string, status models.ServiceTicketStatus) (*models.Claim, error)
	UpdateClaim(ctx context.Context, id string, req UpdateClaimRequest) (*models.Claim, error)
	EvaluateClaim(ctx context.Context, claimID string, req EvaluateClaimRequest) (*models.Claim, error)
}

type service struct {
	queryRepo   QueryRepository
	commandRepo CommandRepository
	txManager   database.TxManager
}

func NewService(
	queryRepo QueryRepository,
	commandRepo CommandRepository,
	txManager database.TxManager,
) Service {
	return &service{
		queryRepo:   queryRepo,
		commandRepo: commandRepo,
		txManager:   txManager,
	}
}

func (s *service) CreateWarranty(ctx context.Context, ticketID string, warrantyDays int) (*models.Warranty, error) {
	if ticketID == "" || warrantyDays <= 0 {
		return nil, fmt.Errorf("%w: ticketID and warrantyDays > 0 are required", ErrInvalidInput)
	}

	wID, err := s.generateWarrantyID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	w := &models.Warranty{
		ID:        wID,
		TicketID:  ticketID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, warrantyDays),
		Status:    models.WarrantyStatusActive,
	}

	if err := s.commandRepo.CreateWarranty(ctx, w); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty created successfully",
		slog.String("warranty_id", w.ID),
		slog.String("ticket_id", w.TicketID),
		slog.Time("end_date", w.EndDate),
	)

	return w, nil
}

func (s *service) GetWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error) {
	w, err := s.queryRepo.FindWarrantyByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	// Dynamic check if warranty has expired
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

	return w, nil
}

func (s *service) UpdateWarrantyStatus(ctx context.Context, id string, req UpdateWarrantyStatusRequest) (*models.Warranty, error) {
	if req.Status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	switch req.Status {
	case models.WarrantyStatusActive, models.WarrantyStatusExpired, models.WarrantyStatusVoid:
		// Valid
	default:
		return nil, fmt.Errorf("%w: invalid warranty status", ErrInvalidInput)
	}

	trimmedNotes := strings.TrimSpace(req.Notes)
	if req.Status == models.WarrantyStatusVoid && trimmedNotes == "" {
		return nil, fmt.Errorf("%w: notes/reason is required when voiding a warranty", ErrInvalidInput)
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
	if req.WarrantyID == "" || req.IssueDescription == "" {
		return nil, fmt.Errorf("%w: warranty_id and issue_description are required", ErrInvalidInput)
	}

	w, err := s.queryRepo.FindWarrantyByID(ctx, req.WarrantyID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	// Check if active
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
		Status:           models.StatusReceived,
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

func (s *service) GetClaims(ctx context.Context, status, search string, limit int, cursor string) ([]models.Claim, string, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > utils.MaxLimit {
		limit = utils.MaxLimit
	}
	return s.queryRepo.FindAllClaims(ctx, status, search, limit, cursor)
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

func (s *service) UpdateClaimStatus(ctx context.Context, id string, status models.ServiceTicketStatus) (*models.Claim, error) {
	if status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	// Validate status
	switch status {
	case models.StatusReceived, models.StatusRepairing, models.StatusPendingConfirmation,
		models.StatusFixed, models.StatusCompleted, models.StatusCancelled, models.StatusReturned:
		// Valid
	default:
		return nil, fmt.Errorf("%w: invalid claim status", ErrInvalidInput)
	}

	c, err := s.queryRepo.FindClaimByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}

	c.Status = status
	if err := s.commandRepo.UpdateClaim(ctx, c); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Claim status updated",
		slog.String("claim_id", id),
		slog.String("status", string(status)),
	)

	return c, nil
}

func (s *service) UpdateClaim(ctx context.Context, id string, req UpdateClaimRequest) (*models.Claim, error) {
	if req.IssueDescription == "" {
		return nil, fmt.Errorf("%w: issue_description is required", ErrInvalidInput)
	}

	c, err := s.queryRepo.FindClaimByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}

	c.IssueDescription = req.IssueDescription
	if req.RepairAction != "" {
		c.RepairAction = &req.RepairAction
	} else {
		c.RepairAction = nil
	}
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
	if req.Status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrInvalidInput)
	}

	switch req.Status {
	case models.ClaimEvaluationAccepted, models.ClaimEvaluationRejected, models.ClaimEvaluationVoid:
		// Valid
	default:
		return nil, fmt.Errorf("%w: invalid claim evaluation status", ErrInvalidInput)
	}

	trimmedNotes := strings.TrimSpace(req.Notes)
	if (req.Status == models.ClaimEvaluationRejected || req.Status == models.ClaimEvaluationVoid) && trimmedNotes == "" {
		return nil, fmt.Errorf("%w: notes/reason is required when status is REJECTED or VOID", ErrInvalidInput)
	}

	c, err := s.queryRepo.FindClaimByID(ctx, claimID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}

	var repairStatus models.ServiceTicketStatus
	switch req.Status {
	case models.ClaimEvaluationAccepted:
		repairStatus = models.StatusRepairing
	case models.ClaimEvaluationRejected, models.ClaimEvaluationVoid:
		repairStatus = models.StatusCancelled
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
		// 1. Update claim evaluation
		err = s.commandRepo.UpdateClaimEvaluation(txCtx, c.ID, repairStatus, req.Status, notesPtr)
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
		slog.String("repair_status", string(repairStatus)),
	)

	return s.queryRepo.FindClaimByID(ctx, claimID)
}

// Helpers
func (s *service) generateWarrantyID() (string, error) {
	return uuid.New().String(), nil
}

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
