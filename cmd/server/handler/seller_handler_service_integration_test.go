package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/seller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetAllSellerIntegration function from Seller domain.
// Test the GetAll function of the handler in integration with the service and a mock of repository.
// It uses the real handler and service with a mock of the repository (the handler and service are not mocked).
// Two cases of test:
// 1. Success: Get all the sellers available.
// 2. Error: Get an error when try to get all the sellers because an internal server error occurs.
func TestGetAllSellerIntegration(t *testing.T) {
	// Edge case: integration_get_all_sellers_ok.
	// Summary: Get all the sellers available.
	t.Run("should return all the sellers available and a 200 OK Status code", func(t *testing.T) {
		// Arrange
		expectedStatusCode := http.StatusOK

		//Expected body in JSON format.
		expectedBody := `{
			"data": [
				{
					"id": 1,
					"cid": 10001,
					"company_name": "FreshFoods Ltd.",
					"address": "123 Green St, Veggieville",
					"telephone": "555-0101",
					"locality_id": 1
				},
				{
					"id": 2,
					"cid": 22,
					"company_name": "The company",
					"address": "calle falsa 66",
					"telephone": "44593",
					"locality_id": 3
				}
				]
			}`

		//Config the mock of the repository.
		mockRepo := seller.NewMockRepository()
		mockRepo.On("GetAll", mock.Anything).Return([]domain.Seller{
			{ID: 1, CID: 10001, CompanyName: "FreshFoods Ltd.", Address: "123 Green St, Veggieville", Telephone: "555-0101", IDLocality: 1},
			{ID: 2, CID: 22, CompanyName: "The company", Address: "calle falsa 66", Telephone: "44593", IDLocality: 3},
		}, nil)

		//Config the service with the mock of the repository.
		service := seller.NewService(mockRepo)

		//Config the handler with the service.
		sellerHandler := NewSeller(service)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller"
		router.GET(route, sellerHandler.GetAll())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/seller", nil)

		//Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		//Check if the status code is the correct one (200 OK in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		//Check if the list of sellers was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		//Check if the mock of the service was called.
		mockRepo.AssertExpectations(t)
	})

	//Edge case: integration_get_all_sellers_error.
	//Summary: Get an error when try to get all the sellers because an internal server error occurs.
	t.Run("should return an error and a 500 Internal Server Error Status code", func(t *testing.T) {
		// Arrange
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error

		//Expected body in JSON format.
		expectedBody := `{
			"code":"internal_server_error", 
			"message":"internal server error"
		}`

		expectedError := errors.New("internal server error")

		//Config the mock of the repository.
		mockRepo := seller.NewMockRepository()
		mockRepo.On("GetAll", mock.Anything).Return([]domain.Seller{}, expectedError)

		//Config the service with the mock of the repository.
		service := seller.NewService(mockRepo)

		//Config the handler with the service.
		sellerHandler := NewSeller(service)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller"
		router.GET(route, sellerHandler.GetAll())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/seller", nil)

		//Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		//Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		//Check if the error was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		//Check if the mock of the service was called.
		mockRepo.AssertExpectations(t)
	})
}
