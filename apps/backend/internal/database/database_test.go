//go:build integration

package database_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type DatabaseConnectionTestSuite struct {
	testutil.IntegrationSuite
}

func TestDatabaseConnectionSuite(t *testing.T) {
	suite.Run(t, new(DatabaseConnectionTestSuite))
}

func (s *DatabaseConnectionTestSuite) TestDatabaseConnection() {
	db := s.DB
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var result int
	err := db.DB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	s.Require().NoError(err)
	s.Equal(1, result)
}

func TestMain(m *testing.M) {
	tdb, err := testutil.SetupTestDB()
	if err != nil {
		log.Fatalf("Failed to setup integration test database: %v", err)
	}

	code := m.Run()

	tdb.Terminate()

	os.Exit(code)
}
