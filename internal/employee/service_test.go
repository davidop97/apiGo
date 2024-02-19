package employee

import (
	"context"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestService_Read(t *testing.T) {
	t.Run("find_all", func(t *testing.T) {
		//Given
		ctx := context.Background()
		expectedEmployees := []domain.Employee{
			{
				ID:           1,
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			},
			{
				ID:           2,
				CardNumberID: "C987E012F",
				FirstName:    "George",
				LastName:     "Smith",
				WarehouseID:  2,
			},
		}
		repository := &RepositoryMock{}
		repository.On("GetAll", ctx).Return(expectedEmployees, nil)
		service := NewService(repository)

		//When
		obtainedEmployees, obtainedError := service.GetAllEmployees(ctx)

		//Then
		assert.NoError(t, obtainedError)
		assert.Equal(t, expectedEmployees, obtainedEmployees)
		repository.AssertExpectations(t)
	})
	t.Run("find_by_id_existent", func(t *testing.T) {
		//Given
		ctx := context.Background()
		id := 1
		expectedEmployee := domain.Employee{
			ID:           1,
			CardNumberID: "D789E012F",
			FirstName:    "Harold",
			LastName:     "Doe",
			WarehouseID:  1,
		}
		repository := &RepositoryMock{}
		repository.On("Get", ctx, id).Return(expectedEmployee, nil)
		service := NewService(repository)

		//When
		obtainedEmployee, obtainedError := service.GetEmployeeByID(ctx, id)

		//Then
		assert.NoError(t, obtainedError)
		assert.Equal(t, expectedEmployee, obtainedEmployee)
		repository.AssertExpectations(t)
	})
	t.Run("find_by_id_non_existent", func(t *testing.T) {
		//Given
		ctx := context.Background()
		id := 1
		expectedEmployee := domain.Employee{}
		expectedError := ErrNotFound
		repository := &RepositoryMock{}
		repository.On("Get", ctx, id).Return(expectedEmployee, expectedError)
		service := NewService(repository)

		//When
		obtainedEmployee, obtainedError := service.GetEmployeeByID(ctx, id)

		//Then
		assert.ErrorIs(t, obtainedError, expectedError)
		assert.Equal(t, expectedEmployee, obtainedEmployee)
		repository.AssertExpectations(t)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		//Given
		ctx := context.Background()
		employee := domain.Employee{
			ID:           1,
			CardNumberID: "D789E012F",
			FirstName:    "Harold",
			LastName:     "Doe",
			WarehouseID:  1,
		}
		exists := false
		cardNumberID := "D789E012F"
		expectedId := 1
		repository := &RepositoryMock{}
		repository.On("Exists", ctx, cardNumberID).Return(exists)
		repository.On("Save", ctx, employee).Return(expectedId, nil)
		service := NewService(repository)

		//When
		obtainedId, obtainedError := service.SaveEmployee(ctx, employee)

		//Then
		assert.NoError(t, obtainedError)
		assert.Equal(t, expectedId, obtainedId)
		repository.AssertExpectations(t)
	})
	t.Run("create_conflict", func(t *testing.T) {
		//Given
		ctx := context.Background()
		employee := domain.Employee{
			ID:           1,
			CardNumberID: "D789E012F",
			FirstName:    "Harold",
			LastName:     "Doe",
			WarehouseID:  1,
		}
		exists := true
		cardNumberID := "D789E012F"
		expectedId := 0
		expectedError := ErrEmployeeAlreadyExists
		repository := &RepositoryMock{}
		repository.On("Exists", ctx, cardNumberID).Return(exists)
		service := NewService(repository)

		//When
		obtainedId, obtainedError := service.SaveEmployee(ctx, employee)

		//Then
		assert.ErrorIs(t, obtainedError, expectedError)
		assert.Equal(t, expectedId, obtainedId)
		repository.AssertExpectations(t)
	})
}
func TestService_Update(t *testing.T) {
	t.Run("update_existent", func(t *testing.T) {
		//Given
		ctx := context.Background()
		id := 1
		originalEmployee := domain.Employee{
			ID:           1,
			CardNumberID: "D789E012F",
			FirstName:    "Harold",
			LastName:     "Doe",
			WarehouseID:  1,
		}
		updatedEmployee := domain.Employee{
			ID:           1,
			CardNumberID: "D789E012F",
			FirstName:    "George",
			LastName:     "Smith",
			WarehouseID:  3,
		}
		repository := &RepositoryMock{}
		repository.On("Get", ctx, id).Return(originalEmployee, nil)
		repository.On("Update", ctx, updatedEmployee).Return(nil)
		service := NewService(repository)

		//When
		obtainedError := service.UpdateEmployee(ctx, updatedEmployee)

		//Then
		assert.NoError(t, obtainedError)
		repository.AssertExpectations(t)
	})
	t.Run("update_non_existent", func(t *testing.T) {
		//Given
		ctx := context.Background()
		id := 1
		emptyEmployee := domain.Employee{}
		updatedEmployee := domain.Employee{
			ID:           1,
			CardNumberID: "D789E012F",
			FirstName:    "George",
			LastName:     "Smith",
			WarehouseID:  3,
		}
		expectedError := ErrNotFound
		repository := &RepositoryMock{}
		repository.On("Get", ctx, id).Return(emptyEmployee, expectedError)
		service := NewService(repository)

		//When
		obtainedError := service.UpdateEmployee(ctx, updatedEmployee)

		//Then
		assert.ErrorIs(t, obtainedError, expectedError)
		repository.AssertExpectations(t)
	})
}
func TestService_Delete(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {
		//Given
		ctx := context.Background()
		id := 1
		repository := &RepositoryMock{}
		repository.On("Delete", ctx, id).Return(nil)
		service := NewService(repository)

		//When
		obtainedError := service.DeleteEmployee(ctx, id)

		//Then
		assert.NoError(t, obtainedError)
		repository.AssertExpectations(t)
	})
	t.Run("delete_non_existent", func(t *testing.T) {
		//Given
		ctx := context.Background()
		id := 1
		expectedError := ErrNotFound
		repository := &RepositoryMock{}
		repository.On("Delete", ctx, id).Return(expectedError)
		service := NewService(repository)

		//When
		obtainedError := service.DeleteEmployee(ctx, id)

		//Then
		assert.ErrorIs(t, obtainedError, expectedError)
		repository.AssertExpectations(t)
	})
}
