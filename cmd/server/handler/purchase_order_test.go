package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/purchase_order"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests for Purchase Order handler

// TestHandler_CreatePurchaseOrder
// Test all cases for Create a new purchase order
// Test methods: Create
func TestHandler_CreatePurchaseOrder(t *testing.T) {
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_ok
	// DESCRIPTION: When the data entry is successful, a 201 code will be returned along with the object entered.
	t.Run("it should create a new purchase order", func(t *testing.T) {
		// arrange
		expectedStatusCode := http.StatusCreated
		expectedPurchaseOrder := `{"data":{"id":1,"order_number":"123","order_date":"2022-01-01","tracking_code":"abc","buyer_id":1,"product_record_id":5,"order_status_id":10}}`
		purchaseOrderToCreate := RequestBodyPurchaseCreate{
			OrderNumber:     "123",
			OrderDate:       "2022-01-01",
			TrackingCode:    "abc",
			BuyerID:         1,
			ProductRecordID: 5,
			OrderStatusID:   10,
		}
		// convert purchaseOrderToCreate to json so it can be used as a body request
		jsonPayload, _ := json.Marshal(purchaseOrderToCreate)

		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("Save", mock.Anything, mock.Anything).Return(1, nil)

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(jsonPayload)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedPurchaseOrder, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_incorrect_json_format
	// DESCRIPTION: If json format is not valid, it returns 400 status code.
	t.Run("it should return error when body request is malformed", func(t *testing.T) {
		// arrange
		invalidJSON := `{"some_field": "value"`
		// prepare expected results
		expectedMessageError := `{"error":"Bad request"}`
		expectedStatusCode := http.StatusBadRequest

		// create mock of the service
		mockService := purchase_order.ServiceMock{}

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(invalidJSON)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_missing_fields
	// DESCRIPTION: If body request does not have all requiere fields, it returns 422 status code.
	t.Run("it should return error when missing fields", func(t *testing.T) {
		// arrange
		invalidJSON := `{"some_field": "value"}`
		// prepare expected results
		expectedMessageError := `{"error":"Missing fields"}`
		expectedStatusCode := http.StatusUnprocessableEntity

		// create mock of the service
		mockService := purchase_order.ServiceMock{}

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(invalidJSON)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_incorrect_date_field
	// DESCRIPTION: If order date has incorrect format, it returns 422 status code.
	t.Run("it should return error when date field is in invalid format", func(t *testing.T) {
		// arrange
		purchaseOrderToCreate := RequestBodyPurchaseCreate{
			OrderNumber:     "123",
			OrderDate:       "01-01-2022",
			TrackingCode:    "abc",
			BuyerID:         1,
			ProductRecordID: 5,
			OrderStatusID:   10,
		}
		// convert purchaseOrderToCreate to json so it can be used as a body request
		jsonPayload, _ := json.Marshal(purchaseOrderToCreate)
		// prepare expected results
		expectedMessageError := `{"error":"Invalid date format"}`
		expectedStatusCode := http.StatusUnprocessableEntity

		// create mock of the service
		mockService := purchase_order.ServiceMock{}

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(jsonPayload)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_purchase_order_already_exists
	// DESCRIPTION: If purchase order already exists, it returns 409 status code.
	t.Run("it should return error when purchase order already exists", func(t *testing.T) {
		// arrange
		expectedStatusCode := http.StatusConflict
		expectedMessageError := `{"error":"Purchase order already exists"}`
		purchaseOrderToCreate := RequestBodyPurchaseCreate{
			OrderNumber:     "123",
			OrderDate:       "2022-01-01",
			TrackingCode:    "abc",
			BuyerID:         1,
			ProductRecordID: 5,
			OrderStatusID:   10,
		}
		// convert purchaseOrderToCreate to json so it can be used as a body request
		jsonPayload, _ := json.Marshal(purchaseOrderToCreate)

		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("Save", mock.Anything, mock.Anything).Return(0, purchase_order.ErrPurchaseOrderAlreadyExists)

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(jsonPayload)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_buyer_does_not_exists
	// DESCRIPTION: If buyer does not exists, it returns 409 status code.
	t.Run("it should return error when buyer does not exists", func(t *testing.T) {
		// arrange
		expectedStatusCode := http.StatusConflict
		expectedMessageError := `{"error":"Buyer id not found"}`
		purchaseOrderToCreate := RequestBodyPurchaseCreate{
			OrderNumber:     "123",
			OrderDate:       "2022-01-01",
			TrackingCode:    "abc",
			BuyerID:         1,
			ProductRecordID: 5,
			OrderStatusID:   10,
		}
		// convert purchaseOrderToCreate to json so it can be used as a body request
		jsonPayload, _ := json.Marshal(purchaseOrderToCreate)

		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("Save", mock.Anything, mock.Anything).Return(0, purchase_order.ErrBuyerIDNotExists)

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(jsonPayload)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_product_record_does_not_exists
	// DESCRIPTION: If product record not exists, it returns 409 status code.
	t.Run("it should return error when product record does not exists", func(t *testing.T) {
		// arrange
		expectedStatusCode := http.StatusConflict
		expectedMessageError := `{"error":"Product record id not found"}`
		purchaseOrderToCreate := RequestBodyPurchaseCreate{
			OrderNumber:     "123",
			OrderDate:       "2022-01-01",
			TrackingCode:    "abc",
			BuyerID:         1,
			ProductRecordID: 5,
			OrderStatusID:   10,
		}
		// convert purchaseOrderToCreate to json so it can be used as a body request
		jsonPayload, _ := json.Marshal(purchaseOrderToCreate)

		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("Save", mock.Anything, mock.Anything).Return(0, purchase_order.ErrProductsRecordIDNotExits)

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(jsonPayload)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: error_internal_server
	// DESCRIPTION: If a internal server error occurs creating a new purchase order, it returns 500 status code.
	t.Run("it should return error when internal server occurs processing", func(t *testing.T) {
		// arrange
		expectedStatusCode := http.StatusInternalServerError
		expectedMessageError := `{"error":"Internal server error"}`
		purchaseOrderToCreate := RequestBodyPurchaseCreate{
			OrderNumber:     "123",
			OrderDate:       "2022-01-01",
			TrackingCode:    "abc",
			BuyerID:         1,
			ProductRecordID: 5,
			OrderStatusID:   10,
		}
		// convert purchaseOrderToCreate to json so it can be used as a body request
		jsonPayload, _ := json.Marshal(purchaseOrderToCreate)

		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("Save", mock.Anything, mock.Anything).Return(0, errors.New("internal server error"))

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/purchaseOrders"

		// Add handler to the router
		r.POST(route, handler.Create())

		// create the request
		// create the body to the request
		reqBodyBytes := bytes.NewBuffer([]byte(string(jsonPayload)))
		request, _ := http.NewRequest("POST", route, reqBodyBytes)
		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})
}

// TestHandler_ReportPurchaseOrdersByBuyer
// Test all cases for report Purchase Orders By Buyer
// Test methods: ReportPurchaseOrdersByBuyer
func TestHandler_ReportPurchaseOrdersByBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: ReportPurchaseOrdersByBuyer
	// EDGE CASE: report_ok
	// DESCRIPTION: When report is successfully created, 200 status code will be returned along with the created report.
	t.Run("should return a report with purchase Orders By Buyer ", func(t *testing.T) {
		// arrange
		buyerID := 1
		expectedStatusCode := http.StatusOK
		expectedPurchaseOrder := `{"data":[{"id":1,"card_number_id":"123","first_name":"John","last_name":"Doe","purchase_orders_count":3},{"id":2,"card_number_id":"456","first_name":"Jane","last_name":"Smith","purchase_orders_count":0}]}`

		purchaseOrderReport := []domain.PurchaseOrdersByBuyer{
			{
				ID:                  1,
				CardNumberID:        "123",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 3,
			},
			{
				ID:                  2,
				CardNumberID:        "456",
				FirstName:           "Jane",
				LastName:            "Smith",
				PurchaseOrdersCount: 0,
			},
		}
		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("PurchaseOrdersByBuyer", mock.Anything, buyerID).Return(purchaseOrderReport, nil)

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/reportPurchaseOrders/:id"

		// Add handler to the router
		r.GET(route, handler.ReportPurchaseOrdersByBuyer())

		// create the request
		url := fmt.Sprintf("/api/v1/buyers/reportPurchaseOrders/%d", buyerID)
		request, _ := http.NewRequest("GET", url, nil)

		// Add the "id" query parameter
		query := request.URL.Query()
		query.Add("id", strconv.Itoa(buyerID))
		request.URL.RawQuery = query.Encode()

		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedPurchaseOrder, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: ReportPurchaseOrdersByBuyer
	// EDGE CASE: error_buyer_id_not_found
	// DESCRIPTION: When the buyer id not found in the system, status code 404 is returned.
	t.Run(" should fail with 404 status code when buyer id not found ", func(t *testing.T) {
		buyerID := 1
		expectedStatusCode := http.StatusNotFound
		expectedMessageError := `{"error":"Buyer id not found"}`

		purchaseOrderReport := []domain.PurchaseOrdersByBuyer{}
		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("PurchaseOrdersByBuyer", mock.Anything, buyerID).Return(purchaseOrderReport, purchase_order.ErrBuyerIDNotExists)

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/reportPurchaseOrders/:id"

		// Add handler to the router
		r.GET(route, handler.ReportPurchaseOrdersByBuyer())

		// create the request
		url := fmt.Sprintf("/api/v1/buyers/reportPurchaseOrders/%d", buyerID)
		request, _ := http.NewRequest("GET", url, nil)

		// Add the "id" query parameter
		query := request.URL.Query()
		query.Add("id", strconv.Itoa(buyerID))
		request.URL.RawQuery = query.Encode()

		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: ReportPurchaseOrdersByBuyer
	// EDGE CASE: error_internal_server
	// DESCRIPTION: When is a problem extracting report from the system, status code 500 is returned.
	t.Run("should fail with 500 status code when trying to get report", func(t *testing.T) {
		buyerID := 1
		expectedStatusCode := http.StatusInternalServerError
		expectedMessageError := `{"error":"Internal server error"}`

		purchaseOrderReport := []domain.PurchaseOrdersByBuyer{}
		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("PurchaseOrdersByBuyer", mock.Anything, buyerID).Return(purchaseOrderReport, errors.New("some errors"))

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/reportPurchaseOrders/:id"

		// Add handler to the router
		r.GET(route, handler.ReportPurchaseOrdersByBuyer())

		// create the request
		url := fmt.Sprintf("/api/v1/buyers/reportPurchaseOrders/%d", buyerID)
		request, _ := http.NewRequest("GET", url, nil)

		// Add the "id" query parameter
		query := request.URL.Query()
		query.Add("id", strconv.Itoa(buyerID))
		request.URL.RawQuery = query.Encode()

		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.JSONEq(t, expectedMessageError, response.Body.String())
		mockService.AssertExpectations(t)
	})

	// ASSOCIATED USER STORY: ReportPurchaseOrdersByBuyer
	// EDGE CASE: empty_report
	// DESCRIPTION: When no data found for a buyer id, status code 204 is returned.
	t.Run("should return 204 status code when no data was found and there are no errors", func(t *testing.T) {
		buyerID := 1
		expectedStatusCode := http.StatusNoContent
		expectedPurchaseOrder := ""

		purchaseOrderReport := []domain.PurchaseOrdersByBuyer{}
		//jsonReport, _ := json.Marshal(purchaseOrderReport)
		// create mock of the service
		mockService := purchase_order.ServiceMock{}
		// set service mock with the expected value
		mockService.On("PurchaseOrdersByBuyer", mock.Anything, buyerID).Return(purchaseOrderReport, nil)

		// create handler using mock service
		handler := NewPurchaseOrder(&mockService)

		// prepare server to call endpoint GetAll
		// create gin router
		r := gin.New()
		route := "/api/v1/buyers/reportPurchaseOrders/:id"

		// Add handler to the router
		r.GET(route, handler.ReportPurchaseOrdersByBuyer())

		// create the request
		url := fmt.Sprintf("/api/v1/buyers/reportPurchaseOrders/%d", buyerID)
		request, _ := http.NewRequest("GET", url, nil)

		// Add the "id" query parameter
		query := request.URL.Query()
		query.Add("id", strconv.Itoa(buyerID))
		request.URL.RawQuery = query.Encode()

		// create the response recorder
		response := httptest.NewRecorder()

		// act
		r.ServeHTTP(response, request)

		// assert
		// assert results
		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, expectedPurchaseOrder, response.Body.String())
		mockService.AssertExpectations(t)
	})
}
