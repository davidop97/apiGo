package employee

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) GetAllEmployees(ctx context.Context) ([]domain.Employee, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.Employee), args.Error(1)
}

func (s *ServiceMock) GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (s *ServiceMock) SaveEmployee(ctx context.Context, employee domain.Employee) (int, error) {
	args := s.Called(ctx, employee)
	return args.Int(0), args.Error(1)
}

func (s *ServiceMock) UpdateEmployee(ctx context.Context, employee domain.Employee) error {
	args := s.Called(ctx, employee)
	return args.Error(0)
}

func (s *ServiceMock) DeleteEmployee(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
