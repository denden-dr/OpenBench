package ticket

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, t *models.ServiceTicket) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *mockRepository) FindAll(ctx context.Context, status string, search string, limit, offset int) ([]models.ServiceTicket, int, error) {
	args := m.Called(ctx, status, search, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]models.ServiceTicket), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *mockRepository) FindByID(ctx context.Context, id string) (*models.ServiceTicket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.ServiceTicket), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) Update(ctx context.Context, t *models.ServiceTicket) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

var errDb = errors.New("db connection failure")

func TestService_CreateTicket(t *testing.T) {
	tests := []struct {
		name        string
		req         CreateTicketRequest
		mockErr     error
		expectedErr error
	}{
		{
			name: "Success - complete payload",
			req: CreateTicketRequest{
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				DevicePasscode:   "1234",
				IssueDescription: "Layar pecah",
				RepairAction:     "Ganti LCD",
				Cost:             1500000,
				WarrantyDays:     30,
			},
			mockErr:     nil,
			expectedErr: nil,
		},
		{
			name: "Failure - empty customer name",
			req: CreateTicketRequest{
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
			},
			mockErr:     nil,
			expectedErr: ErrInvalidInput,
		},
		{
			name: "Failure - empty customer phone",
			req: CreateTicketRequest{
				CustomerName:     "Budi Santoso",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
			},
			mockErr:     nil,
			expectedErr: ErrInvalidInput,
		},
		{
			name: "Failure - repo error",
			req: CreateTicketRequest{
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
			},
			mockErr:     errDb,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			if tt.expectedErr != nil {
				if tt.mockErr != nil {
					repo.On("Create", mock.Anything, mock.AnythingOfType("*models.ServiceTicket")).Return(tt.mockErr)
				}
			} else {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.ServiceTicket")).Return(nil)
			}
			svc := NewService(repo)

			res, err := svc.CreateTicket(context.Background(), tt.req)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
				is.NotEmpty(res.TicketID)
				is.True(strings.HasPrefix(res.TicketNumber, "TKT-"))
				is.Equal(models.StatusReceived, res.Status)
				is.Equal(tt.req.CustomerName, res.CustomerName)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_UpdateTicketStatus(t *testing.T) {
	existingTicket := &models.ServiceTicket{
		ID:           "ticket-1",
		TicketNumber: "TKT-20260708-ABCD",
		Status:       models.StatusReceived,
	}

	tests := []struct {
		name        string
		ticketID    string
		req         ChangeStatusRequest
		setupMock   func(repo *mockRepository)
		expectedErr error
		checkStatus models.ServiceTicketStatus
	}{
		{
			name:     "Success - Change status to REPAIRING",
			ticketID: "ticket-1",
			req: ChangeStatusRequest{
				Status: models.StatusRepairing,
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindByID", mock.Anything, "ticket-1").Return(existingTicket, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(t *models.ServiceTicket) bool {
					return t.Status == models.StatusRepairing
				})).Return(nil)
			},
			expectedErr: nil,
			checkStatus: models.StatusRepairing,
		},
		{
			name:     "Failure - Ticket Not Found",
			ticketID: "non-existent",
			req: ChangeStatusRequest{
				Status: models.StatusRepairing,
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
			},
			expectedErr: ErrTicketNotFound,
		},
		{
			name:     "Failure - Empty Status",
			ticketID: "ticket-1",
			req: ChangeStatusRequest{
				Status: "",
			},
			setupMock:   func(repo *mockRepository) {},
			expectedErr: ErrInvalidInput,
		},
		{
			name:     "Failure - Invalid Status Enum",
			ticketID: "ticket-1",
			req: ChangeStatusRequest{
				Status: "INVALID_STATUS",
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

			res, err := svc.UpdateTicketStatus(context.Background(), tt.ticketID, tt.req)
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

func TestService_UpdateTicketDetails(t *testing.T) {
	existingTicket := &models.ServiceTicket{
		ID:               "ticket-1",
		TicketNumber:     "TKT-20260708-ABCD",
		Status:           models.StatusReceived,
		CustomerName:     "Budi Santoso",
		CustomerPhone:    "081234567890",
		IssueDescription: "Layar pecah",
	}

	tests := []struct {
		name        string
		ticketID    string
		req         UpdateTicketRequest
		setupMock   func(repo *mockRepository)
		expectedErr error
	}{
		{
			name:     "Success - Update details",
			ticketID: "ticket-1",
			req: UpdateTicketRequest{
				CustomerName:     "Budi Santoso Baru",
				CustomerPhone:    "081234567899",
				IssueDescription: "Layar pecah & touch error",
				RepairAction:     "Ganti LCD Fullset",
				Cost:             1800000,
				WarrantyDays:     60,
				Notes:            "LCD Original",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindByID", mock.Anything, "ticket-1").Return(existingTicket, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(t *models.ServiceTicket) bool {
					return t.CustomerName == "Budi Santoso Baru" &&
						t.CustomerPhone == "081234567899" &&
						t.IssueDescription == "Layar pecah & touch error" &&
						*t.RepairAction == "Ganti LCD Fullset" &&
						t.Cost == 1800000 &&
						t.WarrantyDays == 60 &&
						*t.Notes == "LCD Original"
				})).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "Failure - Empty Customer Name",
			ticketID: "ticket-1",
			req: UpdateTicketRequest{
				CustomerPhone:    "081234567899",
				IssueDescription: "Layar pecah",
			},
			setupMock:   func(repo *mockRepository) {},
			expectedErr: ErrInvalidInput,
		},
		{
			name:     "Failure - Ticket Not Found",
			ticketID: "non-existent",
			req: UpdateTicketRequest{
				CustomerName:     "Budi Santoso Baru",
				CustomerPhone:    "081234567899",
				IssueDescription: "Layar pecah",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
			},
			expectedErr: ErrTicketNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			tt.setupMock(repo)
			svc := NewService(repo)

			res, err := svc.UpdateTicketDetails(context.Background(), tt.ticketID, tt.req)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
				is.Equal(tt.req.CustomerName, res.CustomerName)
				is.Equal(tt.req.CustomerPhone, res.CustomerPhone)
				is.Equal(tt.req.RepairAction, *res.RepairAction)
				is.Equal(tt.req.Cost, res.Cost)
				is.Equal(tt.req.WarrantyDays, res.WarrantyDays)
				is.Equal(tt.req.Notes, *res.Notes)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_EmergencyUpdateTicket(t *testing.T) {
	existingTicket := &models.ServiceTicket{
		ID:               "ticket-1",
		TicketNumber:     "TKT-20260708-ABCD",
		Status:           models.StatusReceived,
		CustomerName:     "Budi Santoso",
		CustomerPhone:    "081234567890",
		DeviceBrand:      "Samsung",
		DeviceModel:      "Galaxy S23",
		IssueDescription: "Layar pecah",
	}

	tests := []struct {
		name        string
		ticketID    string
		req         EmergencyUpdateTicketRequest
		setupMock   func(repo *mockRepository)
		expectedErr error
	}{
		{
			name:     "Success - Emergency Update status and brand/model",
			ticketID: "ticket-1",
			req: EmergencyUpdateTicketRequest{
				CustomerName:     "Budi Santoso Baru",
				CustomerPhone:    "081234567899",
				DeviceBrand:      "Apple",
				DeviceModel:      "iPhone 15 Pro",
				DevicePasscode:   "9999",
				Status:           models.StatusCompleted,
				IssueDescription: "Layar pecah & touch error",
				RepairAction:     "Ganti LCD Fullset",
				Cost:             3000000,
				WarrantyDays:     90,
				Notes:            "LCD Original Apple",
			},
			setupMock: func(repo *mockRepository) {
				repo.On("FindByID", mock.Anything, "ticket-1").Return(existingTicket, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(t *models.ServiceTicket) bool {
					return t.DeviceBrand == "Apple" &&
						t.DeviceModel == "iPhone 15 Pro" &&
						t.Status == models.StatusCompleted
				})).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "Failure - Empty Status",
			ticketID: "ticket-1",
			req: EmergencyUpdateTicketRequest{
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				Status:           "",
				IssueDescription: "Layar pecah",
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

			res, err := svc.EmergencyUpdateTicket(context.Background(), tt.ticketID, tt.req)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(res)
				is.Equal(tt.req.DeviceBrand, res.DeviceBrand)
				is.Equal(tt.req.DeviceModel, res.DeviceModel)
				is.Equal(tt.req.Status, res.Status)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestService_GetTicketsAndByID(t *testing.T) {
	ticket1 := &models.ServiceTicket{
		ID:            "ticket-1",
		TicketNumber:  "TKT-20260708-ABCD",
		Status:        models.StatusReceived,
		CustomerName:  "Budi Santoso",
		CustomerPhone: "081234567890",
		DeviceBrand:   "Samsung",
		DeviceModel:   "Galaxy S23",
	}
	ticket2 := &models.ServiceTicket{
		ID:            "ticket-2",
		TicketNumber:  "TKT-20260708-EFGH",
		Status:        models.StatusRepairing,
		CustomerName:  "Joko Widodo",
		CustomerPhone: "089999999",
		DeviceBrand:   "Apple",
		DeviceModel:   "iPhone 15",
	}

	t.Run("Get Ticket By ID - Success", func(t *testing.T) {
		is := assert.New(t)
		must := require.New(t)

		repo := &mockRepository{}
		repo.On("FindByID", mock.Anything, "ticket-1").Return(ticket1, nil)
		svc := NewService(repo)

		res, err := svc.GetTicketByID(context.Background(), "ticket-1")
		must.NoError(err)
		must.NotNil(res)
		is.Equal(ticket1.TicketNumber, res.TicketNumber)
		repo.AssertExpectations(t)
	})

	t.Run("Get Ticket By ID - Not Found", func(t *testing.T) {
		is := assert.New(t)

		repo := &mockRepository{}
		repo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
		svc := NewService(repo)

		_, err := svc.GetTicketByID(context.Background(), "non-existent")
		is.ErrorIs(err, ErrTicketNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("Get Tickets - Filter Status", func(t *testing.T) {
		is := assert.New(t)
		must := require.New(t)

		repo := &mockRepository{}
		repo.On("FindAll", mock.Anything, "REPAIRING", "", 10, 0).Return([]models.ServiceTicket{*ticket2}, 1, nil)
		svc := NewService(repo)

		res, total, err := svc.GetTickets(context.Background(), "REPAIRING", "", 10, 0)
		must.NoError(err)
		is.Equal(1, total)
		must.Len(res, 1)
		is.Equal("ticket-2", res[0].TicketID)
		repo.AssertExpectations(t)
	})

	t.Run("Get Tickets - Search Query", func(t *testing.T) {
		is := assert.New(t)
		must := require.New(t)

		repo := &mockRepository{}
		repo.On("FindAll", mock.Anything, "", "Joko", 10, 0).Return([]models.ServiceTicket{*ticket2}, 1, nil)
		svc := NewService(repo)

		res, total, err := svc.GetTickets(context.Background(), "", "Joko", 10, 0)
		must.NoError(err)
		is.Equal(1, total)
		must.Len(res, 1)
		is.Equal("Joko Widodo", res[0].CustomerName)
		repo.AssertExpectations(t)
	})
}
