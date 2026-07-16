package ticket

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/events"
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

func (m *mockRepository) FindAll(ctx context.Context, status string, search string, limit int, cursor string) ([]models.ServiceTicket, string, error) {
	args := m.Called(ctx, status, search, limit, cursor)
	if args.Get(0) != nil {
		return args.Get(0).([]models.ServiceTicket), args.String(1), args.Error(2)
	}
	return nil, args.String(1), args.Error(2)
}

func (m *mockRepository) Search(ctx context.Context, req TicketSearchRequest) ([]models.ServiceTicket, string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).([]models.ServiceTicket), args.String(1), args.Error(2)
	}
	return nil, args.String(1), args.Error(2)
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

type mockEventBus struct {
	mock.Mock
}

func (m *mockEventBus) Publish(ctx context.Context, event events.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockEventBus) Subscribe(eventType events.EventType, handler events.EventHandler) {
	m.Called(eventType, handler)
}

type mockWarrantyGenerator struct {
	mock.Mock
}

func (m *mockWarrantyGenerator) CreateWarranty(ctx context.Context, ticketID string, warrantyDays int) (*models.Warranty, error) {
	args := m.Called(ctx, ticketID, warrantyDays)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Warranty), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockTxManager struct{}

func (m *mockTxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
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
		{
			name: "Failure - negative cost",
			req: CreateTicketRequest{
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
				Cost:             -1000,
			},
			mockErr:     nil,
			expectedErr: ErrInvalidInput,
		},
		{
			name: "Failure - negative warranty days",
			req: CreateTicketRequest{
				CustomerName:     "Budi Santoso",
				CustomerPhone:    "081234567890",
				DeviceBrand:      "Samsung",
				DeviceModel:      "Galaxy S23",
				IssueDescription: "Layar pecah",
				WarrantyDays:     -5,
			},
			mockErr:     nil,
			expectedErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			bus := &mockEventBus{}
			wgen := &mockWarrantyGenerator{}
			if tt.expectedErr != nil {
				if tt.mockErr != nil {
					repo.On("Create", mock.Anything, mock.AnythingOfType("*models.ServiceTicket")).Return(tt.mockErr)
				}
			} else {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.ServiceTicket")).Return(nil)
			}
			svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

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
			bus.AssertExpectations(t)
			wgen.AssertExpectations(t)
		})
	}
}

