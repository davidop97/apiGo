package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	inboudorder "github.com/davidop97/apiGo/internal/inboudOrder"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetAllReports(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//Given
		var (
			expectedReports = []inboudorder.Report{
				{
					Employee: &domain.Employee{
						ID:           1,
						CardNumberID: "402323",
						FirstName:    "Harold",
						LastName:     "Doe",
						WarehouseID:  1,
					},
					InboudOrdersCount: 1,
				},
				{
					Employee: &domain.Employee{
						ID:           2,
						CardNumberID: "402324",
						FirstName:    "Jane",
						LastName:     "Doe",
						WarehouseID:  2,
					},
					InboudOrdersCount: 2,
				},
			}
			expectedStatusCode = http.StatusOK
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"data":[{"id":1,"card_number_id":"402323","first_name":"Harold","last_name":"Doe","warehouse_id":1, "inboud_orders_count":1},
		{"id":2,"card_number_id":"402324","first_name":"Jane","last_name":"Doe","warehouse_id":2, "inboud_orders_count":2}]}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("GetAllReports", mock.Anything).Return(expectedReports, nil)
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/employees/reportInboundOrder"
		engine.GET(route, handler.GetAllReports())
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
	t.Run("error", func(t *testing.T) {
		//Given
		var (
			expectedReports    = []inboudorder.Report{}
			expectedStatusCode = http.StatusInternalServerError
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"message":"Internal Server Error"}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("GetAllReports", mock.Anything).Return(expectedReports, errors.New("internal server error"))
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/employees/reportInboundOrder"
		engine.GET(route, handler.GetAllReports())
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
}

func TestHandler_GenerateReport(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//Given
		var (
			id             = 1
			expectedReport = inboudorder.Report{
				Employee: &domain.Employee{
					ID:           1,
					CardNumberID: "A123B456C",
					FirstName:    "John",
					LastName:     "Doe",
					WarehouseID:  1,
				},
				InboudOrdersCount: 3,
			}
			expectedStatusCode = http.StatusOK
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"data":{"id":1,"card_number_id":"A123B456C","first_name":"John","last_name":"Doe","warehouse_id":1, "inboud_orders_count":3}}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("GenerateReport", mock.Anything, id).Return(expectedReport, nil)
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/employees/reportInboundOrders"
		engine.GET(route, handler.GenerateReport())
		url := fmt.Sprintf("/employees/reportInboundOrders?id=%d", id)
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
	t.Run("error_employee_id_required", func(t *testing.T) {
		//Given
		var (
			id                 = 1
			expectedReport     = inboudorder.Report{}
			expectedStatusCode = http.StatusBadRequest
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"error": "Employee ID is required"}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("GenerateReport", mock.Anything, id).Return(expectedReport, errors.New("Employee ID is required"))
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/employees/reportInboundOrders"
		engine.GET(route, handler.GenerateReport())
		request, _ := http.NewRequest("GET", route, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		//service.AssertExpectations(t)
	})
	t.Run("error_invalid_id", func(t *testing.T) {
		//Given
		var (
			id                 = "dsasdds"
			expectedReport     = inboudorder.Report{}
			expectedStatusCode = http.StatusBadRequest
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"error": "Invalid employee ID"}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("GenerateReport", mock.Anything, id).Return(expectedReport, errors.New("Invalid employee ID"))
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/employees/reportInboundOrders"
		engine.GET(route, handler.GenerateReport())
		url := fmt.Sprintf("/employees/reportInboundOrders?id=%s", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		//service.AssertExpectations(t)
	})
	t.Run("employe_not_found", func(t *testing.T) {
		//Given
		var (
			id                 = 40
			expectedReport     = inboudorder.Report{}
			expectedStatusCode = http.StatusNotFound
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"message":"employee not found"}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("GenerateReport", mock.Anything, id).Return(expectedReport, inboudorder.ErrEmployeeNotFound)
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/employees/reportInboundOrders"
		engine.GET(route, handler.GenerateReport())
		url := fmt.Sprintf("/employees/reportInboundOrders?id=%d", id)
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
	t.Run("error", func(t *testing.T) {
		//Given
		var (
			id                 = 1
			expectedReport     = inboudorder.Report{}
			expectedStatusCode = http.StatusInternalServerError
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			expectedBody = `{"message":"internal error"}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("GenerateReport", mock.Anything, id).Return(expectedReport, errors.New("internal server error"))
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/employees/reportInboundOrders"
		engine.GET(route, handler.GenerateReport())
		url := fmt.Sprintf("/employees/reportInboundOrders?id=%d", id)
		request, _ := http.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		//service.AssertExpectations(t)
	})
}
func TestHandler_CreateInboundOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//Given
		var (
			id                   = 1
			expectedInboundOrder = domain.InboudOrder{
				OrderDate:      "2024-02-09",
				OrderNumber:    "order#1",
				EmployeeID:     4,
				ProductBatchID: 1,
				WarehouseID:    1,
			}
			expectedStatusCode = http.StatusCreated
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"order_date":"2024-02-09","order_number":"order#1","employee_id":4, "product_batch_id":1, "warehouse_id":1}`
			expectedBody = `{"data":{"id":1,"order_date":"2024-02-09","order_number":"order#1","employee_id":4, "product_batch_id":1, "warehouse_id":1}}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("CreateInboundOrder", mock.Anything, expectedInboundOrder).Return(id, nil)
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/inboundOrders"
		engine.POST(route, handler.CreateInboundOrder())
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
	t.Run("error_empty_values", func(t *testing.T) {
		var (
			id                   = 0
			expectedInboundOrder = domain.InboudOrder{
				OrderDate:      "2024-02-09",
				EmployeeID:     4,
				ProductBatchID: 1,
				WarehouseID:    1,
			}
			expectedStatusCode = http.StatusUnprocessableEntity
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"order_date":"2024-02-09","employee_id":4, "product_batch_id":1, "warehouse_id":1}`
			expectedBody = `{"message": "field order_number is empty"}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("CreateInboundOrder", mock.Anything, expectedInboundOrder).Return(id, errors.New("bad request"))
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/inboundOrders"
		engine.POST(route, handler.CreateInboundOrder())
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		// service.AssertExpectations(t)
	})
	t.Run("error_check_negative", func(t *testing.T) {
		var (
			id                   = 0
			expectedInboundOrder = domain.InboudOrder{
				OrderDate:      "2024-02-09",
				OrderNumber:    "order#1",
				EmployeeID:     4,
				ProductBatchID: 1,
				WarehouseID:    -1,
			}
			expectedStatusCode = http.StatusBadRequest
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest  = `{"order_date":"2024-02-09","order_number":"order#1","employee_id":4, "product_batch_id":1, "warehouse_id":-1}`
			expectedBody = `{"message": "negative warehouse_id"}`
		)
		service := &inboudorder.ServiceMock{}
		service.On("CreateInboundOrder", mock.Anything, expectedInboundOrder).Return(id, errors.New("negative warehouse_id"))
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/inboundOrders"
		engine.POST(route, handler.CreateInboundOrder())
		body := strings.NewReader(bodyRequest)
		request, _ := http.NewRequest("POST", route, body)
		response := httptest.NewRecorder()

		//When
		engine.ServeHTTP(response, request)

		//Then
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedHeaders, response.Header())
		assert.JSONEq(t, expectedBody, response.Body.String())
		//service.AssertExpectations(t)
	})
	t.Run("error_duplicate_inbound_order", func(t *testing.T) {
		var (
			id                   = 0
			expectedInboundOrder = domain.InboudOrder{
				OrderDate:      "2024-02-09",
				OrderNumber:    "order#1",
				EmployeeID:     4,
				ProductBatchID: 1,
				WarehouseID:    1,
			}
			expectedStatusCode = http.StatusConflict
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest   = `{"order_date":"2024-02-09","order_number":"order#1","employee_id":4, "product_batch_id":1, "warehouse_id":1}`
			expectedBody  = `{"message": "duplicate inbound order number"}`
			expectedError = inboudorder.ErrInboundOrderAlreadyExists
		)
		service := &inboudorder.ServiceMock{}
		service.On("CreateInboundOrder", mock.Anything, expectedInboundOrder).Return(id, expectedError)
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/inboundOrders"
		engine.POST(route, handler.CreateInboundOrder())
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
	t.Run("error_employee_does_not_exists", func(t *testing.T) {
		var (
			id                   = 0
			expectedInboundOrder = domain.InboudOrder{
				OrderDate:      "2024-02-09",
				OrderNumber:    "order#1",
				EmployeeID:     4,
				ProductBatchID: 1,
				WarehouseID:    1,
			}
			expectedStatusCode = http.StatusConflict
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest   = `{"order_date":"2024-02-09","order_number":"order#1","employee_id":4, "product_batch_id":1, "warehouse_id":1}`
			expectedBody  = `{"message": "Employee does not exists"}`
			expectedError = inboudorder.ErrEmployeeDoesNotExists
		)
		service := &inboudorder.ServiceMock{}
		service.On("CreateInboundOrder", mock.Anything, expectedInboundOrder).Return(id, expectedError)
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/inboundOrders"
		engine.POST(route, handler.CreateInboundOrder())
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
	t.Run("error_warehouse_does_not_exists", func(t *testing.T) {
		var (
			id                   = 0
			expectedInboundOrder = domain.InboudOrder{
				OrderDate:      "2024-02-09",
				OrderNumber:    "order#1",
				EmployeeID:     4,
				ProductBatchID: 1,
				WarehouseID:    1,
			}
			expectedStatusCode = http.StatusConflict
			expectedHeaders    = http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			}
			bodyRequest   = `{"order_date":"2024-02-09","order_number":"order#1","employee_id":4, "product_batch_id":1, "warehouse_id":1}`
			expectedBody  = `{"message": "Warehouse does not exists"}`
			expectedError = inboudorder.ErrWarehouseDoesNotExists
		)
		service := &inboudorder.ServiceMock{}
		service.On("CreateInboundOrder", mock.Anything, expectedInboundOrder).Return(id, expectedError)
		handler := NewInboudOrder(service)
		engine := gin.New()
		route := "/inboundOrders"
		engine.POST(route, handler.CreateInboundOrder())
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
}
