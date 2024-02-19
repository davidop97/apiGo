package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/employee"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIntegration_ReadEmployee(t *testing.T) {
	t.Run("should return a list of employees", func(t *testing.T) {
		// Given
		var (
			expectedEmployees = []domain.Employee{
				{
					ID:           1,
					CardNumberID: "D789E012F",
					FirstName:    "Harold",
					LastName:     "Doe",
					WarehouseID:  1,
				},
				{
					ID:           2,
					CardNumberID: "C987E012F",
					FirstName:    "George",
					LastName:     "Smith",
					WarehouseID:  2,
				},
			}
			expectedStatusCode = http.StatusOK
			expectedBody       = `{"data":[{"id":1,"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1},
								  {"id":2,"card_number_id":"C987E012F","first_name":"George","last_name":"Smith","warehouse_id":2}]}`
		)
		repository := &employee.RepositoryMock{}
		repository.On("GetAll", mock.Anything).Return(expectedEmployees, nil)
		service := employee.NewService(repository)
		handler := NewEmployee(service)
		r := gin.New()
		route := "/api/v1/employees"
		r.GET(route, handler.GetAll())
		// - create http GET request
		request, _ := http.NewRequest("GET", route, nil)
		// - set up response recorder to capture handler's response
		response := httptest.NewRecorder()

		// When
		r.ServeHTTP(response, request)

		// Assert
		// - check if the obtained status code matches expected
		assert.Equal(t, expectedStatusCode, response.Code)
		// - check if the obtained response body macthes expected
		assert.JSONEq(t, expectedBody, response.Body.String())
		// - verify if expectations set on the repository mock were met
		repository.AssertExpectations(t)
	})
	t.Run("it should return an error if repository is unable to connect to the database when calling GetAll", func(t *testing.T) {
		// Arrange
		// - Preparing an empty section list to simulate no data returned from the service.
		emptySectionList := []domain.Employee{}

		// - Specifying the expected HTTP status code for an internal server error.
		expectedStatusCode := http.StatusInternalServerError

		// - Creating an error to simulate a database connection issue.
		err := errors.New("database connection error")

		// - Declaring the expected response body the handler should return.
		expectedBody := `{}`

		// - Creating a mock of the repository layer.
		repository := &employee.RepositoryMock{}
		// - Setting up the mock response for the GetAll method to return an error.
		repository.On("GetAll", mock.Anything).Return(emptySectionList, err)

		// - Instantiate service with mocked repository
		service := employee.NewService(repository)

		// - Instantiating the handler with the mocked service.
		handler := NewEmployee(service)

		// - Creating a new Gin router and registering the handler for the route.
		r := gin.New()
		route := "/api/v1/employees"
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

		// - Verifying that all expectations set on the repositoory mock were met.
		repository.AssertExpectations(t)
	})
}

func TestIntegration_CreateEmployee(t *testing.T) {
	t.Run("it should create a employee and return its id", func(t *testing.T) {
		// Arrange
		// - simulate a section creation request
		id := 1
		employeeRequest := domain.Employee{
			ID:           0,
			CardNumberID: "D789E012F",
			FirstName:    "Harold",
			LastName:     "Doe",
			WarehouseID:  1,
		}
		// - request body as JSON
		bodyRequest := `{"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}`
		// - declare expected status code and response body
		expectedStatusCode := http.StatusCreated
		expectedBody := `{"data":{"id":1,"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}}`
		// - create a repository mock
		repository := &employee.RepositoryMock{}
		repository.On("Exists", mock.Anything, employeeRequest.CardNumberID).Return(false)
		repository.On("Save", mock.Anything, employeeRequest).Return(id, nil)
		// - instantiate service with mocked repository
		service := employee.NewService(repository)
		// - instantiate handler
		handler := NewEmployee(service)
		// - create router and register handler
		r := gin.New()
		route := "/api/v1/employees"
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
		// - verify if expectations on the repository mock were met
		repository.AssertExpectations(t)
	})

	t.Run("it should return an error if the employee number already exists", func(t *testing.T) {
		// When
		CardNumberID := "D789E012F"
		// - Request body as JSON
		bodyRequest := `{"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}`
		// - Declare expected status code and response body for a conflict error
		expectedStatusCode := http.StatusConflict
		expectedBody := `{"message": "duplicate employee number"}`
		// - Create a repository mock
		repository := &employee.RepositoryMock{}
		repository.On("Exists", mock.Anything, CardNumberID).Return(true)
		// - Instantiate service with mocked repository
		service := employee.NewService(repository)
		// - Instantiate handler
		handler := NewEmployee(service)
		// - Create router and register handler
		r := gin.New()
		route := "/api/v1/employees"
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
		repository.AssertExpectations(t)
	})
}
