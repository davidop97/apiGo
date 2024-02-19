package section

import (
	"context"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestService_Read(t *testing.T) {
	// First test case: it checks if the service correctly returns a list of all sections.
	t.Run("it should return a list of all sections", func(t *testing.T) {
		// Arrange
		ctx := context.Background()           // Creating a context for the test.
		expectedSections := []domain.Section{ // Setting up expected result.
			// First mock section.
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 1,
				MinimumTemperature: 1,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
			// Second mock section.
			{
				ID:                 2,
				SectionNumber:      2,
				CurrentTemperature: 2,
				MinimumTemperature: 2,
				CurrentCapacity:    2,
				MinimumCapacity:    2,
				MaximumCapacity:    2,
				WarehouseID:        2,
				ProductTypeID:      2,
			},
		}

		repository := &RepositoryMock{}                            // Creating a mock repository.
		repository.On("GetAll", ctx).Return(expectedSections, nil) // Setting up the mock response for GetAll method.
		service := NewService(repository)                          // Creating a service with the mock repository.

		// Act
		obtainedProducts, obtainedError := service.GetAll(ctx) // Calling the method to test.

		// Assert
		assert.NoError(t, obtainedError)                    // Verifying no error was returned.
		assert.Equal(t, expectedSections, obtainedProducts) // Verifying the result is as expected.
		repository.AssertExpectations(t)                    // Ensuring all expectations on the mock were met.
	})

	// Second test case: it checks if the service can return a specific section by its ID.
	t.Run("it should return the section corresponding to the given id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		id := 1                            // The ID to query for.
		expectedSection := domain.Section{ // The expected section to be returned.
			// Mock section data.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}

		repository := &RepositoryMock{}                            // Creating a mock repository.
		repository.On("Get", ctx, id).Return(expectedSection, nil) // Setting up the mock response for Get method.
		service := NewService(repository)                          // Creating the service.

		// Act
		obtainedSection, obtainedError := service.Get(ctx, id) // Calling the method with the test ID.

		// Assert
		assert.NoError(t, obtainedError)                  // Verifying no error was returned.
		assert.Equal(t, expectedSection, obtainedSection) // Verifying the section is as expected.
		repository.AssertExpectations(t)                  // Ensuring all expectations on the mock were met.
	})

	// Third test case: it checks the behavior of the service when a requested section does not exist.
	t.Run("it should return an error when the section does not exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		id := 1                             // The ID to query for.
		expectedSection := domain.Section{} // The expected result (empty section).
		expectedError := ErrNotFound        // The expected error to be returned.

		repository := &RepositoryMock{}                                    // Creating a mock repository.
		repository.On("Get", ctx, id).Return(expectedSection, ErrNotFound) // Setting up the mock response for Get method.
		service := NewService(repository)                                  // Creating the service.

		// Act
		obtainedSection, obtainedError := service.Get(ctx, id) // Calling the method with the test ID.

		// Assert
		assert.ErrorIs(t, obtainedError, expectedError)   // Verifying the correct error was returned.
		assert.Equal(t, expectedSection, obtainedSection) // Verifying the section is as expected (empty).
		repository.AssertExpectations(t)                  // Ensuring all expectations on the mock were met.
	})
}

