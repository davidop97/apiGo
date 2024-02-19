package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davidop97/apiGo/internal/batch"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ReadBatch(t *testing.T) {
	t.Run("it should return a handler function that returns all product batches", func(t *testing.T) {
		// Arange
		// - Set service mock response values
		batches := []domain.ProductBatch{
			{
				ID:                 1,
				BatchNumber:        1,
				CurrentQuantity:    1,
				CurrentTemperature: 1,
				DueDate:            "2023-11-10",
				InitialQuantity:    1,
				ManufacturingDate:  "2023-11-10",
				ManufacturingHour:  1,
				MinimumTemperature: 1,
				ProductID:          1,
				SectionID:          1,
			},
			{
				ID:                 2,
				BatchNumber:        2,
				CurrentQuantity:    2,
				CurrentTemperature: 2,
				DueDate:            "2023-11-10",
				InitialQuantity:    2,
				ManufacturingDate:  "2023-11-10",
				ManufacturingHour:  2,
				MinimumTemperature: 2,
				ProductID:          2,
				SectionID:          2,
			},
		}
		// - Set expected test results
		expectedStatusCode := http.StatusOK
		expectedBody := `{"data":[{"id":1,"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1},
								  {"id":2,"batch_number":2,"current_quantity":2,"current_temperature":2,"due_date":"2023-11-10","initial_quantity":2,"manufacturing_date":"2023-11-10","manufacturing_hour":2,"minimum_temperature":2,"product_id":2,"section_id":2}]}`
		// - Mock service
		service := &batch.ServiceMock{}
		service.On("GetAll", mock.Anything).Return(batches, nil)
		// - Create handler with mocked service
		handler := NewProductBatch(service)
		// - Create router with the handler
		r := gin.New()
		// - Route endpoint
		route := "/api/v1/productBatches"
		r.GET(route, handler.GetAll())
		// - Create http request and a response recorder
		request, _ := http.NewRequest("GET", route, nil)
		response := httptest.NewRecorder()

		// Act
		// - Serve http request and record it
		r.ServeHTTP(response, request)

		// Assert
		// - Check if obtained status code matches expected one
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check if obtained response body matches expected one
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Check if service mock expectations were met
		service.AssertExpectations(t)
	})
}

