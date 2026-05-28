package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	mockrepo "github.com/denden-dr/openbench/apps/backend/mocks/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWarrantyClaimService_CreateClaim(t *testing.T) {
	now := time.Now().UTC()
	pastExitDate := now.AddDate(0, 0, -35)
	futureExitDate := now.AddDate(0, 0, -1)

	tests := []struct {
		name             string
		req              *dto.CreateWarrantyClaimRequest
		setupMock        func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository)
		expectedError    error
		expectedContains string
		expectedAssert   func(t *testing.T, res *dto.WarrantyClaimResponse)
	}{
		{
			name: "success with additional description",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID:              "11111111-1111-1111-1111-111111111111",
				Issue:                 "Layar flicker",
				AdditionalDescription: "Flicker di bagian bawah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "11111111-1111-1111-1111-111111111111").Return(&model.Ticket{
					ID:           "11111111-1111-1111-1111-111111111111",
					CustomerName: "Budi",
					Brand:        "Apple",
					Model:        "iPhone 13",
					Issue:        "LCD Mati",
					Price:        decimal.NewFromInt(1500000),
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "11111111-1111-1111-1111-111111111111").Return(nil, nil).Once()
				mClaim.On("Create", mock.Anything, mock.MatchedBy(func(claim *model.WarrantyClaim) bool {
					return claim.TicketID == "11111111-1111-1111-1111-111111111111" &&
						claim.Issue == "Layar flicker" &&
						claim.AdditionalDescription != nil &&
						*claim.AdditionalDescription == "Flicker di bagian bawah" &&
						claim.Status == model.ClaimWaitingInspection
				})).Run(func(args mock.Arguments) {
					claim := args.Get(1).(*model.WarrantyClaim)
					claim.ID = "claim-123"
					claim.CreatedAt = now
					claim.UpdatedAt = now
				}).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.WarrantyClaimResponse) {
				assert.Equal(t, "claim-123", res.ID)
				assert.Equal(t, "11111111-1111-1111-1111-111111111111", res.TicketID)
				assert.Equal(t, "waiting_inspection", res.Status)
				assert.Equal(t, "Layar flicker", res.Issue)
				assert.NotNil(t, res.AdditionalDescription)
				assert.Equal(t, "Flicker di bagian bawah", *res.AdditionalDescription)
				assert.NotNil(t, res.OriginalTicket)
				assert.Equal(t, "Budi", res.OriginalTicket.CustomerName)
			},
		},
		{
			name: "success without additional description",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "22222222-2222-2222-2222-222222222222",
				Issue:    "Baterai drop",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "22222222-2222-2222-2222-222222222222").Return(&model.Ticket{
					ID:           "22222222-2222-2222-2222-222222222222",
					CustomerName: "Andi",
					Brand:        "Samsung",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "22222222-2222-2222-2222-222222222222").Return(nil, nil).Once()
				mClaim.On("Create", mock.Anything, mock.MatchedBy(func(claim *model.WarrantyClaim) bool {
					return claim.AdditionalDescription == nil
				})).Run(func(args mock.Arguments) {
					claim := args.Get(1).(*model.WarrantyClaim)
					claim.ID = "claim-456"
					claim.CreatedAt = now
					claim.UpdatedAt = now
				}).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.WarrantyClaimResponse) {
				assert.Equal(t, "claim-456", res.ID)
				assert.Equal(t, "22222222-2222-2222-2222-222222222222", res.TicketID)
				assert.Equal(t, "Baterai drop", res.Issue)
				assert.Nil(t, res.AdditionalDescription)
			},
		},
		{
			name: "ticket not picked up",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "33333333-3333-3333-3333-333333333333",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "33333333-3333-3333-3333-333333333333").Return(&model.Ticket{
					ID:           "33333333-3333-3333-3333-333333333333",
					CustomerName: "Budi",
					Status:       model.StatusServiceIn,
					WarrantyDays: 30,
					ExitDate:     nil,
				}, nil).Once()
			},
			expectedError: ErrTicketNotPickedUp,
		},
		{
			name: "warranty expired",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "44444444-4444-4444-4444-444444444444",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "44444444-4444-4444-4444-444444444444").Return(&model.Ticket{
					ID:           "44444444-4444-4444-4444-444444444444",
					CustomerName: "Budi",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &pastExitDate,
				}, nil).Once()
			},
			expectedError: ErrWarrantyExpired,
		},
		{
			name: "ticket not found",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "55555555-5555-5555-5555-555555555555",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "55555555-5555-5555-5555-555555555555").Return(nil, repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
		{
			name: "repository error on ticket get",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "66666666-6666-6666-6666-666666666666",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "66666666-6666-6666-6666-666666666666").Return(nil, errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
		{
			name: "repository error on claim create",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "11111111-1111-1111-1111-111111111111",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "11111111-1111-1111-1111-111111111111").Return(&model.Ticket{
					ID:           "11111111-1111-1111-1111-111111111111",
					CustomerName: "Budi",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "11111111-1111-1111-1111-111111111111").Return(nil, nil).Once()
				mClaim.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
		{
			name: "validation failure - empty ticket id",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
			},
			expectedContains: "TicketID",
		},
		{
			name: "validation failure - invalid ticket uuid",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "invalid-uuid-format",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
			},
			expectedContains: "TicketID",
		},
		{
			name: "duplicate open claim rejected",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "77777777-7777-7777-7777-777777777777",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "77777777-7777-7777-7777-777777777777").Return(&model.Ticket{
					ID:           "77777777-7777-7777-7777-777777777777",
					CustomerName: "Budi",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "77777777-7777-7777-7777-777777777777").Return(&model.WarrantyClaim{
					ID:       "existing-claim",
					TicketID: "77777777-7777-7777-7777-777777777777",
					Status:   model.ClaimWaitingInspection,
				}, nil).Once()
			},
			expectedError: ErrDuplicateWarrantyClaim,
		},
		{
			name: "duplicate check repository error",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "88888888-8888-8888-8888-888888888888",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "88888888-8888-8888-8888-888888888888").Return(&model.Ticket{
					ID:           "88888888-8888-8888-8888-888888888888",
					CustomerName: "Budi",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "88888888-8888-8888-8888-888888888888").Return(nil, errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClaimRepo := new(mockrepo.MockWarrantyClaimRepository)
			mockTicketRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockClaimRepo, mockTicketRepo)

			s := NewWarrantyClaimService(mockClaimRepo, mockTicketRepo)
			res, err := s.CreateClaim(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var appErr *AppError
				if errors.As(tt.expectedError, &appErr) {
					assert.ErrorIs(t, err, tt.expectedError)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
				assert.Nil(t, res)
			} else if tt.expectedContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedContains)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if tt.expectedAssert != nil {
					tt.expectedAssert(t, res)
				}
			}

			mockClaimRepo.AssertExpectations(t)
			mockTicketRepo.AssertExpectations(t)
		})
	}
}

