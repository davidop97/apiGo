package carries

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (r *ServiceMock) GetAll(ctx context.Context) ([]domain.Carries, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Carries), args.Error(1)
}

func (r *ServiceMock) GetAllCarriesByLocality(ctx context.Context) ([]domain.LocalityCarries, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.LocalityCarries), args.Error(1)
}

func (r *ServiceMock) GetAllCarriesByLocalityID(ctx context.Context, localityID int) (domain.LocalityCarries, error) {
	args := r.Called(ctx, localityID)
	return args.Get(0).(domain.LocalityCarries), args.Error(1)
}

func (r *ServiceMock) Save(ctx context.Context, c domain.Carries) (int, error) {
	args := r.Called(ctx, c)
	return args.Int(0), args.Error(1)
}
