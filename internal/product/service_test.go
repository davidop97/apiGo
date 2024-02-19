package product

import (
	"context"
	"errors"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Test for create product method
// User story: CREATE
// Cases: create_ok, create_conflict
func Test_Service_Create(t *testing.T) {
	// create_ok
	t.Run("should return a domain.Product and nil if there is no error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "123456",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Exists", ctx, expectedProduct.ProductCode).Return(false)
		repositoryMock.On("Save", ctx, expectedProduct).Return(expectedProduct.ID, nil)
		service := NewService(repositoryMock)

		//Act
		productID, err := service.Save(ctx, expectedProduct)
		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct.ID, productID)
		repositoryMock.AssertExpectations(t)
	})

	// create_conflict
	t.Run("should return error if the product already exists", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "123456",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Exists", ctx, expectedProduct.ProductCode).Return(true)
		service := NewService(repositoryMock)

		//Act
		productID, err := service.Save(ctx, expectedProduct)
		//Assert
		assert.Error(t, err)
		assert.Equal(t, 0, productID)
		repositoryMock.AssertExpectations(t)
	})

	// Err saving product
	t.Run("should return error if there is an error saving the product", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "123456",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Exists", ctx, expectedProduct.ProductCode).Return(false)
		repositoryMock.On("Save", ctx, expectedProduct).Return(0, ErrorSavingProduct)
		service := NewService(repositoryMock)

		//Act
		productID, err := service.Save(ctx, expectedProduct)
		//Assert
		assert.Error(t, err)
		assert.Equal(t, 0, productID)
		repositoryMock.AssertExpectations(t)
	})
}

// Test for Get method
// User story: READ
// Cases: find_all, find_by_id_non_existent, find_by_id_existent
func Test_Service_Read(t *testing.T) {
	// find_all
	t.Run("should return a slice of domain.Product and nil if there is no error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProducts := []domain.Product{
			{
				ID:             1,
				Description:    "Fresh Milk",
				ExpirationRate: 0.1,
				FreezingRate:   0.05,
				Height:         25,
				Length:         10,
				Netweight:      1,
				ProductCode:    "123456",
				RecomFreezTemp: -4,
				Width:          10,
				ProductTypeID:  1,
				SellerID:       1,
			},
			{
				ID:             2,
				Description:    "Apple Juice",
				ExpirationRate: 0.2,
				FreezingRate:   0.1,
				Height:         25,
				Length:         10,
				Netweight:      1,
				ProductCode:    "654321",
				RecomFreezTemp: 0,
				Width:          10,
				ProductTypeID:  2,
				SellerID:       1,
			},
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("GetAll", ctx).Return(expectedProducts, nil)
		service := NewService(repositoryMock)

		//Act
		products, err := service.GetAll(ctx)
		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
		repositoryMock.AssertExpectations(t)
	})

	// find_by_id_non_existent
	t.Run("should return a domain.Product and nil if there is no error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		nonExistID := 999
		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, nonExistID).Return(domain.Product{}, ErrNotFound)
		service := NewService(repositoryMock)

		//Act
		products, err := service.Get(ctx, nonExistID)

		//Assert
		assert.Error(t, err)
		assert.Equal(t, domain.Product{}, products)
		repositoryMock.AssertExpectations(t)
	})

	// find_by_id_existent
	t.Run("should return a domain.Product and nil if there is no error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		existID := 1
		expectedProduct := domain.Product{
			ID: existID,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, existID).Return(expectedProduct, nil)
		service := NewService(repositoryMock)

		//Act
		products, err := service.Get(ctx, existID)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, products)
		repositoryMock.AssertExpectations(t)
	})

	// Err not found product
	t.Run("should return a error if the product not found in GetAll", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("GetAll", ctx).Return([]domain.Product{}, ErrNotFound)
		service := NewService(repositoryMock)

		//Act
		products, err := service.GetAll(ctx)

		//Assert
		assert.Error(t, err)
		assert.Equal(t, []domain.Product(nil), products)
		repositoryMock.AssertExpectations(t)
	})
}

// Test fot update method
// User story: UPDATE
// update_exist, update_non_exist
func Test_Service_Update(t *testing.T) {
	//update_exist
	t.Run("should return a domain.Product and nil if there is no error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "ABC123",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, expectedProduct.ID).Return(expectedProduct, nil)
		repositoryMock.On("Update", ctx, expectedProduct).Return(nil)
		service := NewService(repositoryMock)

		//Act
		err := service.Update(ctx, expectedProduct)

		//Assert
		assert.NoError(t, err)
		repositoryMock.AssertExpectations(t)
	})

	//update_non_exist
	t.Run("should return err if the product does not exist", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProduct := domain.Product{
			ID:             999,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "123456",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, expectedProduct.ID).Return(domain.Product{}, ErrNotFound)
		service := NewService(repositoryMock)

		//Act
		err := service.Update(ctx, expectedProduct)

		//Assert
		assert.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

	//err ErrProductCodeExists
	t.Run("should return err if the product code exists", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		oldProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "123456",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		newProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "12345",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, oldProduct.ID).Return(oldProduct, nil)
		repositoryMock.On("Exists", ctx, newProduct.ProductCode).Return(true)
		service := NewService(repositoryMock)

		//Act
		err := service.Update(ctx, newProduct)

		//Assert
		assert.Error(t, err)
		assert.Equal(t, ErrProductCodeExists, err)
		repositoryMock.AssertExpectations(t)

	})

	//Err if update fail
	t.Run("should return err if the update fail", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		oldProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "123456",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		newProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "12345",
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, oldProduct.ID).Return(oldProduct, nil)
		repositoryMock.On("Exists", ctx, newProduct.ProductCode).Return(false)
		repositoryMock.On("Update", ctx, newProduct).Return(errors.New("update fail"))
		service := NewService(repositoryMock)

		//Act
		err := service.Update(ctx, newProduct)

		//Assert
		assert.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})
}