func TestWarrantyClaimService_ListClaims(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		page           int
		limit          int
		setupMock      func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository)
		expectedError  error
		expectedAssert func(t *testing.T, res *dto.PaginatedWarrantyClaimsResponse)
	}{
		{
			name:   "success list all claims",
			status: "all",
			page:   1,
			limit:  10,
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("CountPaginated", mock.Anything, "all").Return(int64(2), nil).Once()
				mClaim.On("ListPaginated", mock.Anything, "all", 10, 0).Return([]*model.WarrantyClaim{
					{
						ID:       "claim-1",
						TicketID: "ticket-1",
						Issue:    "Layar rusak",
						Status:   model.ClaimWaitingInspection,
					},
					{
						ID:       "claim-2",
						TicketID: "ticket-2",
						Issue:    "Baterai drop",
						Status:   model.ClaimApproved,
					},
				}, nil).Once()
				mTicket.On("GetByIDs", mock.Anything, []string{"ticket-1", "ticket-2"}).Return([]model.Ticket{
					{
						ID:           "ticket-1",
						CustomerName: "Budi",
						Brand:        "Apple",
					},
					{
						ID:           "ticket-2",
						CustomerName: "Andi",
						Brand:        "Samsung",
					},
				}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.PaginatedWarrantyClaimsResponse) {
				assert.Equal(t, int64(2), res.Total)
				assert.Equal(t, int64(1), res.TotalPages)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, "claim-1", res.Data[0].ID)
				assert.Equal(t, "waiting_inspection", res.Data[0].Status)
				assert.Equal(t, "Budi", res.Data[0].OriginalTicket.CustomerName)
				assert.Equal(t, "claim-2", res.Data[1].ID)
				assert.Equal(t, "approved", res.Data[1].Status)
				assert.Equal(t, "Andi", res.Data[1].OriginalTicket.CustomerName)
			},
		},
		{
			name:   "success filter by status",
			status: "waiting_inspection",
			page:   1,
			limit:  10,
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("CountPaginated", mock.Anything, "waiting_inspection").Return(int64(1), nil).Once()
				mClaim.On("ListPaginated", mock.Anything, "waiting_inspection", 10, 0).Return([]*model.WarrantyClaim{
					{
						ID:       "claim-1",
						TicketID: "ticket-1",
						Issue:    "Layar rusak",
						Status:   model.ClaimWaitingInspection,
					},
				}, nil).Once()
				mTicket.On("GetByIDs", mock.Anything, []string{"ticket-1"}).Return([]model.Ticket{
					{
						ID:           "ticket-1",
						CustomerName: "Budi",
					},
				}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.PaginatedWarrantyClaimsResponse) {
				assert.Equal(t, int64(1), res.Total)
				assert.Len(t, res.Data, 1)
				assert.Equal(t, "waiting_inspection", res.Data[0].Status)
			},
		},
		{
			name:   "success empty list",
			status: "all",
			page:   1,
			limit:  10,
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("CountPaginated", mock.Anything, "all").Return(int64(0), nil).Once()
				mClaim.On("ListPaginated", mock.Anything, "all", 10, 0).Return([]*model.WarrantyClaim{}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.PaginatedWarrantyClaimsResponse) {
				assert.Equal(t, int64(0), res.Total)
				assert.Len(t, res.Data, 0)
			},
		},
		{
			name:   "invalid status value",
			status: "invalid_status",
			page:   1,
			limit:  10,
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				// Status validation happens first, no calls to database
			},
			expectedError: ErrInvalidStatus,
		},
		{
			name:   "repository count error",
			status: "all",
			page:   1,
			limit:  10,
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("CountPaginated", mock.Anything, "all").Return(int64(0), errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClaimRepo := new(mockrepo.MockWarrantyClaimRepository)
			mockTicketRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockClaimRepo, mockTicketRepo)

			s := NewWarrantyClaimService(mockClaimRepo, mockTicketRepo)
			res, err := s.ListClaims(context.Background(), tt.status, tt.page, tt.limit)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var appErr *AppError
				if errors.As(tt.expectedError, &appErr) {
					assert.ErrorIs(t, err, tt.expectedError)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if tt.expectedAssert != nil {
					tt.expectedAssert(t, res)
				}
			}

			mockClaimRepo.AssertExpectations(t)
			mockTicketRepo.AssertExpectations(t)
		})
	}
}

