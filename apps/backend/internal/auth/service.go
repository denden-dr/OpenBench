package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrSessionRevoked      = errors.New("session is revoked")
	ErrTokenCompromised    = errors.New("token compromise detected, session revoked")
	ErrTokenExpired        = errors.New("refresh token expired")
)

type Service interface {
	SignIn(ctx context.Context, email, password string, accessExpiry, refreshExpiry time.Duration) (*SignInResult, error)
	Refresh(ctx context.Context, rawRefreshToken string, accessExpiry, refreshExpiry time.Duration) (*RefreshResult, error)
	SignOut(ctx context.Context, rawRefreshToken string) error
	GetUserByID(ctx context.Context, userID string) (*User, error)
}

type authService struct {
	repo      Repository
	db        *database.Database
	jwtSecret string
}

func NewService(repo Repository, db *database.Database, jwtSecret string) Service {
	return &authService{
		repo:      repo,
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) SignIn(ctx context.Context, email, password string, accessExpiry, refreshExpiry time.Duration) (*SignInResult, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := GenerateAccessToken(user.ID, user.Role, s.jwtSecret, accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	hashedRefreshToken := HashSha256(rawRefreshToken)
	tokenID := uuid.New().String()
	familyID := uuid.New().String()
	expiresAt := time.Now().Add(refreshExpiry)

	record := &RefreshTokenRecord{
		ID:        tokenID,
		FamilyID:  familyID,
		UserID:    user.ID,
		TokenHash: hashedRefreshToken,
		ExpiresAt: expiresAt,
	}

	err = s.repo.CreateRefreshToken(ctx, nil, record)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &SignInResult{
		User:            user,
		AccessToken:     accessToken,
		RawRefreshToken: rawRefreshToken,
		ExpiresAt:       expiresAt,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, rawRefreshToken string, accessExpiry, refreshExpiry time.Duration) (*RefreshResult, error) {
	hashedToken := HashSha256(rawRefreshToken)

	tx, err := s.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	t, err := s.repo.GetRefreshTokenWithLock(ctx, tx, hashedToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if t.IsRevoked {
		return nil, ErrSessionRevoked
	}

	if t.IsUsed {
		return s.handleUsedToken(ctx, tx, t, accessExpiry)
	}

	if time.Now().After(t.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	return s.rotateToken(ctx, tx, t, accessExpiry, refreshExpiry)
}

func (s *authService) handleUsedToken(ctx context.Context, tx *sqlx.Tx, t *RefreshTokenRecord, accessExpiry time.Duration) (*RefreshResult, error) {
	// Check grace period: 5 seconds
	if t.UsedAt != nil && time.Since(*t.UsedAt) < 5*time.Second && t.ReplacedByTokenID != nil {
		succ, err := s.repo.GetRefreshTokenByID(ctx, tx, *t.ReplacedByTokenID)
		if err == nil && !succ.IsRevoked {
			user, err := s.repo.GetUserByID(ctx, t.UserID)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve user: %w", err)
			}

			accessToken, err := GenerateAccessToken(t.UserID, user.Role, s.jwtSecret, accessExpiry)
			if err != nil {
				return nil, fmt.Errorf("failed to generate access token: %w", err)
			}

			if err := tx.Commit(); err != nil {
				return nil, fmt.Errorf("failed to commit transaction: %w", err)
			}

			return &RefreshResult{
				AccessToken:    accessToken,
				GraceRefreshed: true,
			}, nil
		}
	}

	// Compromise detected: revoke entire token family
	_ = s.repo.RevokeTokenFamily(ctx, tx, t.FamilyID)
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction after compromise: %w", err)
	}

	return nil, ErrTokenCompromised
}

func (s *authService) rotateToken(ctx context.Context, tx *sqlx.Tx, t *RefreshTokenRecord, accessExpiry, refreshExpiry time.Duration) (*RefreshResult, error) {
	newRawRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	newHashedRefreshToken := HashSha256(newRawRefreshToken)
	newTokenID := uuid.New().String()
	newExpiresAt := time.Now().Add(refreshExpiry)

	t.IsUsed = true
	now := time.Now()
	t.UsedAt = &now
	t.ReplacedByTokenID = &newTokenID

	err = s.repo.UpdateRefreshToken(ctx, tx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to update current refresh token: %w", err)
	}

	newRecord := &RefreshTokenRecord{
		ID:        newTokenID,
		FamilyID:  t.FamilyID,
		UserID:    t.UserID,
		TokenHash: newHashedRefreshToken,
		ExpiresAt: newExpiresAt,
	}
	err = s.repo.CreateRefreshToken(ctx, tx, newRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, t.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	newAccessToken, err := GenerateAccessToken(t.UserID, user.Role, s.jwtSecret, accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return &RefreshResult{
		AccessToken:     newAccessToken,
		RawRefreshToken: newRawRefreshToken,
		ExpiresAt:       newExpiresAt,
		GraceRefreshed:  false,
	}, nil
}

func (s *authService) SignOut(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}
	hashedToken := HashSha256(rawRefreshToken)
	return s.repo.RevokeTokenByHash(ctx, hashedToken)
}

func (s *authService) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
