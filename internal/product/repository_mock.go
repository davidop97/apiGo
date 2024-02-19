package product

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Product, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, productCode string) bool {
	args := r.Called(ctx, productCode)
	return args.Bool(0)
}

func (r *RepositoryMock) Save(ctx context.Context, p domain.Product) (int, error) {
	args := r.Called(ctx, p)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, p domain.Product) error {
	args := r.Called(ctx, p)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *RepositoryMock) CreateProductRecord(ctx context.Context, p domain.ProductRecordCreate) (int, error) {
	args := r.Called(ctx, p)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) GetProductRecord(ctx context.Context, idProduct int) ([]domain.ProductRecordGet, error) {
	args := r.Called(ctx, idProduct)
	return args.Get(0).([]domain.ProductRecordGet), args.Error(1)
}
