package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// User story: READ
// As a user I want to be able to read a product
// find_all, find_by_id_non_existent, find_by_id_existent

func Test_Handler_Ping(t *testing.T) {
	t.Run("When the petition is Ping, it should return 'Pong'", func(t *testing.T) {
		// Arrange
		expectedResponse := "pong"

		// Config of mock
		route := "/api/v1/ping"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock)                     // Instance of handler
		req := httptest.NewRequest(http.MethodGet, route, nil) // Request
		w := httptest.NewRecorder()                            // Instance of response

		// Act
		gin.SetMode(gin.TestMode) // Config gin to test mode
		router := gin.New()
		router.GET(route, handler.Ping())
		router.ServeHTTP(w, req) // Execute request

		// Assert
		assert.Equal(t, http.StatusOK, w.Code) // Check status code 200

		var response string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("could not unmarshal response: %v", err)
		}
		assert.Equal(t, expectedResponse, response) // Check body
	})
}

func Test_Handler_Get(t *testing.T) {
	//find_all
	t.Run("when the petition is correct, it should return a list of products. Code 200", func(t *testing.T) {
		//Arrange
		// Config of products to return
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

		// Config of mock
		route := "/api/v1/products"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("GetAll", mock.Anything).Return(expectedProducts, nil)
		handler := NewProduct(handlerMock)                     // Instance of handler
		req := httptest.NewRequest(http.MethodGet, route, nil) // Request
		w := httptest.NewRecorder()                            // Instance of response

		//Act
		gin.SetMode(gin.TestMode) // Config gin to test mode
		router := gin.New()
		router.GET(route, handler.GetAll())
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusOK, w.Code)                                    // Check status code 200
		assert.JSONEq(t, expectedResponseBody(expectedProducts), w.Body.String()) // Check body
		handlerMock.AssertExpectations(t)                                         // Check if mock was called
	})

	// find_by_id_non_existent
	t.Run("when the id does not exist, it should return a code 404", func(t *testing.T) {
		//Arrange
		id := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, id).Return(domain.Product{}, product.ErrNotFound)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.Get())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/%d", id), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusNotFound, w.Code) // Check status code 404
		handlerMock.AssertExpectations(t)            // Check if mock was called

	})

	//find_by_id_existent
	t.Run("when the id exists, it should return a code 200", func(t *testing.T) {
		//Arrange
		id := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, id).Return(domain.Product{}, nil)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.Get())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/%d", id), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusOK, w.Code) // Check status code 200
		handlerMock.AssertExpectations(t)      // Check if mock was called

	})

	// ErrInternalServer -- GetAll
	t.Run("when the service returns an error, it should return a code 500", func(t *testing.T) {
		//Arrange
		route := "/api/v1/products"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("GetAll", mock.Anything).Return([]domain.Product{}, errors.New("error"))
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.GetAll())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
		handlerMock.AssertExpectations(t)                       // Check if mock was called
	})

	// ErrInternalServer -- Get
	t.Run("when the service returns an error, it should return a code 500", func(t *testing.T) {
		//Arrange
		id := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, id).Return(domain.Product{}, errors.New("error"))
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.Get())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/%d", id), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
		handlerMock.AssertExpectations(t)                       // Check if mock was called
	})

	//Err id is not a number
	t.Run("when the id is not a number, it should return a StatusBadRequest", func(t *testing.T) {
		//Arrange
		id := "a"
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.Get())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/%s", id), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusBadRequest, w.Code) // Check status code 400
		handlerMock.AssertExpectations(t)              // Check if mock was called
	})

}

