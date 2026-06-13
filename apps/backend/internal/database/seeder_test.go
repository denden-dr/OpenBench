//go:build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type SeederTestSuite struct {
	testutil.IntegrationSuite
}

func TestSeederTestSuite(t *testing.T) {
	suite.Run(t, new(SeederTestSuite))
}

func (s *SeederTestSuite) TestSeedDevUsers() {
	db := s.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clean up any existing admin user to test creation
	_, err := db.DB.ExecContext(ctx, "DELETE FROM users WHERE email = $1", "admin@openbench.dev")
	s.Require().NoError(err)

	// First seeding (Creation)
	err = database.SeedDevUsers(db)
	s.Require().NoError(err)

	// Verify details
	var u struct {
		ID           string `db:"id"`
		Email        string `db:"email"`
		PasswordHash string `db:"password_hash"`
		Role         string `db:"role"`
	}
	err = db.DB.GetContext(ctx, &u, "SELECT id, email, password_hash, role FROM users WHERE email = $1", "admin@openbench.dev")
	s.Require().NoError(err)

	s.Equal("admin@openbench.dev", u.Email)
	s.Equal("admin", u.Role)
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("SecureAdminPassword123!"))
	s.NoError(err)

	// Second seeding (Idempotency check)
	err = database.SeedDevUsers(db)
	s.NoError(err)
}
