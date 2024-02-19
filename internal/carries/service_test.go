package carries

import (
	"context"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestService_Read(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: return all carries
	// DESCRIPTION: Should return all carries
	t.Run("it should return all carries", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedCarries := []domain.Carries{
			{
				ID:          1,
				CID:         "Test Name",
				CompanyName: "Test Company Name",
				Address:     "Test Address",
				Telephone:   "Test Telephone",
				LocalityID:  1,
			},
			{
				ID:          2,
				CID:         "Test Name 2",
				CompanyName: "Test Company Name 2",
				Address:     "Test Address 2",
				Telephone:   "Test Telephone 2",
				LocalityID:  2,
			},
		}

		repository := &RepositoryMock{}
		repository.On("GetAll", ctx).Return(expectedCarries, nil)

		service := NewService(repository)

		// Act.
		obtainedCarries, err := service.GetAll(ctx)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, expectedCarries, obtainedCarries)
		repository.AssertExpectations(t)
	})
}

func TestService_Create(t *testing.T) {
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create a carry
	// DESCRIPTION: Should create a carry
	t.Run("it should create a carry", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedCarry := domain.Carries{
			ID:          1,
			CID:         "Test Name",
			CompanyName: "Test Company Name",
			Address:     "Test Address",
			Telephone:   "Test Telephone",
			LocalityID:  1,
		}

		repository := &RepositoryMock{}
		repository.On("Save", ctx, expectedCarry).Return(expectedCarry.ID, nil)

		service := NewService(repository)

		// Act.
		obtainedCarryID, err := service.Save(ctx, expectedCarry)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, expectedCarry.ID, obtainedCarryID)
		repository.AssertExpectations(t)
	})

	// EDGE CASE: create a carry with incorrect data
	// DESCRIPTION: Should return error
	t.Run("it should return error when one of the fields is empty", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedCarry := domain.Carries{
			ID:          1,
			CID:         "Test Name",
			CompanyName: "",
			Address:     "",
			Telephone:   "Test Telephone",
			LocalityID:  1,
		}

		repository := &RepositoryMock{}

		service := NewService(repository)

		// Act.
		obtainedCarryID, err := service.Save(ctx, expectedCarry)

		// Assert.
		assert.EqualError(t, err, ErrIncorrectData.Error())
		assert.Equal(t, 0, obtainedCarryID)
		repository.AssertExpectations(t)
	})
	t.Run("it should return error when the locality ID is negative", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedCarry := domain.Carries{
			ID:          1,
			CID:         "Test Name",
			CompanyName: "Test Company Name",
			Address:     "Test Address",
			Telephone:   "Test Telephone",
			LocalityID:  -1,
		}

		repository := &RepositoryMock{}

		service := NewService(repository)

		// Act.
		obtainedCarryID, err := service.Save(ctx, expectedCarry)

		// Assert.
		assert.EqualError(t, err, ErrIncorrectData.Error())
		assert.Equal(t, 0, obtainedCarryID)
		repository.AssertExpectations(t)
	})
	t.Run("it should return error when the carry already exists", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedCarry := domain.Carries{
			ID:          1,
			CID:         "Test Name",
			CompanyName: "Test Company Name",
			Address:     "Test Address",
			Telephone:   "Test Telephone",
			LocalityID:  1,
		}

		repository := &RepositoryMock{}
		repository.On("Save", ctx, expectedCarry).Return(0, ErrDuplicateCarry)

		service := NewService(repository)

		// Act.
		obtainedCarryID, err := service.Save(ctx, expectedCarry)

		// Assert.
		assert.EqualError(t, err, ErrDuplicateCarry.Error())
		assert.Equal(t, 0, obtainedCarryID)
		repository.AssertExpectations(t)
	})
}

func TestService_GetAllCarriesByLocality(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: return all carries by locality
	// DESCRIPTION: Should return all carries by locality
	t.Run("it should return all carries by locality", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedLocalityCarries := []domain.LocalityCarries{
			{
				LocalityID:   "1",
				LocalityName: "Test Locality",
				CarriesCount: 2,
			},
			{
				LocalityID:   "2",
				LocalityName: "Test Locality 2",
				CarriesCount: 20,
			},
		}

		repository := &RepositoryMock{}
		repository.On("GetAllCarriesByLocality", ctx).Return(expectedLocalityCarries, nil)

		service := NewService(repository)

		// Act.
		obtainedLocalityCarries, err := service.GetAllCarriesByLocality(ctx)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, expectedLocalityCarries, obtainedLocalityCarries)
		repository.AssertExpectations(t)
	})
}

func TestService_GetAllCarriesByLocalityID(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: return all carries by locality ID
	// DESCRIPTION: Should return all carries by locality ID
	t.Run("it should return all carries by locality ID", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()

		expectedLocalityCarries := domain.LocalityCarries{
			LocalityID:   "1",
			LocalityName: "Test Locality",
			CarriesCount: 2,
		}

		repository := &RepositoryMock{}
		repository.On("GetAllCarriesByLocalityID", ctx, 1).Return(expectedLocalityCarries, nil)

		service := NewService(repository)

		// Act.
		obtainedLocalityCarries, err := service.GetAllCarriesByLocalityID(ctx, 1)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, expectedLocalityCarries, obtainedLocalityCarries)
		repository.AssertExpectations(t)
	})
}
