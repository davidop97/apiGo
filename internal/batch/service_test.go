package batch

import (
	"context"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestService_GetAll(t *testing.T) {
	t.Run("it should return a slice of all product batches", func(t *testing.T) {
		// Arrange
		// - Initialize a context for the test. This context is passed through the service to the repository.
		ctx := context.Background()
		// - Define the expected result: a slice of ProductBatch instances to simulate the expected behavior
		//   of the repository when the service requests all product batches.
		expectedBatches := []domain.ProductBatch{
			{
				ID:                 1,
				BatchNumber:        1,
				CurrentQuantity:    1,
				CurrentTemperature: 1,
				DueDate:            "2023-11-10",
				InitialQuantity:    1,
				ManufacturingDate:  "2023-11-10",
				ManufacturingHour:  1,
				MinimumTemperature: 1,
				ProductID:          1,
				SectionID:          1,
			},
			{
				ID:                 2,
				BatchNumber:        2,
				CurrentQuantity:    2,
				CurrentTemperature: 2,
				DueDate:            "2023-11-10",
				InitialQuantity:    2,
				ManufacturingDate:  "2023-11-10",
				ManufacturingHour:  2,
				MinimumTemperature: 2,
				ProductID:          2,
				SectionID:          2,
			},
		}
		// - Mock the repository to return the expected slice of product batches when GetAll is called.
		//   This simulates the repository's behavior without needing to interact with the actual data source.
		repository := &RepositoryMock{}
		repository.On("GetAll", ctx).Return(expectedBatches, nil)
		// - Instantiate the service with the mocked repository. This setup allows the test to focus
		//   on the service's ability to process and return the data correctly.
		service := NewService(repository)

		// Act
		// - Call the GetAll method on the service, capturing the batches returned and any error.
		obtainedBatches, obtainedError := service.GetAll(ctx)

		// Assert
		// - Verify that no error was returned. This ensures that the service's retrieval process
		//   operates as expected under normal conditions.
		assert.NoError(t, obtainedError)
		// - Check that the slice of product batches returned by the service matches the expected slice.
		//   This confirms that the service correctly processes and relays the repository's data.
		assert.Equal(t, expectedBatches, obtainedBatches)
		// - Confirm that the repository's expectations (i.e., a call to GetAll with the specified context)
		//   were met. This ensures the service correctly utilizes the repository in retrieving the data.
		repository.AssertExpectations(t)
	})

}

func TestService_Save(t *testing.T) {
	t.Run("it should save a product batch and return its id", func(t *testing.T) {
		// Arrange
		// - Create a context to be used in the test. This simulates the passing of context through the service layer to the repository.
		ctx := context.Background()
		// - Define the ID expected to be returned after saving the product batch. This ID simulates the database generated ID for the new batch.
		expectedID := 1
		// - Construct a ProductBatch instance to be saved. This object contains all necessary details of the product batch, including a predefined ID.
		batch := domain.ProductBatch{
			ID:                 1,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2023-11-10",
			InitialQuantity:    1,
			ManufacturingDate:  "2001-11-10",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductID:          1,
			SectionID:          1,
		}
		// - Mock the repository to first ensure that a batch with the expected ID does not already exist, and then to simulate the saving of the new batch.
		repository := &RepositoryMock{}
		repository.On("Exists", ctx, expectedID).Return(false)    // Simulate the batch does not already exist.
		repository.On("Save", ctx, batch).Return(expectedID, nil) // Simulate successful saving of the batch.
		// - Instantiate the service with the mocked repository, allowing the service's save functionality to be tested independently of database operations.
		service := NewService(repository)

		// Act
		// - Call the Save method on the service with the new batch, capturing the returned ID and any error.
		obtainedID, obtainedError := service.Save(ctx, batch)

		// Assert
		// - Verify that no error was returned during the save operation. This check ensures that the service can save a batch without encountering issues.
		assert.NoError(t, obtainedError)
		// - Ensure the ID returned from the save operation matches the expected ID. This confirms that the service correctly returns the identifier of the saved batch.
		assert.Equal(t, expectedID, obtainedID)
		// - Confirm that the mock repository's expectations (i.e., calls to Exists and Save with the specified context and batch) were fulfilled.
		//   This ensures that the service interacts with the repository as expected.
		repository.AssertExpectations(t)
	})

	t.Run("it should return an error if the product batch already exists", func(t *testing.T) {
		// Arrange
		// - Initialize a context for the test. This context simulates the passing of request-scoped information through the service and repository layers.
		ctx := context.Background()
		// - Define the ID of the product batch that is being tested for duplication.
		id := 1
		// - Create a ProductBatch instance that will be attempted to save. This batch is set up to simulate a duplicate entry scenario.
		batch := domain.ProductBatch{
			ID:                 id,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2023-11-10",
			InitialQuantity:    1,
			ManufacturingDate:  "2001-11-10",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductID:          1,
			SectionID:          1,
		}
		// - Specify the expected error when attempting to save a duplicate product batch. This error simulates the database or repository layer rejecting the duplicate entry.
		expectedError := ErrDuplicateBatchNumber
		// - Mock the repository to indicate that a batch with the given ID already exists when the Exists method is called.
		repository := &RepositoryMock{}
		repository.On("Exists", ctx, id).Return(true) // Simulate the batch already exists.
		// - Instantiate the service with the mocked repository to test the service's behavior in handling duplicate entries.
		service := NewService(repository)

		// Act
		// - Attempt to save the duplicate batch through the service, capturing the returned ID and any error.
		obtainedID, obtainedError := service.Save(ctx, batch)

		// Assert
		// - Verify that the correct error is returned when attempting to save a duplicate product batch.
		//   This ensures the service correctly identifies and handles duplicate entries according to business rules.
		assert.ErrorIs(t, obtainedError, expectedError)
		// - Check that the ID returned is 0, indicating that no new record was created due to the duplicate entry.
		assert.Equal(t, 0, obtainedID)
		// - Confirm that the repository's expectations, specifically the call to Exists with the specified context and ID, were met.
		//   This check ensures that the service properly consults the repository to determine the existence of the batch before attempting to save.
		repository.AssertExpectations(t)
	})

}
