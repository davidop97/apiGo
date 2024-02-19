package inboudorder

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Exists(ctx context.Context, employeeID int) (*domain.Employee, error) {
	args := r.Called(ctx, employeeID)
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (r *RepositoryMock) GetAllReports(ctx context.Context) ([]Report, error) {
	args := r.Called(ctx)
	return args.Get(0).([]Report), args.Error(1)
}

func (r *RepositoryMock) GenerateReport(ctx context.Context, employeeID int) (Report, error) {
	args := r.Called(ctx, employeeID)
	return args.Get(0).(Report), args.Error(1)
}

func (r *RepositoryMock) ExistsEmployee(ctx context.Context, employeeID int) bool {
	args := r.Called(ctx, employeeID)
	return args.Bool(0)
}

func (r *RepositoryMock) ExistsInboundOrder(ctx context.Context, orderNumber string) bool {
	args := r.Called(ctx, orderNumber)
	return args.Bool(0)
}

func (r *RepositoryMock) ExistsWarehouse(ctx context.Context, warehouseID int) bool {
	args := r.Called(ctx, warehouseID)
	return args.Bool(0)
}

func (r *RepositoryMock) Save(ctx context.Context, i domain.InboudOrder) (int, error) {
	args := r.Called(ctx, i)
	return args.Int(0), args.Error(1)
}
