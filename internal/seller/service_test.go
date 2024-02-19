package seller

import (
	"context"
	"errors"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestService_Read function.
// Test all the Get functions from seller service.
// - GetAllSellers
// - GetSellerByID
// - GetLocalityIdFromSeller
func TestService_Read(t *testing.T) {
	//Associated User Story: READ.
	//Edge case: find_all.
	//Summary: find_all case from GetAllSellers. Get all the sellers available in the database.
	t.Run("It should return all the sellers", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//Expected results
		expectedSellers := []domain.Seller{
			{
				ID:          1,
				CID:         1,
				CompanyName: "Seller 1",
				Address:     "Address 1",
				Telephone:   "Telephone 1",
				IDLocality:  1,
			},
			{
				ID:          2,
				CID:         2,
				CompanyName: "Seller 2",
				Address:     "Address 2",
				Telephone:   "Telephone 2",
				IDLocality:  2,
			},
		}

		//Create a new mock repository and mock the Get function.
		repository := NewMockRepository()
		repository.On("GetAll", ctx).Return(expectedSellers, nil)

		service := NewService(repository)

		//Act
		obtainedSellers, err := service.GetAllSellers(ctx)

		//Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the if obtained sellers are equal to the expected sellers.
		assert.Equal(t, expectedSellers, obtainedSellers)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_by_id_existent.
	//Summary: find_by_id_existent case from Get. Get a seller by its ID.
	t.Run("It should return the seller requested in the ID param ", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//Expected seller
		expectedSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Seller 1",
			Address:     "Address 1",
			Telephone:   "Telephone 1",
			IDLocality:  1,
		}

		//Create a new mock repository and mock the Get function.
		repository := NewMockRepository()
		repository.On("Get", ctx, 1).Return(expectedSeller, nil)

		service := NewService(repository)

		//Act
		obtainedSeller, err := service.GetSellerByID(ctx, 1)

		//Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the obtained seller is the expected seller.
		assert.Equal(t, expectedSeller, obtainedSeller)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_by_id_non_existent.
	//Summary: find_by_id_non_existent case from Get. If the ID of the seller not exist,
	//this method return an empty seller.
	t.Run("It should return an empty seller if the ID not exist ", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//Expected seller (empty seller).
		expectedSeller := domain.Seller{}

		//Create a new mock repository and mock the Get function.
		repository := NewMockRepository()
		nonexistentID := 1500
		repository.On("Get", ctx, nonexistentID).Return(expectedSeller, nil)

		service := NewService(repository)

		//Act
		obtainedSeller, err := service.GetSellerByID(ctx, nonexistentID)

		//Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the obtained seller is the expected (empty seller).
		assert.Equal(t, expectedSeller, obtainedSeller)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_locality_existent_from_seller (extra case for Sprint3).
	//Summary: find_locality_existent_from_seller case from GetLocalityIdFromSeller. Check if locality_id from seller exists
	// when trying to create a new seller.
	t.Run("It should return true if the locality exists ", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//Expected result
		expectedResult := true
		requestLocalityID := 1

		//Create a new mock repository and mock the function.
		repository := NewMockRepository()
		repository.On("GetLocalityIdFromSeller", ctx, requestLocalityID).Return(expectedResult)

		service := NewService(repository)

		//Act
		obtainedResult := service.GetLocalityIdFromSeller(ctx, requestLocalityID)

		//Assert
		//Check the obtained result is the expected.
		assert.Equal(t, expectedResult, obtainedResult)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_locality_non_existent_from_seller (extra case for Sprint3).
	//Summary: find_locality_non_existent_from_seller case from GetLocalityIdFromSeller. Return
	//false due the locality does not exist.
	t.Run("It should return false because the locality does not exist", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//Expected result
		expectedResult := false
		requestLocalityID := 1500

		//Create a new mock repository and mock the function.
		repository := NewMockRepository()
		repository.On("GetLocalityIdFromSeller", ctx, requestLocalityID).Return(expectedResult)

		service := NewService(repository)

		//Act
		obtainedResult := service.GetLocalityIdFromSeller(ctx, requestLocalityID)

		//Assert
		//Check the obtained result is the expected.
		assert.Equal(t, expectedResult, obtainedResult)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
	})
}

// TestService_Create function.
// Test the Save function from seller service.
func TestService_Create(t *testing.T) {
	//Associated User Story: CREATE.
	//Edge case: create_ok.
	//Summary: Creates a new seller if all the fields are correct.
	t.Run("It should save the seller in the database and return the ID of the new seller and nil", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//Expected seller
		expectedSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Seller 1",
			Address:     "Address 1",
			Telephone:   "Telephone 1",
			IDLocality:  1,
		}

		//Create a new mock repository and mock the Save function.
		repository := NewMockRepository()
		//Execute the Exists methods first to check if the seller already exists (using CID).
		//In this case, the seller does not exist (return false).
		repository.On("Exists", ctx, expectedSeller.CID).Return(false)
		repository.On("Save", ctx, expectedSeller).Return(expectedSeller.ID, nil)

		service := NewService(repository)

		//Act
		obtainedID, err := service.Save(ctx, expectedSeller)

		//Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the obtained ID is the expected ID.
		assert.Equal(t, expectedSeller.ID, obtainedID)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})
	//Associated User Story: CREATE.
	//Edge case: create_conflict.
	//Summary: Can not create a new seller because the seller already exists.
	//(using CID to check the existence of the seller).
	t.Run("It should return 0 and error: 'seller already exists'", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//This seller is trying to be created but it already exists.
		expectedSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Seller 1",
			Address:     "Address 1",
			Telephone:   "Telephone 1",
			IDLocality:  1,
		}

		//Expected result 0 and error "seller already exists"
		expectedResult := 0
		expectedError := errors.New("seller already exists")

		//Create a new mock repository and mock the Save function.
		repository := NewMockRepository()
		//Execute the Exists methods first to check if the seller already exists (using CID).
		//In this case, the seller already exists then return true.
		repository.On("Exists", ctx, expectedSeller.CID).Return(true)

		service := NewService(repository)

		//Act
		obtainedID, err := service.Save(ctx, expectedSeller)

		//Assert
		//Check if an error occurred.
		assert.Error(t, err)
		//Check the obtained error is the expected error.
		assert.Equal(t, expectedError, err)
		//Check the obtained ID is the expected ID.
		assert.Equal(t, expectedResult, obtainedID)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})
	//Associated User Story: CREATE.
	//Edge case: create_error.
	//Summary: Can not create a new seller because an internal error occurred.
	t.Run("It should return 0 and error: 'an error occurred'", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		//This seller is trying to be created but can not
		//be save due an error occurred.
		expectedSeller := domain.Seller{
			CID:         1,
			CompanyName: "Seller 1",
			Address:     "Address 1",
			Telephone:   "Telephone 1",
			IDLocality:  1,
		}

		//Expected result 0 and error "an error occurred".
		expectedResult := 0
		expectedError := errors.New("an error occurred")

		//Create a new mock repository and mock the Save function.
		repository := NewMockRepository()
		//Execute the Exists methods first, to check if the seller already exists (using CID).
		//In this case, the seller does not exist then return false.
		repository.On("Exists", ctx, expectedSeller.CID).Return(false)
		repository.On("Save", ctx, expectedSeller).Return(expectedResult, expectedError)

		service := NewService(repository)

		//Act
		obtainedID, err := service.Save(ctx, expectedSeller)

		//Assert
		//Check if an error occurred.
		assert.Error(t, err)
		//Check if the obtained error is the expected error.
		assert.Equal(t, expectedError, err)
		//Check if the obtained ID is the expected ID.
		assert.Equal(t, expectedResult, obtainedID)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

}

