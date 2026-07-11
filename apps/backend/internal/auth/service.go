package auth

import (
	"context"
	"errors"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
)

type Service interface {
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, error)
}

type service struct {
	queryRepo QueryRepository
	cfg       *config.Config
}

func NewService(queryRepo QueryRepository, cfg *config.Config) Service {
	return &service{
		queryRepo: queryRepo,
		cfg:       cfg,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.queryRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		slog.WarnContext(ctx, "Failed login attempt - user not found", slog.String("email", email))
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		slog.WarnContext(ctx, "Failed login attempt - invalid password", slog.String("email", email))
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "User logged in successfully",
		slog.String("email", user.Email),
		slog.String("user_id", user.ID),
	)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.cfg.Auth.AccessExpiry),
		User: UserProfileResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.cfg.Auth.RefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	user, err := s.queryRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	newAccessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Token refreshed successfully",
		slog.String("email", user.Email),
		slog.String("user_id", user.ID),
	)

	return &RefreshResponse{
		AccessToken: newAccessToken,
		ExpiresAt:   time.Now().Add(s.cfg.Auth.AccessExpiry),
	}, nil
}

func (s *service) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.cfg.Auth.AccessExpiry).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.AccessSecret))
}

func (s *service) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(s.cfg.Auth.RefreshExpiry).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.RefreshSecret))
}
