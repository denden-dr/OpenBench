package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	mockrepo "github.com/denden-dr/openbench/apps/backend/mocks/repository"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type stubResult struct{}

func (s *stubResult) LastInsertId() (int64, error) { return 0, nil }
func (s *stubResult) RowsAffected() (int64, error) { return 1, nil }

type stubTx struct {
	committed  bool
	rolledBack bool
}

func (s *stubTx) Commit() error                                                         { s.committed = true; return nil }
func (s *stubTx) Rollback() error                                                       { s.rolledBack = true; return nil }
func (s *stubTx) GetContext(context.Context, interface{}, string, ...interface{}) error { return nil }
func (s *stubTx) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return &stubResult{}, nil
}
func (s *stubTx) QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row {
	panic("stubTx.QueryRowxContext should not be called — repository methods are mocked")
}

func ptrInt(v int) *int {
	return &v
}

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
				CustomerPhone:  "08123456789",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   ptrInt(30),
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(ticket *model.Ticket) bool {
					return ticket.CustomerName == "Budi" &&
						ticket.CustomerPhone == "08123456789" &&
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
				assert.Equal(t, "08123456789", res.CustomerPhone)
				assert.Equal(t, "service_in", res.Status)
				assert.Equal(t, "unpaid", res.PaymentStatus)
			},
		},
		{
			name: "success with default warranty days",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Andi",
				CustomerPhone:  "08123456789",
				CustomerGender: "Male",
				Brand:          "Samsung",
				Model:          "Galaxy S21",
				Issue:          "Baterai kembung",
				Price:          decimal.NewFromInt(500000),
				WarrantyDays:   nil,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(ticket *model.Ticket) bool {
					return ticket.WarrantyDays == 30 && ticket.CustomerPhone == "08123456789"
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
				CustomerPhone:  "08123456789",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(-1),
				WarrantyDays:   ptrInt(30),
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: ErrNegativePrice,
		},
		{
			name: "negative warranty error",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Budi",
				CustomerPhone:  "08123456789",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   ptrInt(-5),
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: ErrNegativeWarranty,
		},
		{
			name: "validation failure - empty customer name",
			req: &dto.CreateTicketRequest{
				CustomerName:   "",
				CustomerPhone:  "08123456789",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   ptrInt(30),
			},
			setupMock:     func(m *mockrepo.MockTicketRepository) {},
			expectedError: errors.New("validation failed"),
		},
		{
			name: "repository error",
			req: &dto.CreateTicketRequest{
				CustomerName:   "Budi",
				CustomerPhone:  "08123456789",
				CustomerGender: "Male",
				Brand:          "Apple",
				Model:          "iPhone 13",
				Issue:          "LCD Mati",
				Price:          decimal.NewFromInt(1500000),
				WarrantyDays:   ptrInt(30),
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
	statusWaiting := "waiting_confirmation"
	statusCancelled := "cancelled"

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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Status:       "service_in",
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.MatchedBy(func(t *model.Ticket) bool {
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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-not-found").Return(nil, repository.ErrNotFound).Once()
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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-error").Return(nil, errors.New("db error")).Once()
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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					CustomerName:  "Budi",
					Status:        "fixed",
					PaymentStatus: "unpaid",
					WarrantyDays:  30,
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.MatchedBy(func(t *model.Ticket) bool {
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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "fixed",
					PaymentStatus: "unpaid",
					WarrantyDays:  30,
				}, nil).Once()
			},
			expectedError: ErrInvalidPaymentStatus,
		},
		{
			name: "error transition out of picked_up blocked",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Status: &fixed,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				now := time.Now()
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "picked_up",
					PaymentStatus: "paid",
					WarrantyDays:  30,
					ExitDate:      &now,
				}, nil).Once()
			},
			expectedError: ErrInvalidStatusTransition,
		},
		{
			name: "repository update error",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				CustomerName: &newName,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.Anything).Return(errors.New("update db error")).Once()
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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "picked_up",
					PaymentStatus: "paid",
					WarrantyDays:  30,
					ExitDate:      &now,
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.MatchedBy(func(t *model.Ticket) bool {
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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:            "ticket-123",
					Status:        "picked_up",
					PaymentStatus: "paid",
					WarrantyDays:  30,
					ExitDate:      &now,
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.MatchedBy(func(t *model.Ticket) bool {
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
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:     "ticket-123",
					Status: "on_process",
				}, nil).Once()
			},
			expectedError: ErrNonPickedUpWithExitDate,
		},
		{
			name: "error repository update returns not found",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				CustomerName: &newName,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.Anything).Return(repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
		{
			name: "success transition to waiting_confirmation",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Status: &statusWaiting,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Status:       "on_process",
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.Status == "waiting_confirmation"
				})).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "waiting_confirmation", res.Status)
			},
		},
		{
			name: "success transition to cancelled",
			id:   "ticket-123",
			req: &dto.UpdateTicketRequest{
				Status: &statusCancelled,
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				tx := &stubTx{}
				m.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				m.On("GetByIDForUpdateTx", mock.Anything, tx, "ticket-123").Return(&model.Ticket{
					ID:           "ticket-123",
					CustomerName: "Budi",
					Status:       "waiting_confirmation",
				}, nil).Once()
				m.On("UpdateTx", mock.Anything, tx, mock.MatchedBy(func(t *model.Ticket) bool {
					return t.Status == "cancelled"
				})).Return(nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.TicketResponse) {
				assert.Equal(t, "cancelled", res.Status)
			},
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
		page           int
		limit          int
		search         string
		status         string
		setupMock      func(m *mockrepo.MockTicketRepository)
		expectedError  error
		expectedAssert func(t *testing.T, res *dto.PaginatedTicketsResult)
	}{
		{
			name:   "success list tickets",
			page:   1,
			limit:  20,
			search: "",
			status: "all",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("CountPaginated", mock.Anything, "", "all").Return(int64(2), nil).Once()
				m.On("ListPaginated", mock.Anything, "", "all", 20, 0).Return([]model.Ticket{
					{
						ID:           "ticket-1",
						CustomerName: "Andi",
					},
					{
						ID:           "ticket-2",
						CustomerName: "Siti",
					},
				}, nil).Once()
				m.On("GetStatusCounts", mock.Anything, "").Return(map[string]int64{
					"service_in": 2,
				}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.PaginatedTicketsResult) {
				assert.Equal(t, int64(2), res.Total)
				assert.Equal(t, int64(1), res.TotalPages)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, "ticket-1", res.Data[0].ID)
				assert.Equal(t, "ticket-2", res.Data[1].ID)
				assert.Equal(t, int64(2), res.StatusCounts["all"])
				assert.Equal(t, int64(2), res.StatusCounts["service_in"])
			},
		},
		{
			name:   "success empty list",
			page:   1,
			limit:  20,
			search: "",
			status: "all",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("CountPaginated", mock.Anything, "", "all").Return(int64(0), nil).Once()
				m.On("ListPaginated", mock.Anything, "", "all", 20, 0).Return([]model.Ticket{}, nil).Once()
				m.On("GetStatusCounts", mock.Anything, "").Return(map[string]int64{}, nil).Once()
			},
			expectedError: nil,
			expectedAssert: func(t *testing.T, res *dto.PaginatedTicketsResult) {
				assert.Equal(t, int64(0), res.Total)
				assert.Len(t, res.Data, 0)
			},
		},
		{
			name:   "invalid status",
			page:   1,
			limit:  20,
			search: "",
			status: "invalid_status",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				// No calls to repo since status validation happens first in service
			},
			expectedError: ErrInvalidStatus,
		},
		{
			name:   "repository count error",
			page:   1,
			limit:  20,
			search: "",
			status: "all",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("CountPaginated", mock.Anything, "", "all").Return(int64(0), errors.New("count db error")).Once()
				m.On("ListPaginated", mock.Anything, "", "all", 20, 0).Return([]model.Ticket{}, nil).Maybe()
				m.On("GetStatusCounts", mock.Anything, "").Return(map[string]int64{}, nil).Maybe()
			},
			expectedError: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			res, err := s.ListTickets(context.Background(), tt.page, tt.limit, tt.search, tt.status)

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

