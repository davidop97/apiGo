package handler

import (
	"errors"
	"fmt"
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

func TestHandler_Read(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		//Given
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
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"data":[{"id":1,"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1},
								  {"id":2,"card_number_id":"C987E012F","first_name":"George","last_name":"Smith","warehouse_id":2}]}`
		)
		service := &employee.ServiceMock{}
		service.On("GetAllEmployees", mock.Anything).Return(expectedEmployees, nil)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees"
		engine.GET(route, handler.GetAll())
		request, _ := http.NewRequest("GET", route, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("find_by_id_existent", func(t *testing.T) {
		//Given
		var (
			id               = 1
			expectedEmployee = domain.Employee{
				ID:           1,
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			}
			expectedStatusCode = http.StatusOK
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"data":{"id":1,"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}}`
		)
		service := &employee.ServiceMock{}
		service.On("GetEmployeeByID", mock.Anything, id).Return(expectedEmployee, nil)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.GET(route, handler.Get())
		url := fmt.Sprintf("/api/v1/employees/%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("find_by_id_non_existent", func(t *testing.T) {
		//Given
		var (
			id                 = 1
			expectedEmployee   = domain.Employee{}
			expectedStatusCode = http.StatusNotFound
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"message":"employee not found"}`
		)
		service := &employee.ServiceMock{}
		service.On("GetEmployeeByID", mock.Anything, id).Return(expectedEmployee, employee.ErrNotFound)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.GET(route, handler.Get())
		url := fmt.Sprintf("/api/v1/employees/%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
}
func TestHandler_Create(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		//Given
		var (
			id              = 1
			employeeRequest = domain.Employee{
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			}
			expectedStatusCode = http.StatusCreated
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}`
			expectedBody = `{"data":{"id":1,"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}}`
		)
		service := &employee.ServiceMock{}
		service.On("SaveEmployee", mock.Anything, employeeRequest).Return(id, nil)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees"
		engine.POST(route, handler.Create())
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("create_conflict", func(t *testing.T) {
		//Given
		var (
			employeeRequest = domain.Employee{
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			}
			expectedStatusCode = http.StatusConflict
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}`
			expectedBody = `{"message": "duplicate employee number"}`
			err          = employee.ErrEmployeeAlreadyExists
		)
		service := &employee.ServiceMock{}
		service.On("SaveEmployee", mock.Anything, employeeRequest).Return(0, err)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees"
		engine.POST(route, handler.Create())
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("create_fail", func(t *testing.T) {
		//Given
		var (
			employeeRequest = domain.Employee{
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			}
			expectedStatusCode = http.StatusUnprocessableEntity
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest        = `{"first_name":"Harold","last_name":"Doe","warehouse_id":1}`
			expectedBody       = `{"message":"field %s is empty"}`
			ErremptyFieldError = errors.New("field card_number_id is empty")
			err                = ErremptyFieldError
		)
		service := &employee.ServiceMock{}
		service.On("SaveEmployee", mock.Anything, employeeRequest).Return(0, err)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees"
		engine.POST(route, handler.Create())
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, fmt.Sprintf(expectedBody, "card_number_id"), response.Body.String())
		// service.AssertExpectations(t)
	})
}
func TestHandler_Update(t *testing.T) {
	t.Run("update_ok", func(t *testing.T) {
		//Given
		var (
			id               = 1
			originalEmployee = domain.Employee{
				ID:           1,
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			}
			updatedEmployee = domain.Employee{
				ID:           1,
				CardNumberID: "D789E012F",
				FirstName:    "George",
				LastName:     "Smith",
				WarehouseID:  3,
			}
			expectedStatusCode = http.StatusOK
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"card_number_id":"D789E012F","first_name":"George","last_name":"Smith","warehouse_id":3}`
			expectedBody = `{"data":{"id":1,"card_number_id":"D789E012F","first_name":"George","last_name":"Smith","warehouse_id":3}}`
		)
		service := &employee.ServiceMock{}
		service.On("GetEmployeeByID", mock.Anything, id).Return(originalEmployee, nil)
		service.On("UpdateEmployee", mock.Anything, updatedEmployee).Return(nil)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.PATCH(route, handler.Update())
		url := fmt.Sprintf("/api/v1/employees/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("update_duplicate_employee", func(t *testing.T) {
		//Given
		var (
			id               = 1
			originalEmployee = domain.Employee{
				ID:           1,
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			}
			updatedEmployee = domain.Employee{
				ID:           1,
				CardNumberID: "D789E012F",
				FirstName:    "Harold",
				LastName:     "Doe",
				WarehouseID:  1,
			}
			expectedStatusCode = http.StatusConflict
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"card_number_id":"D789E012F","first_name":"Harold","last_name":"Doe","warehouse_id":1}`
			expectedBody = `{"message": "duplicate section number"}`
		)
		err := employee.ErrEmployeeAlreadyExists
		service := &employee.ServiceMock{}
		service.On("GetEmployeeByID", mock.Anything, id).Return(originalEmployee, nil)
		service.On("UpdateEmployee", mock.Anything, updatedEmployee).Return(err)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.PATCH(route, handler.Update())
		url := fmt.Sprintf("/api/v1/employees/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("update_non_existent", func(t *testing.T) {
		//Given
		var (
			id                 = 1
			emptyEmployee      = domain.Employee{}
			expectedStatusCode = http.StatusNotFound
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"card_number_id":"D789E012F","first_name":"George","last_name":"Smith","warehouse_id":3}`
			expectedBody = `{"message":"employee not found"}`
			err          = employee.ErrNotFound
		)
		service := &employee.ServiceMock{}
		service.On("GetEmployeeByID", mock.Anything, id).Return(emptyEmployee, err)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.PATCH(route, handler.Update())
		url := fmt.Sprintf("/api/v1/employees/%d", id)
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("PATCH", url, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("invalid_id", func(t *testing.T) {
		//Given
		var (
			id                 = "invalid id"
			expectedStatusCode = http.StatusBadRequest
			expectedBody       = `{"message":"bad id"}`
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
		)
		service := &employee.ServiceMock{}
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.PATCH(route, handler.Update())
		url := fmt.Sprintf("/api/v1/employees/%s", id)
		request, _ := http.NewRequest("PATCH", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
}
func TestHandler_Delete(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {
		//Given
		var (
			id              = 1
			expectedHeaders = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedStatusCode = http.StatusNoContent
		)
		service := &employee.ServiceMock{}
		service.On("DeleteEmployee", mock.Anything, id).Return(nil)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.DELETE(route, handler.Delete())
		url := fmt.Sprintf("/api/v1/employees/%d", id)
		request, _ := http.NewRequest("DELETE", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.Empty(t, response.Body)
		service.AssertExpectations(t)
	})
	t.Run("delete_non_existent", func(t *testing.T) {
		//Given
		var (
			id              = 1
			expectedHeaders = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedStatusCode = http.StatusNotFound
			err                = employee.ErrNotFound
			expectedBody       = `{"message":"employee not found"}`
		)
		service := &employee.ServiceMock{}
		service.On("DeleteEmployee", mock.Anything, id).Return(err)
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.DELETE(route, handler.Delete())
		url := fmt.Sprintf("/api/v1/employees/%d", id)
		request, _ := http.NewRequest("DELETE", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("invalid_id", func(t *testing.T) {
		//Given
		var (
			id                 = "invalid id"
			expectedStatusCode = http.StatusBadRequest
			expectedBody       = `{"message":"bad id"}`
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
		)
		service := &employee.ServiceMock{}
		handler := NewEmployee(service)
		engine := gin.New()
		route := "/api/v1/employees/:id"
		engine.DELETE(route, handler.Delete())
		url := fmt.Sprintf("/api/v1/employees/%s", id)
		request, _ := http.NewRequest("DELETE", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		service.AssertExpectations(t)
	})
}
