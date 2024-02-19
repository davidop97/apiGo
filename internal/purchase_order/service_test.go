package purchase_order

import (
	"context"
	"errors"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Tests for Purchase Order services

// TestService_SavePurchaseOrder
// Test all cases for Save functions
// Tested methods: Save
func TestService_SavePurchaseOrder(t *testing.T) {
	// ASSOCIATED USER STORY: SAVE
	// EDGE CASE: save_ok
	// DESCRIPTION: If it contains the required fields it will be save
	t.Run("Successfully save purchase orders", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedPurchaseOrderID := 1
		buyerID := 2
		productRecordID := 3
		poToSave := domain.PurchaseOrder{
			ID:              expectedPurchaseOrderID,
			OrderNumber:     "123456",
			OrderDate:       "2024-01-01",
			TrackingCode:    "654321",
			BuyerID:         buyerID,
			ProductRecordID: productRecordID,
			OrderStatusID:   1,
		}

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsPurchaseOrder", ctx, expectedPurchaseOrderID).Return(false)
		repoMock.On("ExistsBuyer", ctx, buyerID).Return(true)
		repoMock.On("ExistsProductsRecord", ctx, productRecordID).Return(true)
		repoMock.On("Save", ctx, poToSave).Return(expectedPurchaseOrderID, nil)

		// call to service interface
		service := NewService(repoMock)

		// act
		idResult, err := service.Save(ctx, poToSave)

		// assert
		// check no error occurs
		assert.NoError(t, err)
		// check that the id returned is equal to the one we expect
		assert.Equal(t, expectedPurchaseOrderID, idResult)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: SAVE
	// EDGE CASE: save_purchase_order_exists
	// DESCRIPTION: If purchase order already exists it cannot be saved.
	t.Run("should return zero as a id when purchase order already exists ", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedErr := ErrPurchaseOrderAlreadyExists
		expectedPurchaseOrderID := 0
		buyerID := 2
		productRecordID := 3
		poToSave := domain.PurchaseOrder{
			ID:              expectedPurchaseOrderID,
			OrderNumber:     "123456",
			OrderDate:       "2024-01-01",
			TrackingCode:    "654321",
			BuyerID:         buyerID,
			ProductRecordID: productRecordID,
			OrderStatusID:   1,
		}

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsPurchaseOrder", ctx, expectedPurchaseOrderID).Return(true)

		// call to service interface
		service := NewService(repoMock)

		// act
		idResult, err := service.Save(ctx, poToSave)

		// assert
		// check no error occurs
		assert.EqualError(t, expectedErr, err.Error())
		// check that the id returned is equal to the one we expect
		assert.Equal(t, expectedPurchaseOrderID, idResult)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: SAVE
	// EDGE CASE: save_buyer_does_not_exists
	// DESCRIPTION: If buyer does not exists it cannot be saved.
	t.Run("should return zero as a id when buyer does not exists ", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedErr := ErrBuyerIDNotExists
		expectedPurchaseOrderID := 0
		buyerID := 2
		productRecordID := 3
		poToSave := domain.PurchaseOrder{
			ID:              expectedPurchaseOrderID,
			OrderNumber:     "123456",
			OrderDate:       "2024-01-01",
			TrackingCode:    "654321",
			BuyerID:         buyerID,
			ProductRecordID: productRecordID,
			OrderStatusID:   1,
		}

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsPurchaseOrder", ctx, expectedPurchaseOrderID).Return(false)
		repoMock.On("ExistsBuyer", ctx, buyerID).Return(false)

		// call to service interface
		service := NewService(repoMock)

		// act
		idResult, err := service.Save(ctx, poToSave)

		// assert
		// check no error occurs
		assert.EqualError(t, expectedErr, err.Error())
		// check that the id returned is equal to the one we expect
		assert.Equal(t, expectedPurchaseOrderID, idResult)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: SAVE
	// EDGE CASE: save_product_record_does_not_exists
	// DESCRIPTION: If product record does not exists it cannot be saved.
	t.Run("should return zero as a id when product record does not exists ", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedErr := ErrProductsRecordIDNotExits
		expectedPurchaseOrderID := 0
		buyerID := 2
		productRecordID := 3
		poToSave := domain.PurchaseOrder{
			ID:              expectedPurchaseOrderID,
			OrderNumber:     "123456",
			OrderDate:       "2024-01-01",
			TrackingCode:    "654321",
			BuyerID:         buyerID,
			ProductRecordID: productRecordID,
			OrderStatusID:   1,
		}

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsPurchaseOrder", ctx, expectedPurchaseOrderID).Return(false)
		repoMock.On("ExistsBuyer", ctx, buyerID).Return(true)
		repoMock.On("ExistsProductsRecord", ctx, productRecordID).Return(false)

		// call to service interface
		service := NewService(repoMock)

		// act
		idResult, err := service.Save(ctx, poToSave)

		// assert
		// check no error occurs
		assert.EqualError(t, expectedErr, err.Error())
		// check that the id returned is equal to the one we expect
		assert.Equal(t, expectedPurchaseOrderID, idResult)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: SAVE
	// EDGE CASE: save_not_ok
	// DESCRIPTION: If it contains the required fields it will be save
	t.Run("should return zero as a id when save function return error", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// expected result
		expectedErr := errors.New("some errors")
		expectedPurchaseOrderID := 0
		buyerID := 2
		productRecordID := 3
		poToSave := domain.PurchaseOrder{
			ID:              expectedPurchaseOrderID,
			OrderNumber:     "123456",
			OrderDate:       "2024-01-01",
			TrackingCode:    "654321",
			BuyerID:         buyerID,
			ProductRecordID: productRecordID,
			OrderStatusID:   1,
		}

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsPurchaseOrder", ctx, expectedPurchaseOrderID).Return(false)
		repoMock.On("ExistsBuyer", ctx, buyerID).Return(true)
		repoMock.On("ExistsProductsRecord", ctx, productRecordID).Return(true)
		repoMock.On("Save", ctx, poToSave).Return(expectedPurchaseOrderID, errors.New("some errors"))

		// call to service interface
		service := NewService(repoMock)

		// act
		idResult, err := service.Save(ctx, poToSave)

		// assert
		// check no error occurs
		assert.EqualError(t, expectedErr, err.Error())
		// check that the id returned is equal to the one we expect
		assert.Equal(t, expectedPurchaseOrderID, idResult)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})
}

// TestService_PurchaseOrdersByBuyer
// Test all cases for PurchaseOrdersByBuyers functions
// Tested methods: PurchaseOrdersByBuyers
func Test_PurchaseOrdersByBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: PurchaseOrdersByBuyer
	// EDGE CASE: purchaseOrdersByBuyer_ok
	// DESCRIPTION: If buyer exists, a report is created with amount of purchase order by buyer
	t.Run("should return  purchase orders list when buyer exists", func(t *testing.T) {
		// arrange
		buyerID := 1

		ctx := context.Background()

		expectReport := []domain.PurchaseOrdersByBuyer{
			{
				ID:                  1,
				CardNumberID:        "CARD1",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 3,
			},
			{
				ID:                  2,
				CardNumberID:        "CARD2",
				FirstName:           "Jane",
				LastName:            "Smith",
				PurchaseOrdersCount: 0,
			},
		}

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsBuyer", ctx, buyerID).Return(true)
		repoMock.On("PurchaseOrdersByBuyers", ctx, buyerID).Return(expectReport, nil)

		// call to service interface
		service := NewService(repoMock)

		// act
		idResult, err := service.PurchaseOrdersByBuyer(ctx, buyerID)

		// assert
		// check no error occurs
		assert.NoError(t, err)
		// check that the id returned is equal to the one we expect
		assert.Equal(t, expectReport, idResult)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: PurchaseOrdersByBuyer
	// EDGE CASE: purchaseOrdersByBuyer_buyer_does_not_exists
	// DESCRIPTION: If buyer does not exists, a report is created with amount of purchase order by buyer
	t.Run("should return  purchase orders list when buyer exists", func(t *testing.T) {
		// arrange
		buyerID := 1

		ctx := context.Background()
		expectedError := ErrBuyerIDNotExists

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsBuyer", ctx, buyerID).Return(false)

		// call to service interface
		service := NewService(repoMock)

		// act
		report, err := service.PurchaseOrdersByBuyer(ctx, buyerID)

		// assert
		// check no error occurs
		assert.EqualError(t, expectedError, err.Error())
		// check that the id returned is equal to the one we expect
		assert.Nil(t, report)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: PurchaseOrdersByBuyer
	// EDGE CASE: purchaseOrdersByBuyer_fail
	// DESCRIPTION: If an error occurs when obtain the purchase order by buyer report, return an nil, err
	t.Run("should return an error when try to get purchase orders list without buyers", func(t *testing.T) {
		// arrange
		buyerID := 1
		ctx := context.Background()
		expectedError := errors.New("some errors were encountered")
		expectedReport := []domain.PurchaseOrdersByBuyer{}

		// create mock for repository
		repoMock := &RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsBuyer", ctx, buyerID).Return(true)
		repoMock.On("PurchaseOrdersByBuyers", ctx, buyerID).Return(expectedReport, expectedError)

		// call to service interface
		service := NewService(repoMock)

		// act
		report, err := service.PurchaseOrdersByBuyer(ctx, buyerID)

		// assert
		// check no error occurs
		assert.EqualError(t, expectedError, err.Error())
		// check that the id returned is equal to the one we expect
		assert.Nil(t, report)
		// check if the repository was called with the expected parameters
		repoMock.AssertExpectations(t)
	})
}
