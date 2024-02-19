package product

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

// Mock for test in handler

func (m *ServiceMock) Get(ctx context.Context, productID int) (domain.Product, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *ServiceMock) GetAll(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *ServiceMock) Save(ctx context.Context, product domain.Product) (int, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(int), args.Error(1)
}

func (m *ServiceMock) Update(ctx context.Context, product domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ServiceMock) Delete(ctx context.Context, productID int) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

func (m *ServiceMock) Exists(ctx context.Context, productCode string) bool {
	args := m.Called(ctx, productCode)
	return args.Bool(0)
}

func (m *ServiceMock) CreateProductRecord(ctx context.Context, p domain.ProductRecordCreate) (int, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(int), args.Error(1)
}

func (m *ServiceMock) GetProductRecord(ctx context.Context, idProduct int) ([]domain.ProductRecordGet, error) {
	args := m.Called(ctx, idProduct)
	return args.Get(0).([]domain.ProductRecordGet), args.Error(1)
}