// TestService_Delete function.
// Test the Delete function from seller service.
func TestService_Delete(t *testing.T) {
	//Associated User Story: DELETE.
	//Edge case: delete_ok.
	//Summary: Delete a seller using its ID.
	t.Run("It should delete the seller and return nil", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		//ID from the seller to be deleted.
		IDSeller := 1

		//Create a new mock repository and mock the functions.
		repository := NewMockRepository()
		//Execute the Get methods first, to check if the seller already exists (using ID).
		//In this case, the seller exists then return err=nil.
		repository.On("Get", ctx, IDSeller).Return(domain.Seller{}, nil)
		//Then try to delete the seller. If everything is OK, return err=nil.
		repository.On("Delete", ctx, IDSeller).Return(nil)

		service := NewService(repository)

		//Act
		obtainedErr := service.Delete(ctx, IDSeller)

		//Assert
		//Check no error occurs.
		assert.NoError(t, obtainedErr)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

	//Associated User Story: DELETE.
	//Edge case: delete_non_existent.
	//Summary: Can not delete a user because the seller does not exist
	//(using ID for checking existence).
	t.Run("It should return an 'seller not found' error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		//ID from the seller to be deleted.
		IDSeller := 1500

		expectedError := errors.New("seller not found")

		//Create a new mock repository and mock the Get function.
		repository := NewMockRepository()
		//Execute the Get methods first, to check if the seller already exists (using ID).
		//In this case, the seller does not exist then return err=seller not found
		//And an empty seller.
		repository.On("Get", ctx, IDSeller).Return(domain.Seller{}, expectedError)

		service := NewService(repository)

		//Act
		obtainedErr := service.Delete(ctx, IDSeller)

		//Assert
		//Check if an error occurred.
		assert.Error(t, obtainedErr)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

	//Associated User Story: DELETE.
	//Edge case: delete_error (extra case for Sprint3).
	//Summary: Can not delete a user because an error occurred.
	t.Run("It should return an error: 'an error occurred'", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		//ID from the seller to be deleted.
		IDSeller := 1500

		expectedError := errors.New("an error occurred")

		//Create a new mock repository and mock the functions.
		repository := NewMockRepository()
		//Execute the Get methods first, to check if the seller already exists (using ID).
		//In this case, the seller exists therefore return err=nil
		repository.On("Get", ctx, IDSeller).Return(domain.Seller{}, nil)
		//Then try to delete the seller. If an error occurs, return err=an error occurred.
		repository.On("Delete", ctx, IDSeller).Return(expectedError)

		service := NewService(repository)

		//Act
		obtainedErr := service.Delete(ctx, IDSeller)

		//Assert
		//Check if an error occurred.
		assert.Error(t, obtainedErr)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)

	})

	//Associated User Story: DELETE.
	//Edge case: delete_error2 (extra case for Sprint3).
	//Summary: Can not delete a user because an error occurred
	//(ErrNotFound in Delete method). This means that there were no rows affected.
	//RowsAffected(function from Repository) returns a value < 1.
	t.Run("It should return an error: 'seller not found'", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		//ID from the seller to be deleted.
		IDSeller := 1500

		ErrNotFound = errors.New("seller not found")

		//Create a new mock repository and mock the functions.
		repository := NewMockRepository()
		//Execute the Get methods first, to check if the seller already exists (using ID).
		repository.On("Get", ctx, IDSeller).Return(domain.Seller{}, nil)
		//Then try to delete the seller. If an error occurs, return it.
		repository.On("Delete", ctx, IDSeller).Return(ErrNotFound)

		service := NewService(repository)

		//Act
		obtainedErr := service.Delete(ctx, IDSeller)

		//Assert
		//Check if an error occurred.
		assert.Error(t, obtainedErr)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
		//Check if the obtained error is the expected error.
		assert.Equal(t, ErrNotFound, obtainedErr)

	})

}

