package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestHandler_Read is a test for Read handler
func TestHandler_WarehouseRead(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: Read all warehouses
	// DESCRIPTION: A code 200 should be returned with a list of warehouses
	t.Run("it should return all warehouses", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedWarehouses := []domain.Warehouse{
			{
				ID:                 1,
				Address:            "Test Address",
				Telephone:          "Test Telephone",
				WarehouseCode:      "Test WarehouseCode",
				MinimumCapacity:    100,
				MinimumTemperature: 10,
			},
			{
				ID:                 2,
				Address:            "Test Address 2",
				Telephone:          "Test Telephone 2",
				WarehouseCode:      "Test WarehouseCode 2",
				MinimumCapacity:    200,
				MinimumTemperature: 20,
			},
		}
		expectedStatus := 200
		expectedBody := `{"data":[{"id":1,"address":"Test Address","telephone":"Test Telephone","warehouse_code":"Test WarehouseCode",
						"minimum_capacity":100,"minimum_temperature":10},
						{"id":2,"address":"Test Address 2","telephone":"Test Telephone 2","warehouse_code":"Test WarehouseCode 2",
						"minimum_capacity":200,"minimum_temperature":20}]}`

		//Create service mock
		service := &warehouse.ServiceMock{}
		service.On("GetAll", mock.Anything).Return(expectedWarehouses, nil)

		//Create handler
		handler := NewWarehouse(service)
		server.GET("/api/v1//warehouses", handler.GetAll())

		//Create request and response
		req, err := http.NewRequest(http.MethodGet, "/api/v1/warehouses", nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: get a warehouse by id
	// DESCRIPTION: A code 200 should be returned with a warehouse
	t.Run("it should return a warehouse by id", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}
		expectedStatus := 200
		expectedBody := `{"data":{"id":1,"address":"Test Address","telephone":"Test Telephone","warehouse_code":"Test WarehouseCode",
						"minimum_capacity":100,"minimum_temperature":10}}`

		service := &warehouse.ServiceMock{}
		service.On("Get", mock.Anything, 1).Return(expectedWarehouse, nil)

		handler := NewWarehouse(service)
		server.GET("/api/v1/warehouses/:id", handler.Get())

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/warehouses/%d", expectedWarehouse.ID), nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: get a warehouse by id that does not exist
	// DESCRIPTION: A code 404 should be returned with a not found error
	t.Run("it should return a 404 not found error when the warehouse does not exist", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedStatus := 404
		expectedBody := `{"message": "Not found"}`
		warehouseID := 1

		service := &warehouse.ServiceMock{}
		service.On("Get", mock.Anything, 1).Return(domain.Warehouse{}, warehouse.ErrNotFound)

		handler := NewWarehouse(service)
		server.GET("/api/v1/warehouses/:id", handler.Get())
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/warehouses/%d", warehouseID), nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return a 500 internal server error when the id is not a number", func(t *testing.T) {
		server := gin.New()

		expectedStatus := 500
		expectedBody := `{"message": "Internal server error"}`
		warehouseID := "test"

		service := &warehouse.ServiceMock{}

		handler := NewWarehouse(service)
		server.GET("/api/v1/warehouses/:id", handler.Get())
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/warehouses/"+warehouseID), nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
}

// TestHandler_Delete is a test for Delete handler
func TestHandler_WarehouseDelete(t *testing.T) {
	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete a warehouse by id
	// DESCRIPTION: A code 204 should be returned with no content
	t.Run("it should delete a warehouse by id", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedStatus := 204
		warehouseID := 1

		service := &warehouse.ServiceMock{}
		service.On("Delete", mock.Anything, 1).Return(nil)

		handler := NewWarehouse(service)
		server.DELETE("/api/v1/warehouses/:id", handler.Delete())

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/warehouses/%d", warehouseID), nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		service.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete a warehouse by id that does not exist
	// DESCRIPTION: A code 404 should be returned with a not found error
	t.Run("it should return a 404 not found error when the warehouse does not exist", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedStatus := 404
		expectedBody := `{"message": "Warehouse not found"}`
		warehouseID := 1

		service := &warehouse.ServiceMock{}
		service.On("Delete", mock.Anything, 1).Return(warehouse.ErrNotFound)

		handler := NewWarehouse(service)
		server.DELETE("/api/v1/warehouses/:id", handler.Delete())

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/warehouses/%d", warehouseID), nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return a 500 internal server error when the id is not a number", func(t *testing.T) {
		server := gin.New()

		expectedStatus := 500
		expectedBody := `{"message": "Internal server error"}`
		warehouseID := "test"

		service := &warehouse.ServiceMock{}

		handler := NewWarehouse(service)
		server.DELETE("/api/v1/warehouses/:id", handler.Delete())
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/warehouses/"+warehouseID), nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
}

// TestHandler_Create is a test for Create handler
func TestHandler_WarehouseCreate(t *testing.T) {
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create a warehouse
	// DESCRIPTION: A code 201 should be returned with a warehouse
	t.Run("it should create a warehouse", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}
		expectedStatus := 201
		expectedBody := `{"data":{"id":1,"address":"Test Address","telephone":"Test Telephone","warehouse_code":"Test WarehouseCode",
						"minimum_capacity":100,"minimum_temperature":10}}`

		service := &warehouse.ServiceMock{}
		service.On("Save", mock.Anything, expectedWarehouse).Return(expectedWarehouse.ID, nil)

		handler := NewWarehouse(service)
		server.POST("/api/v1/warehouses", handler.Create())

		body, _ := json.Marshal(expectedWarehouse)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer([]byte(body)))
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create a warehouse with an invalid request body
	// DESCRIPTION: A code 422 should be returned with an invalid request body error
	t.Run("it should return a 422 error when the request body is invalid", func(t *testing.T) {
		//Arrange
		server := gin.New()

		whouse := map[string]interface{}{"telephone": "Test Telephone", "warehouse_code": "Test WarehouseCode", "minimum_capacity": 100, "minimum_temperature": 10}
		expectedStatus := 422
		expectedBody := `{"message": "Invalid request body"}`

		service := &warehouse.ServiceMock{}

		handler := NewWarehouse(service)
		server.POST("/api/v1/warehouses", handler.Create())

		body, _ := json.Marshal(whouse)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer([]byte(body)))
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create a warehouse that already exists
	// DESCRIPTION: A code 409 should be returned with a warehouse already exists error
	t.Run("it should return a 409 error when the warehouse already exists", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}
		expectedStatus := 409
		expectedBody := `{"message": "Warehouse already exists"}`

		service := &warehouse.ServiceMock{}
		service.On("Save", mock.Anything, expectedWarehouse).Return(0, warehouse.ErrDuplicateWarehouse)

		handler := NewWarehouse(service)
		server.POST("/api/v1/warehouses", handler.Create())

		body, _ := json.Marshal(expectedWarehouse)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer([]byte(body)))
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create a warehouse with an invalid field value
	// DESCRIPTION: A code 422 should be returned with an incorrect data error
	t.Run("it should return a code 422 when the body has an invalid field value", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}
		expectedStatus := 422
		expectedBody := `{"message": "Incorrect data"}`

		service := &warehouse.ServiceMock{}
		service.On("Save", mock.Anything, expectedWarehouse).Return(0, warehouse.ErrIncorrectData)

		handler := NewWarehouse(service)
		server.POST("/api/v1/warehouses", handler.Create())

		body, _ := json.Marshal(expectedWarehouse)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer([]byte(body)))
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
}

// TestHandler_Update is a test for Update handler
func TestHandler_WarehouseUpdate(t *testing.T) {
	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: update a warehouse
	// DESCRIPTION: A code 200 should be returned with a warehouse
	t.Run("it should update a warehouse", func(t *testing.T) {
		//Arrange
		server := gin.New()

		whouse := map[string]interface{}{"minimum_capacity": 150, "minimum_temperature": 20}
		oldWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    100,
			MinimumTemperature: 10,
		}
		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			Address:            "Test Address",
			Telephone:          "Test Telephone",
			WarehouseCode:      "Test WarehouseCode",
			MinimumCapacity:    150,
			MinimumTemperature: 20,
		}
		expectedStatus := 200
		expectedBody := `{"data":{"id":1,"address":"Test Address","telephone":"Test Telephone","warehouse_code":"Test WarehouseCode",
						"minimum_capacity":150,"minimum_temperature":20}}`

		service := &warehouse.ServiceMock{}
		service.On("Update", mock.Anything, expectedWarehouse).Return(nil)
		service.On("Get", mock.Anything, 1).Return(oldWarehouse, nil)

		handler := NewWarehouse(service)
		server.PUT("/api/v1/warehouses/:id", handler.Update())

		body, _ := json.Marshal(whouse)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/warehouses/%d", expectedWarehouse.ID), bytes.NewBuffer([]byte(body)))
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: update a warehouse that does not exist
	// DESCRIPTION: A code 404 should be returned with a warehouse not found error
	t.Run("it should return a 404 not found error when the warehouse does not exist", func(t *testing.T) {
		//Arrange
		server := gin.New()

		whouse := map[string]interface{}{"minimum_capacity": 150, "minimum_temperature": 20}
		expectedStatus := 404
		expectedBody := `{"message": "Warehouse not found"}`
		warehouseID := 1

		service := &warehouse.ServiceMock{}
		service.On("Get", mock.Anything, warehouseID).Return(domain.Warehouse{}, warehouse.ErrNotFound)

		handler := NewWarehouse(service)
		server.PUT("/api/v1/warehouses/:id", handler.Update())

		body, _ := json.Marshal(whouse)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/warehouses/%d", warehouseID), bytes.NewBuffer([]byte(body)))
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return a 500 internal server error when the id is not a number", func(t *testing.T) {
		server := gin.New()

		expectedStatus := 500
		expectedBody := `{"message": "Internal server error"}`
		warehouseID := "test"

		service := &warehouse.ServiceMock{}

		handler := NewWarehouse(service)
		server.PUT("/api/v1/warehouses/:id", handler.Update())
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/warehouses/"+warehouseID), nil)
		res := httptest.NewRecorder()

		//Act
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
}
