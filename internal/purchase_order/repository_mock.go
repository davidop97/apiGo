package purchase_order

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Save(ctx context.Context, po domain.PurchaseOrder) (int, error) {
	args := m.Called(ctx, po)
	return args.Get(0).(int), args.Error(1)
}

func (m *RepositoryMock) ExistsPurchaseOrder(ctx context.Context, purchaseOrder int) bool {
	args := m.Called(ctx, purchaseOrder)
	return args.Get(0).(bool)
}

func (m *RepositoryMock) ExistsBuyer(ctx context.Context, id int) bool {
	args := m.Called(ctx, id)
	return args.Bool(0)
}

func (m *RepositoryMock) ExistsProductsRecord(ctx context.Context, id int) bool {
	args := m.Called(ctx, id)
	return args.Bool(0)
}

func (m *RepositoryMock) PurchaseOrdersByBuyers(ctx context.Context, buyerID int) ([]domain.PurchaseOrdersByBuyer, error) {
	args := m.Called(ctx, buyerID)
	return args.Get(0).([]domain.PurchaseOrdersByBuyer), args.Error(1)
}
