package buyer

import (
	"context"
	"errors"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Tests for Buyer services

// TestService_ReadBuyer
// Test all cases for Get functions
// Tested methods: Get and GetAll
func TestService_ReadBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: find_all
	// DESCRIPTION: If the list has "n" elements it will return the number of total elements.
	t.Run("it should returns all buyers", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedbuyers := []domain.Buyer{{
			ID:           123456789,
			FirstName:    "John Doe",
			LastName:     "Doe",
			CardNumberID: "12345678",
		},
			{
				ID:           987654321,
				FirstName:    "SomeName",
				LastName:     "SomeLastName",
				CardNumberID: "987654321",
			},
		}
		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("GetAll", ctx).Return(expectedbuyers, nil)

		// call to service interface
		service := NewService(repository)

		// act
		obtainedBuyers, err := service.GetAll(ctx)

		// assert
		// Check no error occurs.
		assert.NoError(t, err)
		// Check if obtained buyers are equal to the expected buyers.
		assert.Equal(t, expectedbuyers, obtainedBuyers)
		// Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: READ
	// EDGE CASE: find_by_id_existent
	// DESCRIPTION: If the element searched by id exists it will return the information of the requested element.
	t.Run("it should return a buyer by id", func(t *testing.T) {
		// arrange
		id := 1
		ctx := context.Background()

		// expected result
		expectedBuyer := domain.Buyer{
			ID:           1,
			FirstName:    "Jane Doe",
			LastName:     "Doe",
			CardNumberID: "0987654321",
		}

		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Get", ctx, id).Return(expectedBuyer, nil)

		// call to service interface
		service := NewService(repository)

		// act
		obtainedBuyer, err := service.Get(ctx, id)

		// assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the if obtained buyers are equal to the expected buyers.
		assert.Equal(t, expectedBuyer, obtainedBuyer)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: READ
	// EDGE CASE: find_by_id_non_existent
	// DESCRIPTION: If the element searched by id does not exist, it returns error.
	t.Run("it should return an error when there is no such buyer", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedBuyerErr := ErrNotFound
		expectedBuyer := domain.Buyer{}
		id := 1

		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Get", ctx, id).Return(domain.Buyer{}, ErrNotFound)

		// call to service interface
		service := NewService(repository)

		// act
		obtainedBuyer, err := service.Get(ctx, id)

		// assert
		// Check for empty result because buyer not found
		assert.Empty(t, obtainedBuyer)
		// Check the if obtained buyers are equal to the expected buyers.
		assert.Equal(t, expectedBuyer, obtainedBuyer)
		// Check the if obtained error are equal to the expected error.
		assert.EqualError(t, err, expectedBuyerErr.Error())
		// Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

}

// TestService_DeleteBuyer
// Test all cases for Delete functions
// Tested methods: Delete
func TestService_DeleteBuyer(t *testing.T) {

	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete_ok
	// DESCRIPTION: If the deletion is successful the item will not appear in the list.
	t.Run("it should return nil when buyer succesfully deleted", func(t *testing.T) {
		// arrange
		id := 1
		ctx := context.Background()

		// expected result
		buyer := domain.Buyer{
			ID:           1,
			FirstName:    "Jane Doe",
			LastName:     "Doe",
			CardNumberID: "0987654321",
		}
		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Get", ctx, id).Return(buyer, nil)
		repository.On("Delete", ctx, id).Return(nil)

		// call to service interface
		service := NewService(repository)

		// act
		obtainedResult := service.Delete(ctx, id)
		// assert
		// Check for nil because buyer can be deleted
		assert.Nil(t, obtainedResult)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete_non_existent
	// DESCRIPTION: If the buyer does not exist, null will be returned.
	t.Run("it should return an error when buyer does not exists", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		// expected result
		id := 1
		expectedErr := ErrNotFound

		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Get", ctx, id).Return(domain.Buyer{}, ErrNotFound)
		// when the buyer does not exists, delete function it is not called.

		// call to service interface
		service := NewService(repository)

		// act
		obtainedResult := service.Delete(ctx, id)

		// assert
		// Check the if obtained error are equal to the expected error.
		assert.EqualError(t, obtainedResult, expectedErr.Error())
		//Check is the error expected are type that actual error
		assert.IsType(t, expectedErr, obtainedResult)
		// Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: error_calling_delete
	// DESCRIPTION: If there an error deleted by id, it returns error.
	t.Run("it should return errors when there is an error when calling delete function", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		id := 1
		buyer := domain.Buyer{
			ID:           1,
			FirstName:    "Jane Doe",
			LastName:     "Doe",
			CardNumberID: "0987654321",
		}
		expectedError := errors.New("")

		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Get", ctx, id).Return(buyer, nil)
		repository.On("Delete", ctx, id).Return(errors.New(""))

		// call to service interface
		service := NewService(repository)

		// act
		obtainedResult := service.Delete(ctx, id)

		// assert
		//Check is the error expected are type that actual error
		assert.EqualError(t, obtainedResult, expectedError.Error())
		//Check is the error expected are type that actual error
		assert.IsType(t, expectedError, obtainedResult)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

}

// TestService_CreateBuyer
// Test all cases for Save functions
// Tested methods: Save
func TestService_CreateBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_ok
	// DESCRIPTION: If it contains the required fields it will be created
	t.Run("it should return an id of the created buyer", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		// expected result
		cardNumberID := "1234567890"
		buyer := domain.Buyer{
			FirstName:    "John Doe",
			LastName:     "Doe",
			CardNumberID: cardNumberID,
		}
		expectedID := 1

		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Exists", ctx, cardNumberID).Return(false)
		repository.On("Save", ctx, buyer).Return(expectedID, nil)

		// call to service interface
		service := NewService(repository)

		// act
		buyerIDObtained, err := service.Save(ctx, buyer)

		// assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check if obtained buyers are equal to the expected buyers.
		assert.Equal(t, expectedID, buyerIDObtained)
		// Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_conflict
	// DESCRIPTION: If there card_number_id already exists it cannot be created.
	t.Run("it should return an error when buyer does already exists", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		cardNumberID := "1234567890"
		buyer := domain.Buyer{
			FirstName:    "Jane Doe",
			LastName:     "Doe",
			CardNumberID: cardNumberID,
		}

		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Exists", ctx, cardNumberID).Return(true)

		// call to service interface
		service := NewService(repository)

		// act
		buyerIDObtained, err := service.Save(ctx, buyer)
		// assert
		// Check if error exists
		assert.Error(t, err)
		//Check is the error expected are type that actual error
		assert.Equal(t, 0, buyerIDObtained)
		// check that save function it is not called because buyer exists
		repository.AssertNotCalled(t, "Save")
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_calling_create
	// DESCRIPTION: If there an error created a buyer, it returns error.
	t.Run("it should return errors when there is a error when calling to save function", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedError := errors.New("")
		cardNumberID := "1234567890"

		buyer := domain.Buyer{
			FirstName:    "John",
			LastName:     "Doe",
			CardNumberID: cardNumberID,
		}

		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Exists", ctx, cardNumberID).Return(false)
		repository.On("Save", ctx, buyer).Return(0, errors.New(""))

		// call to service interface
		service := NewService(repository)

		// act
		_, err := service.Save(ctx, buyer)

		// assert
		//Check is the error expected are type that actual error
		assert.EqualError(t, err, expectedError.Error())
		// Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

}

// TestService_UpdateBuyer
// Test all cases for Update functions
// Tested methods: Update
func TestService_UpdateBuyer(t *testing.T) {

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: update_existent
	// DESCRIPTION: When the data update is successful, the buyer will be returned with the updated information.
	t.Run("it should return nil if buyer it can updated", func(t *testing.T) {
		// arrange
		id := 1
		cardNumberID := "1234567890"

		buyerInDatabase := domain.Buyer{
			ID:           id,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: cardNumberID,
		}

		buyerToUpdate := domain.Buyer{
			ID:           id,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: cardNumberID,
		}

		ctx := context.Background()
		repository := &RepositoryMock{}
		repository.On("Exists", ctx, buyerToUpdate.CardNumberID).Return(false)
		repository.On("Update", ctx, buyerToUpdate).Return(nil)

		service := NewService(repository)

		// act
		result := service.Update(ctx, id, buyerToUpdate, &buyerInDatabase)

		// assert
		assert.Equal(t, nil, result)
		repository.AssertExpectations(t)

	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: update_non_existent
	// DESCRIPTION: If the buyer to be updated does not exist, null will be returned.
	t.Run("it should return an error when buyer does not exists to update", func(t *testing.T) {
		// arrange

		// expected result
		id := 1
		cardNumberID := "1234567890"
		buyerInDatabase := domain.Buyer{
			ID:           2,
			FirstName:    "John",
			LastName:     "Doe",
			CardNumberID: cardNumberID,
		}
		buyerToUpdate := domain.Buyer{
			ID:           id,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: cardNumberID,
		}
		expectedErr := errors.New("buyer already exists")

		// prepare environment to call mock
		ctx := context.Background()
		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Exists", ctx, buyerToUpdate.CardNumberID).Return(true)

		// call to service interface
		service := NewService(repository)
		// act
		result := service.Update(ctx, id, buyerToUpdate, &buyerInDatabase)
		// assert
		// Check the if obtained error are equal to the expected error.
		assert.EqualError(t, expectedErr, result.Error())
		// Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: error_calling_update
	// DESCRIPTION: If there an error created a buyer, it returns error.
	t.Run("it should return errors when there is a error when calling to update function", func(t *testing.T) {
		// arrange

		// expected result
		id := 1
		cardNumberID := "1234567890"
		expectedError := errors.New("some db error")

		buyerInDatabase := domain.Buyer{
			ID:           id,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: cardNumberID,
		}

		buyerToUpdate := domain.Buyer{
			ID:           id,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: cardNumberID,
		}

		// prepare environment to call mock
		ctx := context.Background()
		// create a mock for repository
		repository := &RepositoryMock{}
		// set mock with the expected value
		repository.On("Exists", ctx, buyerToUpdate.CardNumberID).Return(false)
		repository.On("Update", ctx, buyerToUpdate).Return(expectedError)

		// call to service interface
		service := NewService(repository)

		// act
		result := service.Update(ctx, id, buyerToUpdate, &buyerInDatabase)

		// assert
		// Check the if obtained error are equal to the expected error.
		assert.EqualError(t, expectedError, result.Error())
		// Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

}
