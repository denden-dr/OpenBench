package warranty

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateWarranty(ctx context.Context, w *models.Warranty) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *mockRepository) FindWarrantyByID(ctx context.Context, id string) (*models.Warranty, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Warranty), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) FindWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error) {
	args := m.Called(ctx, ticketID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Warranty), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) UpdateWarrantyStatus(ctx context.Context, id string, status models.WarrantyStatus, notes *string) error {
	args := m.Called(ctx, id, status, notes)
	return args.Error(0)
}

func (m *mockRepository) CreateClaim(ctx context.Context, c *models.Claim) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *mockRepository) FindClaimByID(ctx context.Context, id string) (*models.Claim, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Claim), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) FindAllClaims(ctx context.Context, status string, search string, limit, offset int) ([]models.Claim, int, error) {
	args := m.Called(ctx, status, search, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Claim), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *mockRepository) UpdateClaim(ctx context.Context, c *models.Claim) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *mockRepository) EvaluateClaimTx(ctx context.Context, claimID string, evalStatus models.ClaimEvaluationStatus, evalNotes *string, repairStatus models.ServiceTicketStatus, isVoidWarranty bool, warrantyID string, warrantyNotes *string) error {
	args := m.Called(ctx, claimID, evalStatus, evalNotes, repairStatus, isVoidWarranty, warrantyID, warrantyNotes)
	return args.Error(0)
}

