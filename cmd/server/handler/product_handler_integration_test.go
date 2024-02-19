package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/mock"
)

func TestIntegrationCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Test for product
	t.Run("should create a new product", func(t *testing.T) {
		// arrange
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
		// Convert product to json
		jsonProduct, _ := json.Marshal(expectedProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		//configure router
		router := gin.Default()
		mockRepo := &product.RepositoryMock{}
		mockRepo.On("Exists", mock.Anything, expectedProduct.ProductCode).Return(false)
		mockRepo.On("Save", mock.Anything, expectedProduct).Return(1, nil)
		service := product.NewService(mockRepo)
		productHandler := NewProduct(service)
		router.POST("/api/v1/products", productHandler.Create())

		// act
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/products", reader)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		// assert
		if response.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, response.Code)
		}
	})
	t.Run("should return error when product code already exists", func(t *testing.T) {
		// arrange
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
		// Convert product to json
		jsonProduct, _ := json.Marshal(expectedProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		//configure router
		router := gin.Default()
		mockRepo := &product.RepositoryMock{}
		mockRepo.On("Exists", mock.Anything, expectedProduct.ProductCode).Return(true)
		service := product.NewService(mockRepo)
		productHandler := NewProduct(service)
		router.POST("/api/v1/products", productHandler.Create())

		// act
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/products", reader)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		// assert
		if response.Code != http.StatusConflict {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
		}
	})

	// Test for productRecord
	t.Run("should create a new product record", func(t *testing.T) {
		// arrange
		expectedProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}
		// Convert product to json
		jsonProductRecord, _ := json.Marshal(expectedProductRecord)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProductRecord)

		//configure router
		router := gin.Default()
		mockRepo := &product.RepositoryMock{}
		service := product.NewService(mockRepo)
		productHandler := NewProduct(service)
		mockRepo.On("Get", mock.Anything, expectedProductRecord.ProductID).Return(domain.Product{}, nil)
		mockRepo.On("CreateProductRecord", mock.Anything, expectedProductRecord).Return(1, nil)
		router.POST("/api/v1/productRecords", productHandler.CreateProductRecord())

		// act
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/productRecords", reader)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		// assert
		if response.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, response.Code)
		}
	})
	t.Run("should return a error if product not exists", func(t *testing.T) {
		// arrange
		expectedProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}
		// Convert product to json
		jsonProductRecord, _ := json.Marshal(expectedProductRecord)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProductRecord)

		//configure router
		router := gin.Default()
		mockRepo := &product.RepositoryMock{}
		service := product.NewService(mockRepo)
		productHandler := NewProduct(service)
		mockRepo.On("Get", mock.Anything, expectedProductRecord.ProductID).Return(domain.Product{}, product.ErrNotFound)
		//mockRepo.On("CreateProductRecord", mock.Anything, expectedProductRecord).Return(1, nil)
		router.POST("/api/v1/productRecords", productHandler.CreateProductRecord())

		// act
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/productRecords", reader)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		// assert
		if response.Code != http.StatusConflict {
			t.Errorf("Expected status code %d, got %d", http.StatusConflict, response.Code)
		}
	})
}
