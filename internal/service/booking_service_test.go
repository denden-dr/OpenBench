package service

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingService_CreateBooking(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	deviceName := "iPhone 13"
	issue := "Broken screen"

	repo := new(MockBookingRepository)
	svc := NewBookingService(repo)

	expectedBooking := &domain.Booking{
		UserID:           userID,
		DeviceName:       deviceName,
		IssueDescription: issue,
		Status:           domain.StatusPendingDiagnosis,
	}

	repo.On("Create", ctx, mock.MatchedBy(func(b *domain.Booking) bool {
		return b.UserID == userID && b.DeviceName == deviceName && b.Status == domain.StatusPendingDiagnosis
	})).Return(expectedBooking, nil)

	b, err := svc.CreateBooking(ctx, userID, deviceName, issue)

	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, domain.StatusPendingDiagnosis, b.Status)
	repo.AssertExpectations(t)
}

func TestBookingService_ApproveRepair(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(repo *MockBookingRepository)
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func(repo *MockBookingRepository) {
				repo.On("UpdateStatus", ctx, bookingID, userID, domain.StatusDiagnosisComplete, domain.StatusApproved).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Invalid State",
			mockSetup: func(repo *MockBookingRepository) {
				repo.On("UpdateStatus", ctx, bookingID, userID, domain.StatusDiagnosisComplete, domain.StatusApproved).Return(domain.ErrInvalidStateTransition)
			},
			expectedError: domain.ErrInvalidStateTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockBookingRepository)
			svc := NewBookingService(repo)
			tt.mockSetup(repo)
			err := svc.ApproveRepair(ctx, userID, bookingID)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestBookingService_CancelRepair(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	bookingID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(repo *MockBookingRepository)
		expectedError error
	}{
		{
			name: "Success from Pending",
			mockSetup: func(repo *MockBookingRepository) {
				repo.On("FindByIDAndUser", ctx, bookingID, userID).Return(&domain.Booking{Status: domain.StatusPendingDiagnosis}, nil).Once()
				repo.On("UpdateStatus", ctx, bookingID, userID, domain.StatusPendingDiagnosis, domain.StatusCanceled).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Cannot cancel In Progress",
			mockSetup: func(repo *MockBookingRepository) {
				repo.On("FindByIDAndUser", ctx, bookingID, userID).Return(&domain.Booking{Status: domain.StatusInProgress}, nil).Once()
			},
			expectedError: domain.ErrInvalidStateTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockBookingRepository)
			svc := NewBookingService(repo)
			tt.mockSetup(repo)
			err := svc.CancelRepair(ctx, userID, bookingID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
