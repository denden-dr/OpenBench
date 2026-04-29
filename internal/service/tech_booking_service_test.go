package service

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBookingRepository is a mock implementation of the BookingRepository interface.
type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error) {
	args := m.Called(ctx, b)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) FindByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*domain.Booking, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) UpdateStatus(ctx context.Context, id, userID uuid.UUID, currentStatus, newStatus domain.BookingStatus) error {
	args := m.Called(ctx, id, userID, currentStatus, newStatus)
	return args.Error(0)
}

func (m *MockBookingRepository) GetAvailableBookings(ctx context.Context) ([]*domain.Booking, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetBookingsByTechID(ctx context.Context, techID uuid.UUID) ([]*domain.Booking, error) {
	args := m.Called(ctx, techID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) UpdateTechnician(ctx context.Context, id, techID uuid.UUID) error {
	args := m.Called(ctx, id, techID)
	return args.Error(0)
}

func (m *MockBookingRepository) UpdateDiagnosis(ctx context.Context, id, techID uuid.UUID, diagnosis string, cost float64, repairTime string) error {
	args := m.Called(ctx, id, techID, diagnosis, cost, repairTime)
	return args.Error(0)
}

func (m *MockBookingRepository) UpdateTechStatus(ctx context.Context, id, techID uuid.UUID, currentStatus, newStatus domain.BookingStatus) error {
	args := m.Called(ctx, id, techID, currentStatus, newStatus)
	return args.Error(0)
}

func TestTechBookingService_AssignBooking(t *testing.T) {
	ctx := context.Background()
	techID := uuid.New()
	bookingID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(m *MockBookingRepository)
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func(m *MockBookingRepository) {
				m.On("UpdateTechnician", ctx, bookingID, techID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Already Claimed",
			mockSetup: func(m *MockBookingRepository) {
				m.On("UpdateTechnician", ctx, bookingID, techID).Return(domain.ErrConflict)
			},
			expectedError: domain.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockBookingRepository)
			tt.mockSetup(repo)

			svc := NewTechBookingService(repo)
			err := svc.AssignBooking(ctx, techID, bookingID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestTechBookingService_DiagnoseBooking(t *testing.T) {
	ctx := context.Background()
	techID := uuid.New()
	bookingID := uuid.New()
	input := DiagnoseInput{
		Diagnosis:  "Broken Screen",
		Cost:       150.0,
		RepairTime: "2 hours",
	}

	tests := []struct {
		name          string
		input         DiagnoseInput
		mockSetup     func(m *MockBookingRepository)
		expectedError error
	}{
		{
			name:  "Success",
			input: input,
			mockSetup: func(m *MockBookingRepository) {
				m.On("UpdateDiagnosis", ctx, bookingID, techID, input.Diagnosis, input.Cost, input.RepairTime).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Invalid Cost",
			input: DiagnoseInput{
				Diagnosis:  "Broken Screen",
				Cost:       -10.0,
				RepairTime: "2 hours",
			},
			mockSetup:     func(m *MockBookingRepository) {},
			expectedError: domain.ErrInvalidInput,
		},
		{
			name: "Empty Repair Time",
			input: DiagnoseInput{
				Diagnosis:  "Broken Screen",
				Cost:       150.0,
				RepairTime: "",
			},
			mockSetup:     func(m *MockBookingRepository) {},
			expectedError: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockBookingRepository)
			tt.mockSetup(repo)

			svc := NewTechBookingService(repo)
			err := svc.DiagnoseBooking(ctx, techID, bookingID, tt.input)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestTechBookingService_UpdateBookingStatus(t *testing.T) {
	ctx := context.Background()
	techID := uuid.New()
	bookingID := uuid.New()

	tests := []struct {
		name          string
		newStatus     domain.BookingStatus
		mockSetup     func(m *MockBookingRepository)
		expectedError error
	}{
		{
			name:      "Move to In Progress",
			newStatus: domain.StatusInProgress,
			mockSetup: func(m *MockBookingRepository) {
				m.On("UpdateTechStatus", ctx, bookingID, techID, domain.StatusApproved, domain.StatusInProgress).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Move to Completed",
			newStatus: domain.StatusCompleted,
			mockSetup: func(m *MockBookingRepository) {
				m.On("UpdateTechStatus", ctx, bookingID, techID, domain.StatusInProgress, domain.StatusCompleted).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "Invalid Transition",
			newStatus:     domain.StatusApproved,
			mockSetup:     func(m *MockBookingRepository) {},
			expectedError: domain.ErrInvalidStateTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockBookingRepository)
			tt.mockSetup(repo)

			svc := NewTechBookingService(repo)
			err := svc.UpdateBookingStatus(ctx, techID, bookingID, tt.newStatus)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
