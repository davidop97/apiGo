package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/seller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestHandler_READ function.
// Test all the Get functions from seller handler.
// - GetAllSellers
// - GetSellerByID
func TestHandler_READ(t *testing.T) {
	//Associated User Story: READ.
	//Edge case: find_all.
	//Summary: return a list of sellers.
	t.Run("should return a list of sellers", func(t *testing.T) {

		//Arrange
		//Create a struct to unmarshal the response
		//Uses the same struct as the response from the API included the 'data' word
		//in the JSON response.
		//For example:
		// {"data":
		// 	[
		// 		{
		// 			"id": 1,
		// 			"cid": 10001,
		// 			"company_name": "Empresa1.",
		// 			"address": "123 Green St",
		// 			"telephone": "111-1111",
		// 			"locality_id": 1
		// 		},
		// 		{
		// 			"id": 2,
		// 			"cid": 10002,
		// 			"company_name": "Empresa2.",
		// 			"address": "456 EverGreen St",
		// 			"telephone": "555-2222",
		// 			"locality_id": 1
		// 		}
		// 	]
		// }
		type APIResponse struct {
			Data []domain.Seller `json:"data"`
		}

		// Create APIResponse instead of []domain.Seller
		// to unmarshal the response.
		var actualSellers APIResponse

		//Expected results
		expectedStatusCode := http.StatusOK //200 OK
		expectedSellers := []domain.Seller{
			{
				ID:          1,
				CID:         1,
				CompanyName: "Seller 1",
				Address:     "Address 1",
				Telephone:   "Telephone 1",
				IDLocality:  1,
			},
			{
				ID:          2,
				CID:         2,
				CompanyName: "Seller 2",
				Address:     "Address 2",
				Telephone:   "Telephone 2",
				IDLocality:  2,
			},
		}

		// Create the mock of the service
		serviceMock := seller.NewMockService()
		// Mock the GetAllSellers function of the service to return the expectedSellers
		// and nil as error. GetAllSellers not need a parameter so a mock.Anything is used.
		serviceMock.On("GetAllSellers", mock.Anything).Return(expectedSellers, nil)
		// Create the handler with the mock of the service
		handler := NewSeller(serviceMock)

		// Create the gin router. And the route to test.
		r := gin.New()
		route := "/api/v1/seller"

		// Add the handler to the router.
		r.GET(route, handler.GetAll())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, route, nil)
		//Create the response recorder.
		response := httptest.NewRecorder()

		//Act
		r.ServeHTTP(response, request)

		// Unmarshal the response into the actualSellers
		err := json.Unmarshal([]byte(response.Body.Bytes()), &actualSellers)
		if err != nil {
			t.Fatal(err)
		}

		// Assert
		// Check if the status code is the correct one (200 OK in this case).
		assert.Equal(t, expectedStatusCode, response.Code)
		// Check if the obtained results are equals to the expected ones.
		// Use actualSellers.Data because the APIResponse struct has the 'data' word.
		assert.Equal(t, expectedSellers, actualSellers.Data)
		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_all_not_found.
	//Summary: return an empty list of sellers because there are not sellers available.
	t.Run("should return a 404 Not Found status code when no sellers are found", func(t *testing.T) {

		// Arrange
		mockSellerService := &seller.ServiceMock{}
		sellerHandler := NewSeller(mockSellerService)

		//Expected results.
		expectedError := errors.New("Sellers not found")
		expectedStatusCode := http.StatusNotFound // 404 Not Found.

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "Sellers not found"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		// Mock the service to return an empty list of sellers and an error.
		mockSellerService.On("GetAllSellers", mock.Anything).Return([]domain.Seller{}, seller.ErrNotFound)

		request, _ := http.NewRequest("GET", "/seller", nil)

		// Create the response recorder, the router and the route to test.
		response := httptest.NewRecorder()
		router := gin.New()
		router.GET("/seller", sellerHandler.GetAll())

		// Act
		router.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "Sellers not found"}
		//and the expected is a string: "Sellers not found".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the message from the actualErrorResponse.
		actualMessage := actualErrorResponse["message"]

		// Assert
		//Check if the status code is the correct one (404 Not Found in this case).
		assert.Equal(t, expectedStatusCode, response.Code)
		//Check if the message is the correct one ("Sellers not found" in this case).
		assert.Equal(t, expectedError.Error(), actualMessage)
		//Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_all_error.
	//Summary: return an empty list of sellers due an internal server error.
	t.Run("should return an internal server error", func(t *testing.T) {

		//Arrange

		//Expected results
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error
		expectedSeller := []domain.Seller{}                  //Empty seller
		expectedError := errors.New("internal server error")
		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "internal server error"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		// Create the mock of the service.
		serviceMock := seller.NewMockService()

		// Mock the GetAllSellers function of the service to return the expectedSellers (empty list)
		// and an error. GetAllSellers not need a parameter so a mock.Anything is used.
		serviceMock.On("GetAllSellers", mock.Anything).Return(expectedSeller, expectedError)

		// Create the handler with the mock of the service.
		handler := NewSeller(serviceMock)

		// Create the gin router and the route to test.
		r := gin.New()
		route := "/api/v1/seller"

		// Add the handler to the router.
		r.GET(route, handler.GetAll())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, route, nil)

		//Create the response recorder.
		response := httptest.NewRecorder()

		//Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "internal server error"}
		//and the expected is a string: "internal server error".
		_ = json.Unmarshal([]byte(response.Body.Bytes()), &actualErrorResponse)

		//Get the error message from the map.
		actualErrorMessage := actualErrorResponse["message"]

		// Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualErrorMessage)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_by_id_existent.
	//Summary: return a seller by id.
	t.Run("should return a seller using its ID", func(t *testing.T) {

		//Arrange
		// ctx := context.Background() //Not necessary to send the context in these tests.

		//Create a struct to unmarshal the response.
		//Uses the same struct as the response from the API included the 'data' word
		//in the JSON response.
		//For example:
		// {"data":
		// 		{
		// 			"id": 1,
		// 			"cid": 10001,
		// 			"company_name": "Empresa1.",
		// 			"address": "123 Green St",
		// 			"telephone": "111-1111",
		// 			"locality_id": 1
		// 		}
		// }
		type APIResponse struct {
			Data domain.Seller `json:"data"`
		}
		// Create actualSeller using the APIResponse struct.
		var actualSeller APIResponse

		//Expected results
		expectedStatusCode := http.StatusOK //200 OK
		expectedSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Seller 1",
			Address:     "Address 1",
			Telephone:   "Telephone 1",
			IDLocality:  1,
		}

		sellerID := 1

		// Create the mock of the service
		serviceMock := seller.NewMockService()
		// Mock the GetSellerByID function of the service to return the expectedSeller
		// and nil as error. GetSellerByID function needs a parameter, the ID of the seller.
		serviceMock.On("GetSellerByID", mock.Anything, sellerID).Return(expectedSeller, nil)
		// Create the handler with the mock of the service
		handler := NewSeller(serviceMock)

		// Create the gin router and the route to test.
		r := gin.New()
		route := "/api/v1/seller/:id"

		// Add the handler to the router.
		r.GET(route, handler.Get())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/%d", sellerID), nil)
		//Create the response recorder.
		response := httptest.NewRecorder()

		//Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualSeller
		err := json.Unmarshal([]byte(response.Body.Bytes()), &actualSeller) //The []byte is not necessary.
		if err != nil {
			t.Fatal(err)
		}

		// Assert
		// Check if the status code is the correct one (200 OK in this case).
		assert.Equal(t, expectedStatusCode, response.Code)
		// Check if the obtained result is the same than the expected one.
		// Use actualSeller.Data because the APIResponse struct has the 'data' word.
		assert.Equal(t, expectedSeller, actualSeller.Data)
		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})
	//Associated User Story: READ.
	//Edge case: find_by_id_non_existent.
	//Summary: return an error due to the seller does not exist.
	t.Run("should return an error not found", func(t *testing.T) {

		//Arrange

		//Expected results
		expectedStatusCode := http.StatusNotFound //404 Not Found
		expectedSeller := domain.Seller{}         //Empty seller
		// expectedError := "seller not found"
		expectedError := seller.ErrNotFound
		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "seller not found"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		sellerID := 1500 //Non existent seller ID

		// Create the mock of the service.
		serviceMock := seller.NewMockService()

		// Mock the GetSellerByID function of the service to return an empty seller
		// and error: seller not found. GetSellerByID function needs a parameter, the ID of the seller.
		serviceMock.On("GetSellerByID", mock.Anything, sellerID).Return(expectedSeller, expectedError)

		// Create the handler with the mock of the service.
		handler := NewSeller(serviceMock)

		// Create the gin router and the route to test.
		r := gin.New()
		route := "/api/v1/seller/:id"

		// Add the handler to the router.
		r.GET(route, handler.Get())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/%d", sellerID), nil)

		//Create the response recorder.
		response := httptest.NewRecorder()

		//Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "seller not found"}
		//and the expected is a string: "seller not found".
		_ = json.Unmarshal([]byte(response.Body.Bytes()), &actualErrorResponse)

		//Get the error message from the map.
		actualErrorMessage := actualErrorResponse["message"]

		// Assert
		// Check if the status code is the correct one (404 Not Found in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualErrorMessage)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})
	//Associated User Story: READ.
	//Edge case: find_by_id_error.
	//Summary: return an error due an internal server error.
	t.Run("should return an error: internal server error", func(t *testing.T) {

		//Arrange

		//Expected results
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error
		expectedSeller := domain.Seller{}                    //Empty seller

		expectedError := errors.New("internal server error")
		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "error message"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		sellerID := 1

		// Create the mock of the service.
		serviceMock := seller.NewMockService()

		// Mock the GetSellerByID function of the service to return an empty seller
		// and error. GetSellerByID function needs a parameter, the ID of the seller.
		serviceMock.On("GetSellerByID", mock.Anything, sellerID).Return(expectedSeller, expectedError)

		// Create the handler with the mock of the service.
		handler := NewSeller(serviceMock)

		// Create the gin router and the route to test.
		r := gin.New()
		route := "/api/v1/seller/:id"

		// Add the handler to the router.
		r.GET(route, handler.Get())

		//Create the request
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/seller/%d", sellerID), nil)

		//Create the response recorder.
		response := httptest.NewRecorder()

		//Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "internal server error"}
		//and the expected is a string: "internal server error".
		_ = json.Unmarshal([]byte(response.Body.Bytes()), &actualErrorResponse)

		//Get the error message from the map.
		actualErrorMessage := actualErrorResponse["message"]

		// Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualErrorMessage)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

	//Associated User Story: READ.
	//Edge case: find_by_id_invalidID.
	//Summary: return an error due to the ID is invalid.
	t.Run("should return an error: 'invalid id'", func(t *testing.T) {
		// Arrange
		expectedStatusCode := http.StatusBadRequest
		expectedError := errors.New("invalid id")
		actualErrorResponse := make(map[string]string)

		serviceMock := seller.NewMockService()

		handler := NewSeller(serviceMock)

		r := gin.New()
		r.GET("/api/v1/seller/:id", handler.Get())

		// Use an invalid id in the request.
		request, _ := http.NewRequest("GET", "/api/v1/seller/invalid", nil)
		response := httptest.NewRecorder()

		// Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "invalid id"}
		//and the expected is a string: "invalid id".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		actualMessage := actualErrorResponse["message"]

		// Assert
		// Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})
}

// TestHandler_DELETE function.
// Test the delete function from seller handler.
func TestHandler_DELETE(t *testing.T) {
	//Associated User Story: DELETE.
	//Edge case: delete_ok.
	//Summary: return a status code 204 when the seller was deleted.
	t.Run("should return a status code 204 when the seller was deleted", func(t *testing.T) {
		// Arrange
		expectedStatusCode := http.StatusNoContent //204 No Content
		sellerID := 1

		// Create the mock of the service.
		serviceMock := seller.NewMockService()

		// Mock the DeleteSeller function of the service to return nil.
		// DeleteSeller function needs a parameter, the ID of the seller.
		serviceMock.On("Delete", mock.Anything, sellerID).Return(nil)

		// Create the handler with the mock of the service.
		handler := NewSeller(serviceMock)

		// Create the gin router and the route to test.
		r := gin.New()
		route := "/api/v1/seller/:id"

		// Add the handler to the router.
		r.DELETE(route, handler.Delete())

		//Create the request
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/seller/%d", sellerID), nil)

		//Create the response recorder.
		response := httptest.NewRecorder()

		// Act
		r.ServeHTTP(response, request)

		// Assert
		// Check if the status code is the correct one (204 No Content in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the response body is empty.
		assert.Equal(t, "", response.Body.String())

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

	//Associated User Story: DELETE.
	//Edge case: delete_error.
	//Summary: return a status code 500 due an internal server error.
	t.Run("should return a status code 404 due to the seller does not exist", func(t *testing.T) {
		// Arrange
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error
		sellerID := 1

		// Create the mock of the service.
		serviceMock := seller.NewMockService()

		// Mock the DeleteSeller function of the service to return an error.
		// DeleteSeller function needs a parameter, the ID of the seller.
		serviceMock.On("Delete", mock.Anything, sellerID).Return(errors.New("internal server error"))

		// Create the handler with the mock of the service.
		handler := NewSeller(serviceMock)

		// Create the gin router and the route to test.
		r := gin.New()
		route := "/api/v1/seller/:id"

		// Add the handler to the router.
		r.DELETE(route, handler.Delete())

		//Create the request
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/seller/%d", sellerID), nil)

		//Create the response recorder.
		response := httptest.NewRecorder()

		// Act
		r.ServeHTTP(response, request)

		// Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

	//Associated User Story: DELETE.
	//Edge case: delete_non_existent.
	//Summary: return an error due the seller does not exist when trying to delete it.
	t.Run("should return an error not found when trying to delete it", func(t *testing.T) {

		//Arrange

		//Expected results
		expectedStatusCode := http.StatusNotFound //404 Not Found
		expectedError := seller.ErrNotFound

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "seller not found"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		sellerID := 1500 //Non existent seller ID

		// Create the mock of the service.
		serviceMock := seller.NewMockService()

		// Mock the Delete function of the service to return an error: 'seller
		// not found'. Delete function needs a parameter, the ID of the seller.
		serviceMock.On("Delete", mock.Anything, sellerID).Return(expectedError)

		// Create the handler with the mock of the service.
		handler := NewSeller(serviceMock)

		// Create the gin router and the route to test.
		r := gin.New()
		route := "/api/v1/seller/:id"

		// Add the handler to the router.
		r.DELETE(route, handler.Delete())

		//Create the request
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/seller/%d", sellerID), nil)

		//Create the response recorder.
		response := httptest.NewRecorder()

		//Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "seller not found"}
		//and the expected is a string: "seller not found".
		_ = json.Unmarshal([]byte(response.Body.Bytes()), &actualErrorResponse)

		//Get the error message from the map.
		actualErrorMessage := actualErrorResponse["message"]

		// Assert
		// Check if the status code is the correct one (404 Not Found in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualErrorMessage)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

	//Associated User Story: DELETE.
	//Edge case: delete_error_invalidID.
	//Summary: return an error due to the ID is invalid when trying to delete a seller.
	t.Run("should return an error: 'invalid id'", func(t *testing.T) {

		// Arrange
		expectedStatusCode := http.StatusBadRequest //400 Bad Request
		expectedError := errors.New("invalid id")

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "invalid id"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		serviceMock := seller.NewMockService()

		handler := NewSeller(serviceMock)

		r := gin.New()
		r.DELETE("/api/v1/seller/:id", handler.Delete())

		// Use an invalid id in the request.
		request, _ := http.NewRequest(http.MethodDelete, "/api/v1/seller/invalid", nil)
		response := httptest.NewRecorder()

		// Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "invalid id"}
		//and the expected is a string: "invalid id".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the error message from the map.
		actualMessage := actualErrorResponse["message"]

		// Assert
		// Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

	//Associated User Story: DELETE.
	//Edge case: delete_error_id_must_be_1_or_greater.
	//Summary: return an error due to the ID is not 1 or greater when trying to delete a seller.
	t.Run("should return an error: 'id must be 1 or greater'", func(t *testing.T) {

		// Arrange
		expectedStatusCode := http.StatusBadRequest //400 Bad Request
		expectedError := errors.New("id must be 1 or greater")

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "id must be 1 or greater"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		serviceMock := seller.NewMockService()

		handler := NewSeller(serviceMock)

		r := gin.New()
		r.DELETE("/api/v1/seller/:id", handler.Delete())

		// Use a negative id ('-1') in the request.
		request, _ := http.NewRequest(http.MethodDelete, "/api/v1/seller/-1", nil)
		response := httptest.NewRecorder()

		// Act
		r.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "id must be 1 or greater"}
		//and the expected is a string: "id must be 1 or greater".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the error message from the map.
		actualMessage := actualErrorResponse["message"]

		// Assert
		// Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		//Check if the mock of the service was called.
		serviceMock.AssertExpectations(t)
	})

}

// TestHandler_CREATE function.
// Test all the Create function from seller handler.
func TestHandler_CREATE(t *testing.T) {
	//Associated User Story: CREATE.
	//Edge case: create_ok.
	//Summary: return the new created seller and a 201 status code.
	t.Run("should return the new created seller and a 201 status code", func(t *testing.T) {

		//Arrange
		mockSellerService := &seller.ServiceMock{}

		// Crear un nuevo handler con el mock service
		sellerHandler := NewSeller(mockSellerService)

		// The seller to be created.
		seller := domain.Seller{
			CID:         1,
			CompanyName: "TestCompany",
			Address:     "TestAddress",
			Telephone:   "TestTelephone",
			IDLocality:  1,
		}

		// Prepare the functions of the mock of the service.
		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, 1).Return(true)
		mockSellerService.On("Save", mock.Anything, mock.Anything).Return(1, nil)

		// Create a new POST request with the seller as the body
		sellerJson, _ := json.Marshal(seller)
		request, err := http.NewRequest("POST", "/seller", bytes.NewBuffer(sellerJson))

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/seller", sellerHandler.Create())

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (201 Created in this case).
		assert.Equal(t, http.StatusCreated, response.Code)

		// Check if the error is nil.
		assert.NoError(t, err)

	})

	//Associated User Story: CREATE.
	//Edge case: create_conflict.
	//Summary: return a 409 conflict status code because the seller already exists.
	//(using CID field to check existence).
	t.Run("should return a 409 (Conflict) when the seller already exists", func(t *testing.T) {

		// Arrange
		mockSellerService := &seller.ServiceMock{}
		sellerHandler := NewSeller(mockSellerService)

		expectedError := seller.ErrSellerAlreadyExists
		expectedStatusCode := http.StatusConflict

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "seller already exists"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		sellerRequest := domain.Seller{
			CID:         1,
			CompanyName: "TestCompany",
			Address:     "TestAddress",
			Telephone:   "TestPhone",
			IDLocality:  1,
		}

		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, sellerRequest.IDLocality).Return(true)
		mockSellerService.On("Save", mock.Anything, sellerRequest).Return(0, expectedError)

		sellerJson, _ := json.Marshal(sellerRequest)
		request, _ := http.NewRequest("POST", "/seller", bytes.NewBuffer(sellerJson))

		response := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/seller", sellerHandler.Create())

		// Act
		router.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "seller already exists"}
		//and the expected is a string: "seller already exists".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the error message from the map.
		actualMessage := actualErrorResponse["message"]

		// Assert
		//Check if the status code is the correct one (409 Conflict in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: CREATE.
	//Edge case: create_fail
	//Summary: return a 422 Unprocessable Entity status code because the request not contains the required fields.
	//In this case, the request not contains the company_name field.
	t.Run("should return return a 422 Unprocessable Entity status code because the request not contains the required fields", func(t *testing.T) {

		// Arrange
		mockSellerService := &seller.ServiceMock{}
		sellerHandler := NewSeller(mockSellerService)

		expectedError := errors.New("Invalid or missing 'company_name'")
		expectedStatusCode := http.StatusUnprocessableEntity //422 Unprocessable Entity.

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "Invalid or missing 'company_name'"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		//Seller request without the company_name field
		//for this test case.
		sellerRequest := domain.Seller{
			CID:        1,
			Address:    "TestAddress",
			Telephone:  "TestPhone",
			IDLocality: 1,
		}

		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, sellerRequest.IDLocality).Return(true)

		sellerJson, _ := json.Marshal(sellerRequest)
		request, _ := http.NewRequest("POST", "/seller", bytes.NewBuffer(sellerJson))

		response := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/seller", sellerHandler.Create())

		// Act
		router.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "Invalid or missing 'company_name'"}
		//and the expected is a string: "Invalid or missing 'company_name'".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the error message from the map.
		actualMessage := actualErrorResponse["message"]

		// Assert
		//Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: CREATE.
	//Edge case: create_bad_request.
	//Summary: Return a 400 Bad Request status code because the request JSON is invalid.
	t.Run("should return a 400 Bad Request status code because the request JSON is invalid", func(t *testing.T) {

		// Arrange
		mockSellerService := &seller.ServiceMock{}
		sellerHandler := NewSeller(mockSellerService)

		expectedError := errors.New("invalid json")
		expectedStatusCode := http.StatusBadRequest //400 Bad Request.

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "invalid json"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		//Create a string with an incorrect JSON structure and convert it to an io.Reader
		//to use it as the body of the request.
		incorrectJson := `{"this is: an incorrect JSON"}`
		reader := strings.NewReader(incorrectJson)

		//Create the request with the incorrect JSON as the body.
		request, _ := http.NewRequest("POST", "/seller", reader)

		// Create the response recorder and the router.
		response := httptest.NewRecorder()
		router := gin.Default()

		//Register the route.
		router.POST("/seller", sellerHandler.Create())

		// Act
		router.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "invalid json"}
		//and the expected is a string: "invalid json".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the error message from the map.
		actualMessage := actualErrorResponse["message"]

		// Assert
		//Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error message is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: CREATE.
	//Edge case: create_error_invalid_CID.
	//Summary: Return a 422 Unprocessable Entity status code because the CID is invalid.
	//In this case the CID is a string and not an int.
	t.Run("should return a 422 Unprocessable Entity status code because the CID is invalid", func(t *testing.T) {

		// Arrange
		mockSellerService := &seller.ServiceMock{}
		sellerHandler := NewSeller(mockSellerService)

		expectedError := errors.New("invalid or missing cid. CID must be 1 or greater")
		expectedStatusCode := http.StatusUnprocessableEntity //422 Unprocessable Entity.

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "invalid or missing cid. CID must be 1 or greater"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		//Create a string with a string CID instead of an int and convert it to an io.Reader
		//to use it as the body of the request.
		incorrectJson := `{
			"cid":"22",
			"company_name": "Empresa",
			"address": "calle siempreviva",
			"telephone":"555",
			"locality_id": 1
		  }`
		reader := strings.NewReader(incorrectJson)

		//Create the request with the string CID in the JSON as the body.
		request, _ := http.NewRequest("POST", "/seller", reader)

		// Create the response recorder and the router.
		response := httptest.NewRecorder()
		router := gin.Default()

		//Register the route.
		router.POST("/seller", sellerHandler.Create())

		// Act
		router.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "invalid or missing cid. CID must be 1 or greater"}
		//and the expected is a string: "invalid or missing cid. CID must be 1 or greater".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the error message from the map.
		actualMessage := actualErrorResponse["message"]

		// Assert
		//Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error message is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: CREATE.
	//Edge case: create_error_invalid_Locality_ID.
	//Summary: Return a 422 Unprocessable Entity status code because the Locality_id is invalid.
	//In this case the Locality_id is a string and not an int.
	t.Run("should return a 422 Unprocessable Entity status code because the Locality_id is invalid", func(t *testing.T) {

		// Arrange
		mockSellerService := &seller.ServiceMock{}
		sellerHandler := NewSeller(mockSellerService)

		expectedError := errors.New("invalid or missing locality. Locality_id must be 1 or greater")
		expectedStatusCode := http.StatusUnprocessableEntity //422 Unprocessable Entity.

		//Use a map to unmarshal the response. The response is a JSON with the following structure:
		// {
		// 	"error": "invalid locality_id"
		// }
		// Create the map to unmarshal the response.
		actualErrorResponse := make(map[string]string)

		//Create a string JSON with a string Locality_id instead of an int and
		//convert it to an io.Reader to use it as the body of the request.
		incorrectJson := `{
			"cid":22,
			"company_name": "Empresa",
			"address": "calle siempreviva",
			"telephone":"555",
			"locality_id": "1"
		  }`
		reader := strings.NewReader(incorrectJson)

		//Create the request with the string Locality_id in the JSON as the body.
		request, _ := http.NewRequest("POST", "/seller", reader)

		// Create the response recorder and the router.
		response := httptest.NewRecorder()
		router := gin.Default()

		//Register the route.
		router.POST("/seller", sellerHandler.Create())

		// Act
		router.ServeHTTP(response, request)

		//Unmarshal the response into the actualErrorResponse
		//because the response is a JSON with the format: {"message": "nvalid or missing locality. Locality_id must be 1 or greater"}
		//and the expected is a string: "invalid or missing locality. Locality_id must be 1 or greater".
		_ = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)

		//Get the error message from the map.
		actualMessage := actualErrorResponse["message"]

		// Assert
		//Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the obtained error message is the same than the expected one.
		assert.Equal(t, expectedError.Error(), actualMessage)

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: CREATE.
	//Edge case: create_error_non_existent_Locality_ID.
	//Summary: return a 422 status code because the 'locality_id' does not exist.
	t.Run("should return a 422 status code beacuse the 'locality_id' does not exist", func(t *testing.T) {

		//Arrange
		mockSellerService := &seller.ServiceMock{}

		// Crear un nuevo handler con el mock service
		sellerHandler := NewSeller(mockSellerService)

		// The seller to be created.
		seller := domain.Seller{
			CID:         1,
			CompanyName: "TestCompany",
			Address:     "TestAddress",
			Telephone:   "TestTelephone",
			IDLocality:  1500, //This locality_id does not exist.
		}

		expectedCode := http.StatusUnprocessableEntity //422 Unprocessable Entity.

		//The expected body (id locality not exists error message)
		//of the response in JSON format.
		expectedBody := `{"code":"unprocessable_entity", "message":"id locality not exists"}`

		// Prepare the function of the service_mock.
		// The function GetLocalityIdFromSeller returns false because the locality_id does not exist. This function needs two parameters: the context and the locality_id.
		//In this case, the context is not important, so it is mocked with mock.Anything.
		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, 1500).Return(false)

		// Create a new POST request with the seller as the body
		sellerJson, _ := json.Marshal(seller)
		request, _ := http.NewRequest("POST", "/seller", bytes.NewBuffer(sellerJson))

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/seller", sellerHandler.Create())

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedCode, response.Code)

		// Check if the error message is the correct one
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)

	})

	//Associated User Story: CREATE.
	//Edge case: create_error_500.
	//Summary: return a 500 status code due an internal server error occurs
	// when the 'sellerService' Save function was executed
	t.Run("should return a 500 status code", func(t *testing.T) {

		//Arrange

		//Expected results
		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error.
		expectedError := errors.New("internal server error")

		//The expected body (internal server error)
		//of the response in JSON format.
		expectedBody := `{"code":"internal_server_error", "message":"internal server error"}`

		mockSellerService := &seller.ServiceMock{}

		// Crear un nuevo handler con el mock service
		sellerHandler := NewSeller(mockSellerService)

		// The seller to be created.
		seller := domain.Seller{
			CID:         1,
			CompanyName: "TestCompany",
			Address:     "TestAddress",
			Telephone:   "TestTelephone",
			IDLocality:  1,
		}

		// Prepare the functions of the mock of the service.
		// The function GetLocalityIdFromSeller returns true because the locality_id exists.
		//This function needs two parameters: the context and the locality_id. In this case,
		//the context is not important, so it is mocked with mock.Anything.
		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, 1).Return(true)
		// The function Save returns 0 and an error because an internal server error occurs.
		mockSellerService.On("Save", mock.Anything, mock.Anything).Return(0, expectedError)

		// Create a new POST request with the seller as the body
		sellerJson, _ := json.Marshal(seller)
		request, _ := http.NewRequest("POST", "/seller", bytes.NewBuffer(sellerJson))

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/seller", sellerHandler.Create())

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message is the correct one
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)

	})

}

// TestHandler_UPDATE function.
// Test the Update function from seller handler.
func TestHandler_UPDATE(t *testing.T) {

	//Associated User Story: UPDATE.
	//Edge case: update_ok.
	//Summary: return the updated seller and a 200 status code.
	t.Run("should update seller and return a 200 status code", func(t *testing.T) {

		// Arrange
		sellerID := 1

		// The original seller to be updated.
		originalSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Company",
			Address:     "Address",
			Telephone:   "Telephone",
			IDLocality:  1,
		}

		//The original seller with the modifications in CompanyName, Address,
		// Telephone and IDLocality.
		updatedSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "UpdatedCompany",
			Address:     "UpdatedAddress",
			Telephone:   "UpdatedTelephone",
			IDLocality:  2,
		}

		expectedStatusCode := http.StatusOK //200 OK.

		//The body of the request in JSON format.
		bodyRequestJSON := `{
		"cid":1,
		"company_name": "UpdatedCompany",
		"address": "UpdatedAddress",
		"telephone":"UpdatedTelephone",
		"locality_id": 2
	   }`

		//The expected body (the seller with modifications)
		//of the response in JSON format.
		expectedBody := `{
		"data": {
			"id": 1,
			"cid": 1,
			"company_name": "UpdatedCompany",
			"address": "UpdatedAddress",
			"telephone": "UpdatedTelephone",
			"locality_id": 2
			}
		}`

		// Create a mock of the service.
		mockSellerService := &seller.ServiceMock{}

		// Create a new handler with the mock service
		sellerHandler := NewSeller(mockSellerService)

		// Prepare the mock service with the expected inputs and outputs.
		mockSellerService.On("GetSellerByID", mock.Anything, sellerID).Return(originalSeller, nil)
		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, 2).Return(true)
		mockSellerService.On("Update", mock.Anything, updatedSeller).Return(nil)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller/:id"
		router.PATCH(route, sellerHandler.Update())

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(bodyRequestJSON)

		//Create the request with the method PATCH, the route and the body of the request.
		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/%d", sellerID), bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (200 OK in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the seller data was correctly updated.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: UPDATE.
	//Edge case: update_error
	//Summary: return an 500 status code because an internal server error occurs.
	t.Run("should return an 500 status code because an internal server error occurs", func(t *testing.T) {

		// Arrange
		sellerID := 1500 //The seller does not exist.

		bodyRequestJSON := `{
			"cid":1,
			"company_name": "UpdatedCompany",
			"address": "UpdatedAddress",
			"telephone":"UpdatedTelephone",
			"locality_id": 2
		   }`

		expectedBody := `{"code":"internal_server_error", "message":"internal server error"}`

		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error.
		expectedSeller := domain.Seller{}                    //Empty seller
		expectedError := errors.New("internal server error")

		// Create a mock of the service.
		mockSellerService := &seller.ServiceMock{}

		// Create a new handler with the mock service
		sellerHandler := NewSeller(mockSellerService)

		// Prepare the mock service with the expected inputs and outputs.
		mockSellerService.On("GetSellerByID", mock.Anything, sellerID).Return(expectedSeller, expectedError)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller/:id"
		router.PATCH(route, sellerHandler.Update())

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(bodyRequestJSON)

		//Create the request with the method PATCH, the route and the body of the request.
		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/%d", sellerID), bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message is the correct one.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: UPDATE.
	//Edge case: update_non_existent.
	//Summary: return an 404 status code because the seller does not exist.
	t.Run("should return an 404 status code because the seller does not exist", func(t *testing.T) {

		// Arrange
		sellerID := 1500 //The seller does not exist.

		//The body of the request in JSON format.
		bodyRequestJSON := `{
		"cid":1,
		"company_name": "UpdatedCompany",
		"address": "UpdatedAddress",
		"telephone":"UpdatedTelephone",
		"locality_id": 2
	   }`

		//The expected body of the response in JSON format.
		expectedBody := `{"code":"not_found", "message":"seller not found"}`

		//The expected results.
		expectedStatusCode := http.StatusNotFound //404 Not Found.
		expectedSeller := domain.Seller{}         //Empty seller
		expectedError := seller.ErrNotFound       // message: "seller not found"

		// Create a mock of the service.
		mockSellerService := &seller.ServiceMock{}

		// Create a new handler with the mock service.
		sellerHandler := NewSeller(mockSellerService)

		// Prepare the mock service with the expected inputs and outputs.
		// This function return an empty seller and an error due the seller does not exist.
		mockSellerService.On("GetSellerByID", mock.Anything, sellerID).Return(expectedSeller, expectedError)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller/:id"
		router.PATCH(route, sellerHandler.Update())

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(bodyRequestJSON)

		//Create the request with the method PATCH, the route and the body of the request.
		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/%d", sellerID), bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (404 Not Found in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message is the correct one.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: UPDATE.
	//Edge case: update_error_invalidID.
	//Summary: return an 400 status code because the seller id is invalid.
	//In this case, seller id is a string.
	t.Run("should return an 400 status code because the seller id is invalid", func(t *testing.T) {

		// Arrange
		sellerID := "invalidID" //The seller id is a string instead of a valid integer.

		//The body of the request in JSON format.
		bodyRequestJSON := `{
		"cid":1,
		"company_name": "UpdatedCompany",
		"address": "UpdatedAddress",
		"telephone":"UpdatedTelephone",
		"locality_id": 2
	   }`

		//The expected body of the response in JSON format.
		expectedBody := `{"code":"bad_request", "message":"id must be 1 or greater"}`

		//The expected results.
		expectedStatusCode := http.StatusBadRequest //400 Bad Request.

		// Create a mock of the service.
		mockSellerService := &seller.ServiceMock{}

		// Create a new handler with the mock service.
		sellerHandler := NewSeller(mockSellerService)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller/:id"
		router.PATCH(route, sellerHandler.Update())

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(bodyRequestJSON)

		//Create the request with the method PATCH, the route and the body of the request.
		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/%s", sellerID), bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message is the correct one.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: UPDATE.
	//Edge case: update_error_invalidJSON.
	//Summary: return a 400 status code because the JSON of the request is invalid.
	//In this case, the JSON has an extra comma in the field 'locality_id'.
	t.Run("should return return a 400 status code because the JSON of the request is invalid", func(t *testing.T) {

		// Arrange
		sellerID := 1

		expectedSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Company",
			Address:     "Address",
			Telephone:   "Telephone",
			IDLocality:  1,
		}

		//The body of the request in JSON format.
		//In this case, the JSON has an extra comma in the field 'locality_id'.
		invalidBodyRequestJSON := `{
		"cid":1,
		"company_name": "UpdatedCompany",
		"address": "UpdatedAddress",
		"telephone":"UpdatedTelephone",
		"locality_id": 2,
	   }`

		//The expected body of the response in JSON format.
		expectedBody := `{"code":"bad_request", "message":"invalid json"}`

		//The expected results.
		expectedStatusCode := http.StatusBadRequest //400 Bad Request.

		// Create a mock of the service.
		mockSellerService := &seller.ServiceMock{}

		// Create a new handler with the mock service.
		sellerHandler := NewSeller(mockSellerService)

		// Prepare the mock service with the expected inputs and outputs.
		// This function return a seller and nil due the seller exists.
		mockSellerService.On("GetSellerByID", mock.Anything, sellerID).Return(expectedSeller, nil)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller/:id"
		router.PATCH(route, sellerHandler.Update())

		//Create the body for the request with the data in JSON format.
		//Use the invalid JSON for this test.
		bodyRequest := strings.NewReader(invalidBodyRequestJSON)

		//Create the request with the method PATCH, the route and the body of the request.
		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/%d", sellerID), bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (400 Bad Request in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message is the correct one.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: UPDATE.
	//Edge case: update_error_LocalityID_non_existent.
	//Summary: return a 422 status code because the 'locality_id', to update a seller, does not exist.
	t.Run("should return a 422 status code because the 'locality_id' of the seller does not exist", func(t *testing.T) {

		// Arrange
		sellerID := 1

		//This seller has a locality_id that exists.
		originalSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Company",
			Address:     "Address",
			Telephone:   "Telephone",
			IDLocality:  1,
		}

		//The body of the request, to update
		//the seller, in JSON format
		//with a locality_id that does not exist
		BodyRequestJSON := `{
		"cid":1,
		"company_name": "UpdatedCompany",
		"address": "UpdatedAddress",
		"telephone":"UpdatedTelephone",
		"locality_id": 1500
	   }`

		//The expected body of the response in JSON format.
		expectedBody := `{"code":"unprocessable_entity", "message":"id locality not exists"}`

		//422 Unprocessable Entity.
		expectedStatusCode := http.StatusUnprocessableEntity

		// Create a mock of the service.
		mockSellerService := &seller.ServiceMock{}

		// Create a new handler with the mock service.
		sellerHandler := NewSeller(mockSellerService)

		// Prepare the mock service with the expected inputs and outputs.
		// This function return a seller and nil due the seller exists.
		mockSellerService.On("GetSellerByID", mock.Anything, sellerID).Return(originalSeller, nil)
		// This function return false due the locality_id does not exist.
		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, 1500).Return(false)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller/:id"
		router.PATCH(route, sellerHandler.Update())

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(BodyRequestJSON)

		//Create the request with the method PATCH, the route and the body of the request.
		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/%d", sellerID), bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (422 Unprocessable Entity in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message is the correct one.
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})

	//Associated User Story: UPDATE.
	//Edge case: update_error_500.
	//Summary: return a 500 status code after execute the Update function from 'sellerService'.
	t.Run("should eturn a 500 status code after execute the Update function from 'sellerService'", func(t *testing.T) {

		// Arrange
		sellerID := 1

		// The original seller to be updated.
		originalSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Company",
			Address:     "Address",
			Telephone:   "Telephone",
			IDLocality:  1,
		}

		//The original seller with the modifications in CompanyName, Address,
		// Telephone and IDLocality.
		updatedSeller := domain.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "UpdatedCompany",
			Address:     "UpdatedAddress",
			Telephone:   "UpdatedTelephone",
			IDLocality:  2,
		}

		expectedStatusCode := http.StatusInternalServerError //500 Internal Server Error.
		expectedError := errors.New("internal server error")

		//The body of the request in JSON format.
		bodyRequestJSON := `{
		"cid":1,
		"company_name": "UpdatedCompany",
		"address": "UpdatedAddress",
		"telephone":"UpdatedTelephone",
		"locality_id": 2
	   }`

		//The expected body (internal server error)
		//of the response in JSON format.
		expectedBody := `{"code":"internal_server_error", "message":"internal server error"}`

		// Create a mock of the service.
		mockSellerService := &seller.ServiceMock{}

		// Create a new handler with the mock service
		sellerHandler := NewSeller(mockSellerService)

		// Prepare the mock service with the expected inputs and outputs.
		// This function return a seller and nil due the seller exists.
		mockSellerService.On("GetSellerByID", mock.Anything, sellerID).Return(originalSeller, nil)
		// This function return true due the locality_id exists.
		mockSellerService.On("GetLocalityIdFromSeller", mock.Anything, 2).Return(true)
		// This function return an error due the Update function from 'sellerService' failed
		//for some internal reason.
		mockSellerService.On("Update", mock.Anything, updatedSeller).Return(expectedError)

		//Create the router and set the route.
		router := gin.New()
		route := "/api/v1/seller/:id"
		router.PATCH(route, sellerHandler.Update())

		//Create the body for the request with the data in JSON format.
		bodyRequest := strings.NewReader(bodyRequestJSON)

		//Create the request with the method PATCH, the route and the body of the request.
		request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/seller/%d", sellerID), bodyRequest)

		// Register the route and create the response recorder.
		response := httptest.NewRecorder()

		//Act
		router.ServeHTTP(response, request)

		//Assert
		// Check if the status code is the correct one (500 Internal Server Error in this case).
		assert.Equal(t, expectedStatusCode, response.Code)

		// Check if the error message is the correct one
		assert.JSONEq(t, expectedBody, response.Body.String())

		// Check if the mock of the service was called.
		mockSellerService.AssertExpectations(t)
	})
}
