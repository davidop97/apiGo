package seller

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func NewMockRepository() *RepositoryMock {
	return &RepositoryMock{}
}

func (r *RepositoryMock) GetAll(ctx context.Context) (l []domain.Seller, err error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Seller), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (i domain.Seller, err error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, cid int) bool {
	args := r.Called(ctx, cid)
	return args.Bool(0)
}

func (r *RepositoryMock) Save(ctx context.Context, s domain.Seller) (id int, err error) {
	args := r.Called(ctx, s)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, s domain.Seller) (err error) {
	args := r.Called(ctx, s)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) (err error) {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *RepositoryMock) GetLocalityIdFromSeller(ctx context.Context, id int) bool {
	args := r.Called(ctx, id)
	return args.Bool(0)
}