// Test for delete method
// User story: DELETE
// delete_non_exist, delete_ok
func TestService_Delete(t *testing.T) {
	// delete_non_exist
	t.Run("should return err if the product does not exist", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		nonExistID := 999
		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, nonExistID).Return(domain.Product{}, ErrNotFound)
		service := NewService(repositoryMock)

		//Act
		err := service.Delete(ctx, nonExistID)

		//Assert
		assert.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})

	// delete_ok
	t.Run("should return nil if there is no error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		existID := 1
		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, existID).Return(domain.Product{}, nil)
		repositoryMock.On("Delete", ctx, existID).Return(nil)
		service := NewService(repositoryMock)

		//Act
		err := service.Delete(ctx, existID)

		//Assert
		assert.NoError(t, err)
		repositoryMock.AssertExpectations(t)
	})

	// Err if delete fail
	t.Run("should return err if the delete fail", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		existID := 1
		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, existID).Return(domain.Product{}, nil)
		repositoryMock.On("Delete", ctx, existID).Return(errors.New("delete fail"))
		service := NewService(repositoryMock)

		//Act
		err := service.Delete(ctx, existID)

		//Assert
		assert.Error(t, err)
		repositoryMock.AssertExpectations(t)
	})
}

// Test for create method
// User story: POST
// create_ok, create_err
func TestService_CreateProductRecord(t *testing.T) {
	// create_ok
	t.Run("should return nil if there is no error", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, expectedProductRecord.ProductID).Return(domain.Product{}, nil)
		repositoryMock.On("CreateProductRecord", ctx, expectedProductRecord).Return(0, nil)
		service := NewService(repositoryMock)

		//Act
		id, err := service.CreateProductRecord(ctx, expectedProductRecord)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, id)
		repositoryMock.AssertExpectations(t)
	})
	// create_err
	t.Run("should return err if the product does not exist", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, expectedProductRecord.ProductID).Return(domain.Product{}, ErrNotFound)
		service := NewService(repositoryMock)

		//Act
		id, err := service.CreateProductRecord(ctx, expectedProductRecord)

		//Assert
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		repositoryMock.AssertExpectations(t)
	})
	//create_err
	t.Run("should return err if the create fail", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, expectedProductRecord.ProductID).Return(domain.Product{}, nil)
		repositoryMock.On("CreateProductRecord", ctx, expectedProductRecord).Return(0, errors.New("create fail"))
		service := NewService(repositoryMock)

		//Act
		id, err := service.CreateProductRecord(ctx, expectedProductRecord)

		//Assert
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		repositoryMock.AssertExpectations(t)
	})
}

// Test for getProductRecord method
// User story: GET
// get_ok, get_err
func TestService_GetProductRecord(t *testing.T) {
	// get_ok
	t.Run("should return a product record", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		expectedProductRecord := []domain.ProductRecordGet{
			{
				ProductID:   44,
				Description: "Test",
				RecordCount: 1,
			},
		}
		idProduct := 44

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, idProduct).Return(domain.Product{}, nil)
		repositoryMock.On("GetProductRecord", ctx, idProduct).Return(expectedProductRecord, nil)
		service := NewService(repositoryMock)

		//Act
		productRecord, err := service.GetProductRecord(ctx, idProduct)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProductRecord, productRecord)
		repositoryMock.AssertExpectations(t)
	})
	// get_err for product not found
	t.Run("should return err if the product does not exist", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		idProduct := 44

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, idProduct).Return(domain.Product{}, ErrNotFound)
		service := NewService(repositoryMock)

		//Act
		productRecord, err := service.GetProductRecord(ctx, idProduct)

		//Assert
		assert.Error(t, err)
		assert.Equal(t, []domain.ProductRecordGet(nil), productRecord)
		repositoryMock.AssertExpectations(t)
	})
	// get_err for get product record fail
	t.Run("should return err if the get product record fail", func(t *testing.T) {
		//Arrange
		ctx := context.Background()

		idProduct := 44

		repositoryMock := &RepositoryMock{}
		repositoryMock.On("Get", ctx, idProduct).Return(domain.Product{}, nil)
		repositoryMock.On("GetProductRecord", ctx, idProduct).Return([]domain.ProductRecordGet(nil), errors.New("get product record fail"))
		service := NewService(repositoryMock)

		//Act
		productRecord, err := service.GetProductRecord(ctx, idProduct)

		//Assert
		assert.Error(t, err)
		assert.Equal(t, []domain.ProductRecordGet(nil), productRecord)
		repositoryMock.AssertExpectations(t)
	})
}