func TestWarrantyClaimService_ApproveClaim(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		setupMock     func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository, mTx *mockrepo.MockTransaction)
		expectedError error
	}{
		{
			name: "success approve claim",
			id:   "claim-123",
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository, mTx *mockrepo.MockTransaction) {
				mClaim.On("BeginTx", mock.Anything).Return(mTx, nil).Once()
				claim := &model.WarrantyClaim{
					ID:       "claim-123",
					TicketID: "ticket-123",
					Issue:    "LCD rusak",
					Status:   model.ClaimWaitingInspection,
				}
				mClaim.On("GetByIDForUpdateTx", mock.Anything, mTx, "claim-123").Return(claim, nil).Once()
				ticket := &model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Brand:        "Apple",
					Model:        "iPhone",
				}
				mTicket.On("GetByIDForUpdateTx", mock.Anything, mTx, "ticket-123").Return(ticket, nil).Once()
				mTicket.On("CreateTx", mock.Anything, mTx, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.ParentTicketID != nil && *t.ParentTicketID == "ticket-123" && t.IsWarranty
				})).Return(nil).Once()
				mClaim.On("UpdateTx", mock.Anything, mTx, mock.Anything).Return(nil).Once()
				mTx.On("Commit").Return(nil).Once()
				mTx.On("Rollback").Return(nil).Maybe()
			},
			expectedError: nil,
		},
		{
			name: "claim already decided",
			id:   "claim-123",
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository, mTx *mockrepo.MockTransaction) {
				mClaim.On("BeginTx", mock.Anything).Return(mTx, nil).Once()
				claim := &model.WarrantyClaim{
					ID:     "claim-123",
					Status: model.ClaimApproved,
				}
				mClaim.On("GetByIDForUpdateTx", mock.Anything, mTx, "claim-123").Return(claim, nil).Once()
				mTx.On("Rollback").Return(nil).Once()
			},
			expectedError: ErrClaimAlreadyDecided,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClaimRepo := new(mockrepo.MockWarrantyClaimRepository)
			mockTicketRepo := new(mockrepo.MockTicketRepository)
			mockTx := new(mockrepo.MockTransaction)
			tt.setupMock(mockClaimRepo, mockTicketRepo, mockTx)

			s := NewWarrantyClaimService(mockClaimRepo, mockTicketRepo)
			res, err := s.ApproveClaim(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}

			mockClaimRepo.AssertExpectations(t)
			mockTicketRepo.AssertExpectations(t)
			mockTx.AssertExpectations(t)
		})
	}
}

