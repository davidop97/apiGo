package warehouse

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (s *ServiceMock) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (s *ServiceMock) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	args := s.Called(ctx, w)
	return args.Int(0), args.Error(1)
}

func (s *ServiceMock) Delete(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

func (s *ServiceMock) Update(ctx context.Context, w domain.Warehouse) error {
	args := s.Called(ctx, w)
	return args.Error(0)
}