func TestService_Create(t *testing.T) {
	// First test case: it checks if the service can correctly save a new section and return its ID.
	t.Run("it should save a section and return its id", func(t *testing.T) {
		// Arrange
		ctx := context.Background() // Creating a context for the test.
		section := domain.Section{  // Creating a mock section to be saved.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		exists := false    // The section does not exist yet.
		sectionNumber := 1 // The section number to be checked.
		expectedId := 1    // The expected ID to be returned after saving.

		repository := &RepositoryMock{}                             // Creating a mock repository.
		repository.On("Exists", ctx, sectionNumber).Return(exists)  // Setting up the mock response for Exists method.
		repository.On("Save", ctx, section).Return(expectedId, nil) // Setting up the mock response for Save method.
		service := NewService(repository)                           // Creating the service.

		// Act
		obtainedId, obtainedError := service.Save(ctx, section) // Calling the method to test.

		// Assert
		assert.NoError(t, obtainedError)        // Verifying no error was returned.
		assert.Equal(t, expectedId, obtainedId) // Verifying the ID is as expected.
		repository.AssertExpectations(t)        // Ensuring all expectations on the mock were met.
	})

	// Second test case: it checks if the service returns an error when trying to save a section with an existing section number.
	t.Run("it should return an error if section_number already exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		section := domain.Section{ // Creating a mock section with the same details as before.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		exists := true // The section number already exists.
		sectionNumber := 1
		expectedId := 0                         // No ID should be returned on failure.
		expectedError := ErrDuplicateSectNumber // The expected error for duplicate section number.

		repository := &RepositoryMock{}                            // Creating a mock repository.
		repository.On("Exists", ctx, sectionNumber).Return(exists) // Setting up the mock response for Exists method.
		service := NewService(repository)                          // Creating the service.

		// Act
		obtainedId, obtainedError := service.Save(ctx, section) // Calling the method with the intention of failing.

		// Assert
		assert.ErrorIs(t, obtainedError, expectedError) // Verifying the correct error was returned.
		assert.Equal(t, expectedId, obtainedId)         // Verifying no ID was returned.
		repository.AssertExpectations(t)                // Ensuring all expectations on the mock were met.
	})
}

func TestService_Delete(t *testing.T) {
	// First test case: it checks if the service can correctly delete a section by its ID.
	t.Run("it should delete a section corresponding to the given id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()                  // Creating a context for the test.
		id := 1                                      // The ID of the section to be deleted.
		repository := &RepositoryMock{}              // Creating a mock repository.
		repository.On("Delete", ctx, id).Return(nil) // Setting up the mock response for Delete method.
		service := NewService(repository)            // Creating the service with the mock repository.

		// Act
		err := service.Delete(ctx, id) // Calling the Delete method on the service.

		// Assert
		assert.NoError(t, err)           // Verifying no error was returned.
		repository.AssertExpectations(t) // Ensuring all expectations on the mock were met.
	})

	// Second test case: it checks if the service returns an error when trying to delete a non-existent section.
	t.Run("it should return an error if section id is doesn't exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		id := 1                                                // The ID of the section that doesn't exist.
		expectedError := ErrNotFound                           // The expected error to be returned.
		repository := &RepositoryMock{}                        // Creating a mock repository.
		repository.On("Delete", ctx, id).Return(expectedError) // Setting up the mock to return an error.
		service := NewService(repository)                      // Creating the service.

		// Act
		obtainedError := service.Delete(ctx, id) // Calling the Delete method, expecting an error.

		// Assert
		assert.ErrorIs(t, obtainedError, expectedError) // Verifying the correct error was returned.
		repository.AssertExpectations(t)                // Ensuring all expectations on the mock were met.
	})
}