func TestService_CreateWarranty(t *testing.T) {
	tests := []struct {
		name         string
		ticketID     string
		warrantyDays int
		mockErr      error
		expectedErr  error
	}{
		{
			name:         "Success",
			ticketID:     "tkt-123",
			warrantyDays: 30,
			mockErr:      nil,
			expectedErr:  nil,
		},
		{
			name:         "Failure - Invalid ticketID",
			ticketID:     "",
			warrantyDays: 30,
			mockErr:      nil,
			expectedErr:  ErrInvalidInput,
		},
		{
			name:         "Failure - Invalid warrantyDays",
			ticketID:     "tkt-123",
			warrantyDays: 0,
			mockErr:      nil,
			expectedErr:  ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			if tt.expectedErr == nil {
				repo.On("CreateWarranty", mock.Anything, mock.AnythingOfType("*models.Warranty")).Return(nil)
			}
			svc := NewService(repo)

			res, err := svc.CreateWarranty(context.Background(), tt.ticketID, tt.warrantyDays)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
				is.Equal(tt.ticketID, res.TicketID)
				is.Equal(models.WarrantyStatusActive, res.Status)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_GetWarrantyByTicketID(t *testing.T) {
	now := time.Now()
	activeWarranty := &models.Warranty{
		ID:        "w-1",
		TicketID:  "tkt-1",
		StartDate: now.AddDate(0, 0, -5),
		EndDate:   now.AddDate(0, 0, 25),
		Status:    models.WarrantyStatusActive,
	}
	expiredWarranty := &models.Warranty{
		ID:        "w-2",
		TicketID:  "tkt-2",
		StartDate: now.AddDate(0, 0, -35),
		EndDate:   now.AddDate(0, 0, -5),
		Status:    models.WarrantyStatusActive, // Needs to be marked EXPIRED dynamically
	}

	tests := []struct {
		name        string
		ticketID    string
		setupMock   func(repo *mockRepository)
		expectedErr error
		checkStatus models.WarrantyStatus
	}{
		{
			name:     "Success - Active",
			ticketID: "tkt-1",
			setupMock: func(repo *mockRepository) {
				repo.On("FindWarrantyByTicketID", mock.Anything, "tkt-1").Return(activeWarranty, nil)
			},
			expectedErr: nil,
			checkStatus: models.WarrantyStatusActive,
		},
		{
			name:     "Success - Expired dynamic update",
			ticketID: "tkt-2",
			setupMock: func(repo *mockRepository) {
				repo.On("FindWarrantyByTicketID", mock.Anything, "tkt-2").Return(expiredWarranty, nil)
				repo.On("UpdateWarrantyStatus", mock.Anything, "w-2", models.WarrantyStatusExpired, mock.Anything).Return(nil)
			},
			expectedErr: nil,
			checkStatus: models.WarrantyStatusExpired,
		},
		{
			name:     "Failure - Not Found",
			ticketID: "non-existent",
			setupMock: func(repo *mockRepository) {
				repo.On("FindWarrantyByTicketID", mock.Anything, "non-existent").Return(nil, nil)
			},
			expectedErr: ErrWarrantyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			tt.setupMock(repo)
			svc := NewService(repo)

			res, err := svc.GetWarrantyByTicketID(context.Background(), tt.ticketID)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
				is.Equal(tt.checkStatus, res.Status)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_CreateClaim(t *testing.T) {
	now := time.Now()
	activeWarranty := &models.Warranty{
		ID:        "w-active",
		TicketID:  "tkt-1",
		StartDate: now.AddDate(0, 0, -5),
		EndDate:   now.AddDate(0, 0, 25),
		Status:    models.WarrantyStatusActive,
	}
	expiredWarranty := &models.Warranty{
		ID:        "w-expired",
		TicketID:  "tkt-2",
		StartDate: now.AddDate(0, 0, -35),
		EndDate:   now.AddDate(0, 0, -5),
		Status:    models.WarrantyStatusExpired,
	}

	tests := []struct {
		name        string
		req         CreateClaimRequest
		setupMock   func(repo *mockRepository)
		expectedErr error
	}{
		{
			name: "Success - Active Warranty",
			req: CreateClaimRequest{
				WarrantyID:       "w-active",
				IssueDescription: "Touchscreen error",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindWarrantyByID", mock.Anything, "w-active").Return(activeWarranty, nil)
				repo.On("CreateClaim", mock.Anything, mock.AnythingOfType("*models.Claim")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Failure - Expired Warranty",
			req: CreateClaimRequest{
				WarrantyID:       "w-expired",
				IssueDescription: "Touchscreen error",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindWarrantyByID", mock.Anything, "w-expired").Return(expiredWarranty, nil)
			},
			expectedErr: ErrWarrantyNotActive,
		},
		{
			name: "Failure - Empty Input",
			req: CreateClaimRequest{
				WarrantyID:       "",
				IssueDescription: "Touchscreen error",
			},
			setupMock:   func(repo *mockRepository) {},
			expectedErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			tt.setupMock(repo)
			svc := NewService(repo)

			res, err := svc.CreateClaim(context.Background(), tt.req)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
				is.Equal(models.StatusReceived, res.Status)
				is.Equal(models.ClaimEvaluationPending, res.EvaluationStatus)
				is.Equal(tt.req.IssueDescription, res.IssueDescription)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_EvaluateClaim(t *testing.T) {
	claim := &models.Claim{
		ID:          "c-1",
		ClaimNumber: "CLM-001",
		WarrantyID:  "w-1",
		Status:      models.StatusReceived,
	}

	tests := []struct {
		name        string
		claimID     string
		req         EvaluateClaimRequest
		setupMock   func(repo *mockRepository)
		expectedErr error
	}{
		{
			name:    "Success - Accept Claim",
			claimID: "c-1",
			req: EvaluateClaimRequest{
				Status: models.ClaimEvaluationAccepted,
				Notes:  "Claim is valid",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindClaimByID", mock.Anything, "c-1").Return(claim, nil).Twice()
				repo.On("EvaluateClaimTx", mock.Anything, "c-1", models.ClaimEvaluationAccepted, mock.Anything, models.StatusRepairing, false, "w-1", mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Success - Reject Claim (Requires Notes)",
			claimID: "c-1",
			req: EvaluateClaimRequest{
				Status: models.ClaimEvaluationRejected,
				Notes:  "Device shows drop impact",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindClaimByID", mock.Anything, "c-1").Return(claim, nil).Twice()
				repo.On("EvaluateClaimTx", mock.Anything, "c-1", models.ClaimEvaluationRejected, mock.Anything, models.StatusCancelled, false, "w-1", mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Success - Void Claim (Requires Notes, voids warranty)",
			claimID: "c-1",
			req: EvaluateClaimRequest{
				Status: models.ClaimEvaluationVoid,
				Notes:  "Unauthorized modification",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindClaimByID", mock.Anything, "c-1").Return(claim, nil).Twice()
				repo.On("EvaluateClaimTx", mock.Anything, "c-1", models.ClaimEvaluationVoid, mock.Anything, models.StatusCancelled, true, "w-1", mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Failure - Reject Claim Without Notes",
			claimID: "c-1",
			req: EvaluateClaimRequest{
				Status: models.ClaimEvaluationRejected,
				Notes:  "   ",
			},
			setupMock:   func(repo *mockRepository) {},
			expectedErr: ErrInvalidInput,
		},
		{
			name:    "Failure - Void Claim Without Notes",
			claimID: "c-1",
			req: EvaluateClaimRequest{
				Status: models.ClaimEvaluationVoid,
				Notes:  "",
			},
			setupMock:   func(repo *mockRepository) {},
			expectedErr: ErrInvalidInput,
		},
		{
			name:    "Failure - Claim Not Found",
			claimID: "c-non-existent",
			req: EvaluateClaimRequest{
				Status: models.ClaimEvaluationAccepted,
				Notes:  "Valid",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindClaimByID", mock.Anything, "c-non-existent").Return(nil, nil)
			},
			expectedErr: ErrClaimNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			tt.setupMock(repo)
			svc := NewService(repo)

			res, err := svc.EvaluateClaim(context.Background(), tt.claimID, tt.req)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_UpdateWarrantyStatus(t *testing.T) {
	warranty := &models.Warranty{
		ID:       "w-1",
		TicketID: "tkt-1",
		Status:   models.WarrantyStatusActive,
	}

	tests := []struct {
		name        string
		warrantyID  string
		req         UpdateWarrantyStatusRequest
		setupMock   func(repo *mockRepository)
		expectedErr error
	}{
		{
			name:       "Success - Direct Update to Void with Notes",
			warrantyID: "w-1",
			req: UpdateWarrantyStatusRequest{
				Status: models.WarrantyStatusVoid,
				Notes:  "Liquid damage detected during diagnostic",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindWarrantyByID", mock.Anything, "w-1").Return(warranty, nil).Twice()
				repo.On("UpdateWarrantyStatus", mock.Anything, "w-1", models.WarrantyStatusVoid, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "Failure - Direct Update to Void without Notes",
			warrantyID: "w-1",
			req: UpdateWarrantyStatusRequest{
				Status: models.WarrantyStatusVoid,
				Notes:  "",
			},
			setupMock:   func(repo *mockRepository) {},
			expectedErr: ErrInvalidInput,
		},
		{
			name:       "Failure - Warranty Not Found",
			warrantyID: "w-non-existent",
			req: UpdateWarrantyStatusRequest{
				Status: models.WarrantyStatusVoid,
				Notes:  "Some reason",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindWarrantyByID", mock.Anything, "w-non-existent").Return(nil, nil)
			},
			expectedErr: ErrWarrantyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			tt.setupMock(repo)
			svc := NewService(repo)

			res, err := svc.UpdateWarrantyStatus(context.Background(), tt.warrantyID, tt.req)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
			}
			repo.AssertExpectations(t)
		})
	}
}
