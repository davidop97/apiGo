package buyer

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Buyer), args.Error(1)

}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Buyer, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, cardNumberID string) bool {
	args := r.Called(ctx, cardNumberID)
	return args.Bool(0)
}

func (r *RepositoryMock) Save(ctx context.Context, b domain.Buyer) (int, error) {
	args := r.Called(ctx, b)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, b domain.Buyer) error {
	args := r.Called(ctx, b)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
