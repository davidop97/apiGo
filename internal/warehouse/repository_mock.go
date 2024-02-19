package warehouse

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, warehouseCode string) bool {
	args := r.Called(ctx, warehouseCode)
	return args.Bool(0)
}

func (r *RepositoryMock) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	args := r.Called(ctx, w)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, w domain.Warehouse) error {
	args := r.Called(ctx, w)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