func TestService_Update(t *testing.T) {
	// Testing the 'Update' function in the service.

	// First test case: Verifying if the service correctly updates a section given its ID.
	t.Run("it should update fields of the section corresponding to the given id", func(t *testing.T) {
		// Arrange
		// Setting up the context and initializing the test data.
		ctx := context.Background()
		id := 1
		updatedSection := domain.Section{ // This is the section data after the update.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 2,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		}
		originalSection := domain.Section{ // This represents the original data before the update.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		repository := &RepositoryMock{}                            // Mocking the repository.
		repository.On("Get", ctx, id).Return(originalSection, nil) // Mocking the 'Get' method to return the original section data.
		repository.On("Update", ctx, updatedSection).Return(nil)   // Mocking the 'Update' method to return a successful operation.
		service := NewService(repository)                          // Creating a new service with the mocked repository.

		// Act
		// Calling the 'Update' method and capturing any returned error.
		obtainedError := service.Update(ctx, updatedSection)

		// Assert
		// Checking that no error was returned and the mock expectations were met.
		assert.NoError(t, obtainedError)
		repository.AssertExpectations(t)
	})

	// Second test case: Verifying the behavior when an invalid section ID is provided.
	t.Run("it should return an error if section id doesn't exists", func(t *testing.T) {
		// Arrange
		// Similar setup as the first test case but with the anticipation of an error.
		ctx := context.Background()
		id := 1
		updatedSection := domain.Section{ // Section data being used for the update.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		emptySection := domain.Section{}                                  // Empty section to simulate a non-existent section.
		expectedError := ErrNotFound                                      // Expected error when the section does not exist.
		repository := &RepositoryMock{}                                   // Mocking the repository again.
		repository.On("Get", ctx, id).Return(emptySection, expectedError) // Mocking 'Get' to return an error for non-existent section.
		service := NewService(repository)

		// Act
		// Attempting to update and expecting an error.
		obtainedError := service.Update(ctx, updatedSection)

		// Assert
		// Checking that the expected error was returned and the mock expectations were met.
		assert.ErrorIs(t, obtainedError, expectedError)
		repository.AssertExpectations(t)
	})

	// Third test case: Verifying the behavior when trying to update to an existing section number.
	t.Run("it should return an error if section number already exists", func(t *testing.T) {
		// Arrange
		// Setup similar to the first test but with an additional check for existing section number.
		ctx := context.Background()
		id := 1                           // ID of the section to be updated.
		sectionNumber := 2                // The new section number which already exists.
		updatedSection := domain.Section{ // Section data with the new section number.
			ID:                 1,
			SectionNumber:      2,
			CurrentTemperature: 2,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		}
		originalSection := domain.Section{ // Original data of the section before update.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		expectedError := ErrDuplicateSectNumber                    // Expected error for duplicate section number.
		repository := &RepositoryMock{}                            // Mocking the repository.
		repository.On("Get", ctx, id).Return(originalSection, nil) // Mocking 'Get' to return original data.
		repository.On("Exists", ctx, sectionNumber).Return(true)   // Mocking 'Exists' to simulate that the section number already exists.
		service := NewService(repository)

		// Act
		// Attempting to update and expecting a duplicate error.
		obtainedError := service.Update(ctx, updatedSection)

		// Assert
		// Checking that the expected error was returned and the mock expectations were met.
		assert.ErrorIs(t, obtainedError, expectedError)
		repository.AssertExpectations(t)
	})
}

func TestService_ProductCount(t *testing.T) {
	t.Run("it should return a product count for the section corresponding to the given id", func(t *testing.T) {
		// Arrange
		// - Initialize a context for the test, simulating the passage of request-scoped data through the application layers.
		ctx := context.Background()
		// - Define the ID of the section for which the product count is being requested.
		id := 1
		// - Prepare the expected response, a slice of ProdCountResponse objects, which will simulate the repository's response for the given section ID.
		expectedProductCount := []ProdCountResponse{
			{
				ID:            1, // Section ID
				SectionNumber: 1, // Number of the section
				ProductCount:  1, // The expected count of products in the section
			},
		}
		// - Mock the repository to return the expected product count when the ProductCount method is called with the specified section ID.
		repository := &RepositoryMock{}
		repository.On("ProductCount", ctx, id).Return(expectedProductCount, nil) // Simulate successful retrieval of product count.
		// - Instantiate the service with the mocked repository to test the service's ability to process and relay product count data.
		service := NewService(repository)

		// Act
		// - Call the ProductCount method on the service with the given section ID, capturing the returned product count and any error.
		obtainedProductCount, obtainedError := service.ProductCount(ctx, id)

		// Assert
		// - Verify that no error was returned. This check ensures that the service can retrieve the product count without encountering issues.
		assert.NoError(t, obtainedError)
		// - Ensure the product count returned by the service matches the expected count. This confirms that the service accurately processes and relays the repository's data.
		assert.Equal(t, expectedProductCount, obtainedProductCount)
		// - Confirm that the mock repository's expectations, specifically the call to ProductCount with the specified context and section ID, were met.
		//   This verification ensures that the service properly utilizes the repository in determining the product count for the specified section.
		repository.AssertExpectations(t)
	})

}
