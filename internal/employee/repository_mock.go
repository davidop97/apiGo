package employee

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Employee, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Employee), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Employee, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, cardNumberID string) bool {
	args := r.Called(ctx, cardNumberID)
	return args.Bool(0)
}

func (r *RepositoryMock) Save(ctx context.Context, e domain.Employee) (int, error) {
	args := r.Called(ctx, e)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, e domain.Employee) error {
	args := r.Called(ctx, e)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
