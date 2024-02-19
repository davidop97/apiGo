package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidop97/apiGo/internal/buyer"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests for Buyer handler

// TestHandler_ReadBuyer
// Test all cases for Get functions
// Test methods: Get and GetAll
func TestHandler_ReadBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: find_all
	// DESCRIPTION: When the request is successful the backend will return a list of all existing buyers.
	t.Run("it should returns all buyers", func(t *testing.T) {
		// arrange

		// buyers that service layer returns
		buyers := []domain.Buyer{
			{
				ID:           1,
				FirstName:    "Jane",
				LastName:     "Doe",
				CardNumberID: "123456789",
			},
			{
				ID:           2,
				FirstName:    "John",
				LastName:     "Smith",
				CardNumberID: "987654321",
			},
		}
		// prepare expected results
		expectedBuyers := `{"data":[{"id":1,"card_number_id":"123456789","first_name":"Jane","last_name":"Doe"},
							{"id":2,"card_number_id":"987654321","first_name":"John","last_name":"Smith"}]}`
		expectedStatusCode := http.StatusOK

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("GetAll", mock.Anything).Return(buyers, nil)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.GET(route, handler.GetAll())

		// create the request
		request, _ := http.NewRequest("GET", route, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedBuyers, response.Body.String())
		serviceMock.AssertExpectations(t)

	})
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: find_by_id_existent
	// DESCRIPTION: When the request is successful the backend will return the requested buyer information.
	t.Run("it should return a buyer by id", func(t *testing.T) {
		// buyer with the id that service returns
		expectedID := "1"
		buyerFromService := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: "123456789",
		}
		// prepare expected results
		expectedBuyers := `{"data":{"id":1,"card_number_id":"123456789","first_name":"Jane","last_name":"Doe"}}`
		expectedStatusCode := http.StatusOK

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, 1).Return(buyerFromService, nil)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.GET(route, handler.Get())

		// create the request
		request, _ := http.NewRequest("GET", "/api/v1/buyers/"+expectedID, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedBuyers, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: READ
	// EDGE CASE: find_by_id_non_existent
	// DESCRIPTION: When the buyer does not exist, a 404 code will be returned.
	t.Run("it should return 404 status code when buyer id is not found", func(t *testing.T) {
		expectedID := "1"
		//expectedError
		expectedMessageError := `{"error":"Buyer not found"}`
		buyerFromService := domain.Buyer{}
		// prepare expected results
		expectedStatusCode := http.StatusNotFound

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, mock.Anything).Return(buyerFromService, buyer.ErrNotFound)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.GET(route, handler.Get())

		// create the request
		request, _ := http.NewRequest("GET", "/api/v1/buyers/"+expectedID, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: READ
	// EDGE CASE: error_calling_get_invalid_id
	// DESCRIPTION: If the element searched by invalid id, it returns error.
	t.Run("it should return 400 status code when buyer id is invalid", func(t *testing.T) {
		//expectedID := "1"
		//expectedError
		expectedMessageError := `{"error":"Invalid id"}`
		//buyerFromService := domain.Buyer{}
		// prepare expected results
		expectedStatusCode := http.StatusBadRequest

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// service mock with the expected value, before launch a errors with invalid id

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.GET(route, handler.Get())

		// create the request
		request, _ := http.NewRequest("GET", "/api/v1/buyers", nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: READ
	// EDGE CASE: error_calling_get
	// DESCRIPTION: If any errors appear in the server, it returns error.
	t.Run("it should return an error when internal server error appears", func(t *testing.T) {
		expectedID := "1"
		//expectedError
		expectedMessageError := `{"error":" Internal error"}`
		buyerFromService := domain.Buyer{}
		// prepare expected results
		expectedStatusCode := http.StatusInternalServerError

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, mock.Anything).Return(buyerFromService, errors.New("some errors"))

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.GET(route, handler.Get())

		// create the request
		request, _ := http.NewRequest("GET", "/api/v1/buyers/"+expectedID, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: READ
	// EDGE CASE: error_calling_getall
	// DESCRIPTION: If any errors appear in the server, it returns error.
	t.Run("it should return an error when internal server error appears calling GetAll function", func(t *testing.T) {
		//expectedError
		expectedMessageError := "Internal server error"
		buyerFromService := []domain.Buyer{}
		// prepare expected results
		expectedStatusCode := http.StatusInternalServerError

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("GetAll", mock.Anything).Return(buyerFromService, errors.New("some errors"))

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.GET(route, handler.GetAll())

		// create the request
		request, _ := http.NewRequest("GET", route, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		// Check if the response body contains the expected error message
		assert.Contains(t, response.Body.String(), expectedMessageError)
		serviceMock.AssertExpectations(t)
	})

}

// TestService_DeleteBuyer
// Test all cases for Delete functions in handler layer
// Tested methods: Delete
func TestService_DeleteBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete_ok
	// DESCRIPTION: A code 204 will be returned when the deletion is successful.
	t.Run("it should returns 204 when the deletion is successful", func(t *testing.T) {
		// buyer with the id that service returns
		expectedID := "1"

		// prepare expected results
		expectedStatusCode := http.StatusNoContent

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Delete", mock.Anything, 1).Return(nil)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.DELETE(route, handler.Delete())

		// create the request
		request, _ := http.NewRequest("DELETE", "/api/v1/buyers/"+expectedID, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, "", response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete_non_existent
	// DESCRIPTION: When the buyer does not exist, a 404 code will be returned.
	t.Run("it should returns 404 when the buyer does not exist", func(t *testing.T) {
		expectedID := "1"
		//expectedError
		expectedMessageError := `{"error":"Buyer not found"}`
		// prepare expected results
		expectedStatusCode := http.StatusNotFound

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Delete", mock.Anything, 1).Return(buyer.ErrNotFound)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.DELETE(route, handler.Delete())

		// create the request
		request, _ := http.NewRequest("DELETE", "/api/v1/buyers/"+expectedID, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: error_calling_delete_invalid_id
	// DESCRIPTION: If the element searched by invalid id, it returns error.
	t.Run("it should return 404 status code when buyer id is invalid", func(t *testing.T) {
		//expectedError
		expectedMessageError := `{"error":"Invalid id"}`
		//buyerFromService := domain.Buyer{}
		// prepare expected results
		expectedStatusCode := http.StatusBadRequest

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// service mock with the expected value, before launch a errors with invalid id

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.DELETE(route, handler.Delete())

		// create the request
		request, _ := http.NewRequest("DELETE", "/api/v1/buyers", nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: error_calling_delete
	// DESCRIPTION: If any errors appear in the server, it returns error.
	t.Run("it should return an error when internal server error appears calling Delete function", func(t *testing.T) {
		//expectedError
		expectedID := "1"
		expectedMessageError := `{"error":" Internal error"}`
		// prepare expected results
		expectedStatusCode := http.StatusInternalServerError

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Delete", mock.Anything, 1).Return(errors.New("some errors"))

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.DELETE(route, handler.Delete())

		// create the request
		request, _ := http.NewRequest("DELETE", "/api/v1/buyers/"+expectedID, nil)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})
}

// TestService_CreateBuyer
// Test all cases for Create functions
// Tested methods: Create
func TestService_CreateBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_ok
	// DESCRIPTION: When the data entry is successful, a 201 code will be returned along with the object entered.
	t.Run("it should return 201 code When the data entry is successful", func(t *testing.T) {
		// arrange

		expectedStatusCode := http.StatusCreated
		expectedBuyers := `{"data":{"id":1,"card_number_id":"123456789","first_name":"Jane","last_name":"Doe"}}`
		buyerToCreate := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: "123456789",
		}

		buyerJSON, _ := json.Marshal(buyerToCreate)
		// configure the expected Buyer object for the Save method

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Save", mock.Anything, mock.Anything).Return(1, nil)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSON)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedBuyers, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_fail
	// DESCRIPTION: If the JSON object does not contain the required fields, a 422 code will be returned.
	t.Run("it should return 422 if the JSON object does not contain the required fields", func(t *testing.T) {
		expectedStatusCode := http.StatusUnprocessableEntity
		expectedMessageError := `{"error":"Missing fields"}`
		// configure the Buyer object to the Save method
		buyerToCreate := domain.Buyer{
			LastName: "Doe",
		}

		buyerJSON, _ := json.Marshal(buyerToCreate)

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		//serviceMock.On("Save", mock.Anything, buyerToCreate).Return(1, buyer.ErrAlreadyExists)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSON)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_conflict
	// DESCRIPTION: If the card_number_id already exists, a 409 Conflict error is returned.
	t.Run("it should return 409 Conflict error if the card_number_id already exists", func(t *testing.T) {
		expectedStatusCode := http.StatusConflict
		expectedMessageError := `{"error":"Buyer already exists"}`
		// configure the Buyer object to the Save method
		buyerToCreate := domain.Buyer{
			ID:           0,
			CardNumberID: "123456789",
			FirstName:    "Jane",
			LastName:     "Doe",
		}

		buyerJSON, _ := json.Marshal(buyerToCreate)

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Save", mock.Anything, buyerToCreate).Return(1, buyer.ErrAlreadyExists)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSON)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_bad_request_json
	// DESCRIPTION: If json body is malformed, a 400 code status bad request is returned.
	t.Run("it should return 400 status code when body json is invalid", func(t *testing.T) {

		invalidJSON := `{"some_field": "value"`
		// prepare expected results
		expectedMessageError := `{"error":"Bad request"}`
		expectedStatusCode := http.StatusBadRequest

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// no mock to configure because there are no external functions to call
		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(invalidJSON)))
		request, _ := http.NewRequest("POST", "/api/v1/buyers", reqBodyBytes)
		request.Header.Set("Content-Type", "application/json")
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_internal_error
	// DESCRIPTION: If the card_number_id already exists, a 409 Conflict error is returned.
	t.Run("it should return 500 status code if there any error", func(t *testing.T) {
		expectedStatusCode := http.StatusInternalServerError
		expectedMessageError := `{"error":"Internal server error"}`
		// configure the Buyer object to the Save method
		buyerToCreate := domain.Buyer{
			ID:           0,
			CardNumberID: "123456789",
			FirstName:    "Jane",
			LastName:     "Doe",
		}

		buyerJSON, _ := json.Marshal(buyerToCreate)

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Save", mock.Anything, buyerToCreate).Return(1, errors.New("some errors"))

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSON)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})
}

// TestService_UpdateBuyer
// Test all cases for Update functions
// Tested methods: Update
func TestService_UpdateBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: update_ok
	// DESCRIPTION: When the data update is successful, the buyer will be returned with the updated information along with a code 200.
	t.Run("it should return 200 when the data update is successful", func(t *testing.T) {
		// buyer with the id that service returns
		buyerIDToUpdate := 1
		buyerFromService := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: "123456789",
		}
		buyertoUpdate := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: "123456789",
		}
		buyerJSONToUpdate, _ := json.Marshal(buyertoUpdate)
		// prepare expected results
		expectedBuyers := `{"data":{"id":1,"card_number_id":"123456789","first_name":"Jane","last_name":"Does"}}`
		expectedStatusCode := http.StatusOK

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, buyerIDToUpdate).Return(buyerFromService, nil)
		serviceMock.On("Update", mock.Anything, 1, buyertoUpdate, &buyerFromService).Return(nil)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.PATCH(route, handler.Update())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSONToUpdate)))
		request, _ := http.NewRequest("PATCH", "/api/v1/buyers/1", reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedBuyers, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: update_non_existent
	// DESCRIPTION: If the buyer to be updated does not exist, a 404 code will be returned.
	t.Run("it should return 404 code when the buyer to be updated does not exist", func(t *testing.T) {
		// buyer with the id that service returns
		buyerIDToUpdate := 1
		buyerFromService := domain.Buyer{
			ID:           2,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: "123456789",
		}
		buyertoUpdate := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: "123456789",
		}
		buyerJSONToUpdate, _ := json.Marshal(buyertoUpdate)
		// prepare expected results
		expectedMessageError := `{"error":"Buyer not found"}`
		expectedStatusCode := http.StatusNotFound

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, buyerIDToUpdate).Return(buyerFromService, buyer.ErrNotFound)
		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.PATCH(route, handler.Update())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSONToUpdate)))
		request, _ := http.NewRequest("PATCH", "/api/v1/buyers/1", reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: error_calling_update_invalid_id
	// DESCRIPTION: If the element update by invalid id, it returns error.
	t.Run("it should return 400 status code when buyer id is invalid", func(t *testing.T) {
		buyertoUpdate := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: "123456789",
		}
		buyerJSONToUpdate, _ := json.Marshal(buyertoUpdate)
		// prepare expected results
		expectedMessageError := `{"error":"Invalid id"}`
		expectedStatusCode := http.StatusBadRequest

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// no mock to configure because there are no external functions to call
		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers"

		// Add handler to the router
		r.PATCH(route, handler.Update())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSONToUpdate)))
		request, _ := http.NewRequest("PATCH", "/api/v1/buyers", reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: error_calling_get_in_update
	// DESCRIPTION: If any errors appear calling to GET method, it returns error.
	t.Run("it should return 500 when called to Get buyer by id and internal server error appears", func(t *testing.T) {
		// buyer with the id that service returns
		buyerIDToUpdate := 1
		buyertoUpdate := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: "123456789",
		}
		buyerJSONToUpdate, _ := json.Marshal(buyertoUpdate)
		// prepare expected results
		expectedMessageError := `{"error":"Internal server error"}`
		expectedStatusCode := http.StatusInternalServerError

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, buyerIDToUpdate).Return(domain.Buyer{}, errors.New("unexpected error in get"))
		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.PATCH(route, handler.Update())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSONToUpdate)))
		request, _ := http.NewRequest("PATCH", "/api/v1/buyers/1", reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: error_calling_update_buyer_already_exists
	// DESCRIPTION: If ErrAlreadyExists appear calling to Update method, it returns 409 code.
	t.Run("it should return 409 when the buyer to update already exists", func(t *testing.T) {
		// buyer with the id that service returns
		buyerIDToUpdate := 1
		buyerFromService := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: "123456789",
		}
		buyertoUpdate := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: "123456789",
		}
		buyerJSONToUpdate, _ := json.Marshal(buyertoUpdate)
		// prepare expected results
		expectedMessageError := `{"error":"Buyer already exists"}`
		expectedStatusCode := http.StatusConflict

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, buyerIDToUpdate).Return(buyerFromService, nil)
		serviceMock.On("Update", mock.Anything, 1, buyertoUpdate, &buyerFromService).Return(buyer.ErrAlreadyExists)
		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.PATCH(route, handler.Update())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSONToUpdate)))
		request, _ := http.NewRequest("PATCH", "/api/v1/buyers/1", reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: error_calling_update_internal_error
	// DESCRIPTION: If internal error appears calling to Update method, it returns 500 code.
	t.Run("it should return 500 when trying to update a buyer", func(t *testing.T) {
		// buyer with the id that service returns
		buyerIDToUpdate := 1
		buyerFromService := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: "123456789",
		}
		buyertoUpdate := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Does",
			CardNumberID: "123456789",
		}
		buyerJSONToUpdate, _ := json.Marshal(buyertoUpdate)
		// prepare expected results
		expectedMessageError := `{"error":"Internal server error"}`
		expectedStatusCode := http.StatusInternalServerError

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, buyerIDToUpdate).Return(buyerFromService, nil)
		serviceMock.On("Update", mock.Anything, buyerIDToUpdate, buyertoUpdate, &buyerFromService).Return(errors.New("error calling update"))
		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.PATCH(route, handler.Update())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(buyerJSONToUpdate)))
		request, _ := http.NewRequest("PATCH", "/api/v1/buyers/1", reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: UPDATE
	// EDGE CASE: create_bad_request_json
	// DESCRIPTION: If json body is malformed, a 400 status bad request error is returned.
	t.Run("it should return 400 status code when body json is invalid", func(t *testing.T) {
		// buyer with the id that service returns
		buyerIDToUpdate := 1
		buyerFromService := domain.Buyer{
			ID:           1,
			FirstName:    "Jane",
			LastName:     "Doe",
			CardNumberID: "123456789",
		}
		invalidJSON := `{"some_field": "value"`
		// prepare expected results
		expectedMessageError := `{"error":"Bad request"}`
		expectedStatusCode := http.StatusBadRequest

		// create mock of the service
		serviceMock := buyer.NewBuyerService()
		// set service mock with the expected value
		serviceMock.On("Get", mock.Anything, buyerIDToUpdate).Return(buyerFromService, nil)

		// create handler using mock service
		handler := NewBuyer(serviceMock)

		// prepare server to call endpoint Update
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/:id"

		// Add handler to the router
		r.PATCH(route, handler.Update())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(invalidJSON)))
		request, _ := http.NewRequest("PATCH", "/api/v1/buyers/1", reqBodyBytes)

		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		serviceMock.AssertExpectations(t)
	})
}