func TestTicketService_GetPublicTicket(t *testing.T) {
	entryDate := time.Now()
	t1 := model.Ticket{
		ID:             "abcdef12-3456-7890-abcd-ef1234567890",
		CustomerName:   "Budi Anto",
		CustomerPhone:  "081234567890",
		CustomerGender: "Male",
		Brand:          "Apple",
		Model:          "iPhone 13",
		Issue:          "LCD Mati",
		Status:         "service_in",
		PaymentStatus:  "unpaid",
		WarrantyDays:   30,
		EntryDate:      entryDate,
	}

	tests := []struct {
		name          string
		id            string
		setupMock     func(m *mockrepo.MockTicketRepository)
		expectedError error
		expectedRes   *dto.PublicTicketResponse
	}{
		{
			name: "success get by full UUID",
			id:   "abcdef12-3456-7890-abcd-ef1234567890",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "abcdef12-3456-7890-abcd-ef1234567890").Return(&t1, nil).Once()
			},
			expectedError: nil,
			expectedRes: &dto.PublicTicketResponse{
				ID:                  "abcdef12-3456-7890-abcd-ef1234567890",
				CustomerNameMasked:  "B*** A***",
				CustomerPhoneMasked: "0812******90",
				Brand:               "Apple",
				Model:               "iPhone 13",
				Issue:               "LCD Mati",
				Status:              "service_in",
				EntryDate:           entryDate,
			},
		},
		{
			name: "fail get by short ID",
			id:   "abcdef12",
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByID", mock.Anything, "abcdef12").Return(nil, repository.ErrNotFound).Once()
			},
			expectedError: ErrTicketNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			res, err := s.GetPublicTicket(context.Background(), tt.id)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes.ID, res.ID)
				assert.Equal(t, tt.expectedRes.CustomerNameMasked, res.CustomerNameMasked)
				assert.Equal(t, tt.expectedRes.CustomerPhoneMasked, res.CustomerPhoneMasked)
				assert.Equal(t, tt.expectedRes.Brand, res.Brand)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTicketService_TrackPublicTicket(t *testing.T) {
	t1 := model.Ticket{
		ID:            "abcdef12-3456-7890-abcd-ef1234567890",
		CustomerPhone: "+62 812-3456-7890",
	}

	tests := []struct {
		name          string
		req           *dto.PublicTrackRequest
		setupMock     func(m *mockrepo.MockTicketRepository)
		expectedError error
		expectedID    string
	}{
		{
			name: "success track match",
			req: &dto.PublicTrackRequest{
				ShortID: "abcdef12",
				Phone:   "081234567890",
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByShortID", mock.Anything, "abcdef12").Return([]model.Ticket{t1}, nil).Once()
			},
			expectedError: nil,
			expectedID:    "abcdef12-3456-7890-abcd-ef1234567890",
		},
		{
			name: "track failed - phone mismatch",
			req: &dto.PublicTrackRequest{
				ShortID: "abcdef12",
				Phone:   "081299999999",
			},
			setupMock: func(m *mockrepo.MockTicketRepository) {
				m.On("GetByShortID", mock.Anything, "abcdef12").Return([]model.Ticket{t1}, nil).Once()
			},
			expectedError: ErrTicketNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockTicketRepository)
			tt.setupMock(mockRepo)

			s := NewTicketService(mockRepo)
			id, err := s.TrackPublicTicket(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNormalizePhone(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"+62 812-3456-789", "08123456789"},
		{"628123456789", "08123456789"},
		{"08123456789", "08123456789"},
		{"  0812 3456  ", "08123456"},
		{"abc-123", "123"},
	}

	for _, tc := range cases {
		result := normalizePhone(tc.input)
		if result != tc.expected {
			t.Errorf("normalizePhone(%q) = %q; expected %q", tc.input, result, tc.expected)
		}
	}
}
