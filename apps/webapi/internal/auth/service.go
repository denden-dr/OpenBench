package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
)

type Service interface {
	Login(ctx context.Context, email, password string) (LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (RefreshResponse, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	Me(ctx context.Context, userID string) (UserProfileResponse, error)
}

type service struct {
	queryRepo   QueryRepository
	commandRepo CommandRepository
	cfg         *config.Config
}

func NewService(queryRepo QueryRepository, commandRepo CommandRepository, cfg *config.Config) Service {
	return &service{
		queryRepo:   queryRepo,
		commandRepo: commandRepo,
		cfg:         cfg,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	user, err := s.queryRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, err
	}
	if user == nil {
		slog.WarnContext(ctx, "Failed login attempt - user not found", slog.String("email", email))
		return LoginResponse{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		slog.WarnContext(ctx, "Failed login attempt - invalid password", slog.String("email", email))
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return LoginResponse{}, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, err
	}

	slog.InfoContext(ctx, "User logged in successfully",
		slog.String("email", user.Email),
		slog.String("user_id", user.ID),
	)

	return LoginResponse{
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

func (s *service) Refresh(ctx context.Context, refreshToken string) (RefreshResponse, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.cfg.Auth.RefreshSecret), nil
	}, jwt.WithIssuer("OpenBench"), jwt.WithAudience("OpenBench-Client"))

	if err != nil || !token.Valid {
		return RefreshResponse{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return RefreshResponse{}, ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		return RefreshResponse{}, ErrInvalidToken
	}

	// Blacklist the old refresh token (Refresh Token Rotation)
	if jti, ok := claims["jti"].(string); ok && jti != "" {
		if expFloat, ok := claims["exp"].(float64); ok {
			expiresAt := time.Unix(int64(expFloat), 0)
			if time.Now().Before(expiresAt) {
				_ = s.commandRepo.BlacklistToken(ctx, jti, expiresAt)
			}
		}
	}

	user, err := s.queryRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return RefreshResponse{}, err
	}
	if user == nil {
		return RefreshResponse{}, ErrUserNotFound
	}

	newAccessToken, err := s.generateAccessToken(user)
	if err != nil {
		return RefreshResponse{}, err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return RefreshResponse{}, err
	}

	slog.InfoContext(ctx, "Token refreshed successfully",
		slog.String("email", user.Email),
		slog.String("user_id", user.ID),
	)

	return RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(s.cfg.Auth.AccessExpiry),
	}, nil
}

func (s *service) Logout(ctx context.Context, accessToken, refreshToken string) error {
	blacklistToken := func(tokenString string, secret string) {
		if tokenString == "" {
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		// We still process the claims even if it's expired
		if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok || jti == "" {
			return
		}

		expFloat, ok := claims["exp"].(float64)
		if !ok {
			return
		}

		expiresAt := time.Unix(int64(expFloat), 0)

		// If it's already expired, no need to blacklist
		if time.Now().After(expiresAt) {
			return
		}

		_ = s.commandRepo.BlacklistToken(ctx, jti, expiresAt)
	}

	blacklistToken(accessToken, s.cfg.Auth.AccessSecret)
	blacklistToken(refreshToken, s.cfg.Auth.RefreshSecret)

	slog.InfoContext(ctx, "User logged out successfully")

	return nil
}

func (s *service) Me(ctx context.Context, userID string) (UserProfileResponse, error) {
	user, err := s.queryRepo.GetUserByID(ctx, userID)
	if err != nil {
		return UserProfileResponse{}, err
	}
	if user == nil {
		return UserProfileResponse{}, ErrUserNotFound
	}

	return UserProfileResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *service) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"jti":   uuid.New().String(),
		"exp":   time.Now().Add(s.cfg.Auth.AccessExpiry).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "OpenBench",
		"aud":   "OpenBench-Client",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.AccessSecret))
}

func (s *service) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"jti":   uuid.New().String(),
		"exp":   time.Now().Add(s.cfg.Auth.RefreshExpiry).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "OpenBench",
		"aud":   "OpenBench-Client",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.RefreshSecret))
}
