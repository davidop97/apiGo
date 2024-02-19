package handler

// package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	// "github.com/davidop97/apiGo/cmd/server/handler"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/locality"
)

// TestGetAllLocalityIntegration function from Locality domain.
// Test the GetAll function of the handler in integration with the service and a mock of repository.
// It uses the real handler and service with a mock of the repository (the handler and service are not mocked).
// Two cases of test:
// 1. Success: Get all the localities available.
// 2. Error: Get an error when try to get all the localities because an internal server error occurs.
func TestGetAllLocalityIntegration(t *testing.T) {
	// Edge case: integration_get_all_localities_ok.
	// Summary: Get all the localities available.
	t.Run("should return all the localities available and a 200 OK Status code", func(t *testing.T) {
		// Arrange
		expectedStatusCode := http.StatusOK

		//Expected body in JSON format.
		expectedBody := `{"data":[{"id":1,"postal_code":5700,"locality_name":"San Luis","province_name":"San Luis","country_name":"Argentina"},{"id":2,"postal_code":5000,"locality_name":"Cordoba","province_name":"Cordoba","country_name":"Argentina"}]}`

		//Config the mock of the repository.
		mockRepo := locality.NewMockRepository()
		mockRepo.On("GetAll", mock.Anything).Return([]domain.Locality{
			{ID: 1, PostalCode: 5700, LocalityName: "San Luis", ProvinceName: "San Luis", CountryName: "Argentina"},
			{ID: 2, PostalCode: 5000, LocalityName: "Cordoba", ProvinceName: "Cordoba", CountryName: "Argentina"},
		}, nil)

		//Config the service with the mock of the repository.
		service := locality.NewService(mockRepo)

		//Config the handler with the service.
		localityHandler := NewLocality(service)
		// localityHandler := handler.NewLocality(service)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.GET(route, localityHandler.GetAll())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/localities", nil)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (200 OK in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the list of localities was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockRepo.AssertExpectations(t)
	})

	//Edge case: integration_get_all_localities_error.
	//Summary: Get an error when try to get all the localities because an internal server error occurs.
	t.Run("should return an error and a 500 Internal Server Error Status code", func(t *testing.T) {
		// Arrange
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error
		// The body of the response in JSON format.
		expectedBody := `{
			"code":"internal_server_error", 
			"message":"internal server error"
		}`

		expectedError := errors.New("internal server error")

		//Config the mock of the repository.
		mockRepo := locality.NewMockRepository()
		mockRepo.On("GetAll", mock.Anything).Return([]domain.Locality{}, expectedError)

		//Config the service with the mock of the repository.
		service := locality.NewService(mockRepo)

		//Config the handler with the service.
		localityHandler := NewLocality(service)
		// localityHandler := handler.NewLocality(service)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.GET(route, localityHandler.GetAll())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/localities", nil)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the list of localities was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockRepo.AssertExpectations(t)
	})
}
