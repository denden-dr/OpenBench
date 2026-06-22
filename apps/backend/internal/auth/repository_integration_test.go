//go:build integration

package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/auth"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AuthRepositoryTestSuite struct {
	testutil.IntegrationSuite
	repo auth.AuthRepository
}

func TestAuthRepositorySuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

func (s *AuthRepositoryTestSuite) SetupTest() {
	s.IntegrationSuite.SetupTest() // Clean tables dynamically
	s.repo = auth.NewRepository(s.DB)
}

func (s *AuthRepositoryTestSuite) TestUserQueries() {
	ctx := context.Background()

	// 1. Seed user directly
	hashedPassword, err := auth.HashPassword("TestPass123!")
	s.Require().NoError(err)

	userID := uuid.New().String()
	_, err = s.DB.DB.ExecContext(ctx,
		"INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)",
		userID, "test-repo@openbench.dev", hashedPassword, "user",
	)
	s.Require().NoError(err)

	// 2. Test GetUserByID
	s.Run("GetUserByID - Success", func() {
		u, err := s.repo.GetUserByID(ctx, userID)
		s.Require().NoError(err)
		s.Assert().Equal("test-repo@openbench.dev", u.Email)
		s.Assert().Equal("user", u.Role)
	})

	s.Run("GetUserByID - Not Found", func() {
		u, err := s.repo.GetUserByID(ctx, uuid.New().String())
		s.Assert().Error(err)
		s.Assert().Nil(u)
	})

	// 3. Test GetUserByEmail
	s.Run("GetUserByEmail - Success", func() {
		u, err := s.repo.GetUserByEmail(ctx, "test-repo@openbench.dev")
		s.Require().NoError(err)
		s.Assert().Equal(userID, u.ID)
	})

	s.Run("GetUserByEmail - Not Found", func() {
		u, err := s.repo.GetUserByEmail(ctx, "nonexistent@openbench.dev")
		s.Assert().Error(err)
		s.Assert().Nil(u)
	})
}

