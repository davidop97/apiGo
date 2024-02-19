package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/locality"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestHandler_GetLocalityByID function.
// Test the GetLocalityByID function in the following cases:
// - Locality found.
// - Locality not found.
// - Error in the service.
// - Error if the id is invalid.
// - Error if the id is not one or greater.
func TestHandler_GetLocalityByID(t *testing.T) {
	//Edge case: get_locality_by_id_ok.
	//Summary: Get a locality by id if exists.
	t.Run("should return a locality by id if exists", func(t *testing.T) {

		//Arrange
		LocalityID := 1
		Locality := domain.Locality{
			ID:           1,
			PostalCode:   5000,
			LocalityName: "Cordoba",
			ProvinceName: "Cordoba",
			CountryName:  "Argentina",
		}

		expectedStatusCode := http.StatusOK //200 OK.

		//The body of the response in JSON format.
		expectedBody := `{
		"data": {
			"id": 1,
			"postal_code": 5000,
			"locality_name": "Cordoba",
			"province_name": "Cordoba",
			"country_name": "Argentina"
		}
	   }`

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("GetLocalityByID", mock.Anything, LocalityID).Return(Locality, nil)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities/:id"
		router.GET(route, localityHandler.GetLocalityById())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/localities/%d", LocalityID), nil)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (200 OK in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the locality data was correctly updated.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: get_locality_by_id_not_found.
	//Summary: Get a locality by id, but the locality does not exist. Return a 404 Not Found.
	t.Run("should return a 404 Not Found because the locality does not exist", func(t *testing.T) {
		//Arrange
		LocalityID := 1500 //Locality does not exist.

		expectedStatusCode := http.StatusNotFound //404 Not Found.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"not_found", 
			"message":"locality not found"
		}`

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("GetLocalityByID", mock.Anything, LocalityID).Return(domain.Locality{}, locality.ErrLocalityNotFound)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities/:id"
		router.GET(route, localityHandler.GetLocalityById())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/localities/%d", LocalityID), nil)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (404 Not Found in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: get_locality_by_id_error.
	//Summary: Return a 500 internal server error case.
	t.Run("should return a 500 internal server error", func(t *testing.T) {
		//Arrange
		LocalityID := 1

		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"internal_server_error", 
			"message":"internal server error"
		}`

		expectedError := errors.New("internal server error")

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("GetLocalityByID", mock.Anything, LocalityID).Return(domain.Locality{}, expectedError)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities/:id"
		router.GET(route, localityHandler.GetLocalityById())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/localities/%d", LocalityID), nil)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: get_locality_by_id_invalid_id.
	//Summary: Return a 400 status code bad request error because the id is invalid.
	t.Run("should return a 400 status code because the id is invalid", func(t *testing.T) {
		//Arrange
		LocalityID := "invalid_id"

		expectedStatusCode := http.StatusBadRequest //400 Bad Request.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"bad_request", 
			"message":"invalid id"
		}`

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities/:id"
		router.GET(route, localityHandler.GetLocalityById())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/localities/%s", LocalityID), nil)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: get_locality_by_id_not_one_or_greater.
	//Summary: Return a 400 status code bad request error because the id is not one or greater.
	t.Run("should return a 400 status code because the id is not one or greater", func(t *testing.T) {
		//Arrange
		LocalityID := -1

		expectedStatusCode := http.StatusBadRequest //400 Bad Request.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"bad_request", 
			"message":"id must be 1 or greater"
		}`

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities/:id"
		router.GET(route, localityHandler.GetLocalityById())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/localities/%d", LocalityID), nil)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})
}

// TestHandler_GetAll function.
// Test the GetAllfunction in the following cases:
// - Success: return a 200 status code and the list of localities.
// - Error: return a 404 when the list of localities is empty.
// - Error: return a 500 internal server error.
func TestHandler_GetAll(t *testing.T) {
	//Edge case: get_all_localities_success.
	//Summary: Return a 200 status code and the list of localities.
	t.Run("should return a 200 status code and the list of localities", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusOK //200 OK.

		//The body of the response in JSON format.
		expectedBody := `{
				"data": [
					{
						"id": 1,
						"postal_code": 5700,
						"locality_name": "San Luis",
						"province_name": "San Luis",
						"country_name": "Argentina"
					},
					{
						"id": 2,
						"postal_code": 5000,
						"locality_name": "Cordoba",
						"province_name": "Cordoba",
						"country_name": "Argentina"
					}
				]
		   }`

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("GetAll", mock.Anything).Return([]domain.Locality{
			{
				ID:           1,
				PostalCode:   5700,
				LocalityName: "San Luis",
				ProvinceName: "San Luis",
				CountryName:  "Argentina",
			},
			{
				ID:           2,
				PostalCode:   5000,
				LocalityName: "Cordoba",
				ProvinceName: "Cordoba",
				CountryName:  "Argentina",
			},
		}, nil)

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
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: get_all_localities_not_found.
	//Summary: Return a 404 status code not found error because the list of localities is empty.
	t.Run("should return a 404 status code because the list of localities is empty", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusNotFound //404 Not Found.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"not_found", 
			"message":"Localities not found"
		}`

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("GetAll", mock.Anything).Return([]domain.Locality{}, locality.ErrNoRows)

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
		// Check if the status code is the correct one (404 Not Found in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: get_all_localities_internal_server_errpr.
	//Summary: Return a 500 status code internal server error.
	t.Run("should return a a 500 status code internal server error", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"internal_server_error", 
			"message":"internal server error"
		}`

		expectedError := errors.New("internal server error")

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("GetAll", mock.Anything).Return([]domain.Locality{}, expectedError)

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

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})
}

// TestHandler_CreateLocality function.
// Test the Create function from Locality in the following cases:
// - Success: return a 201 status code and the locality created.
// - Error: return a 400 conflict error because the locality already exists.
// - Error: return a 500 internal server error.
func TestHandler_CreateLocality(t *testing.T) {
	//Edge case: create_locality_ok.
	//Summary: Return a 201 status code and the locality created.
	t.Run("should return a 201 status code and the locality created", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusCreated //201 Created.

		expectedLocalityID := 1

		//The body of the response in JSON format.
		expectedBody := `{
			"data": {
				"id": 1,
				"postal_code": 5700,
				"locality_name": "San Luis",
				"province_name": "San Luis",
				"country_name": "Argentina"
			}
		}`

		localityRequest := domain.Locality{
			PostalCode:   5700,
			LocalityName: "San Luis",
			ProvinceName: "San Luis",
			CountryName:  "Argentina",
		}

		localityRequestJSON := `{
			"postal_code": 5700,
			"locality_name": "San Luis",
			"province_name": "San Luis",
			"country_name": "Argentina"
		}`

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(localityRequestJSON)

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("Save", mock.Anything, localityRequest).Return(expectedLocalityID, nil)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.POST(route, localityHandler.Create())

		//Create the request
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/localities", bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (201 Created in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the locality was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: create_locality_error_internal_server_error.
	//Summary: can not add the locality because of an internal server error.
	t.Run("should return a 500 status code internal server error", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"internal_server_error", 
			"message":"internal server error"
		}`

		expectedLocalityID := 0 //The locality ID is 0 because the locality was not created.
		expectedError := errors.New("internal server error")

		localityRequest := domain.Locality{
			PostalCode:   5700,
			LocalityName: "San Luis",
			ProvinceName: "San Luis",
			CountryName:  "Argentina",
		}

		localityRequestJSON := `{
			"postal_code": 5700,
			"locality_name": "San Luis",
			"province_name": "San Luis",
			"country_name": "Argentina"
		}`

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(localityRequestJSON)

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("Save", mock.Anything, localityRequest).Return(expectedLocalityID, expectedError)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.POST(route, localityHandler.Create())

		//Create the request
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/localities", bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: create_locality_error_locality_already_exists.
	//Summary: can not add the locality because locality already exists.
	t.Run("should return a 409 status code conflict because the locality already exists", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusConflict //409 Conflict.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"conflict", 
			"message":"locality already exists"
		}`

		expectedLocalityID := 0 //The locality ID is 0 because the locality was not created.
		expectedError := locality.ErrLocalityAlreadyExists

		localityRequest := domain.Locality{
			PostalCode:   5700,
			LocalityName: "San Luis",
			ProvinceName: "San Luis",
			CountryName:  "Argentina",
		}

		localityRequestJSON := `{
			"postal_code": 5700,
			"locality_name": "San Luis",
			"province_name": "San Luis",
			"country_name": "Argentina"
		}`

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(localityRequestJSON)

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Prepare the mock service with the expected inputs and outputs.
		mockLocalityService.On("Save", mock.Anything, localityRequest).Return(expectedLocalityID, expectedError)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.POST(route, localityHandler.Create())

		//Create the request
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/localities", bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (409 Conflict in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: create_locality_error_invalid_json.
	//Summary: can not add the locality because the request body has a wrong JSON format.
	t.Run("should return a 422 status code Unprocessable Entity because the request body has a wrong JSON format", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusUnprocessableEntity //422 Unprocessable Entity.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"unprocessable_entity", 
			"message":"invalid json"
		}`

		//The body of the request in JSON with wrong format.
		//In this case, it has an extra comma.
		localityRequestJSON := `{
			"postal_code": 5700,
			"locality_name": "San Luis",
			"province_name": "San Luis",
			"country_name": "Argentina",,
		}`

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(localityRequestJSON)

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.POST(route, localityHandler.Create())

		//Create the request
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/localities", bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: create_locality_error_invalid_postal_code.
	//Summary: can not add the locality because the request body has a wrong postal_code.
	t.Run("should return a 422 status code Unprocessable Entity because the request body has a wrong postal_code", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusUnprocessableEntity //422 Unprocessable Entity.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"unprocessable_entity", 
			"message":"Invalid or missing 'postal_code'"
		}`

		//The body of the request in JSON with wrong postal_code.
		//In this case, the postal_code is a negative number.
		localityRequestJSON := `{
			"postal_code": -5700,
			"locality_name": "San Luis",
			"province_name": "San Luis",
			"country_name": "Argentina"
		}`

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(localityRequestJSON)

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.POST(route, localityHandler.Create())

		//Create the request
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/localities", bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

	//Edge case: create_locality_error_missing_field.
	//Summary: can not add the locality because the request body has a missing field.
	t.Run("should return a 422 status code Unprocessable Entity because the request body has a missing field", func(t *testing.T) {
		//Arrange
		expectedStatusCode := http.StatusUnprocessableEntity //422 Unprocessable Entity.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"unprocessable_entity", 
			"message":"Invalid or missing 'locality_name'"
		}`

		//The body of the request in JSON format with a missing field.
		//In this case, the locality_name is missing.
		localityRequestJSON := `{
			"postal_code": 5700,
			"province_name": "San Luis",
			"country_name": "Argentina"
		}`

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(localityRequestJSON)

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities"
		router.POST(route, localityHandler.Create())

		//Create the request
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/localities", bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})

}

// TestHandler_GetReportSellers function
// This function test the function GetReportSellers which obtain a report of the number of sellers per locality.
// the result is a list of localities with the number of sellers. an example of expected result is:
// "data": [
//
//	{
//		"locality_id": 1,
//		"locality_name": "San Luis",
//		"postal_code": 5700,
//		"sellers_count": 5
//	},
//	{
//		"locality_id": 2,
//		"locality_name": "Cordoba",
//		"postal_code": 5000,
//		"sellers_count": 3
//	}
//
// ]
// Test the GetReportSellers (in base a Locality) function of the handler in the following cases:
// - Success: return a 200 status code and the report of sellers.
// - Error: return a 500 status code and an error message because an internal error occurred.
func TestHandler_GetReportSellers(t *testing.T) {
	//Edge case: get_report_sellers_success.
	//Summary: return a 200 status code and the report of sellers.
	t.Run("should return a 200 status code and the report of sellers", func(t *testing.T) {
		//Arrange

		expectedStatusCode := http.StatusOK //200 OK.

		expectedLocality := domain.Locality{
			ID:           1,
			PostalCode:   5000,
			LocalityName: "Cordoba",
			ProvinceName: "Cordoba",
			CountryName:  "Argentina",
		}

		//The body of the response in JSON format.
		expectedBodyJson := `{
			"data": [
				{
					"locality_id": 1,
					"locality_name": "San Luis",
					"postal_code": 5700,
					"sellers_count": 5
				},
				{
					"locality_id": 2,
					"locality_name": "Cordoba",
					"postal_code": 5000,
					"sellers_count": 3
				}
			]
		}`

		expectedBody := []domain.ReportSellers{
			{
				Locality_id:   1,
				Locality_name: "San Luis",
				Postal_code:   5700,
				Sellers_count: 5,
			},
			{
				Locality_id:   2,
				Locality_name: "Cordoba",
				Postal_code:   5000,
				Sellers_count: 3,
			},
		}

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		mockLocalityService.On("GetLocalityByID", mock.Anything, 1).Return(expectedLocality, nil)
		mockLocalityService.On("GetReportSellers", mock.Anything, 1).Return(expectedBody, nil)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities/reportSellers/:id"
		router.GET(route, localityHandler.GetReportSellers())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/localities/reportSellers/1", nil)
		query := request.URL.Query()
		query.Add("id", strconv.Itoa(1))
		request.URL.RawQuery = query.Encode()

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (200 OK in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the report of sellers was correctly returned.
		assert.JSONEq(t, expectedBodyJson, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)

	})

	//Edge case: get_report_sellers_locality_internal_server_error.
	//Summary: return a 500 status code internal server error.
	t.Run("should return a 500 status code internal server error", func(t *testing.T) {
		//Arrange

		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error.

		//The body of the response in JSON format.
		expectedBody := `{
			"code":"internal_server_error", 
			"message":"internal server error"
		}`

		mockLocalityService := &locality.ServiceMock{}

		//Create a new handler with the mock of the service.
		localityHandler := NewLocality(mockLocalityService)

		// Configure the mock to return an error when GetLocalityByID is called.
		mockLocalityService.On("GetLocalityByID", mock.Anything, 1).Return(domain.Locality{}, errors.New("Locality not found"))

		// Configure the mock to return an empty slice and an error when GetReportSellers is called.
		mockLocalityService.On("GetReportSellers", mock.Anything, 1).Return([]domain.ReportSellers{}, errors.New("Sellers not found"))

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/localities/reportSellers/:id"
		router.GET(route, localityHandler.GetReportSellers())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/localities/reportSellers/1", nil)
		query := request.URL.Query()
		query.Add("id", strconv.Itoa(1))
		request.URL.RawQuery = query.Encode()

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message was correctly returned.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockLocalityService.AssertExpectations(t)
	})
}