// User story: CREATE
// As a user I want to be able to create a product
// create_ok, create_fail, create_conflict
func TestProduct_Create(t *testing.T) {
	// create_ok
	t.Run("when the petition is correct, it should return a code 201 with the object", func(t *testing.T) {
		//Arrange
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
		route := "/api/v1/products"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Save", mock.Anything, expectedProduct).Return(1, nil)
		handlerMock.On("Create", mock.Anything, expectedProduct).Return(expectedProduct, nil)
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(expectedProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.Create())
		//Request route with id
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusCreated, w.Code) // Check status code 201
	})
	// create_fail
	t.Run("when the petition is incorrect, it should return a code 422", func(t *testing.T) {
		//Arrange
		invalidProduct := domain.Product{
			ID:             1,
			Description:    "", // Invalid description
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
		route := "/api/v1/products"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(invalidProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.Create())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // Check status code 400
		handlerMock.AssertExpectations(t)                       // Check if mock was called

	})
	// create_conflict
	t.Run("when the product already exists, it should return a code 409", func(t *testing.T) {
		//Arrange
		existingProduct := domain.Product{
			ID:             1,
			Description:    "Fresh Milk",
			ExpirationRate: 0.1,
			FreezingRate:   0.05,
			Height:         25,
			Length:         10,
			Netweight:      1,
			ProductCode:    "123456", // Existing product code
			RecomFreezTemp: -4,
			Width:          10,
			ProductTypeID:  1,
			SellerID:       1,
		}
		route := "/api/v1/products"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Save", mock.Anything, existingProduct).Return(0, product.ErrProductCodeExists)
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(existingProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.Create())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusConflict, w.Code) // Check status code 409
		handlerMock.AssertExpectations(t)
	})
	// Err invalid json
	t.Run("when the json is invalid, it should return StatusUnprocessableEntity", func(t *testing.T) {
		//Arrange
		invalidJSON := []byte(`{`)
		route := "/api/v1/products"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		// Create a new reader with the JSON
		reader := bytes.NewReader(invalidJSON)

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.Create())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // Check status code 422
		handlerMock.AssertExpectations(t)
	})
	//Err internal server
	t.Run("when the service returns an error, it should return StatusInternalServerError", func(t *testing.T) {
		//Arrange
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
		route := "/api/v1/products"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Save", mock.Anything, expectedProduct).Return(0, errors.New("error"))
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(expectedProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.Create())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
		handlerMock.AssertExpectations(t)
	})
}

// User story: UPDATE
// As a user I want to be able to update a product
// update_ok, update_nonexistent
func TestProduct_Update(t *testing.T) {
	// update_ok
	t.Run("when the petition is correct, it should return a code 200 with the object", func(t *testing.T) {
		//Arrange
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
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, expectedProduct.ID).Return(domain.Product{}, nil)
		handlerMock.On("Update", mock.Anything, expectedProduct).Return(nil)
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(expectedProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT(route, handler.Update())
		//Request route with id
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/products/%d", expectedProduct.ID), reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusOK, w.Code) // Check status code 200
		handlerMock.AssertExpectations(t)      // Check if mock was called

	})
	// update_nonexistent
	t.Run("when the id does not exist, it should return a code 404", func(t *testing.T) {
		//Arrange
		nonexistentProduct := domain.Product{
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
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, nonexistentProduct.ID).Return(domain.Product{}, product.ErrNotFound)
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(nonexistentProduct)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT(route, handler.Update())
		//Request route with id
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/products/%d", nonexistentProduct.ID), reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusNotFound, w.Code) // Check status code 404
		handlerMock.AssertExpectations(t)            // Check if mock was called
	})
	//Err if id is not a number
	t.Run("when the id is not a number, it should return a StatusBadRequest", func(t *testing.T) {
		//Arrange
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT(route, handler.Update())
		//Request route with id
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/products/%s", "abc"), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusBadRequest, w.Code) // Check status code 400
		handlerMock.AssertExpectations(t)              // Check if mock was called
	})
	//Err internal server in GET
	t.Run("when the service returns an error, it should return a StatusInternalServerError", func(t *testing.T) {
		//Arrange
		id := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, id).Return(domain.Product{}, errors.New("error"))
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT(route, handler.Update())
		//Request route with id
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/products/%d", id), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
		handlerMock.AssertExpectations(t)                       // Check if mock was called
	})
	//Err invalid json
	t.Run("when the json is invalid, it should return a StatusUnprocessableEntity", func(t *testing.T) {
		//Arrange
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, 1).Return(domain.Product{}, nil)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT(route, handler.Update())
		//Request route with id
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/products/%d", 1), bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // Check status code 422
		handlerMock.AssertExpectations(t)                       // Check if mock was called
	})
	//Err ErrProductCodeExists
	t.Run("when the product code exists, it should return a StatusConflict", func(t *testing.T) {
		//Arrange
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
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, 1).Return(domain.Product{}, nil)
		handlerMock.On("Update", mock.Anything, expectedProduct).Return(product.ErrProductCodeExists)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT(route, handler.Update())
		//Request route with id
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/products/%d", 1), reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusConflict, w.Code) // Check status code 409
		handlerMock.AssertExpectations(t)            // Check if mock was called

	})
	//Err internal server in Update
	t.Run("when the service returns an error, it should return a StatusInternalServerError", func(t *testing.T) {
		//Arrange
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
		id := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Get", mock.Anything, id).Return(domain.Product{}, nil)
		handlerMock.On("Update", mock.Anything, expectedProduct).Return(errors.New(ErrInternalServer))
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT(route, handler.Update())
		//Request route with id
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/products/%d", id), reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
		handlerMock.AssertExpectations(t)                       // Check if mock was called
	})
}