// TestService_Update function.
// Test the Update function from seller service.
func TestService_Update(t *testing.T) {
	//Associated User Story: UPDATE.
	//Edge case: update_ok.
	//Summary: Update a seller using its ID and return nil error.
	t.Run("It should update the seller and return nil error", func(t *testing.T) {

		//Arrange
		ctx := context.Background()
		//ID from the seller to be updated.
		IDSeller := 1

		//Seller with new data.
		updateSeller := domain.Seller{
			CID:         1,
			CompanyName: "Seller 1",
			Address:     "Address new",
			Telephone:   "Telephone new",
			IDLocality:  1,
		}

		//Create a new mock repository and mock the Update function.
		repository := NewMockRepository()

		//Try to update the seller. If everything is OK, return err=nil.
		repository.On("Update", ctx, updateSeller).Return(nil)

		service := NewService(repository)

		//Act
		obtainedErr := service.Update(ctx, updateSeller, IDSeller)

		//Assert
		//Check if an error occurred.
		assert.NoError(t, obtainedErr)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
		//Check if the obtained error is the expected error (nil in this case).
		assert.Equal(t, nil, obtainedErr)
	})

	//Associated User Story: UPDATE.
	//Edge case: update_non_existent.
	//Summary: Can not update a seller because does not exist.
	t.Run("It should return an error when trying to update a seller", func(t *testing.T) {

		//Arrange
		ctx := context.Background()
		//ID from the seller to be updated.
		IDSellerNotExist := 1500

		//Seller with new data.
		updateSeller := domain.Seller{
			CID:         1500,
			CompanyName: "Seller 1500",
			Address:     "Address new",
			Telephone:   "Telephone new",
			IDLocality:  2,
		}
		ExpectedError := errors.New("an error occurred")

		//Create a new mock repository and mock the Update function.
		repository := NewMockRepository()

		//Try to update the seller. In this case, the seller does not exist
		// then return err=an error occurred.
		repository.On("Update", ctx, updateSeller).Return(ExpectedError)

		service := NewService(repository)

		//Act
		obtainedErr := service.Update(ctx, updateSeller, IDSellerNotExist)

		//Assert
		//Check if an error occurred.
		assert.Error(t, obtainedErr)
		//Check the repository was called with the expected parameters.
		repository.AssertExpectations(t)
		//Check if the obtained error is the expected error.
		assert.Equal(t, ExpectedError, obtainedErr)
	})
}
