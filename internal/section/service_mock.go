package section

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) GetAll(ctx context.Context) ([]domain.Section, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.Section), args.Error(1)
}

func (r *ServiceMock) Get(ctx context.Context, id int) (domain.Section, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (r *ServiceMock) Save(ctx context.Context, s domain.Section) (int, error) {
	args := r.Called(ctx, s)
	return args.Int(0), args.Error(1)
}

func (r *ServiceMock) Update(ctx context.Context, s domain.Section) error {
	args := r.Called(ctx, s)
	return args.Error(0)
}

func (r *ServiceMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *ServiceMock) ProductCount(ctx context.Context, id int) ([]ProdCountResponse, error) {
	args := r.Called(ctx, id)
	return args.Get(0).([]ProdCountResponse), args.Error(1)
}
