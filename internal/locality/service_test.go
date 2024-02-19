package locality

import (
	"context"
	"errors"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"

	"github.com/stretchr/testify/assert"
)

// TestService_GetLocality function.
// Test the GetLocality function in the following cases:
// - Get a locality by id if exists.
// - Return an error if the locality does not exist.
func TestService_GetLocality(t *testing.T) {
	// Edge case: get_locality_by_id_ok.
	// Summary: get a locality by id if exists.

	t.Run("It should return the locality requested in the ID param", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		//Expected locality to be returned by the service.
		expectedLocality := domain.Locality{
			ID:           1,
			PostalCode:   5000,
			LocalityName: "Cordoba",
			ProvinceName: "Cordoba",
			CountryName:  "Argentina",
		}

		//Create a new mock of the repository and mock the GetLocality function.
		mockRepository := NewMockRepository()
		mockRepository.On("GetLocality", ctx, 1).Return(expectedLocality, nil)

		service := NewService(mockRepository)

		// Act
		obtainedLocality, err := service.GetLocalityByID(ctx, 1)

		// Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the obtained locality is the expected locality.
		assert.Equal(t, expectedLocality, obtainedLocality)
		//Check the repository was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})

	// Edge case: get_locality_by_id_non_existent.
	// Summary: return an error if the locality does not exist.
	t.Run("It should return an error if the locality does not exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedError := errors.New("locality not found")
		expectedLocality := domain.Locality{}
		nonExistentLocalityID := 1500

		//Create a new mock of the repository and mock the GetLocality function.
		mockRepository := NewMockRepository()
		mockRepository.On("GetLocality", ctx, nonExistentLocalityID).Return(expectedLocality, expectedError)

		service := NewService(mockRepository)

		// Act
		obtainedLocality, err := service.GetLocalityByID(ctx, nonExistentLocalityID)

		// Assert
		//Check the error is the expected error.
		assert.Equal(t, expectedError, err)
		//Check the obtained locality is the expected (empty locality).
		assert.Equal(t, expectedLocality, obtainedLocality)
		//Check the repository mock was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})
}

// TestService_GetAll function.
// Test the GetAll function in the following cases:
// - Get all the localities if they exist.
// - Return an error if the localities were not found.
func TestService_GetAll(t *testing.T) {
	// Edge case: get_all_localities_ok.
	// Summary: get all the localities if they exist.
	t.Run("It should return all the localities if they exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		//Expected localities to be returned by the service.
		expectedLocalities := []domain.Locality{
			{
				ID:           1,
				PostalCode:   5000,
				LocalityName: "Cordoba",
				ProvinceName: "Cordoba",
				CountryName:  "Argentina",
			},
			{
				ID:           2,
				PostalCode:   5700,
				LocalityName: "San Luis",
				ProvinceName: "San Luis",
				CountryName:  "Argentina",
			},
		}

		//Create a new mock of the repository and mock the GetAll function.
		mockRepository := NewMockRepository()
		mockRepository.On("GetAll", ctx).Return(expectedLocalities, nil)

		service := NewService(mockRepository)

		// Act
		obtainedLocalities, err := service.GetAll(ctx)

		// Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the obtained localities are the expected localities.
		assert.Equal(t, expectedLocalities, obtainedLocalities)
		//Check the repository was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})

	// Edge case: get_all_localities_not_found.
	// Summary: return an error if the localities were not found.
	t.Run("It should return an error if the localities were not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedError := errors.New("localities not found")
		expectedLocalities := []domain.Locality{}

		//Create a new mock of the repository and mock the GetAll function.
		mockRepository := NewMockRepository()
		mockRepository.On("GetAll", ctx).Return(expectedLocalities, expectedError)

		service := NewService(mockRepository)

		// Act
		obtainedLocalities, err := service.GetAll(ctx)

		// Assert
		//Check the error is the expected error.
		assert.Equal(t, expectedError, err)
		//Check the obtained localities are the expected (empty localities).
		assert.Equal(t, expectedLocalities, obtainedLocalities)
		//Check the repository mock was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})
}

// TestService_Save function.
// Test the Save function in the following cases:
// - Save a locality if it does not exist.
// - Return an error if the locality already exists.
// - Return an error if the locality could not be saved.
func TestService_Save(t *testing.T) {
	// Edge case: save_locality_ok.
	// Summary: save a locality if it does not exist.
	t.Run("It should save the locality if it does not exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		//Locality to be saved by the service.
		localityToSave := domain.Locality{
			ID:           1,
			PostalCode:   5000,
			LocalityName: "Cordoba",
			ProvinceName: "Cordoba",
			CountryName:  "Argentina",
		}

		//Create a new mock of the repository and mock the Save function.
		mockRepository := NewMockRepository()
		//Execute the Exists methods first to check if the locality already exists (using PostalCode).
		//In this case, the locality does not exist (return false).
		mockRepository.On("Exists", ctx, localityToSave.PostalCode).Return(false)
		mockRepository.On("Save", ctx, localityToSave).Return(localityToSave.ID, nil)

		service := NewService(mockRepository)

		// Act
		obtainedLocalityID, err := service.Save(ctx, localityToSave)

		// Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the obtained ID is the expected ID.
		assert.Equal(t, localityToSave.ID, obtainedLocalityID)
		//Check the repository was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})

	// Edge case: save_locality_already_exists.
	// Summary: return an error if the locality already exists.
	t.Run("It should return an error if the locality already exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedError := errors.New("locality already exists")
		expectedLocalityID := 0

		localityToSave := domain.Locality{
			ID:           1,
			PostalCode:   5000,
			LocalityName: "Cordoba",
			ProvinceName: "Cordoba",
			CountryName:  "Argentina",
		}

		//Create a new mock of the repository and mock the Save function.
		mockRepository := NewMockRepository()
		mockRepository.On("Exists", ctx, localityToSave.PostalCode).Return(true)

		service := NewService(mockRepository)

		// Act
		obtainedLocalityID, err := service.Save(ctx, localityToSave)

		// Assert
		//Check the error is the expected error.
		assert.Equal(t, expectedError, err)
		//Check the obtained ID is the expected ID.
		assert.Equal(t, expectedLocalityID, obtainedLocalityID)
		//Check the repository mock was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})

	// Edge case: save_locality_error.
	// Summary: return an error if the locality could not be saved.
	t.Run("It should return an error if the locality could not be saved", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		//Expected error and id to be returned by the service.
		expectedError := errors.New("error saving locality")
		expectedLocalityID := 0

		//This locality is trying to be created but can not
		//be save due an error occurred.
		localityToSave := domain.Locality{
			ID:           2,
			PostalCode:   5700,
			LocalityName: "San Luis",
			ProvinceName: "San Luis",
			CountryName:  "Argentina",
		}

		//Create a new mock of the repository and mock the Save function.
		mockRepository := NewMockRepository()
		//Execute the Exists methods first to check if the locality already exists (using PostalCode).
		//In this case, the locality does not exist (return false).
		mockRepository.On("Exists", ctx, localityToSave.PostalCode).Return(false)
		//Execute the Save method to save the locality. It should return an error
		//because the locality could not be saved for some reason.
		mockRepository.On("Save", ctx, localityToSave).Return(expectedLocalityID, expectedError)

		service := NewService(mockRepository)

		// Act
		obtainedLocalityID, err := service.Save(ctx, localityToSave)

		// Assert
		//Check the error is the expected error.
		assert.Equal(t, expectedError, err)
		//Check the obtained ID is the expected ID.
		assert.Equal(t, expectedLocalityID, obtainedLocalityID)
		//Check the repository mock was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})
}

// TestService_GetReportSellers function.
// Test the GetReportSellers function in the following cases:
// - Get the report of sellers if the localities exist.
// - Return an error if the localities were not found.
func TestService_GetReportSellers(t *testing.T) {
	// Edge case: get_report_sellers_ok.
	// Summary: get the report of sellers if the localities exist. This case use one locality_id to get the report of that locality.
	t.Run("It should get the report of sellers if the locality exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		localityID := 1

		//Expected report of sellers to be returned by the service.
		expectedReportSellers := []domain.ReportSellers{
			{
				Locality_id:   1,
				Locality_name: "Cordoba",
				Postal_code:   5000,
				Sellers_count: 1,
			},
		}

		//Create a new mock of the repository and mock the GetReportSellers function.
		mockRepository := NewMockRepository()
		mockRepository.On("GetReportSellers", ctx, localityID).Return(expectedReportSellers, nil)

		service := NewService(mockRepository)

		// Act
		obtainedReportSellers, err := service.GetReportSellers(ctx, localityID)

		// Assert
		//Check no error occurs.
		assert.NoError(t, err)
		//Check the obtained report of sellers is the expected report of sellers.
		assert.Equal(t, expectedReportSellers, obtainedReportSellers)
		//Check the repository was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})

	// Edge case: get_report_sellers_error.
	// Summary: return an error if the localities were not found.
	t.Run("It should return an error if the localities were not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		localityID := 1500 //This locality does not exist.

		//Expected error to be returned by the service.
		expectedError := errors.New("localities were not found")
		expectedReportSellers := []domain.ReportSellers{} //Empty report of sellers.

		//Create a new mock of the repository and mock the GetReportSellers function.
		mockRepository := NewMockRepository()
		mockRepository.On("GetReportSellers", ctx, localityID).Return(expectedReportSellers, expectedError)

		service := NewService(mockRepository)

		// Act
		obtainedReportSellers, err := service.GetReportSellers(ctx, localityID)

		// Assert
		//Check the error is the expected error.
		assert.Equal(t, expectedError, err)
		//Check the obtained report of sellers is the expected report of sellers.
		assert.Equal(t, expectedReportSellers, obtainedReportSellers)
		//Check the repository was called with the expected parameters.
		mockRepository.AssertExpectations(t)
	})
}
