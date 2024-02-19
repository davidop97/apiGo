package handler

import (
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

func TestHandlerIntegration_WarehouseRead(t *testing.T) {
	t.Run("it should return a 200 with a warehouse", func(t *testing.T) {
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

		repo := &warehouse.RepositoryMock{}
		repo.On("Get", mock.Anything, expectedWarehouse.ID).Return(expectedWarehouse, nil)

		service := warehouse.NewService(repo)
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
		repo.AssertExpectations(t)
	})
	t.Run("it should return a 404 not found error when the warehouse does not exist", func(t *testing.T) {
		//Arrange
		server := gin.New()

		expectedStatus := 404
		expectedBody := `{"message": "Not found"}`
		warehouseID := 1

		repo := &warehouse.RepositoryMock{}
		repo.On("Get", mock.Anything, warehouseID).Return(domain.Warehouse{}, warehouse.ErrNotFound)

		service := warehouse.NewService(repo)
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
		repo.AssertExpectations(t)
	})
}
