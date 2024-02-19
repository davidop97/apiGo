package inboudorder

import (
	"context"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestService_GetAllReports(t *testing.T) {
	t.Run("should return a list of reports", func(t *testing.T) {
		// Given
		ctx := context.Background()
		expectedInboundOrders := []Report{
			{
				Employee: &domain.Employee{
					ID:           1,
					CardNumberID: "402323",
					FirstName:    "Harold",
					LastName:     "Doe",
					WarehouseID:  1,
				},
				InboudOrdersCount: 1,
			},
			{
				Employee: &domain.Employee{
					ID:           2,
					CardNumberID: "402324",
					FirstName:    "Jane",
					LastName:     "Doe",
					WarehouseID:  2,
				},
				InboudOrdersCount: 2,
			},
		}

		repository := &RepositoryMock{}
		repository.On("GetAllReports", ctx).Return(expectedInboundOrders, nil)
		service := NewService(repository)

		// When
		obtainedInboundOrders, obtainedError := service.GetAllReports(ctx)

		// Then
		assert.NoError(t, obtainedError)
		assert.Equal(t, expectedInboundOrders, obtainedInboundOrders)
		repository.AssertExpectations(t)
	})
}

func TestService_GenerateReport(t *testing.T) {
	t.Run("should return a report", func(t *testing.T) {
		// Given
		ctx := context.Background()
		employeeID := 1
		expectedInboundOrders := Report{
			Employee: &domain.Employee{
				ID:           1,
				CardNumberID: "402323",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			},
			InboudOrdersCount: 1,
		}

		repository := &RepositoryMock{}
		repository.On("GenerateReport", ctx, employeeID).Return(expectedInboundOrders, nil)
		service := NewService(repository)

		// When
		obtainedInboundOrders, obtainedError := service.GenerateReport(ctx, employeeID)

		// Then
		assert.NoError(t, obtainedError)
		assert.Equal(t, expectedInboundOrders, obtainedInboundOrders)
		repository.AssertExpectations(t)
	})
}

func TestService_CreateInboundOrder(t *testing.T) {
	t.Run("should return a id with a new inbound order", func(t *testing.T) {
		// Given
		ctx := context.Background()
		inboundOrder := domain.InboudOrder{
			OrderDate:      "2021-01-01",
			OrderNumber:    "123456",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		expectedId := 1
		existsTrue := true
		existsFalse := false

		repository := &RepositoryMock{}
		repository.On("ExistsEmployee", ctx, inboundOrder.EmployeeID).Return(existsTrue)
		repository.On("ExistsWarehouse", ctx, inboundOrder.WarehouseID).Return(existsTrue)
		repository.On("ExistsInboundOrder", ctx, inboundOrder.OrderNumber).Return(existsFalse)
		repository.On("Save", ctx, inboundOrder).Return(expectedId, nil)
		service := NewService(repository)

		// When
		obtainedId, obtainedError := service.CreateInboundOrder(ctx, inboundOrder)

		// Then
		assert.NoError(t, obtainedError)
		assert.Equal(t, expectedId, obtainedId)
		repository.AssertExpectations(t)
	})
	t.Run("should return a error due to employee doesn`t exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		inboundOrder := domain.InboudOrder{
			OrderDate:      "2021-01-01",
			OrderNumber:    "123456",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		expectedId := 0
		existsFalse := false
		expectedError := ErrEmployeeDoesNotExists

		repository := &RepositoryMock{}
		repository.On("ExistsEmployee", ctx, inboundOrder.EmployeeID).Return(existsFalse)
		service := NewService(repository)

		// When
		obtainedId, obtainedError := service.CreateInboundOrder(ctx, inboundOrder)

		// Then
		assert.ErrorIs(t, obtainedError, expectedError)
		assert.Equal(t, expectedId, obtainedId)
		repository.AssertExpectations(t)
	})
	t.Run("should return a error due to warehouse doesn`t exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		inboundOrder := domain.InboudOrder{
			OrderDate:      "2021-01-01",
			OrderNumber:    "123456",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		expectedId := 0
		existsFalse := false
		existsTrue := true
		expectedError := ErrWarehouseDoesNotExists

		repository := &RepositoryMock{}
		repository.On("ExistsEmployee", ctx, inboundOrder.EmployeeID).Return(existsTrue)
		repository.On("ExistsWarehouse", ctx, inboundOrder.WarehouseID).Return(existsFalse)
		service := NewService(repository)

		// When
		obtainedId, obtainedError := service.CreateInboundOrder(ctx, inboundOrder)

		// Then
		assert.ErrorIs(t, obtainedError, expectedError)
		assert.Equal(t, expectedId, obtainedId)
		repository.AssertExpectations(t)
	})
	t.Run("should return a error due to inbound Order already exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		inboundOrder := domain.InboudOrder{
			OrderDate:      "2021-01-01",
			OrderNumber:    "123456",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		expectedId := 0
		existsTrue := true
		expectedError := ErrInboundOrderAlreadyExists

		repository := &RepositoryMock{}
		repository.On("ExistsEmployee", ctx, inboundOrder.EmployeeID).Return(existsTrue)
		repository.On("ExistsWarehouse", ctx, inboundOrder.WarehouseID).Return(existsTrue)
		repository.On("ExistsInboundOrder", ctx, inboundOrder.OrderNumber).Return(existsTrue)
		service := NewService(repository)

		// When
		obtainedId, obtainedError := service.CreateInboundOrder(ctx, inboundOrder)

		// Then
		assert.ErrorIs(t, obtainedError, expectedError)
		assert.Equal(t, expectedId, obtainedId)
		repository.AssertExpectations(t)
	})
}
