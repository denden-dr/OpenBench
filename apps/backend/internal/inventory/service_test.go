package inventory

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockQueryRepo struct {
	mock.Mock
}

func (m *mockQueryRepo) FindByID(ctx context.Context, id string) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockQueryRepo) FindAll(ctx context.Context, search string, limit int, cursor string) ([]models.Product, string, error) {
	args := m.Called(ctx, search, limit, cursor)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Product), args.String(1), args.Error(2)
	}
	return nil, args.String(1), args.Error(2)
}

type mockCommandRepo struct {
	mock.Mock
}

func (m *mockCommandRepo) Create(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockCommandRepo) Update(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockCommandRepo) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	args := m.Called(ctx, id, quantityChange)
	return args.Error(0)
}

func (m *mockCommandRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_CreateProduct(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	qRepo := &mockQueryRepo{}
	cRepo := &mockCommandRepo{}
	svc := NewService(qRepo, cRepo)

	req := CreateProductRequest{
		Name:  "Charger Type C",
		Price: 75000,
		Stock: 50,
	}

	cRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil)

	res, err := svc.CreateProduct(context.Background(), req)
	must.NoError(err)
	must.NotNil(res)
	is.Equal("Charger Type C", res.Name)
	is.Equal(int64(75000), res.Price)
	is.Equal(50, res.Stock)

	cRepo.AssertExpectations(t)
}

func TestService_AdjustStock(t *testing.T) {
	must := require.New(t)

	qRepo := &mockQueryRepo{}
	cRepo := &mockCommandRepo{}
	svc := NewService(qRepo, cRepo)

	product := &models.Product{
		ID:    "prod-123",
		Name:  "Charger Type C",
		Price: 75000,
		Stock: 10,
	}

	qRepo.On("FindByID", mock.Anything, "prod-123").Return(product, nil)
	cRepo.On("UpdateStock", mock.Anything, "prod-123", -3).Return(nil)

	err := svc.AdjustStock(context.Background(), "prod-123", -3)
	must.NoError(err)

	qRepo.AssertExpectations(t)
	cRepo.AssertExpectations(t)
}

func TestService_AdjustStock_NegativeResult(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	qRepo := &mockQueryRepo{}
	cRepo := &mockCommandRepo{}
	svc := NewService(qRepo, cRepo)

	product := &models.Product{
		ID:    "prod-123",
		Name:  "Charger Type C",
		Price: 75000,
		Stock: 2,
	}

	qRepo.On("FindByID", mock.Anything, "prod-123").Return(product, nil)

	err := svc.AdjustStock(context.Background(), "prod-123", -5)
	must.Error(err)
	is.ErrorIs(err, ErrInvalidInput)

	qRepo.AssertExpectations(t)
	cRepo.AssertExpectations(t)
}