func TestService_UpdateTicketStatus(t *testing.T) {

	tests := []struct {
		name        string
		ticketID    string
		req         ChangeStatusRequest
		setupMock   func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator)
		expectedErr error
		checkStatus models.ServiceTicketStatus
	}{
		{
			name:     "Success - Change status to REPAIRING",
			ticketID: "ticket-1",
			req: ChangeStatusRequest{
				Status: models.StatusRepairing,
			},
			setupMock: func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {
				repo.On("FindByID", mock.Anything, "ticket-1").Return(&models.ServiceTicket{
					ID:           "ticket-1",
					TicketNumber: "TKT-20260708-ABCD",
					Status:       models.StatusReceived,
					WarrantyDays: 0,
				}, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(t *models.ServiceTicket) bool {
					return t.Status == models.StatusRepairing
				})).Return(nil)
			},
			expectedErr: nil,
			checkStatus: models.StatusRepairing,
		},
		{
			name:     "Success - Change status to COMPLETED with warranty triggers synchronous creation",
			ticketID: "ticket-2",
			req: ChangeStatusRequest{
				Status: models.StatusCompleted,
			},
			setupMock: func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {
				repo.On("FindByID", mock.Anything, "ticket-2").Return(&models.ServiceTicket{
					ID:           "ticket-2",
					TicketNumber: "TKT-20260708-WARR",
					Status:       models.StatusReceived,
					WarrantyDays: 30,
				}, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(t *models.ServiceTicket) bool {
					return t.Status == models.StatusCompleted
				})).Return(nil)
				wgen.On("CreateWarranty", mock.Anything, "ticket-2", 30).Return(&models.Warranty{}, nil)
			},
			expectedErr: nil,
			checkStatus: models.StatusCompleted,
		},
		{
			name:     "Failure - Ticket Not Found",
			ticketID: "non-existent",
			req: ChangeStatusRequest{
				Status: models.StatusRepairing,
			},
			setupMock: func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {
				repo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
			},
			expectedErr: ErrTicketNotFound,
		},
		{
			name:     "Failure - Same Status",
			ticketID: "ticket-1",
			req: ChangeStatusRequest{
				Status: models.StatusReceived,
			},
			setupMock: func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {
				repo.On("FindByID", mock.Anything, "ticket-1").Return(&models.ServiceTicket{
					ID:           "ticket-1",
					TicketNumber: "TKT-20260708-ABCD",
					Status:       models.StatusReceived,
					WarrantyDays: 0,
				}, nil)
			},
			expectedErr: ErrInvalidInput,
		},
		{
			name:     "Failure - Empty Status",
			ticketID: "ticket-1",
			req: ChangeStatusRequest{
				Status: "",
			},
			setupMock:   func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {},
			expectedErr: ErrInvalidInput,
		},
		{
			name:     "Failure - Invalid Status Enum",
			ticketID: "ticket-1",
			req: ChangeStatusRequest{
				Status: "INVALID_STATUS",
			},
			setupMock:   func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {},
			expectedErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			bus := &mockEventBus{}
			wgen := &mockWarrantyGenerator{}
			tt.setupMock(repo, bus, wgen)
			svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

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
			bus.AssertExpectations(t)
			wgen.AssertExpectations(t)
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
		setupMock   func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator)
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
			setupMock: func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {
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
			setupMock:   func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {},
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
			setupMock: func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {
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
			bus := &mockEventBus{}
			wgen := &mockWarrantyGenerator{}
			tt.setupMock(repo, bus, wgen)
			svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

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
			bus.AssertExpectations(t)
			wgen.AssertExpectations(t)
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
		WarrantyDays:     30,
	}

	tests := []struct {
		name        string
		ticketID    string
		req         EmergencyUpdateTicketRequest
		setupMock   func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator)
		expectedErr error
	}{
		{
			name:     "Success - Emergency Update status to COMPLETED triggers synchronous creation",
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
			setupMock: func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {
				repo.On("FindByID", mock.Anything, "ticket-1").Return(existingTicket, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(t *models.ServiceTicket) bool {
					return t.DeviceBrand == "Apple" &&
						t.DeviceModel == "iPhone 15 Pro" &&
						t.Status == models.StatusCompleted
				})).Return(nil)
				wgen.On("CreateWarranty", mock.Anything, "ticket-1", 90).Return(&models.Warranty{}, nil)
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
			setupMock:   func(repo *mockRepository, bus *mockEventBus, wgen *mockWarrantyGenerator) {},
			expectedErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			bus := &mockEventBus{}
			wgen := &mockWarrantyGenerator{}
			tt.setupMock(repo, bus, wgen)
			svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

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
			bus.AssertExpectations(t)
			wgen.AssertExpectations(t)
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
		bus := &mockEventBus{}
		wgen := &mockWarrantyGenerator{}
		repo.On("FindByID", mock.Anything, "ticket-1").Return(ticket1, nil)
		svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

		res, err := svc.GetTicketByID(context.Background(), "ticket-1")
		must.NoError(err)
		must.NotNil(res)
		is.Equal(ticket1.TicketNumber, res.TicketNumber)
		repo.AssertExpectations(t)
		bus.AssertExpectations(t)
	})

	t.Run("Get Ticket By ID - Not Found", func(t *testing.T) {
		is := assert.New(t)

		repo := &mockRepository{}
		bus := &mockEventBus{}
		wgen := &mockWarrantyGenerator{}
		repo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
		svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

		_, err := svc.GetTicketByID(context.Background(), "non-existent")
		is.ErrorIs(err, ErrTicketNotFound)
		repo.AssertExpectations(t)
		bus.AssertExpectations(t)
	})

	t.Run("Get Tickets - Filter Status", func(t *testing.T) {
		is := assert.New(t)
		must := require.New(t)

		repo := &mockRepository{}
		bus := &mockEventBus{}
		wgen := &mockWarrantyGenerator{}
		repo.On("FindAll", mock.Anything, "REPAIRING", "", 10, "").Return([]models.ServiceTicket{*ticket2}, "next-cursor-123", nil)
		svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

		res, nextCursor, err := svc.GetTickets(context.Background(), "REPAIRING", "", 10, "")
		must.NoError(err)
		is.Equal("next-cursor-123", nextCursor)
		must.Len(res, 1)
		is.Equal("ticket-2", res[0].TicketID)
		repo.AssertExpectations(t)
		bus.AssertExpectations(t)
	})

	t.Run("Get Tickets - Search Query", func(t *testing.T) {
		is := assert.New(t)
		must := require.New(t)

		repo := &mockRepository{}
		bus := &mockEventBus{}
		wgen := &mockWarrantyGenerator{}
		repo.On("FindAll", mock.Anything, "", "Joko", 10, "").Return([]models.ServiceTicket{*ticket2}, "", nil)
		svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

		res, nextCursor, err := svc.GetTickets(context.Background(), "", "Joko", 10, "")
		must.NoError(err)
		is.Empty(nextCursor)
		must.Len(res, 1)
		is.Equal("Joko Widodo", res[0].CustomerName)
		repo.AssertExpectations(t)
		bus.AssertExpectations(t)
	})
}

