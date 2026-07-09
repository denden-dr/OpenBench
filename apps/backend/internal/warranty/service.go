package warranty

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
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
	GetClaims(ctx context.Context, status, search string, limit, offset int) ([]models.Claim, int, error)
	GetClaimByID(ctx context.Context, id string) (*models.Claim, error)
	UpdateClaimStatus(ctx context.Context, id string, status models.ServiceTicketStatus) (*models.Claim, error)
	UpdateClaim(ctx context.Context, id string, req UpdateClaimRequest) (*models.Claim, error)
	EvaluateClaim(ctx context.Context, claimID string, req EvaluateClaimRequest) (*models.Claim, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
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

	if err := s.repo.CreateWarranty(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (s *service) GetWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error) {
	w, err := s.repo.FindWarrantyByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	// Dynamic check if warranty has expired
	if w.Status == models.WarrantyStatusActive && time.Now().After(w.EndDate) {
		w.Status = models.WarrantyStatusExpired
		_ = s.repo.UpdateWarrantyStatus(ctx, w.ID, models.WarrantyStatusExpired, nil)
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

	w, err := s.repo.FindWarrantyByID(ctx, id)
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

	err = s.repo.UpdateWarrantyStatus(ctx, id, req.Status, notesPtr)
	if err != nil {
		return nil, err
	}

	return s.repo.FindWarrantyByID(ctx, id)
}

func (s *service) CreateClaim(ctx context.Context, req CreateClaimRequest) (*models.Claim, error) {
	if req.WarrantyID == "" || req.IssueDescription == "" {
		return nil, fmt.Errorf("%w: warranty_id and issue_description are required", ErrInvalidInput)
	}

	w, err := s.repo.FindWarrantyByID(ctx, req.WarrantyID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWarrantyNotFound
	}

	// Check if active
	if w.Status == models.WarrantyStatusActive && time.Now().After(w.EndDate) {
		w.Status = models.WarrantyStatusExpired
		_ = s.repo.UpdateWarrantyStatus(ctx, w.ID, models.WarrantyStatusExpired, nil)
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

	if err := s.repo.CreateClaim(ctx, claim); err != nil {
		return nil, err
	}

	return claim, nil
}

func (s *service) GetClaims(ctx context.Context, status, search string, limit, offset int) ([]models.Claim, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.FindAllClaims(ctx, status, search, limit, offset)
}

func (s *service) GetClaimByID(ctx context.Context, id string) (*models.Claim, error) {
	c, err := s.repo.FindClaimByID(ctx, id)
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

	c, err := s.repo.FindClaimByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClaimNotFound
	}

	c.Status = status
	if err := s.repo.UpdateClaim(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *service) UpdateClaim(ctx context.Context, id string, req UpdateClaimRequest) (*models.Claim, error) {
	if req.IssueDescription == "" {
		return nil, fmt.Errorf("%w: issue_description is required", ErrInvalidInput)
	}

	c, err := s.repo.FindClaimByID(ctx, id)
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

	if err := s.repo.UpdateClaim(ctx, c); err != nil {
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

	c, err := s.repo.FindClaimByID(ctx, claimID)
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

	err = s.repo.EvaluateClaimTx(ctx, c.ID, req.Status, notesPtr, repairStatus, isVoidWarranty, c.WarrantyID, warrantyNotes)
	if err != nil {
		return nil, err
	}

	return s.repo.FindClaimByID(ctx, claimID)
}

// Helpers
func (s *service) generateWarrantyID() (string, error) {
	part1, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	part2, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("w-%04d-%04d", part1.Int64(), part2.Int64()), nil
}

func (s *service) generateClaimID() (string, error) {
	part1, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	part2, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("c-%04d-%04d", part1.Int64(), part2.Int64()), nil
}

func (s *service) generateClaimNumber() (string, error) {
	randNum, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	dateStr := time.Now().Format("20060102")
	return fmt.Sprintf("CLM-%s-%04d", dateStr, randNum.Int64()), nil
}
