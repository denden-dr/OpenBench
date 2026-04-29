package service

import (
	"context"
	"fmt"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/dto"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/google/uuid"
)

// UserService defines the business logic for users.
type UserService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetProfile(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user profile: %w", err)
	}

	return toUserResponse(user), nil
}

func toUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
		UpdatedAt: user.UpdatedAt,
	}
}
