package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/section"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ReadSection(t *testing.T) {
	// Test case: Get all sections
	// This test verifies if the handler correctly returns a list of all sections
	t.Run("it should return all sections", func(t *testing.T) {
		// Arrange

		// - expectedSections simulates the data expected from the service layer
		expectedSections := []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 1,
				MinimumTemperature: 1,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
			{
				ID:                 2,
				SectionNumber:      2,
				CurrentTemperature: 2,
				MinimumTemperature: 2,
				CurrentCapacity:    2,
				MinimumCapacity:    2,
				MaximumCapacity:    2,
				WarehouseID:        2,
				ProductTypeID:      2,
			},
		}

		// - expected status code and response body for the API endpoint
		expectedStatusCode := http.StatusOK
		expectedBody := `{"data":[{"id":1,"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1},
								  {"id":2,"section_number":2,"current_temperature":2,"minimum_temperature":2,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}]}`

		// - mock service layer
		service := &section.ServiceMock{}
		service.On("GetAll", mock.Anything).Return(expectedSections, nil)
		// - instantiate handler with the mocked service
		handler := NewSection(service)
		// - set up router and adding handler function to the route
		r := gin.New()
		route := "/api/v1/sections"
		r.GET(route, handler.GetAll())
		// - create http GET request
		request, _ := http.NewRequest("GET", route, nil)
		// - set up response recorder to capture handler's response
		response := httptest.NewRecorder()

		// Act
		// - serve request and record response
		r.ServeHTTP(response, request)

		// Assert
		// - check if the obtained status code matches expected
		assert.Equal(t, expectedStatusCode, response.Code)
		// - check if the obtained response body macthes expected
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - verify if expectations set on the service mock were met
		service.AssertExpectations(t)
	})

	// Test case: Get a section
	// This testcase verifies if the handler returns the section corresponding to the id provided by the user
	t.Run("it should return a section corresponding to the given id", func(t *testing.T) {
		// Arrange
		// - id simulates the id provided by the user
		id := 1
		// - expectedSection simulates the section returned by the service layer
		expectedSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		// - define expected status code and expected response body
		expectedStatusCode := http.StatusOK
		expectedBody := `{"data":{"id":1,"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}}`
		// - create service mock
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(expectedSection, nil)
		// - instantiate handler with the mocked service
		handler := NewSection(service)
		// - create router and register handler
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.GET(route, handler.Get())
		// - create http request and a response recorder
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - serve http GET request and record response
		r.ServeHTTP(response, request)

		// Assert
		// - check if obtained status code matches expected
		assert.Equal(t, expectedStatusCode, response.Code)
		// - check if obtained response body matches expected
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - verify if expectations on the service mock were met
		service.AssertExpectations(t)
	})

	// Test case: ID doesn't exists
	// This test checks if handler returns an error when id provided by user doesn't exists
	t.Run("it should return an error if the given id doesn't exists", func(t *testing.T) {
		// Arrange
		// - id simulates id provided by user
		id := 1
		// - expectedSection simulates service layer return value
		expectedSection := domain.Section{}
		// - expectedStatusCode simulates error returned by service layer
		expectedStatusCode := http.StatusNotFound
		// - declare expected body response
		expectedBody := `{"message":"section not found"}`
		// - create service mock
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(expectedSection, section.ErrNotFound)
		// - instantiate handler with the mocked service
		handler := NewSection(service)
		// - create router and register handler
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.GET(route, handler.Get())
		// - create http GET request and a response recorder
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - serve http request and record response
		r.ServeHTTP(response, request)

		// Assert
		// - check if obtained status code matches expected
		assert.Equal(t, expectedStatusCode, response.Code)
		// - check if obtained response body matches expected
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - verify if expectations on the service mock were met
		service.AssertExpectations(t)
	})

	// Test case: ID is invalid
	// This test checks if handler sends correct error message when ID sent by the user is invalid
	t.Run("it should return an error if the given id is not valid", func(t *testing.T) {
		// Arrange
		// - simulate an invalid id
		// -- id should be an integer, not a string
		id := "invalid id"
		// - declare status code
		expectedStatusCode := http.StatusBadRequest
		// - declare reponse body handler is expected to return
		expectedBody := `{"message":"bad ID"}`
		// - creater a service mock
		service := &section.ServiceMock{}
		// - instantiate handler with the mocked service
		handler := NewSection(service)
		// - create router and register handler
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.GET(route, handler.Get())
		// - create http request and response recorder
		url := fmt.Sprintf("/api/v1/sections/%s", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - serve http GET request and record repsonse
		r.ServeHTTP(response, request)

		// Assert
		// - check if obtained status code matches expected
		assert.Equal(t, expectedStatusCode, response.Code)
		// - check if obtained response body matches expected
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - verufy if expectations on the service mock were met
		service.AssertExpectations(t)
	})

	// Test case: Service is unable to connect to the database when calling GetAll
	// This test checks if the handler sends the correct error message when the service
	// is unable to retrieve sections due to a database connection issue.
	t.Run("it should return an error if service is unable to connect to the database when calling GetAll", func(t *testing.T) {
		// Arrange
		// - Preparing an empty section list to simulate no data returned from the service.
		emptySectionList := []domain.Section{}

		// - Specifying the expected HTTP status code for an internal server error.
		expectedStatusCode := http.StatusInternalServerError

		// - Creating an error to simulate a database connection issue.
		err := errors.New("some internal error")

		// - Declaring the expected response body the handler should return.
		expectedBody := `{"message":"internal error"}`

		// - Creating a mock of the service layer.
		service := &section.ServiceMock{}

		// - Setting up the mock response for the GetAll method to return an error.
		service.On("GetAll", mock.Anything).Return(emptySectionList, err)

		// - Instantiating the handler with the mocked service.
		handler := NewSection(service)

		// - Creating a new Gin router and registering the handler for the route.
		r := gin.New()
		route := "/api/v1/sections/"
		r.GET(route, handler.GetAll())

		// - Creating a new HTTP GET request and a response recorder.
		request, _ := http.NewRequest("GET", route, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serving the HTTP GET request and recording the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Checking if the obtained status code matches the expected code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Checking if the obtained response body matches the expected JSON string.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verifying that all expectations set on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Service is unable to connect to the database when calling Get
	// This test verifies that the handler correctly handles an internal error,
	// such as a failure to connect to the database, when attempting to retrieve a specific section by ID.
	t.Run("it should return an error if service is unable to connect to the database when calling Get", func(t *testing.T) {
		// Arrange
		// - Preparing test data and environment.
		id := 1                             // Simulating an ID for the section to retrieve.
		expectedSection := domain.Section{} // Creating an empty section to simulate no data found.

		// - Specifying the expected HTTP status code for an internal server error.
		expectedStatusCode := http.StatusInternalServerError

		// - Declaring the expected response body for an internal server error.
		expectedBody := `{"message":"internal error"}`

		// - Simulating an internal error, such as a database connection issue.
		err := errors.New("some internal error")

		// - Creating a mock of the service layer.
		service := &section.ServiceMock{}

		// - Setting up the mock response for the Get method to return an error.
		service.On("Get", mock.Anything, id).Return(expectedSection, err)

		// - Instantiating the handler with the mocked service.
		handler := NewSection(service)

		// - Creating a new Gin router and registering the handler for the route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.GET(route, handler.Get())

		// - Creating a new HTTP GET request for the specific section ID and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serving the HTTP GET request and recording the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Checking if the obtained status code matches the expected code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Checking if the obtained response body matches the expected JSON string.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verifying that all expectations set on the service mock were met.
		service.AssertExpectations(t)
	})

}

func TestHandler_CreateSection(t *testing.T) {
	// Test case: Create a section and return its ID
	// This test checks if the handler correctly creates a section and returns the new section's ID.
	t.Run("it should create a section and return its id", func(t *testing.T) {
		// Arrange
		// - simulate a section creation request
		id := 1
		sectionRequest := domain.Section{
			ID:                 0, // ID is zero as it's a new section.
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		// - request body as JSON
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`
		// - declare expected status code and response body
		expectedStatusCode := http.StatusCreated
		expectedBody := `{"data":{"id":1,"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}}`
		// - create a service mock
		service := &section.ServiceMock{}
		service.On("Save", mock.Anything, sectionRequest).Return(id, nil)
		// - instantiate handler with mocked service
		handler := NewSection(service)
		// - create router and register handler
		r := gin.New()
		route := "/api/v1/sections"
		r.POST(route, handler.Create())
		// - create http request and response recorder
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - serve http POST request and record response
		r.ServeHTTP(response, request)

		// Assert
		// - check if obtained status code matches expected
		assert.Equal(t, expectedStatusCode, response.Code)
		// - check if obtained response body matches expected
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - verify if expectations on the service mock were met
		service.AssertExpectations(t)
	})

	// Test case: Section number already exists
	// This test checks if the handler correctly responds with an error when trying to create a section with a duplicate section number.
	t.Run("it should return an error if the section number already exists", func(t *testing.T) {
		// Arrange
		// - Simulate a section creation request with a section number that already exists.
		sectionRequest := domain.Section{
			ID:                 0, // ID is zero as it's a new section.
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		// - Request body as JSON
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`
		// - Declare expected status code and response body for a conflict error
		expectedStatusCode := http.StatusConflict
		expectedBody := `{"message": "duplicate section number"}`
		// - Create an error to simulate a duplicate section number scenario
		err := section.ErrDuplicateSectNumber
		// - Create a service mock
		service := &section.ServiceMock{}
		service.On("Save", mock.Anything, sectionRequest).Return(0, err)
		// - Instantiate handler with mocked service
		handler := NewSection(service)
		// - Create router and register handler
		r := gin.New()
		route := "/api/v1/sections"
		r.POST(route, handler.Create())
		// - Create HTTP request and response recorder
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve HTTP POST request and record response
		r.ServeHTTP(response, request)

		// Assert
		// - Check if obtained status code matches expected
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check if obtained response body matches expected
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Verify if expectations on the service mock were met
		service.AssertExpectations(t)
	})

	// Test case: A section field is missing
	// This test verifies if the handler correctly responds with an error when a required field in the section creation request body is missing.
	t.Run("it should return an error if a section field is missing", func(t *testing.T) {
		// Arrange
		// - Simulate a create request with a missing 'section_number' field.
		bodyRequest := `{"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`

		// - Declare expected status code and response body for a request with a missing field.
		expectedStatusCode := http.StatusUnprocessableEntity
		expectedBody := `{"message": "field section_number is missing"}`

		// - Create a service mock.
		service := &section.ServiceMock{}

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the POST route.
		r := gin.New()
		route := "/api/v1/sections"
		r.POST(route, handler.Create())

		// - Create HTTP POST request with the missing field and a response recorder.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP POST request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected code for unprocessable entity.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates the missing field error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Body request JSON has syntax errors
	// This test checks if the handler correctly identifies and responds with an error when the JSON syntax in the request body is invalid.
	t.Run("it should return an error if the body request json has syntax errors", func(t *testing.T) {
		// Arrange
		// - Simulate a request with a malformed JSON body.
		bodyRequest := `{badjson}`

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: invalid JSON syntax"}`

		// - Create a service mock.
		service := &section.ServiceMock{}

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the POST route.
		r := gin.New()
		route := "/api/v1/sections"
		r.POST(route, handler.Create())

		// - Create HTTP POST request with the malformed JSON body and a response recorder.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve HTTP POST request and record response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if obtained response body correctly indicates a JSON syntax error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: A section field has a wrong data type
	// This test checks if the handler correctly identifies and responds with an error when one of the section fields in the request body has an incorrect data type.
	t.Run("it should return an error if a section field has a wrong data type", func(t *testing.T) {
		// Arrange
		// - Simulate a request where a section field has an incorrect data type.
		//   Here, 'section_number' is given a string value instead of an expected integer.
		bodyRequest := `{"section_number":"string","current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message": "bad request: type string was provided at section_number field, int was expected."}`

		// - Create a service mock
		service := &section.ServiceMock{}

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the POST route.
		r := gin.New()
		route := "/api/v1/sections"
		r.POST(route, handler.Create())

		// - Create HTTP POST request with the malformed body and a response recorder.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve HTTP POST request and record response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if obtained response body correctly indicates a data type error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Warehouse ID is negative
	// This test checks if the handler correctly identifies and responds with an error when the 'warehouse_id' field in the request body is negative.
	t.Run("it should return an error if the warehouse_id is negative", func(t *testing.T) {
		// Arrange
		// - Simulate a request with a negative 'warehouse_id'.
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":-1,"product_type_id":1}`

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message": "negative warehouse_id"}`

		// - Create a service mock.
		service := &section.ServiceMock{}

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the POST route.
		r := gin.New()
		route := "/api/v1/sections"
		r.POST(route, handler.Create())

		// - Create HTTP POST request with the invalid body and a response recorder.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve HTTP POST request and record response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if obtained response body correctly indicates a negative warehouse ID error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Product type ID is negative
	// This test checks if the handler correctly identifies and responds with an error when the 'product_type_id' field in the request body is negative.
	t.Run("it should return an error if the product_type_id is negative", func(t *testing.T) {
		// Arrange
		// - Simulate a request with a negative 'product_type_id'.
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":-1}`

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message": "negative product_type_id"}`

		// - Create a service mock.
		service := &section.ServiceMock{}

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the POST route.
		r := gin.New()
		route := "/api/v1/sections"
		r.POST(route, handler.Create())

		// - Create HTTP POST request with the invalid body and a response recorder.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve HTTP POST request and record response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if obtained response body correctly indicates a negative product type ID error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Service is unable to connect to the database
	// This test checks if the handler correctly handles an internal error, such as a failure to connect to the database, during the section creation process.
	t.Run("it should return an error if service is unable to connect to the database", func(t *testing.T) {
		// Arrange
		// - Simulating a section creation request.
		sectionRequest := domain.Section{
			ID:                 0, // ID is zero as it's a new section.
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		// - Request body as JSON for the new section.
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`

		// - Declare expected status code and response body for an internal server error.
		expectedStatusCode := http.StatusInternalServerError
		expectedBody := `{"message":"internal error"}`

		// - Simulate an internal service error, such as a database connection issue.
		err := errors.New("some internal error")

		// - Create a service mock to simulate the internal error on saving the section.
		service := &section.ServiceMock{}
		service.On("Save", mock.Anything, sectionRequest).Return(0, err)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the POST route.
		r := gin.New()
		route := "/api/v1/sections/"
		r.POST(route, handler.Create())

		// - Create HTTP POST request with the section data and a response recorder.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP POST request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected internal server error code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates an internal error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})
}

func TestHandler_DeleteSection(t *testing.T) {
	// Test case: Delete the section corresponding to the given ID
	// This test verifies if the handler correctly processes the deletion of a section based on its ID.
	t.Run("it should delete the section corresponding to the given id", func(t *testing.T) {
		// Arrange
		// - Simulate a delete request for a specific section ID.
		id := 1

		// - Declare expected status code for a successful delete operation.
		expectedStatusCode := http.StatusNoContent

		// - Create a service mock to simulate the delete operation.
		service := &section.ServiceMock{}
		service.On("Delete", mock.Anything, id).Return(nil) // Mocking the delete method without errors.

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the DELETE route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.DELETE(route, handler.Delete())

		// - Create HTTP DELETE request for the specified section ID and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		request, _ := http.NewRequest("DELETE", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP DELETE request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected success code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the response body is empty as expected for a no-content response.
		assert.Empty(t, response.Body)

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Given ID doesn't exist
	// This test checks if the handler correctly handles the scenario where a delete request is made for a non-existent section ID.
	t.Run("it should return an error if the given id doesn't exists", func(t *testing.T) {
		// Arrange
		// - Simulate a delete request for a non-existent section ID.
		id := 1

		// - Simulating the service layer error for a not found scenario.
		err := section.ErrNotFound

		// - Declare expected status code and response body for a not found error.
		expectedStatusCode := http.StatusNotFound
		expectedBody := `{"message":"section not found"}`

		// - Create a service mock to simulate the not found error response.
		service := &section.ServiceMock{}
		service.On("Delete", mock.Anything, id).Return(err)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the DELETE route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.DELETE(route, handler.Delete())

		// - Create HTTP DELETE request for the specified section ID and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		request, _ := http.NewRequest("DELETE", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP DELETE request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected not found code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates a not found error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Given ID is not valid
	// This test checks if the handler correctly responds with an error when a delete request is made with an invalid ID format.
	t.Run("it should return an error if the given id is not valid", func(t *testing.T) {
		// Arrange
		// - Simulate a delete request with an invalid ID format.
		id := "invalid id" // The ID is a string which is not a valid type.

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad id"}`

		// - Create a service mock.
		service := &section.ServiceMock{}

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the DELETE route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.DELETE(route, handler.Delete())

		// - Create HTTP DELETE request with an invalid ID and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%s", id)
		request, _ := http.NewRequest("DELETE", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP DELETE request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates a bad ID error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Service is unable to connect to the database
	// This test checks if the handler correctly handles an internal error, such as a failure to connect to the database, during a delete operation.
	t.Run("it should return an error if service is unable to connect to the database", func(t *testing.T) {
		// Arrange
		// - Simulating a delete request for a specific section ID.
		id := 1 // The ID of the section to be deleted.

		// - Declare expected status code and response body for an internal server error.
		expectedStatusCode := http.StatusInternalServerError
		expectedBody := `{"message":"internal error"}`

		// - Simulate an internal service error, such as a database connection issue.
		err := errors.New("some internal error")

		// - Create a service mock to simulate the internal error on deleting the section.
		service := &section.ServiceMock{}
		service.On("Delete", mock.Anything, id).Return(err)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the DELETE route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.DELETE(route, handler.Delete())

		// - Create HTTP DELETE request for the specified section ID and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		request, _ := http.NewRequest("DELETE", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP DELETE request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected internal server error code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates an internal error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})
}

func TestHandler_UpdateSection(t *testing.T) {
	// Test case: Update the section corresponding to the given ID
	// This test verifies if the handler correctly updates a section based on the provided ID and new section data.
	t.Run("it should update the section corresponding to the given id", func(t *testing.T) {
		// Arrange
		// - Simulate an update request for a specific section ID.
		id := 1 // The ID of the section to be updated.

		// - Prepare the original section data before the update.
		originalSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}

		// - Prepare the modified section data with updated values.
		modifiedSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 2,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		}

		// - Declare expected status code and response body after a successful update.
		expectedStatusCode := http.StatusOK
		bodyRequest := `{"section_number":1,"current_temperature":2,"minimum_temperature":2,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}`
		expectedBody := `{"data":{"id":1,"section_number":1,"current_temperature":2,"minimum_temperature":2,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}}`

		// - Create a service mock to simulate the get and update operations.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(originalSection, nil) // Mocking the 'Get' method.
		service.On("Update", mock.Anything, modifiedSection).Return(nil)  // Mocking the 'Update' method.

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with the updated section data and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected success code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body matches the expected JSON string of the updated section.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Given ID doesn't exist
	// This test checks if the handler correctly responds with an error when an update request is made for a non-existent section ID.
	t.Run("it should return an error if the given id doesn't exists", func(t *testing.T) {
		// Arrange
		// - Simulating an update request for a non-existent section ID.
		id := 1 // The ID that doesn't exist in the system.

		// - Prepare an empty section to simulate a not found scenario in the service layer.
		emptySection := domain.Section{}

		// - Simulating the service layer error for a not found scenario.
		err := section.ErrNotFound

		// - Declare expected status code and response body for a not found error.
		expectedStatusCode := http.StatusNotFound
		expectedBody := `{"message":"section not found"}`

		// - Create a service mock to simulate the not found error response.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(emptySection, err)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with the update data and a response recorder.
		bodyRequest := `{"section_number":1,"current_temperature":2,"minimum_temperature":2,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}`
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected not found code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates a not found error for the section.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Section number already exists
	// This test verifies if the handler correctly responds with an error when trying to update a section to a number that already exists.
	t.Run("it should return an error if the section number already exists", func(t *testing.T) {
		// Arrange
		// - Simulate an update request for a specific section ID.
		id := 1 // The ID of the section to be updated.

		// - Prepare the original section data before the update.
		originalSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}

		// - Prepare the modified section data with an updated section number that already exists.
		modifiedSection := domain.Section{
			ID:                 1,
			SectionNumber:      2,
			CurrentTemperature: 2,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		}

		// - Declare expected status code and response body for a conflict error.
		expectedStatusCode := http.StatusConflict
		bodyRequest := `{"id":1,"section_number":2,"current_temperature":2,"minimum_temperature":2,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}`
		expectedBody := `{"message": "duplicate section number"}`

		// - Simulate the service layer error for a duplicate section number.
		err := section.ErrDuplicateSectNumber

		// - Create a service mock to simulate the get and update operations.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(originalSection, nil) // Mocking the 'Get' method.
		service.On("Update", mock.Anything, modifiedSection).Return(err)  // Mocking the 'Update' method with a duplicate section number error.

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with the updated section data and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected conflict code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates a duplicate section number error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Given ID is not valid
	// This test checks if the handler correctly responds with an error when an update request is made with an invalid ID format.
	t.Run("it should return an error if the given id is not valid", func(t *testing.T) {
		// Arrange
		// - Simulate an update request with an invalid ID format.
		id := "invalid id" // The ID is a string that is not a valid format for an ID.

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad id"}`

		// - Prepare the body of the request. This body simulates the data that the client is attempting to update.
		bodyRequest := `{"id":1,"section_number":1,"current_temperature":2,"minimum_temperature":2,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}`

		// - Create a service mock. Note: The service layer won't be called in this test,
		//   as the request should fail at the handler level due to invalid ID format.
		service := &section.ServiceMock{}

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with an invalid ID and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%s", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		//   The handler is expected to validate the ID format before processing the request.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates a bad ID error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		//   In this case, no expectations should be set as the service should not be called.
		service.AssertExpectations(t)
	})

	// Test case: Body request JSON has syntax errors
	// This test checks if the handler correctly identifies and responds with an error when the JSON syntax in the request body is invalid.
	t.Run("it should return an error if the body request json has syntax errors", func(t *testing.T) {
		// Arrange
		// - Simulate an update request with a malformed JSON body.
		id := 1                            // The ID for the section to be updated.
		originalSection := domain.Section{ // The original section data.
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}
		bodyRequest := `{badjson}` // Malformed JSON string.

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: invalid JSON syntax"}`

		// - Create a service mock.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(originalSection, nil)

		// - Instantiate handler with mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with malformed JSON and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates an invalid JSON syntax error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: A section field has a wrong data type
	// This test verifies if the handler correctly responds with an error when one of the section fields in the request body has an incorrect data type.
	t.Run("it should return an error if a section field has a wrong data type", func(t *testing.T) {
		// Arrange
		// - Simulate an update request for a specific section ID.
		id := 1 // The ID of the section to be updated.

		// - Prepare the original section data before the update.
		originalSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}

		// - Prepare a request body with an incorrect data type for 'section_number'.
		bodyRequest := `{"section_number":"string","current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`

		// - Declare expected status code and response body for a bad request.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message": "bad request: type string was provided at section_number field, int was expected."}`

		// - Create a service mock.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(originalSection, nil)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with the incorrect data type and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates an incorrect data type error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: warehouse_id field has a negative value
	// This test verifies if the handler correctly responds with an error when the 'warehouse_id' field in the request body has a negative value.
	t.Run("it should return an error if warehouse_id field has a negative value", func(t *testing.T) {
		// Arrange
		// - Simulate an update request for a specific section ID.
		id := 1 // The ID of the section to be updated.

		// - Prepare the original section data.
		originalSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}

		// - Prepare a request body with a negative value for 'warehouse_id'.
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":-1,"product_type_id":1}`

		// - Declare expected status code and response body for a bad request due to negative warehouse_id.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message": "negative warehouse_id"}`

		// - Create a service mock.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(originalSection, nil)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with the incorrect data and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates a negative warehouse_id error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: product_type_id field has a negative value
	// This test verifies if the handler correctly responds with an error when the 'product_type_id' field in the request body has a negative value.
	t.Run("it should return an error if product_type_id field has a negative value", func(t *testing.T) {
		// Arrange
		// - Simulating an update request for a specific section ID.
		id := 1 // The ID of the section to be updated.

		// - Prepare the original section data.
		originalSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}

		// - Prepare a request body with a negative value for 'product_type_id'.
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":-1}`

		// - Declare expected status code and response body for a bad request due to negative product_type_id.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message": "negative product_type_id"}`

		// - Create a service mock.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(originalSection, nil)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with the incorrect data and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected bad request code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates a negative product_type_id error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})

	// Test case: Service is unable to connect to the database
	// This test checks if the handler correctly handles an internal error, such as a failure to connect to the database, during the section update process.
	t.Run("it should return an error if service is unable to connect to the database", func(t *testing.T) {
		// Arrange
		// - Simulate an update request for a specific section ID.
		id := 1 // The ID of the section to be updated.

		// - Prepare the original section data before the update.
		originalSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseID:        1,
			ProductTypeID:      1,
		}

		// - Prepare the modified section data with updated values.
		modifiedSection := domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 2,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		}

		// - Declare expected status code and response body for an internal server error.
		expectedStatusCode := http.StatusInternalServerError
		bodyRequest := `{"section_number":1,"current_temperature":2,"minimum_temperature":2,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}`
		expectedBody := `{"message":"internal error"}`

		// - Simulate an internal service error, such as a database connection issue.
		err := errors.New("some internal error")

		// - Create a service mock to simulate the internal error on updating the section.
		service := &section.ServiceMock{}
		service.On("Get", mock.Anything, id).Return(originalSection, nil)
		service.On("Update", mock.Anything, modifiedSection).Return(err)

		// - Instantiate handler with the mocked service.
		handler := NewSection(service)

		// - Create router and register handler for the PATCH route.
		r := gin.New()
		route := "/api/v1/sections/:id"
		r.PATCH(route, handler.Update())

		// - Create HTTP PATCH request with the updated section data and a response recorder.
		url := fmt.Sprintf("/api/v1/sections/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP PATCH request and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the obtained status code matches the expected internal server error code.
		assert.Equal(t, expectedStatusCode, response.Code)

		// - Check if the obtained response body correctly indicates an internal error.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// - Verify if expectations on the service mock were met.
		service.AssertExpectations(t)
	})
}

func TestHandler_SectionProductCount(t *testing.T) {
	t.Run("it should return the number of products per section, for all sections if an id is not provided", func(t *testing.T) {
		// Arrange
		// - Set up the expected response: a slice of ProdCountResponse objects that represent the product count for each section.
		expectedCount := []section.ProdCountResponse{
			{
				ID:            1,
				SectionNumber: 1,
				ProductCount:  1,
			},
			{
				ID:            2,
				SectionNumber: 2,
				ProductCount:  2,
			},
			{
				ID:            3,
				SectionNumber: 3,
				ProductCount:  3,
			},
		}
		// - An ID of 0 indicates that no specific section ID is provided in this test case, implying a request for all sections.
		id := 0
		// - Define the expected HTTP status code and body for the response. The body is formatted as JSON.
		expectedStatusCode := http.StatusOK
		expectedBody := `
			{
				"data": [
					{"id": 1, "section_number": 1, "product_count": 1},
					{"id": 2, "section_number": 2, "product_count": 2},
					{"id": 3, "section_number": 3, "product_count": 3}
				]
			}
		`
		// - Mock the service to return the expected product count for all sections when the ProductCount method is called without a specific ID.
		service := &section.ServiceMock{}
		service.On("ProductCount", mock.Anything, id).Return(expectedCount, nil)
		// - Instantiate the handler with the mocked service, allowing the test to focus on the handler's response to the service output.
		handler := NewSection(service)
		// - Set up the HTTP server and route for testing the handler.
		r := gin.New()
		route := "/api/v1/sections/reportProducts"
		r.GET(route, handler.ProductCount())

		request, _ := http.NewRequest("GET", route, nil)
		response := httptest.NewRecorder()

		// Act
		// - Execute the HTTP request against the mock server and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the HTTP status code in the response matches the expected status code.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check that the JSON body of the response matches the expected JSON body, confirming the correct data is returned.
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the mocked service's expectations, particularly the call to ProductCount with a non-specified ID, were met.
		service.AssertExpectations(t)
	})

	t.Run("it should return the number of products for the section corresponding to the given id", func(t *testing.T) {
		// Arrange
		// - Define the expected product count response for a specific section. This simulates the expected service output for the given section ID.
		expectedCount := []section.ProdCountResponse{
			{
				ID:            1, // The unique identifier of the section.
				SectionNumber: 1, // The numeric identifier of the section.
				ProductCount:  1, // The total number of products in the section.
			},
		}
		// - Specify the ID of the section for which the product count is being queried.
		id := 1
		// - Set the expected HTTP status code and body for the response. The body is formatted as JSON to match the expected API output.
		expectedStatusCode := http.StatusOK
		expectedBody := `
			{
				"data": [
					{"id": 1, "section_number": 1, "product_count": 1}
				]
			}
		`
		// - Mock the service to return the expected product count when the ProductCount method is called with the specific ID.
		service := &section.ServiceMock{}
		service.On("ProductCount", mock.Anything, id).Return(expectedCount, nil)
		// - Instantiate the handler with the mocked service, setting up the endpoint to be tested.
		handler := NewSection(service)
		r := gin.New()
		route := "/api/v1/sections/reportProducts"
		r.GET(route, handler.ProductCount())

		// - Construct the URL with the query parameter 'id' to simulate a client request for the specific section.
		url := fmt.Sprintf("/api/v1/sections/reportProducts?id=%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request to the handler and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the HTTP status code matches the expected status, indicating the request was processed successfully.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check that the response body matches the expected JSON structure and data, confirming the correct product count was returned.
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the service mock's expectations were met, particularly that the ProductCount method was called with the correct parameters.
		service.AssertExpectations(t)
	})

	t.Run("it should return an error if the section id doesn't exists", func(t *testing.T) {
		// Arrange
		// - Prepare an empty slice of ProdCountResponse to simulate the service response for a non-existing section.
		expectedCount := []section.ProdCountResponse{}
		// - Specify the ID of the non-existing section to test error handling.
		id := 1
		// - Define the error expected to be returned by the service for a non-existing section.
		err := section.ErrNotFound
		// - Set the expected HTTP status code and response body for scenarios where the section ID does not exist.
		expectedStatusCode := http.StatusNotFound
		expectedBody := `{"message": "section not found"}`
		// - Mock the service to return the expected error when the ProductCount method is called with the non-existing section ID.
		service := &section.ServiceMock{}
		service.On("ProductCount", mock.Anything, id).Return(expectedCount, err)
		// - Instantiate the handler with the mocked service, setting up the testable endpoint.
		handler := NewSection(service)
		r := gin.New()
		route := "/api/v1/sections/reportProducts"
		r.GET(route, handler.ProductCount())

		// - Construct the URL with the query parameter 'id' to simulate a client request that includes the section ID.
		url := fmt.Sprintf("/api/v1/sections/reportProducts?id=%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request to the handler and record the response, simulating the endpoint's behavior when the specified section ID does not exist.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the HTTP status code matches the expected status, indicating the proper handling of a non-existing section request.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check that the response body accurately reflects the expected error message, confirming that the error was communicated correctly to the client.
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the service mock's expectations, especially the invocation of ProductCount with the specified ID, were fulfilled, validating the mock's integration with the handler.
		service.AssertExpectations(t)
	})
}