// User story: DELETE
// As a user I want to be able to delete a product
// delete_ok, delete_nonexistent
func TestProduct_Delete(t *testing.T) {
	// delete_ok
	t.Run("when the petition is correct, it should return a code 200 with the object", func(t *testing.T) {
		//Arrange
		sellerID := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Delete", mock.Anything, sellerID).Return(nil)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE(route, handler.Delete())
		//Request route with id
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/products/%d", sellerID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusNoContent, w.Code) // Check status code 200
		handlerMock.AssertExpectations(t)             // Check if mock was called
	})

	// delete_nonexistent
	t.Run("when the id does not exist, it should return a code 404", func(t *testing.T) {
		//Arrange
		nonexistentProductID := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Delete", mock.Anything, nonexistentProductID).Return(product.ErrNotFound)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE(route, handler.Delete())
		//Request route with id
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/products/%d", nonexistentProductID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusNotFound, w.Code) // Check status code 404
		handlerMock.AssertExpectations(t)            // Check if mock was called
	})

	// Err id is not a number
	t.Run("when the id is not a number, it should return a code StatusBadRequest", func(t *testing.T) {
		//Arrange
		nonexistentProductID := "a"
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE(route, handler.Delete())
		//Request route with id
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/products/%s", nonexistentProductID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusBadRequest, w.Code) // Check status code 400
		handlerMock.AssertExpectations(t)              // Check if mock was called
	})

	// Err internal server in Delete
	t.Run("when the service returns an error, it should return a StatusInternalServerError", func(t *testing.T) {
		//Arrange
		sellerID := 1
		route := "/api/v1/products/:id"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("Delete", mock.Anything, sellerID).Return(errors.New(ErrInternalServer))
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE(route, handler.Delete())
		//Request route with id
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/products/%d", sellerID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
		handlerMock.AssertExpectations(t)                       // Check if mock was called
	})

}

func TestProductRecord_Create(t *testing.T) {
	t.Run("when the petition is correct, it should return a code 200 with the object", func(t *testing.T) {
		//Arrange
		expectedProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}
		expectedProductID := 1
		route := "/api/v1/productRecords"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("CreateProductRecord", mock.Anything, expectedProductRecord).Return(expectedProductID, nil)
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(expectedProductRecord)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.CreateProductRecord())
		//Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusCreated, w.Code) // Check status code 201
		handlerMock.AssertExpectations(t)
	})

	t.Run("When the product not exist, it should return StatusConflict", func(t *testing.T) {
		//Arrange
		expectedProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     100,
		}
		//expectedProduct to json
		jsonProduct, _ := json.Marshal(expectedProductRecord)
		route := "/api/v1/productRecords"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("CreateProductRecord", mock.Anything, expectedProductRecord).Return(0, product.ErrNotFound)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.CreateProductRecord())
		//Request route
		req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonProduct))
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusConflict, w.Code) // Check status code 409
		handlerMock.AssertExpectations(t)
	})

	// create_fail_due_to_invalid_json
	t.Run("when the json is invalid, it should return a code 422", func(t *testing.T) {
		//Arrange
		invalidJSON := "invalid json"
		route := "/api/v1/productRecords"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		// Create a new reader with the invalid JSON
		reader := bytes.NewReader([]byte(invalidJSON))

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.CreateProductRecord())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // Check status code 422
	})

	t.Run("when the fields are invalid, it should return a code 422", func(t *testing.T) {
		//Arrange
		invalidProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "", // Invalid last update
			PurchasePrice: 0,  // Invalid purchase price
			SalePrice:     0,  // Invalid sale price
			ProductID:     0,  // Invalid product ID
		}
		route := "/api/v1/productRecords"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(invalidProductRecord)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.CreateProductRecord())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // Check status code 422
	})

	// create_fail_due_to_invalid_date
	t.Run("when the date is invalid, it should return a code 422", func(t *testing.T) {
		//Arrange
		invalidProductRecord := domain.ProductRecordCreate{
			LastUpdate:    "invalid date", // Invalid date
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}
		route := "/api/v1/productRecords"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(invalidProductRecord)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.CreateProductRecord())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // Check status code 422
	})

	// create_fail_due_to_internal_error
	t.Run("when there is an internal server error, it should return a code 500", func(t *testing.T) {
		//Arrange
		productRecord := domain.ProductRecordCreate{
			LastUpdate:    "2021-04-04",
			PurchasePrice: 10,
			SalePrice:     15,
			ProductID:     44,
		}
		route := "/api/v1/productRecords"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("CreateProductRecord", mock.Anything, productRecord).Return(0, errors.New("internal server error"))
		handler := NewProduct(handlerMock) // Instance of handler

		// Convert product to json
		jsonProduct, _ := json.Marshal(productRecord)

		// Create a new reader with the JSON
		reader := bytes.NewReader(jsonProduct)

		// Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST(route, handler.CreateProductRecord())
		// Request route
		req := httptest.NewRequest(http.MethodPost, route, reader)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
	})
}