func (s *AuthRepositoryTestSuite) TestRefreshTokenQueries() {
	ctx := context.Background()

	// 1. Seed user (required for foreign key constraint on refresh_tokens)
	userID := uuid.New().String()
	_, err := s.DB.DB.ExecContext(ctx,
		"INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)",
		userID, "token-owner@openbench.dev", "hashed_pass", "user",
	)
	s.Require().NoError(err)

	tokenID := uuid.New().String()
	familyID := uuid.New().String()
	tokenHash := auth.HashSha256("my_secret_refresh_token")
	expiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	record := &auth.RefreshTokenRecord{
		ID:        tokenID,
		FamilyID:  familyID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	// 2. Test CreateRefreshToken
	s.Run("CreateRefreshToken", func() {
		err = s.repo.CreateRefreshToken(ctx, nil, record)
		s.Require().NoError(err)
	})

	// 3. Test GetRefreshTokenByID
	s.Run("GetRefreshTokenByID", func() {
		r, err := s.repo.GetRefreshTokenByID(ctx, nil, tokenID)
		s.Require().NoError(err)
		s.Assert().Equal(tokenID, r.ID)
		s.Assert().Equal(familyID, r.FamilyID)
		s.Assert().Equal(userID, r.UserID)
		s.Assert().Equal(tokenHash, r.TokenHash)
		s.Assert().False(r.IsUsed)
		s.Assert().False(r.IsRevoked)
	})

	// 4. Test GetRefreshTokenWithLock
	s.Run("GetRefreshTokenWithLock", func() {
		r, err := s.repo.GetRefreshTokenWithLock(ctx, nil, tokenHash)
		s.Require().NoError(err)
		s.Assert().Equal(tokenID, r.ID)
	})

	// 5. Test UpdateRefreshToken
	s.Run("UpdateRefreshToken", func() {
		r, err := s.repo.GetRefreshTokenByID(ctx, nil, tokenID)
		s.Require().NoError(err)

		r.IsUsed = true
		now := time.Now().Truncate(time.Second)
		r.UsedAt = &now
		replacedByID := uuid.New().String()
		r.ReplacedByTokenID = &replacedByID

		err = s.repo.UpdateRefreshToken(ctx, nil, r)
		s.Require().NoError(err)

		updated, err := s.repo.GetRefreshTokenByID(ctx, nil, tokenID)
		s.Require().NoError(err)
		s.Assert().True(updated.IsUsed)
		s.Assert().Equal(now.Unix(), updated.UsedAt.Unix())
		s.Assert().Equal(replacedByID, *updated.ReplacedByTokenID)
	})

	// 6. Test RevokeTokenFamily
	s.Run("RevokeTokenFamily", func() {
		err = s.repo.RevokeTokenFamily(ctx, nil, familyID)
		s.Require().NoError(err)

		updated, err := s.repo.GetRefreshTokenByID(ctx, nil, tokenID)
		s.Require().NoError(err)
		s.Assert().True(updated.IsRevoked)
	})

	// 7. Test RevokeTokenByHash
	s.Run("RevokeTokenByHash", func() {
		newTokenID := uuid.New().String()
		newTokenHash := auth.HashSha256("another_secret_refresh_token")
		newRecord := &auth.RefreshTokenRecord{
			ID:        newTokenID,
			FamilyID:  familyID,
			UserID:    userID,
			TokenHash: newTokenHash,
			ExpiresAt: expiresAt,
		}

		err = s.repo.CreateRefreshToken(ctx, nil, newRecord)
		s.Require().NoError(err)

		err = s.repo.RevokeTokenByHash(ctx, newTokenHash)
		s.Require().NoError(err)

		updated, err := s.repo.GetRefreshTokenByID(ctx, nil, newTokenID)
		s.Require().NoError(err)
		s.Assert().True(updated.IsRevoked)
	})

	// 8. Test PurgeExpiredTokens
	s.Run("PurgeExpiredTokens", func() {
		// Insert one expired token
		expiredTokenID := uuid.New().String()
		expiredTokenHash := auth.HashSha256("expired_token")
		expiredRecord := &auth.RefreshTokenRecord{
			ID:        expiredTokenID,
			FamilyID:  familyID,
			UserID:    userID,
			TokenHash: expiredTokenHash,
			ExpiresAt: time.Now().Add(-5 * time.Minute),
		}
		err = s.repo.CreateRefreshToken(ctx, nil, expiredRecord)
		s.Require().NoError(err)

		// Insert one active unrevoked token
		activeTokenID := uuid.New().String()
		activeTokenHash := auth.HashSha256("active_token")
		activeRecord := &auth.RefreshTokenRecord{
			ID:        activeTokenID,
			FamilyID:  familyID,
			UserID:    userID,
			TokenHash: activeTokenHash,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}
		err = s.repo.CreateRefreshToken(ctx, nil, activeRecord)
		s.Require().NoError(err)

		// Insert one revoked token with expires_at in the future (should NOT be purged)
		revokedFutureTokenID := uuid.New().String()
		revokedFutureTokenHash := auth.HashSha256("revoked_future_token")
		revokedFutureRecord := &auth.RefreshTokenRecord{
			ID:        revokedFutureTokenID,
			FamilyID:  familyID,
			UserID:    userID,
			TokenHash: revokedFutureTokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		err = s.repo.CreateRefreshToken(ctx, nil, revokedFutureRecord)
		s.Require().NoError(err)
		err = s.repo.RevokeTokenByHash(ctx, revokedFutureTokenHash)
		s.Require().NoError(err)

		// Insert one revoked token with expires_at older than 7 days (should be purged)
		revokedOldTokenID := uuid.New().String()
		revokedOldTokenHash := auth.HashSha256("revoked_old_token")
		revokedOldRecord := &auth.RefreshTokenRecord{
			ID:        revokedOldTokenID,
			FamilyID:  familyID,
			UserID:    userID,
			TokenHash: revokedOldTokenHash,
			ExpiresAt: time.Now().Add(-8 * 24 * time.Hour),
		}
		err = s.repo.CreateRefreshToken(ctx, nil, revokedOldRecord)
		s.Require().NoError(err)
		err = s.repo.RevokeTokenByHash(ctx, revokedOldTokenHash)
		s.Require().NoError(err)

		// Run Purge
		err = s.repo.PurgeExpiredTokens(ctx)
		s.Require().NoError(err)

		// The expired one should be gone
		_, err = s.repo.GetRefreshTokenByID(ctx, nil, expiredTokenID)
		s.Assert().Error(err) // should not exist

		// The active one should still be there
		activeGot, err := s.repo.GetRefreshTokenByID(ctx, nil, activeTokenID)
		s.Require().NoError(err)
		s.Assert().Equal(activeTokenID, activeGot.ID)

		// The revoked future token should still be there
		revokedFutureGot, err := s.repo.GetRefreshTokenByID(ctx, nil, revokedFutureTokenID)
		s.Require().NoError(err)
		s.Assert().Equal(revokedFutureTokenID, revokedFutureGot.ID)

		// The revoked old token should be gone
		_, err = s.repo.GetRefreshTokenByID(ctx, nil, revokedOldTokenID)
		s.Assert().Error(err) // should not exist
	})
}
