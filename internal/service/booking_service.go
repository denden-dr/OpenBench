package service

import (
	"context"
	"fmt"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/google/uuid"
)

// BookingService defines the business logic for customer repair bookings.
type BookingService interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, deviceName, issueDescription string) (*domain.Booking, error)
	GetBookingDetails(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error)
	ApproveRepair(ctx context.Context, userID, bookingID uuid.UUID) error
	CancelRepair(ctx context.Context, userID, bookingID uuid.UUID) error
}

type bookingService struct {
	repo repository.BookingRepository
}

// NewBookingService returns a new instance of the booking service.
func NewBookingService(repo repository.BookingRepository) BookingService {
	return &bookingService{repo: repo}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID uuid.UUID, deviceName, issueDescription string) (*domain.Booking, error) {
	b := &domain.Booking{
		UserID:           userID,
		DeviceName:       deviceName,
		IssueDescription: issueDescription,
		Status:           domain.StatusPendingDiagnosis,
	}
	return s.repo.Create(ctx, b)
}

func (s *bookingService) GetBookingDetails(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error) {
	return s.repo.FindByIDAndUser(ctx, bookingID, userID)
}

func (s *bookingService) ApproveRepair(ctx context.Context, userID, bookingID uuid.UUID) error {
	// Transition from DiagnosisComplete to Approved
	err := s.repo.UpdateStatus(ctx, bookingID, userID, domain.StatusDiagnosisComplete, domain.StatusApproved)
	if err != nil {
		return fmt.Errorf("approving repair: %w", err)
	}
	return nil
}

func (s *bookingService) CancelRepair(ctx context.Context, userID, bookingID uuid.UUID) error {
	// Fetch booking to check current status
	b, err := s.repo.FindByIDAndUser(ctx, bookingID, userID)
	if err != nil {
		return err
	}

	// Users can cancel during Pending or Diagnosis Complete
	if b.Status != domain.StatusPendingDiagnosis && b.Status != domain.StatusDiagnosisComplete {
		return domain.ErrInvalidStateTransition
	}

	return s.repo.UpdateStatus(ctx, bookingID, userID, b.Status, domain.StatusCanceled)
}
