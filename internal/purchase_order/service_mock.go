package purchase_order

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Save(ctx context.Context, purchaseOrder domain.PurchaseOrder) (int, error) {
	args := s.Called(ctx, purchaseOrder)
	return args.Get(0).(int), args.Error(1)
}

func (s *ServiceMock) PurchaseOrdersByBuyer(ctx context.Context, buyerID int) ([]domain.PurchaseOrdersByBuyer, error) {
	args := s.Called(ctx, buyerID)
	return args.Get(0).([]domain.PurchaseOrdersByBuyer), args.Error(1)
}