func TestWarrantyClaimService_VoidClaim(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		req           *dto.VoidWarrantyClaimRequest
		setupMock     func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository, mTx *mockrepo.MockTransaction)
		expectedError error
	}{
		{
			name: "success void claim",
			id:   "claim-123",
			req:  &dto.VoidWarrantyClaimRequest{VoidReason: "Physical damage"},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository, mTx *mockrepo.MockTransaction) {
				mClaim.On("BeginTx", mock.Anything).Return(mTx, nil).Once()
				claim := &model.WarrantyClaim{
					ID:       "claim-123",
					TicketID: "ticket-123",
					Issue:    "LCD rusak",
					Status:   model.ClaimWaitingInspection,
				}
				mClaim.On("GetByIDForUpdateTx", mock.Anything, mTx, "claim-123").Return(claim, nil).Once()
				ticket := &model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Brand:        "Apple",
					Model:        "iPhone",
				}
				mTicket.On("GetByIDForUpdateTx", mock.Anything, mTx, "ticket-123").Return(ticket, nil).Once()
				mTicket.On("CreateTx", mock.Anything, mTx, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.ParentTicketID != nil && *t.ParentTicketID == "ticket-123" && t.IsWarranty && t.Status == model.StatusCancelled
				})).Return(nil).Once()
				mClaim.On("UpdateTx", mock.Anything, mTx, mock.Anything).Return(nil).Once()
				mTx.On("Commit").Return(nil).Once()
				mTx.On("Rollback").Return(nil).Maybe()
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClaimRepo := new(mockrepo.MockWarrantyClaimRepository)
			mockTicketRepo := new(mockrepo.MockTicketRepository)
			mockTx := new(mockrepo.MockTransaction)
			tt.setupMock(mockClaimRepo, mockTicketRepo, mockTx)

			s := NewWarrantyClaimService(mockClaimRepo, mockTicketRepo)
			res, err := s.VoidClaim(context.Background(), tt.id, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}

			mockClaimRepo.AssertExpectations(t)
			mockTicketRepo.AssertExpectations(t)
			mockTx.AssertExpectations(t)
		})
	}
}
