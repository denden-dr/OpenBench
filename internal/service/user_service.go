package service

import (
	"context"
	"fmt"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/google/uuid"
)

// UserService defines the business logic for users.
type UserService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error)
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

func (s *userService) GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user profile: %w", err)
	}
	return user, nil
}
