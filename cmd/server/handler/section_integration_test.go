package handler

import (
	"errors"
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

func TestIntegration_ReadSection(t *testing.T) {
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

		// - mock repository layer
		repository := &section.RepositoryMock{}
		repository.On("GetAll", mock.Anything).Return(expectedSections, nil)
		// - instantiate service with the mocked repository
		service := section.NewService(repository)
		// - instantiate handler
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
		// - verify if expectations set on the repository mock were met
		repository.AssertExpectations(t)
	})

	t.Run("it should return an error if repository is unable to connect to the database when calling GetAll", func(t *testing.T) {
		// Arrange
		// - Preparing an empty section list to simulate no data returned from the service.
		emptySectionList := []domain.Section{}

		// - Specifying the expected HTTP status code for an internal server error.
		expectedStatusCode := http.StatusInternalServerError

		// - Creating an error to simulate a database connection issue.
		err := errors.New("some internal error")

		// - Declaring the expected response body the handler should return.
		expectedBody := `{"message":"internal error"}`

		// - Creating a mock of the repository layer.
		repository := &section.RepositoryMock{}
		// - Setting up the mock response for the GetAll method to return an error.
		repository.On("GetAll", mock.Anything).Return(emptySectionList, err)

		// - Instantiate service with mocked repository
		service := section.NewService(repository)

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

		// - Verifying that all expectations set on the repositoory mock were met.
		repository.AssertExpectations(t)
	})
}

func TestIntegration_CreateSection(t *testing.T) {
	t.Run("it should create a section and return its id", func(t *testing.T) {
		// Arrange
		// - simulate a section creation request
		id := 1
		sectionNumber := 1
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
		// - create a repository mock
		repository := &section.RepositoryMock{}
		repository.On("Exists", mock.Anything, sectionNumber).Return(false)
		repository.On("Save", mock.Anything, sectionRequest).Return(id, nil)
		// - instantiate service with mocked repository
		service := section.NewService(repository)
		// - instantiate handler
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
		// - verify if expectations on the repository mock were met
		repository.AssertExpectations(t)
	})

	t.Run("it should return an error if the section number already exists", func(t *testing.T) {
		// Arrange
		// - Simulate a section creation request with a section number that already exists.
		sectionNumber := 1
		// - Request body as JSON
		bodyRequest := `{"section_number":1,"current_temperature":1,"minimum_temperature":1,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`
		// - Declare expected status code and response body for a conflict error
		expectedStatusCode := http.StatusConflict
		expectedBody := `{"message": "duplicate section number"}`
		// - Create a repository mock
		repository := &section.RepositoryMock{}
		repository.On("Exists", mock.Anything, sectionNumber).Return(true)
		// - Instantiate service with mocked repository
		service := section.NewService(repository)
		// - Instantiate handler
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
		repository.AssertExpectations(t)
	})
}
