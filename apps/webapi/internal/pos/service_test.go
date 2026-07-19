package pos

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPosRepo struct {
	mock.Mock
}

func (m *mockPosRepo) FindByID(ctx context.Context, id string) (*models.PosTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.PosTransaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPosRepo) FindAll(ctx context.Context, limit int, cursor string) ([]models.PosTransaction, string, error) {
	args := m.Called(ctx, limit, cursor)
	if args.Get(0) != nil {
		return args.Get(0).([]models.PosTransaction), args.String(1), args.Error(2)
	}
	return nil, args.String(1), args.Error(2)
}

func (m *mockPosRepo) Create(ctx context.Context, t *models.PosTransaction) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

type mockInventoryQueryRepo struct {
	mock.Mock
}

func (m *mockInventoryQueryRepo) FindByID(ctx context.Context, id string) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryQueryRepo) FindAll(ctx context.Context, search string, limit int, cursor string) ([]models.Product, string, error) {
	args := m.Called(ctx, search, limit, cursor)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Product), args.String(1), args.Error(2)
	}
	return nil, args.String(1), args.Error(2)
}

type mockInventoryCommandRepo struct {
	mock.Mock
}

func (m *mockInventoryCommandRepo) Create(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockInventoryCommandRepo) Update(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockInventoryCommandRepo) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	args := m.Called(ctx, id, quantityChange)
	return args.Error(0)
}

func (m *mockInventoryCommandRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockTxManager struct{}

func (m *mockTxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestService_Checkout_Success(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	posRepo := &mockPosRepo{}
	invQueryRepo := &mockInventoryQueryRepo{}
	invCommandRepo := &mockInventoryCommandRepo{}
	txManager := &mockTxManager{}

	svc := NewService(posRepo, posRepo, invQueryRepo, invCommandRepo, txManager)

	product := &models.Product{
		ID:    "prod-1",
		Name:  "Tempered Glass",
		Price: 50000,
		Stock: 10,
	}

	req := models.CheckoutRequest{
		PaymentMethod: models.PaymentMethodCash,
		Items: []models.CheckoutItemRequest{
			{
				ProductID: "prod-1",
				Quantity:  2,
			},
		},
	}

	invQueryRepo.On("FindByID", mock.Anything, "prod-1").Return(product, nil)
	invCommandRepo.On("UpdateStock", mock.Anything, "prod-1", -2).Return(nil)
	posRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.PosTransaction")).Return(nil)

	expectedTx := &models.PosTransaction{
		ID:            "tx-123",
		PaymentMethod: models.PaymentMethodCash,
		TotalAmount:   100000,
	}
	posRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(expectedTx, nil)

	res, err := svc.Checkout(context.Background(), req)
	must.NoError(err)
	must.NotNil(res)
	is.Equal(int64(100000), res.TotalAmount)
	is.Equal(models.PaymentMethodCash, res.PaymentMethod)

	posRepo.AssertExpectations(t)
	invQueryRepo.AssertExpectations(t)
	invCommandRepo.AssertExpectations(t)
}

func TestService_Checkout_InsufficientStock(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	posRepo := &mockPosRepo{}
	invQueryRepo := &mockInventoryQueryRepo{}
	invCommandRepo := &mockInventoryCommandRepo{}
	txManager := &mockTxManager{}

	svc := NewService(posRepo, posRepo, invQueryRepo, invCommandRepo, txManager)

	product := &models.Product{
		ID:    "prod-1",
		Name:  "Tempered Glass",
		Price: 50000,
		Stock: 1, // Only 1 in stock
	}

	req := models.CheckoutRequest{
		PaymentMethod: models.PaymentMethodCash,
		Items: []models.CheckoutItemRequest{
			{
				ProductID: "prod-1",
				Quantity:  2, // Request 2
			},
		},
	}

	invQueryRepo.On("FindByID", mock.Anything, "prod-1").Return(product, nil)

	res, err := svc.Checkout(context.Background(), req)
	must.Error(err)
	is.ErrorIs(err, ErrInsufficientStock)
	is.Nil(res)

	posRepo.AssertExpectations(t)
	invQueryRepo.AssertExpectations(t)
	invCommandRepo.AssertExpectations(t)
}

func TestService_Checkout_ProductNotFound(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	posRepo := &mockPosRepo{}
	invQueryRepo := &mockInventoryQueryRepo{}
	invCommandRepo := &mockInventoryCommandRepo{}
	txManager := &mockTxManager{}

	svc := NewService(posRepo, posRepo, invQueryRepo, invCommandRepo, txManager)

	req := models.CheckoutRequest{
		PaymentMethod: models.PaymentMethodCash,
		Items: []models.CheckoutItemRequest{
			{
				ProductID: "non-existent",
				Quantity:  2,
			},
		},
	}

	invQueryRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)

	res, err := svc.Checkout(context.Background(), req)
	must.Error(err)
	is.ErrorIs(err, ErrInvalidInput)
	is.Nil(res)

	posRepo.AssertExpectations(t)
	invQueryRepo.AssertExpectations(t)
	invCommandRepo.AssertExpectations(t)
}

func BenchmarkCheckout(b *testing.B) {
	posRepo := &mockPosRepo{}
	invQueryRepo := &mockInventoryQueryRepo{}
	invCommandRepo := &mockInventoryCommandRepo{}
	txManager := &mockTxManager{}

	svc := NewService(posRepo, posRepo, invQueryRepo, invCommandRepo, txManager)

	product := &models.Product{
		ID:    "prod-1",
		Name:  "Tempered Glass",
		Price: 50000,
		Stock: 10,
	}

	req := models.CheckoutRequest{
		PaymentMethod: models.PaymentMethodCash,
		Items: []models.CheckoutItemRequest{
			{
				ProductID: "prod-1",
				Quantity:  2,
			},
		},
	}

	invQueryRepo.On("FindByID", mock.Anything, "prod-1").Return(product, nil)
	invCommandRepo.On("UpdateStock", mock.Anything, "prod-1", -2).Return(nil)
	posRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.PosTransaction")).Return(nil)

	expectedTx := &models.PosTransaction{
		ID:            "tx-123",
		PaymentMethod: models.PaymentMethodCash,
		TotalAmount:   100000,
	}
	posRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(expectedTx, nil)

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = svc.Checkout(ctx, req)
	}
}
