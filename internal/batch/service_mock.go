package batch

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) GetAll(ctx context.Context) ([]domain.ProductBatch, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.ProductBatch), args.Error(1)
}

func (s *ServiceMock) Save(ctx context.Context, b domain.ProductBatch) (int, error) {
	args := s.Called(ctx, b)
	return args.Int(0), args.Error(1)
}
