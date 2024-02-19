package section

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Section, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Section), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Section, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (r *RepositoryMock) Save(ctx context.Context, s domain.Section) (int, error) {
	args := r.Called(ctx, s)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, s domain.Section) error {
	args := r.Called(ctx, s)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *RepositoryMock) Exists(ctx context.Context, sectionNumber int) bool {
	args := r.Called(ctx, sectionNumber)
	return args.Bool(0)
}

func (r *RepositoryMock) ProductCount(ctx context.Context, id int) ([]ProdCountResponse, error) {
	args := r.Called(ctx, id)
	return args.Get(0).([]ProdCountResponse), args.Error(1)
}
