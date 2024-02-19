package batch

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.ProductBatch, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.ProductBatch), args.Error(1)
}

func (r *RepositoryMock) Save(ctx context.Context, b domain.ProductBatch) (int, error) {
	args := r.Called(ctx, b)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, batchNumber int) bool {
	args := r.Called(ctx, batchNumber)
	return args.Bool(0)
}
