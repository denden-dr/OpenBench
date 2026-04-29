package service

import (
	"context"
	"errors"
	"testing"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/dto"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpsertFromAuth(ctx context.Context, id uuid.UUID, email string, fullName, avatarURL *string) (*domain.User, error) {
	args := m.Called(ctx, id, email, fullName, avatarURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestUserService_GetProfile(t *testing.T) {
	userID := uuid.New()
	mockUser := &domain.User{ID: userID, Email: "test@example.com"}
	expectedDTO := &dto.UserResponse{ID: userID, Email: "test@example.com"}
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func(m *MockUserRepository)
		expectedUser  *dto.UserResponse
		expectedError string
	}{
		{
			name:   "Success",
			userID: userID,
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(mockUser, nil)
			},
			expectedUser:  expectedDTO,
			expectedError: "",
		},
		{
			name:   "User Not Found",
			userID: userID,
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, repository.ErrUserNotFound)
			},
			expectedUser:  nil,
			expectedError: repository.ErrUserNotFound.Error(),
		},
		{
			name:   "Repository Failure",
			userID: userID,
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, errors.New("db failure"))
			},
			expectedUser:  nil,
			expectedError: "getting user profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepository)
			tt.mockSetup(repo)

			svc := NewUserService(repo)
			user, err := svc.GetProfile(ctx, tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
			repo.AssertExpectations(t)
		})
	}
}