func TestService_SearchTickets(t *testing.T) {
	ticket1 := &models.ServiceTicket{
		ID:            "ticket-1",
		TicketNumber:  "TKT-20260708-ABCD",
		Status:        models.StatusReceived,
		CustomerName:  "Budi Santoso",
		CustomerPhone: "081234567890",
		DeviceBrand:   "Samsung",
		DeviceModel:   "Galaxy S23",
	}

	t.Run("Search Tickets - Text Match", func(t *testing.T) {
		is := assert.New(t)
		must := require.New(t)

		repo := &mockRepository{}
		bus := &mockEventBus{}
		wgen := &mockWarrantyGenerator{}
		req := TicketSearchRequest{
			Search: "Samsung",
			Limit:  10,
			Cursor: "",
		}
		repo.On("Search", mock.Anything, req).Return([]models.ServiceTicket{*ticket1}, "next-cursor-456", nil)
		svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

		res, nextCursor, err := svc.SearchTickets(context.Background(), req)
		must.NoError(err)
		is.Equal("next-cursor-456", nextCursor)
		must.Len(res, 1)
		is.Equal("Samsung", res[0].DeviceBrand)
		repo.AssertExpectations(t)
	})

	t.Run("Search Tickets - Active Status", func(t *testing.T) {
		is := assert.New(t)
		must := require.New(t)

		repo := &mockRepository{}
		bus := &mockEventBus{}
		wgen := &mockWarrantyGenerator{}
		isActive := true
		req := TicketSearchRequest{
			IsActive: &isActive,
			Limit:    10,
			Cursor:   "",
		}
		repo.On("Search", mock.Anything, req).Return([]models.ServiceTicket{*ticket1}, "", nil)
		svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

		res, nextCursor, err := svc.SearchTickets(context.Background(), req)
		must.NoError(err)
		is.Empty(nextCursor)
		must.Len(res, 1)
		repo.AssertExpectations(t)
	})
}

func TestService_CreateTicket_Encryption(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	key := "this_is_a_secret_key_32_chars_ok"
	rawPasscode := "my-secret-passcode-123"

	repo := &mockRepository{}
	bus := &mockEventBus{}
	wgen := &mockWarrantyGenerator{}

	// Verify that the model passed to repo.Create has an encrypted passcode
	repo.On("Create", mock.Anything, mock.MatchedBy(func(tick *models.ServiceTicket) bool {
		return tick.DevicePasscode != rawPasscode && len(tick.DevicePasscode) > len(rawPasscode)
	})).Return(nil)

	svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, key)

	req := CreateTicketRequest{
		CustomerName:     "Budi Santoso",
		CustomerPhone:    "081234567890",
		DeviceBrand:      "Samsung",
		DeviceModel:      "Galaxy S23",
		DevicePasscode:   rawPasscode,
		IssueDescription: "Layar pecah",
	}

	res, err := svc.CreateTicket(context.Background(), req)
	must.NoError(err)
	must.NotNil(res)

	// Verify that the returned passcode to client is decrypted back to raw
	is.Equal(rawPasscode, res.DevicePasscode)
	repo.AssertExpectations(t)
}

func TestService_UpdateTicketStatus_Rollback(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	existingTicket := &models.ServiceTicket{
		ID:           "ticket-rollback",
		TicketNumber: "TKT-ROLLBACK",
		Status:       models.StatusReceived,
		WarrantyDays: 30,
	}

	repo := &mockRepository{}
	bus := &mockEventBus{}
	wgen := &mockWarrantyGenerator{}

	// Mock getting the ticket initially
	repo.On("FindByID", mock.Anything, "ticket-rollback").Return(existingTicket, nil).Once()

	// Mock the update operation in the transaction (which will be executed first)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(t *models.ServiceTicket) bool {
		return t.Status == models.StatusCompleted
	})).Return(nil)

	// Mock the warranty generation to fail
	errWarranty := errors.New("warranty generation failed")
	wgen.On("CreateWarranty", mock.Anything, "ticket-rollback", 30).Return(nil, errWarranty)

	svc := NewService(repo, repo, &mockTxManager{}, wgen, bus, "this_is_a_secret_key_32_chars_ok")

	// Call UpdateTicketStatus
	res, err := svc.UpdateTicketStatus(context.Background(), "ticket-rollback", ChangeStatusRequest{Status: models.StatusCompleted})

	// Verify that error from CreateWarranty is propagated and res is zero
	must.Error(err)
	is.Contains(err.Error(), "failed to create warranty within transaction")
	is.Zero(res)

	repo.AssertExpectations(t)
	wgen.AssertExpectations(t)
}
