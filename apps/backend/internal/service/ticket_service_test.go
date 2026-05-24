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

func TestTicketService_CreateTicket(t *testing.T) {
	tests := []struct {
		name           string
		req            *dto.CreateTicketRequest
		setupMock      func(m *mockrepo.MockTicketRepository)
		expectedError  error
		expectedAssert func(t *testing.T, res *dto.TicketResponse)
	}{
		{
			name: "success with all fields",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Budi",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   30,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(ticket *model.Ticket) bool {
					return ticket.CustomerName == "Budi" &&
						ticket.CustomerGender == "Male" &&
						ticket.Brand == "Apple" &&
						ticket.Model == "iPhone 13" &&
						ticket.Issue == "LCD Mati" &&
						ticket.Price.Equal(decimal.NewFromInt(1500000)) &&
						ticket.WarrantyDays == 30
				})).Run(func(args mock.Arguments) {
					ticket := args.Get(1).(*model.Ticket)
					ticket.ID = "ticket-123"
					ticket.EntryDate = time.Now()
					ticket.Status = "service_in"
					ticket.PaymentStatus = "unpaid"
				}).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "ticket-123", res.ID)
				assert.Equal(t, "Budi", res.CustomerName)
				assert.Equal(t, "service_in", res.Status)
				assert.Equal(t, "unpaid", res.PaymentStatus)
			},
		},
		{
			name: "success with default warranty days",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Andi",
				CustomerGender: "Male",
				Brand:          "Samsung",
				Model:          "Galaxy S21",
				Issue:          "Baterai kembung",
				Price:          decimal.NewFromInt(500000),
				WarrantyDays:   0,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(ticket *model.Ticket) bool {
					return ticket.WarrantyDays == 30
				})).Run(func(args mock.Arguments) {
					ticket := args.Get(1).(*model.Ticket)
					ticket.ID = "ticket-456"
				}).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, 30, res.WarrantyDays)
			},
		},
		{
			name: "negative price error",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Budi",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(-1),
				WarrantyDays:   30,
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: ErrNegativePrice,
		},
		{
			name: "negative warranty error",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Budi",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   -5,
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: ErrNegativeWarranty,
		},
		{
			name: "validation failure - empty customer name",
			req: &dto.CreateTicketRequest{
				CustomerName:   "",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   30,
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: errors.New("validation failed"),
		},
		{
			name: "repository error",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Budi",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   30,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			res, err := s.CreateTicket(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError.Error() == "validation failed" {
					// Validation error checked by presence of error
				} else {
					var appErr *AppError
					if errors.As(tt.expectedError, &appErr) {
						assert.ErrorIs(t, err, tt.expectedError)
					} else {
						assert.EqualError(t, err, tt.expectedError.Error())
					}
				}
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if tt.expectedAssert != nil {
					tt.expectedAssert(t, res)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTicketService_GetTicket(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		setupMock      func(m *mockrepo.MockTicketRepository)
		expectedError  error
		expectedAssert func(t *testing.T, res *dto.TicketResponse)
	}{
		{
			name: "success get ticket",
			id:   "ticket-123",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:             "ticket-123",
					CustomerName:   "Budi",
					CustomerGender: "Male",
					Brand:          "Apple",
					Model:          "iPhone 13",
					Issue:          "LCD Mati",
					Price:          decimal.NewFromInt(1500000),
					Status:         "service_in",
					PaymentStatus:  "unpaid",
					WarrantyDays:   30,
				}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "ticket-123", res.ID)
				assert.Equal(t, "Budi", res.CustomerName)
			},
		},
		{
			name: "ticket not found",
			id:   "ticket-not-found",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-not-found").Return(nil, repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
		{
			name: "repository error",
			id:   "ticket-error",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-error").Return(nil, errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			res, err := s.GetTicket(context.Background(), tt.id)

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
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTicketService_UpdateTicket(t *testing.T) {
	newName := "Budi Baru"
	pickedUp := "picked_up"
	unpaid := "unpaid"
	fixed := "fixed"
	negPrice := decimal.NewFromInt(-10)
	negWarranty := -1
	invalidGender := "InvalidGender"
	newWarrantyDays := 60
	customExitDate := time.Now().Add(-48 * time.Hour)

	tests := []struct {
		name           string
		id             string
		req            *dto.UpdateTicketRequest
		setupMock      func(m *mockrepo.MockTicketRepository)
		expectedError  error
		expectedAssert func(t *testing.T, res *dto.TicketResponse)
	}{
		{
			name: "success updating fields",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				CustomerName: &newName,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Status:       "service_in",
				}, nil).Once()
				m.On("Update", mock.Anything, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.CustomerName == "Budi Baru"
				})).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "Budi Baru", res.CustomerName)
			},
		},
		{
			name: "error negative price",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Price: &negPrice,
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: ErrNegativePrice,
		},
		{
			name: "error negative warranty days",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				WarrantyDays: &negWarranty,
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: ErrNegativeWarranty,
		},
		{
			name: "validation failure - invalid gender",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				CustomerGender: &invalidGender,
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: errors.New("validation failed"),
		},
		{
			name: "error ticket not found",
			id:   "ticket-not-found",
			req: &dto.UpdateTicketRequest{
				CustomerName: &newName,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-not-found").Return(nil, repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
		{
			name: "repository get error",
			id:   "ticket-error",
			req: &dto.UpdateTicketRequest{
				CustomerName: &newName,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-error").Return(nil, errors.New("db error")).Once()
			},
			expectedError: ErrInternal,
		},
		{
			name: "success transition to picked_up (auto-paid, dates set)",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Status: &pickedUp,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					CustomerName:  "Budi",
					Status:        "service_in",
					PaymentStatus: "unpaid",
					WarrantyDays:  30,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.Status == "picked_up" && t.PaymentStatus == "paid" && t.ExitDate != nil
				})).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "picked_up", res.Status)
				assert.Equal(t, "paid", res.PaymentStatus)
				assert.NotNil(t, res.ExitDate)
				assert.NotNil(t, res.WarrantyExpiryDate)
			},
		},
		{
			name: "error picked_up with explicit unpaid status rejected",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Status:        &pickedUp,
				PaymentStatus: &unpaid,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "service_in",
					PaymentStatus: "unpaid",
					WarrantyDays:  30,
				}, nil).Once()
			},
			expectedError: ErrInvalidPaymentStatus,
		},
		{
			name: "success transition out of picked_up clears dates",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Status: &fixed,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				now := time.Now()
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "picked_up",
					PaymentStatus: "paid",
					WarrantyDays:  30,
					ExitDate:      &now,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.Status == "fixed" && t.ExitDate == nil
				})).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "fixed", res.Status)
				assert.Nil(t, res.ExitDate)
				assert.Nil(t, res.WarrantyExpiryDate)
			},
		},
		{
			name: "repository update error",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				CustomerName: &newName,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
				}, nil).Once()
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("update db error")).Once()
			},
			expectedError: ErrInternal,
		},
		{
			name: "success update warranty_days on already picked_up ticket (recomputes expiry)",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				WarrantyDays: &newWarrantyDays,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				now := time.Now().Add(-24 * time.Hour)
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "picked_up",
					PaymentStatus: "paid",
					WarrantyDays:  30,
					ExitDate:      &now,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.Status == "picked_up" && t.WarrantyDays == 60
				})).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "picked_up", res.Status)
				assert.Equal(t, 60, res.WarrantyDays)
				assert.NotNil(t, res.WarrantyExpiryDate)
			},
		},
		{
			name: "success update exit_date on already picked_up ticket (recomputes expiry)",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				ExitDate: &customExitDate,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				now := time.Now().Add(-24 * time.Hour)
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "picked_up",
					PaymentStatus: "paid",
					WarrantyDays:  30,
					ExitDate:      &now,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.Status == "picked_up" && t.ExitDate.Equal(customExitDate)
				})).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "picked_up", res.Status)
				assert.True(t, res.ExitDate.Equal(customExitDate))
				assert.NotNil(t, res.WarrantyExpiryDate)
			},
		},
		{
			name: "error non-picked_up with exit_date",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Status:   &fixed,
				ExitDate: &customExitDate,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:     "ticket-123",
					Status: "service_in",
				}, nil).Once()
			},
			expectedError: ErrNonPickedUpWithDates,
		},
		{
			name: "error repository update returns not found",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				CustomerName: &newName,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
				}, nil).Once()
				m.On("Update", mock.Anything, mock.Anything).Return(repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			res, err := s.UpdateTicket(context.Background(), tt.id, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError.Error() == "validation failed" {
					// Validation error checked by presence of error
				} else {
					var appErr *AppError
					if errors.As(tt.expectedError, &appErr) {
						assert.ErrorIs(t, err, tt.expectedError)
					} else {
						assert.EqualError(t, err, tt.expectedError.Error())
					}
				}
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if tt.expectedAssert != nil {
					tt.expectedAssert(t, res)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTicketService_ListTickets(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mockrepo.MockTicketRepository)
		expectedError  error
		expectedAssert func(t *testing.T, res []dto.TicketResponse)
	}{
		{
			name: "success list tickets",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("List", mock.Anything).Return([]model.Ticket{
					{
						ID:           "ticket-1",
						CustomerName: "Andi",
					},
					{
						ID:           "ticket-2",
						CustomerName: "Siti",
					},
				}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res []dto.TicketResponse) {
				assert.Len(t, res, 2)
				assert.Equal(t, "ticket-1", res[0].ID)
				assert.Equal(t, "ticket-2", res[1].ID)
			},
		},
		{
			name: "success empty list",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("List", mock.Anything).Return([]model.Ticket{}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res []dto.TicketResponse) {
				assert.Len(t, res, 0)
			},
		},
		{
			name: "repository list error",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("List", mock.Anything).Return(nil, errors.New("list db error")).Once()
			},
			expectedError: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			res, err := s.ListTickets(context.Background())

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
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTicketService_DeleteTicket(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		setupMock     func(m *mockrepo.MockTicketRepository)
		expectedError error
	}{
		{
			name: "success delete",
			id:   "ticket-123",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Delete", mock.Anything, "ticket-123").Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "error ticket not found",
			id:   "ticket-not-found",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Delete", mock.Anything, "ticket-not-found").Return(repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
		{
			name: "repository delete error",
			id:   "ticket-error",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Delete", mock.Anything, "ticket-error").Return(errors.New("delete db error")).Once()
			},
			expectedError: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			err := s.DeleteTicket(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var appErr *AppError
				if errors.As(tt.expectedError, &appErr) {
					assert.ErrorIs(t, err, tt.expectedError)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
