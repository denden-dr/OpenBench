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

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func()
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:   "Success",
			userID: userID,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "full_name", "avatar_url", "updated_at"}).
					AddRow(userID, "test@example.com", "Test User", nil, time.Now())

				mock.ExpectQuery("SELECT id, email, full_name, avatar_url, updated_at FROM users WHERE id = \\$1").
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedUser: &domain.User{
				ID:    userID,
				Email: "test@example.com",
			},
			expectedError: nil,
		},
		{
			name:   "Not Found",
			userID: userID,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, email, full_name, avatar_url, updated_at FROM users WHERE id = \\$1").
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.FindByID(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}
		})
	}
}

func TestUserRepository_UpsertFromAuth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()
	userID := uuid.New()
	email := "test@example.com"
	fullName := "Test User"

	tests := []struct {
		name          string
		userID        uuid.UUID
		email         string
		fullName      *string
		mockSetup     func()
		expectedUser  *domain.User
		expectedError bool
	}{
		{
			name:     "Success",
			userID:   userID,
			email:    email,
			fullName: &fullName,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "full_name", "avatar_url", "updated_at"}).
					AddRow(userID, email, fullName, nil, time.Now())

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(userID, email, &fullName, nil).
					WillReturnRows(rows)
			},
			expectedUser: &domain.User{
				ID:    userID,
				Email: email,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.UpsertFromAuth(ctx, tt.userID, tt.email, tt.fullName, nil)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}
		})
	}
}
