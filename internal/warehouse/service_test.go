package warehouse

import (
	"context"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestService_Read(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: return all warehouses
	// DESCRIPTION: Should return all warehouses
	t.Run("it should return all warehouses", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedWarehouses := []domain.Warehouse{
			{
				ID:                 1,
				Address:            "Test Address",
				Telephone:          "Test Telephone",
				WarehouseCode:      "Test WarehouseCode",
				MinimumCapacity:    100,
				MinimumTemperature: 10,
			},
			{
				ID:                 2,
				Address:            "Test Address 2",
				Telephone:          "Test Telephone 2",
				WarehouseCode:      "Test WarehouseCode 2",
				MinimumCapacity:    200,
				MinimumTemperature: 20,
			},
		}

		repository := &RepositoryMock{}
		repository.On("GetAll", ctx).Return(expectedWarehouses, nil)

		service := NewService(repository)

		// Act.
		obtainedWarehouses, err := service.GetAll(ctx)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouses, obtainedWarehouses)
		repository.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: return a warehouse by id
	// DESCRIPTION: should return a warehouse by id
	t.Run("it should return a warehouse by id", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}

		repository := &RepositoryMock{}
		repository.On("Get", ctx, 1).Return(expectedWarehouse, nil)

		service := NewService(repository)

		// Act.
		obtainedWarehouse, err := service.Get(ctx, 1)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, obtainedWarehouse)
		repository.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: return an error if the warehouse does not exist
	// DESCRIPTION: should return an error if the warehouse does not exist
	t.Run("it should return an error if the warehouse does not exist", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedError := ErrNotFound
		expectedWarehouse := domain.Warehouse{}

		repository := &RepositoryMock{}
		repository.On("Get", ctx, 1).Return(domain.Warehouse{}, ErrNotFound)

		service := NewService(repository)

		// Act.
		obtainedWarehouse, err := service.Get(ctx, 1)

		// Assert.
		assert.EqualError(t, err, expectedError.Error())
		assert.Equal(t, expectedWarehouse, obtainedWarehouse)
		repository.AssertExpectations(t)
	})
}

func TestService_Create(t *testing.T) {
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create a warehouse
	// DESCRIPTION: should create a warehouse
	t.Run("it should create a warehouse", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}

		repository := &RepositoryMock{}
		repository.On("Save", ctx, expectedWarehouse).Return(1, nil)

		service := NewService(repository)

		// Act.
		obtainedWarehouseID, err := service.Save(ctx, expectedWarehouse)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse.ID, obtainedWarehouseID)
		repository.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: return an error if the warehouse already exists
	// DESCRIPTION: should return an error if the warehouse already exists
	t.Run("it should return an error if the warehouse already exists", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedError := ErrDuplicateWarehouse
		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}
		expectedId := 0

		repository := &RepositoryMock{}
		repository.On("Save", ctx, expectedWarehouse).Return(0, ErrDuplicateWarehouse)

		service := NewService(repository)

		// Act.
		obtainedWarehouseID, err := service.Save(ctx, expectedWarehouse)

		// Assert.
		assert.EqualError(t, err, expectedError.Error())
		assert.Equal(t, expectedId, obtainedWarehouseID)
		repository.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete a warehouse
	// DESCRIPTION: should delete a warehouse
	t.Run("it should delete a warehouse", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		warehouseID := 1

		repository := &RepositoryMock{}
		repository.On("Delete", ctx, warehouseID).Return(nil)

		service := NewService(repository)

		// Act.
		err := service.Delete(ctx, warehouseID)

		// Assert.
		assert.NoError(t, err)
		repository.AssertExpectations(t)

	})
	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: return an error if the warehouse does not exist
	// DESCRIPTION: should return an error if the warehouse does not exist
	t.Run("it should return an error if the warehouse does not exist", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedError := ErrNotFound

		warehouseID := 1

		repository := &RepositoryMock{}
		repository.On("Delete", ctx, warehouseID).Return(ErrNotFound)

		service := NewService(repository)

		// Act.
		err := service.Delete(ctx, warehouseID)

		// Assert.
		assert.EqualError(t, err, expectedError.Error())
		repository.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: update a warehouse
	// DESCRIPTION: should update a warehouse
	t.Run("it should update a warehouse", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		warehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}

		repository := &RepositoryMock{}
		repository.On("Update", ctx, warehouse).Return(nil)

		service := NewService(repository)

		// Act.
		err := service.Update(ctx, warehouse)

		// Assert.
		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: return an error if the warehouse does not exist
	// DESCRIPTION: should return an error if the warehouse does not exist
	t.Run("it should return an error if the warehouse does not exists", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		warehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}
		expectedError := ErrNotFound

		repository := &RepositoryMock{}
		repository.On("Update", ctx, warehouse).Return(ErrNotFound)

		service := NewService(repository)

		// Act.
		err := service.Update(ctx, warehouse)

		// Assert.
		assert.EqualError(t, err, expectedError.Error())
		repository.AssertExpectations(t)
	})
}
