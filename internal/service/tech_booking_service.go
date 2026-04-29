package service

import (
	"context"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/google/uuid"
)

// DiagnoseInput contains the data for a technician diagnosis.
type DiagnoseInput struct {
	Diagnosis  string
	Cost       float64
	RepairTime string
}

// TechBookingService defines the business logic for technician booking management.
type TechBookingService interface {
	GetAvailableBookings(ctx context.Context) ([]*domain.Booking, error)
	GetMyBookings(ctx context.Context, techID uuid.UUID) ([]*domain.Booking, error)
	AssignBooking(ctx context.Context, techID, bookingID uuid.UUID) error
	DiagnoseBooking(ctx context.Context, techID, bookingID uuid.UUID, input DiagnoseInput) error
	UpdateBookingStatus(ctx context.Context, techID, bookingID uuid.UUID, newStatus domain.BookingStatus) error
}

type techBookingService struct {
	bookingRepo repository.BookingRepository
}

// NewTechBookingService returns a new instance of the technician booking service.
func NewTechBookingService(bookingRepo repository.BookingRepository) TechBookingService {
	return &techBookingService{
		bookingRepo: bookingRepo,
	}
}

func (s *techBookingService) GetAvailableBookings(ctx context.Context) ([]*domain.Booking, error) {
	return s.bookingRepo.GetAvailableBookings(ctx)
}

func (s *techBookingService) GetMyBookings(ctx context.Context, techID uuid.UUID) ([]*domain.Booking, error) {
	return s.bookingRepo.GetBookingsByTechID(ctx, techID)
}

func (s *techBookingService) AssignBooking(ctx context.Context, techID, bookingID uuid.UUID) error {
	return s.bookingRepo.UpdateTechnician(ctx, bookingID, techID)
}

func (s *techBookingService) DiagnoseBooking(ctx context.Context, techID, bookingID uuid.UUID, input DiagnoseInput) error {
	// Validate input
	if input.Cost < 0 {
		return domain.ErrInvalidInput
	}
	if input.RepairTime == "" {
		return domain.ErrInvalidInput
	}

	return s.bookingRepo.UpdateDiagnosis(ctx, bookingID, techID, input.Diagnosis, input.Cost, input.RepairTime)
}

func (s *techBookingService) UpdateBookingStatus(ctx context.Context, techID, bookingID uuid.UUID, newStatus domain.BookingStatus) error {
	var currentStatus domain.BookingStatus

	// Enforce state machine rules for technicians
	switch newStatus {
	case domain.StatusInProgress:
		currentStatus = domain.StatusApproved
	case domain.StatusCompleted:
		currentStatus = domain.StatusInProgress
	default:
		return domain.ErrInvalidStateTransition
	}

	return s.bookingRepo.UpdateTechStatus(ctx, bookingID, techID, currentStatus, newStatus)
}