func TestProductRecord_Get(t *testing.T) {
	t.Run("when the petition is correct, it should return a code 200 with the object", func(t *testing.T) {
		//Arrange
		expectedProductRecord := []domain.ProductRecordGet{
			{
				ProductID:   44,
				Description: "Test",
				RecordCount: 1,
			},
		}
		expectedProductID := 1
		route := "/api/v1/products/reportRecords"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("GetProductRecord", mock.Anything, expectedProductID).Return(expectedProductRecord, nil)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.GetProductRecord())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/reportRecords?id=%d", expectedProductID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusOK, w.Code) // Check status code 200
		handlerMock.AssertExpectations(t)      // Check if mock was called
	})

	t.Run("when the id is string, it should return StatusBadRequest and invalidID error", func(t *testing.T) {
		//Arrange
		expectedProductID := "string"
		route := "/api/v1/products/reportRecords"
		handlerMock := &product.ServiceMock{}
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.GetProductRecord())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/reportRecords?id=%s", expectedProductID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusBadRequest, w.Code) // Check status code 400
		handlerMock.AssertExpectations(t)              // Check if mock was called

	})

	t.Run("when the product record does not exist, it should return a code 404", func(t *testing.T) {
		//Arrange
		expectedProductRecord := []domain.ProductRecordGet{
			{
				ProductID:   44,
				Description: "Test",
				RecordCount: 1,
			},
		}
		nonexistentProductID := 1
		route := "/api/v1/products/reportRecords"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("GetProductRecord", mock.Anything, nonexistentProductID).Return(expectedProductRecord, product.ErrNotFound)
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.GetProductRecord())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/reportRecords?id=%d", nonexistentProductID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusNotFound, w.Code) // Check status code 404
		handlerMock.AssertExpectations(t)            // Check if mock was called
	})

	t.Run("when there is an internal server error, it should return a code 500", func(t *testing.T) {
		//Arrange
		expectedProductRecord := []domain.ProductRecordGet{
			{
				ProductID:   44,
				Description: "Test",
				RecordCount: 1,
			},
		}
		productID := 1
		route := "/api/v1/products/reportRecords"
		handlerMock := &product.ServiceMock{}
		handlerMock.On("GetProductRecord", mock.Anything, productID).Return(expectedProductRecord, errors.New("internal server error"))
		handler := NewProduct(handlerMock) // Instance of handler

		//Config gin to test mode
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET(route, handler.GetProductRecord())
		//Request route with id
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/reportRecords?id=%d", productID), nil)
		w := httptest.NewRecorder() // Instance of response

		//Act
		router.ServeHTTP(w, req) // Execute request

		//Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code) // Check status code 500
		handlerMock.AssertExpectations(t)                       // Check if mock was called
	})
}

func expectedResponseBody(products []domain.Product) string {
	// Create map with data key
	reponseBody := map[string]interface{}{"data": products}
	// Convert products to json
	jsonProducts, err := json.Marshal(reponseBody)
	if err != nil {
		panic(err) // Panic if error
	}
	return string(jsonProducts)
}
