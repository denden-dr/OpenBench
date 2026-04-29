package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTechnicianRepository_FindByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewTechnicianRepository(sqlxDB)
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func()
		expectedTech  *domain.Technician
		expectedError error
	}{
		{
			name:   "Success",
			userID: userID,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"user_id", "bio", "specialties", "rating", "total_repairs_completed", "is_active", "created_at", "updated_at"}).
					AddRow(userID, "Expert tech", "Screens", 4.5, 10, true, time.Now(), time.Now())

				mock.ExpectQuery("SELECT \\* FROM technicians WHERE user_id = \\$1").
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedTech: &domain.Technician{
				UserID: userID,
				Rating: 4.5,
			},
			expectedError: nil,
		},
		{
			name:   "Not Found",
			userID: userID,
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM technicians WHERE user_id = \\$1").
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedTech:  nil,
			expectedError: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			tech, err := repo.FindByUserID(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, tech)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tech)
				assert.Equal(t, tt.expectedTech.UserID, tech.UserID)
				assert.Equal(t, tt.expectedTech.Rating, tech.Rating)
			}
		})
	}
}
