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
				TicketID:              "ticket-123",
				Issue:                 "Layar flicker",
				AdditionalDescription: "Flicker di bagian bawah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Brand:        "Apple",
					Model:        "iPhone 13",
					Issue:        "LCD Mati",
					Price:        decimal.NewFromInt(1500000),
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "ticket-123").Return(nil, nil).Once()
				mClaim.On("Create", mock.Anything, mock.MatchedBy(func(claim *model.WarrantyClaim) bool {
					return claim.TicketID == "ticket-123" &&
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
				assert.Equal(t, "ticket-123", res.TicketID)
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
				TicketID: "ticket-456",
				Issue:    "Baterai drop",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-456").Return(&model.Ticket{
					ID:           "ticket-456",
					CustomerName: "Andi",
					Brand:        "Samsung",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "ticket-456").Return(nil, nil).Once()
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
				assert.Equal(t, "ticket-456", res.TicketID)
				assert.Equal(t, "Baterai drop", res.Issue)
				assert.Nil(t, res.AdditionalDescription)
			},
		},
		{
			name: "ticket not picked up",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "ticket-not-picked",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-not-picked").Return(&model.Ticket{
					ID:           "ticket-not-picked",
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
				TicketID: "ticket-expired",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-expired").Return(&model.Ticket{
					ID:           "ticket-expired",
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
				TicketID: "ticket-not-found",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-not-found").Return(nil, repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
		{
			name: "repository error on ticket get",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "ticket-error",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-error").Return(nil, errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
		{
			name: "repository error on claim create",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "ticket-123",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "ticket-123").Return(nil, nil).Once()
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
			name: "duplicate open claim rejected",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "ticket-dup",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-dup").Return(&model.Ticket{
					ID:           "ticket-dup",
					CustomerName: "Budi",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "ticket-dup").Return(&model.WarrantyClaim{
					ID:       "existing-claim",
					TicketID: "ticket-dup",
					Status:   model.ClaimWaitingInspection,
				}, nil).Once()
			},
			expectedError: ErrDuplicateWarrantyClaim,
		},
		{
			name: "duplicate check repository error",
			req: &dto.CreateWarrantyClaimRequest{
				TicketID: "ticket-duperr",
				Issue:    "Masalah",
			},
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mTicket.On("GetByID", mock.Anything, "ticket-duperr").Return(&model.Ticket{
					ID:           "ticket-duperr",
					CustomerName: "Budi",
					Status:       model.StatusPickedUp,
					WarrantyDays: 30,
					ExitDate:     &futureExitDate,
				}, nil).Once()
				mClaim.On("GetOpenClaimByTicketID", mock.Anything, "ticket-duperr").Return(nil, errors.New("db error")).Once()
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
		setupMock      func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository)
		expectedError  error
		expectedAssert func(t *testing.T, res []*dto.WarrantyClaimResponse)
	}{
		{
			name:   "success list all claims",
			status: "",
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("List", mock.Anything, "").Return([]*model.WarrantyClaim{
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
				mTicket.On("GetByID", mock.Anything, "ticket-1").Return(&model.Ticket{
					ID:           "ticket-1",
					CustomerName: "Budi",
					Brand:        "Apple",
				}, nil).Once()
				mTicket.On("GetByID", mock.Anything, "ticket-2").Return(&model.Ticket{
					ID:           "ticket-2",
					CustomerName: "Andi",
					Brand:        "Samsung",
				}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res []*dto.WarrantyClaimResponse) {
				assert.Len(t, res, 2)
				assert.Equal(t, "claim-1", res[0].ID)
				assert.Equal(t, "waiting_inspection", res[0].Status)
				assert.Equal(t, "Budi", res[0].OriginalTicket.CustomerName)
				assert.Equal(t, "claim-2", res[1].ID)
				assert.Equal(t, "approved", res[1].Status)
				assert.Equal(t, "Andi", res[1].OriginalTicket.CustomerName)
			},
		},
		{
			name:   "success filter by status",
			status: "waiting_inspection",
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("List", mock.Anything, "waiting_inspection").Return([]*model.WarrantyClaim{
					{
						ID:       "claim-1",
						TicketID: "ticket-1",
						Issue:    "Layar rusak",
						Status:   model.ClaimWaitingInspection,
					},
				}, nil).Once()
				mTicket.On("GetByID", mock.Anything, "ticket-1").Return(&model.Ticket{
					ID:           "ticket-1",
					CustomerName: "Budi",
				}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res []*dto.WarrantyClaimResponse) {
				assert.Len(t, res, 1)
				assert.Equal(t, "waiting_inspection", res[0].Status)
			},
		},
		{
			name:   "success empty list",
			status: "",
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("List", mock.Anything, "").Return([]*model.WarrantyClaim{}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res []*dto.WarrantyClaimResponse) {
				assert.Len(t, res, 0)
			},
		},
		{
			name:   "repository list error",
			status: "",
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("List", mock.Anything, "").Return(nil, errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
		{
			name:   "ticket enrichment error",
			status: "",
			setupMock: func(mClaim *mockrepo.MockWarrantyClaimRepository, mTicket *mockrepo.MockTicketRepository) {
				mClaim.On("List", mock.Anything, "").Return([]*model.WarrantyClaim{
					{
						ID:       "claim-1",
						TicketID: "ticket-1",
						Issue:    "Layar rusak",
						Status:   model.ClaimWaitingInspection,
					},
				}, nil).Once()
				mTicket.On("GetByID", mock.Anything, "ticket-1").Return(nil, errors.New("db error")).Once()
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
			res, err := s.ListClaims(context.Background(), tt.status)

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
