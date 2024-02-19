package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidop97/apiGo/internal/buyer"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test integration for Buyer domain in handler layer

// TestIntegrationHandler_ReadBuyer
// Integration test for Read functions in handler layer
// Tested methods: Read
func TestIntegrationHandler_ReadBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: READ
	// EDGE CASE: find_all
	// DESCRIPTION: When the request is successful the backend will return a list of all existing buyers.
	t.Run("it should returns all buyers", func(t *testing.T) {
		// arrange

		// prepare expected results
		expectedBuyers := `{"data":[{"id":1,"card_number_id":"123456789","first_name":"Jane","last_name":"Doe"},
							{"id":2,"card_number_id":"987654321","first_name":"John","last_name":"Smith"}]}`
		expectedStatusCode := http.StatusOK

		// expected result
		expectedbuyers := []domain.Buyer{{
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
		// create a mock for repository
		repositoryMock := &buyer.RepositoryMock{}
		// set mock with the expected value
		repositoryMock.On("GetAll", mock.Anything).Return(expectedbuyers, nil)

		// create mock of the service
		service := buyer.NewService(repositoryMock)

		// create handler using mock service
		handler := NewBuyer(service)

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
		repositoryMock.AssertExpectations(t)

	})
}

// TestIntegrationHandler_DeleteBuyer
// Integration test for Delete functions in handler layer
// Tested methods: Delete
func TestIntegrationHandler_DeleteBuyer(t *testing.T) {

	// ASSOCIATED USER STORY: DELETE
	// EDGE CASE: delete_non_existent
	// DESCRIPTION: When the buyer does not exist, a 404 code will be returned.
	t.Run("it should returns 404 when the buyer does not exist", func(t *testing.T) {
		expectedID := "1"
		//expectedError
		expectedMessageError := `{"error":"Buyer not found"}`
		// prepare expected results
		expectedStatusCode := http.StatusNotFound

		// create a mock for repository
		repositoryMock := &buyer.RepositoryMock{}
		// set mock with the expected value
		repositoryMock.On("Get", mock.Anything, 1).Return(domain.Buyer{}, buyer.ErrNotFound)
		// when the buyer does not exists, delete function it is not called.

		// call to service interface
		service := buyer.NewService(repositoryMock)

		// create handler using mock service
		handler := NewBuyer(service)

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
		repositoryMock.AssertExpectations(t)
	})
}
