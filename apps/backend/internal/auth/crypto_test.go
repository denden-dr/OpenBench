package auth_test

import (
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	password := "MySecretPassword123!"

	hash, err := auth.HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Valid password check
	assert.True(t, auth.CheckPasswordHash(password, hash))

	// Invalid password check
	assert.False(t, auth.CheckPasswordHash("wrong_password", hash))
}

func TestAccessTokenLifecycle(t *testing.T) {
	userID := "user-12345"
	role := "admin"
	secret := "my_jwt_super_secret_key"
	expiry := 5 * time.Minute

	// Generate token
	tokenStr, err := auth.GenerateAccessToken(userID, role, secret, expiry)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// Parse valid token
	claims, err := auth.ParseAccessToken(tokenStr, secret)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.WithinDuration(t, time.Now().Add(expiry), claims.ExpiresAt.Time, 5*time.Second)

	// Parse with wrong secret
	_, err = auth.ParseAccessToken(tokenStr, "wrong_secret")
	assert.Error(t, err)
}

func TestRefreshTokenUniqueness(t *testing.T) {
	token1, err := auth.GenerateRefreshToken()
	require.NoError(t, err)
	assert.Len(t, token1, 64) // Hex-encoded 32 bytes = 64 characters

	token2, err := auth.GenerateRefreshToken()
	require.NoError(t, err)
	assert.Len(t, token2, 64)

	assert.NotEqual(t, token1, token2)
}