func TestHandler_CreateBatch(t *testing.T) {
	t.Run("it should return a handler function that creates a new service and returns it", func(t *testing.T) {
		// Arrange
		// - Define test input and mock response
		id := 1
		newBatch := domain.ProductBatch{
			ID:                 0, // ID is 0 to simulate new entry
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2023-11-10",
			InitialQuantity:    1,
			ManufacturingDate:  "2023-11-10",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductID:          1,
			SectionID:          1,
		}
		// - Mock HTTP request body
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Define expected response status code and body for successful creation
		expectedStatusCode := http.StatusCreated
		expectedBody := `{"data":{"id":1,"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}}`
		// - Set up service mock to expect a Save call with newBatch and return predefined ID
		service := &batch.ServiceMock{}
		service.On("Save", mock.Anything, newBatch).Return(id, nil)
		// - Initialize handler with the mocked service
		handler := NewProductBatch(service)
		// - Create a new Gin router and register the handler for the test route
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Create a new HTTP POST request with the mock request body
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Initialize a response recorder to capture the handler's response
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request and capture the response
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the response status code matches the expected status code
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Verify that the response body matches the expected JSON body
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the service mock's expectations (the Save call) were met
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if there are any missing fields in the request", func(t *testing.T) {
		// Arrange
		// - Define a request body with a missing "current_quantity" field to simulate a bad request
		bodyRequest := `{"batch_number":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Set the expected response for a request with missing fields
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: field current_quantity is missing"}`
		// - Initialize the service mock without specifying behavior, since it should not be called due to request validation failure
		service := &batch.ServiceMock{}
		// - Initialize handler with the service mock
		handler := NewProductBatch(service)
		// - Create a new Gin router and register the handler function for testing endpoint
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Create a new HTTP POST request with the incomplete body
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Initialize a response recorder to capture the response of the handler
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request to the handler and record the response
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the status code is as expected for a bad request
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check that the response body accurately reports the missing field
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that no expectations were set on the mock service, as the request should fail before service interaction
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if there is json syntax errors at the request", func(t *testing.T) {
		// Arrange
		// - Define a request body with a JSON syntax error to simulate a bad request.
		//   Specifically, missing ':' after "batch_number" introduces a syntax error.
		bodyRequest := `{"batch_number"1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Set the expected response for a request with invalid JSON syntax.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: invalid syntax at position 16"}`
		// - Initialize the service mock without specifying behavior, as the service should not be called due to the JSON parsing error.
		service := &batch.ServiceMock{}
		// - Initialize the handler with the mocked service.
		handler := NewProductBatch(service)
		// - Create a new Gin router and register the handler for the test route.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Create a new HTTP POST request with the malformed JSON body.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Initialize a response recorder to capture the handler's response.
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request to the handler and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the response status code matches the expected status code for a request with JSON syntax error.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check that the response body correctly identifies the JSON syntax error.
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the service mock's expectations (such as not being called) were met.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if incorrect data type is provided for a field", func(t *testing.T) {
		// Arrange
		// - Define a request body where "initial_quantity" is incorrectly set as a string instead of an int
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":"not a number","manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Set the expected response for a request with a data type mismatch
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: type string was provided at initial_quantity field, int was expected."}`
		// - Initialize the service mock without specifying behavior, as the service should not be called due to the data type validation failure
		service := &batch.ServiceMock{}
		// - Initialize the handler with the mocked service
		handler := NewProductBatch(service)
		// - Create a new Gin router and register the handler for the test route
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Create a new HTTP POST request with the body containing the data type error
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Initialize a response recorder to capture the handler's response
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request to the handler and record the response
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the status code matches the expected status code for a data type mismatch error
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Check that the response body accurately reports the data type error
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the service mock's expectations (such as not being called) were met
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if current_quantity field has a negative value", func(t *testing.T) {
		// Arrange
		// - Define a request body with a negative value for "current_quantity" to simulate an invalid input scenario.
		//   A nonnegative value is expected for this field according to business rules.
		bodyRequest := `{"batch_number":1,"current_quantity":-5,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Set the expected HTTP status code and response body for this validation failure scenario.
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: current_quantity must be equal or greater than 0"}`
		// - Initialize the service mock without defining specific behavior, as the focus is on testing request validation.
		service := &batch.ServiceMock{}
		// - Initialize the handler with the service mock to process the request.
		handler := NewProductBatch(service)
		// - Create a Gin router and register the handler for the endpoint being tested.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP request with the body containing the invalid value.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Initialize a response recorder to capture the handler's response.
		response := httptest.NewRecorder()

		// Act
		// - Dispatch the HTTP request to the handler and record the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the response status code is as expected for a request violating business logic.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Verify that the response body accurately reflects the specific validation error regarding "current_quantity".
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Confirm that the mock's expectations, such as not being called due to request validation failure, were met.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if initial_quantity field has a negative value", func(t *testing.T) {
		// Arrange
		// - Define a request body where "initial_quantity" is set to a negative value,
		//   violating the business rule that expects a nonnegative value.
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":-5,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Set the expected HTTP status code and response body for a request that fails
		//   due to providing a negative value for "initial_quantity".
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: initial_quantity must be equal or greater than 0"}`
		// - Initialize the service mock. Specific behavior is not defined since the request validation is expected to fail.
		service := &batch.ServiceMock{}
		// - Initialize the handler with the mocked service to test the endpoint.
		handler := NewProductBatch(service)
		// - Create a new Gin router and register the handler for the testing endpoint.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP request with the body that contains the negative value.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Use a response recorder to capture the handler's response to the invalid request.
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request through the router to trigger the handler function,
		//   capturing its response for validation.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the status code matches the expected code for validation failure.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Ensure the response body correctly identifies the issue with "initial_quantity".
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Confirm that the mock service's expectations are met, which in this case means
		//   it should not have been called due to the validation error.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if product_id field has a negative value", func(t *testing.T) {
		// Arrange
		// - Define a request body where "product_id" is set to a negative value,
		//   which violates the expectation for this field to be nonnegative.
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":-1,"section_id":1}`
		// - Set the expected HTTP status code and response body for a request
		//   that fails due to a negative value for "product_id".
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: product_id must be equal or greater than 0"}`
		// - Initialize the service mock without specific behavior as the focus is on
		//   testing the request validation logic rather than the service logic.
		service := &batch.ServiceMock{}
		// - Initialize the handler with the mocked service to process the request.
		handler := NewProductBatch(service)
		// - Create a Gin router and register the handler for the endpoint being tested.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP request with the body that contains the invalid "product_id".
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Use a response recorder to capture the handler's response.
		response := httptest.NewRecorder()

		// Act
		// - Dispatch the HTTP request through the router to the handler, capturing the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Ensure the response status code is as expected for a validation error.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Confirm the response body accurately reflects the validation error related to "product_id".
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Verify that the mock service's expectations are met, which, in this scenario,
		//   means it should not have been invoked due to the request validation failure.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if section_id field has a negative value", func(t *testing.T) {
		// Arrange
		// - Define a request body where "section_id" is incorrectly set to a negative value,
		//   contravening the requirement for this field to have a nonnegative value.
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":-1}`
		// - Set the expected HTTP status code and response body for the scenario
		//   where the request fails due to a negative "section_id".
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: section_id must be equal or greater than 0"}`
		// - Initialize the service mock without specifying behavior, as the focus is on
		//   testing the request validation rather than the service logic.
		service := &batch.ServiceMock{}
		// - Initialize the handler with the mocked service for processing the request.
		handler := NewProductBatch(service)
		// - Create a new Gin router and register the handler for the endpoint being tested.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP request with the body containing the invalid "section_id".
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Use a response recorder to capture the handler's response to the invalid request.
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request through the router to the handler, recording the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the response status code is as expected for a request violating the validation rule.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Ensure the response body correctly identifies the validation error with "section_id".
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Confirm that the mock service's expectations are met, which in this case means
		//   it should not have been called due to the request validation failure.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if manufacturing_hour is not within the range [0 - 23]", func(t *testing.T) {
		// Arrange
		// - Define a request body where "manufacturing_hour" is set to an out-of-range value,
		//   which violates the expectation for this field to be within the range [0 - 23].
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":25,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Set the expected HTTP status code and response body for a request
		//   that fails due to an out-of-range value for "manufacturing_hour".
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: manufacturing_hour value must be within range [0 - 23]"}`
		// - Initialize the service mock without defining specific behavior as the focus is on
		//   testing the request validation logic rather than the service logic.
		service := &batch.ServiceMock{}
		// - Initialize the handler with the mocked service to process the request.
		handler := NewProductBatch(service)
		// - Create a Gin router and register the handler for the testing endpoint.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP request with the body that contains the invalid "manufacturing_hour".
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Use a response recorder to capture the handler's response to the invalid request.
		response := httptest.NewRecorder()

		// Act
		// - Dispatch the HTTP request through the router to the handler, capturing the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the response status code matches the expected code for a validation error.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Confirm the response body accurately reflects the validation error regarding "manufacturing_hour".
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the mock service's expectations are met, which, in this scenario,
		//   means it should not have been invoked due to the request validation failure.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if due_date has a bad format", func(t *testing.T) {
		// Arrange
		// - Define a request body where "due_date" is provided in an incorrect format.
		//   The expected format is "YYYY-MM-DD", but a slash-separated format is given instead.
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023/11/10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Set the expected HTTP status code and response body for a request
		//   that fails due to a formatting issue with "due_date".
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: due_date should match format YYYY-MM-DD"}`
		// - Initialize the service mock. Specific behavior is not defined since the focus is on
		//   testing the request validation logic rather than the service logic.
		service := &batch.ServiceMock{}
		// - Initialize the handler with the mocked service for processing the request.
		handler := NewProductBatch(service)
		// - Create a Gin router and register the handler for the testing endpoint.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP request with the body that contains the incorrectly formatted "due_date".
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Use a response recorder to capture the handler's response to the invalid request.
		response := httptest.NewRecorder()

		// Act
		// - Serve the HTTP request through the router to the handler, capturing the response.
		r.ServeHTTP(response, request)

		// Assert
		// - Ensure the response status code is as expected for a request violating date format expectations.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Confirm the response body accurately describes the specific formatting error with "due_date".
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Verify that the mock service's expectations are met, which, in this case,
		//   means it should not have been called due to the request validation failure.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if manufacturing_date has a bad format", func(t *testing.T) {
		// Arrange
		// - Prepare a request body with "manufacturing_date" in an incorrect format.
		//   The expected format is "YYYY-MM-DD", but the provided format uses slashes ("/") instead of hyphens ("-").
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023/11/10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Define the expected HTTP status code and response body for a request failing
		//   due to an improperly formatted "manufacturing_date".
		expectedStatusCode := http.StatusBadRequest
		expectedBody := `{"message":"bad request: manufacturing_date should match format YYYY-MM-DD"}`
		// - Initialize the service mock. The specific behavior is not defined as the focus
		//   is on testing the request validation rather than the service's business logic.
		service := &batch.ServiceMock{}
		// - Setup the handler with the mocked service to process the request.
		handler := NewProductBatch(service)
		// - Create a new Gin router and register the handler for the endpoint being tested.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Construct the HTTP request with the body containing the incorrectly formatted "manufacturing_date".
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Utilize a response recorder to capture the handler's response to the invalid request.
		response := httptest.NewRecorder()

		// Act
		// - Process the HTTP request through the router, invoking the handler and recording its response.
		r.ServeHTTP(response, request)

		// Assert
		// - Check if the response status code aligns with the expected code for a date format validation error.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Ensure the response body accurately communicates the specific formatting issue with "manufacturing_date".
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Confirm that the mock service's expectations are satisfied, indicating that
		//   it should not have been invoked due to the preliminary request validation failure.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if batch_number already exists", func(t *testing.T) {
		// Arrange
		// - Prepare a new batch with a batch number that simulates an already existing entry.
		//   The expectation is that batch numbers are unique within the system.
		newBatch := domain.ProductBatch{
			ID:                 0,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2023-11-10",
			InitialQuantity:    1,
			ManufacturingDate:  "2023-11-10",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductID:          1,
			SectionID:          1,
		}
		// - Define the error returned when attempting to save a batch with a duplicate batch number.
		err := batch.ErrDuplicateBatchNumber
		// - Construct the request body matching the newBatch object to simulate a client submission.
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Specify the expected HTTP status code and response body for a conflict due to a duplicate batch number.
		expectedStatusCode := http.StatusConflict
		expectedBody := `{"message": "batch number must be unique, provided already exists"}`
		// - Initialize the service mock to return the predefined error when the Save method is called with a duplicate batch number.
		service := &batch.ServiceMock{}
		service.On("Save", mock.Anything, newBatch).Return(0, err) // Service returns 0, indicating no new ID was generated.
		// - Setup the handler with the mocked service to process the POST request.
		handler := NewProductBatch(service)
		// - Create a Gin router and register the handler for the endpoint being tested.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP POST request with the body containing the new batch data.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Utilize a response recorder to capture the handler's response to the duplicate batch number.
		response := httptest.NewRecorder()

		// Act
		// - Dispatch the HTTP request through the router to the handler, capturing the response for validation.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the response status code matches the expected code for a duplicate batch number conflict.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Confirm the response body accurately communicates the issue with the duplicate batch number.
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the mock service's expectations are satisfied, indicating that
		//   it correctly identified the duplicate batch number and returned the expected error.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if product_id doesn't exist", func(t *testing.T) {
		// Arrange
		// - Define a new batch that references a product_id which does not exist in the system.
		//   This simulates the situation where a batch is being created for a non-existent product.
		newBatch := domain.ProductBatch{
			ID:                 0,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2023-11-10",
			InitialQuantity:    1,
			ManufacturingDate:  "2023-11-10",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductID:          1, // Assumed non-existent for this test case.
			SectionID:          1,
		}
		// - Specify the error to be returned when a batch is attempted to be saved with a non-existent product_id.
		err := batch.ErrProductNotFound
		// - Construct the request body to simulate a client attempting to create a batch for the non-existent product.
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Define the expected HTTP status code and response body for the error scenario where the provided product_id was not found.
		expectedStatusCode := http.StatusConflict
		expectedBody := `{"message": "can't create batch, provided product id was not found"}`
		// - Initialize the service mock to return the predefined error when the Save method is called with a non-existent product_id.
		service := &batch.ServiceMock{}
		service.On("Save", mock.Anything, newBatch).Return(0, err) // Indicates no new ID generated due to error.
		// - Setup the handler with the mocked service to process the POST request.
		handler := NewProductBatch(service)
		// - Create a Gin router and register the handler for the endpoint being tested.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP POST request with the body containing the new batch data.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Utilize a response recorder to capture the handler's response to the attempted creation of a batch with a non-existent product_id.
		response := httptest.NewRecorder()

		// Act
		// - Dispatch the HTTP request through the router to the handler, capturing the response for validation.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the response status code matches the expected code for an attempt to create a batch for a non-existent product.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Confirm the response body accurately communicates the issue with the provided product_id.
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the mock service's expectations are met, indicating that
		//   it correctly identified the non-existent product_id and returned the expected error.
		service.AssertExpectations(t)
	})

	t.Run("it should return a handler function that returns an error if section_id doesn't exist", func(t *testing.T) {
		// Arrange
		// - Define a new batch that references a section_id which does not exist in the system,
		//   simulating the situation where a batch is being created for a non-existent section.
		newBatch := domain.ProductBatch{
			ID:                 0,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2023-11-10",
			InitialQuantity:    1,
			ManufacturingDate:  "2023-11-10",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductID:          1,
			SectionID:          1, // Assumed non-existent for this test case.
		}
		// - Specify the error to be returned when a batch is attempted to be saved with a non-existent section_id.
		err := batch.ErrSectionNotFound
		// - Construct the request body to simulate a client attempting to create a batch for the non-existent section.
		bodyRequest := `{"batch_number":1,"current_quantity":1,"current_temperature":1,"due_date":"2023-11-10","initial_quantity":1,"manufacturing_date":"2023-11-10","manufacturing_hour":1,"minimum_temperature":1,"product_id":1,"section_id":1}`
		// - Define the expected HTTP status code and response body for the error scenario where the provided section_id was not found.
		expectedStatusCode := http.StatusConflict
		expectedBody := `{"message": "can't create batch, provided section id was not found"}`
		// - Initialize the service mock to return the predefined error when the Save method is called with a non-existent section_id.
		service := &batch.ServiceMock{}
		service.On("Save", mock.Anything, newBatch).Return(0, err) // Indicates no new ID generated due to error.
		// - Setup the handler with the mocked service to process the POST request.
		handler := NewProductBatch(service)
		// - Create a Gin router and register the handler for the endpoint being tested.
		r := gin.New()
		route := "/api/v1/productBatches"
		r.POST(route, handler.Create())
		// - Prepare the HTTP POST request with the body containing the new batch data.
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		// - Utilize a response recorder to capture the handler's response to the attempted creation of a batch with a non-existent section_id.
		response := httptest.NewRecorder()

		// Act
		// - Dispatch the HTTP request through the router to the handler, capturing the response for validation.
		r.ServeHTTP(response, request)

		// Assert
		// - Verify that the response status code matches the expected code for an attempt to create a batch for a non-existent section.
		assert.Equal(t, expectedStatusCode, response.Code)
		// - Confirm the response body accurately communicates the issue with the provided section_id.
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - Ensure that the mock service's expectations are satisfied, indicating that
		//   it correctly identified the non-existent section_id and returned the expected error.
		service.AssertExpectations(t)
	})
}
