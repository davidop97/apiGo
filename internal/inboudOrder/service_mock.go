package inboudorder

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) GetAllReports(ctx context.Context) ([]Report, error) {
	args := s.Called(ctx)
	return args.Get(0).([]Report), args.Error(1)
}

func (s *ServiceMock) GenerateReport(ctx context.Context, employeeID int) (Report, error) {
	args := s.Called(ctx, employeeID)
	return args.Get(0).(Report), args.Error(1)
}

func (s *ServiceMock) CreateInboundOrder(ctx context.Context, i domain.InboudOrder) (int, error) {
	args := s.Called(ctx, i)
	return args.Int(0), args.Error(1)
}
