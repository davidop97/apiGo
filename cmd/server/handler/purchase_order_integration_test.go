package handler

import (
	"bytes"
	"encoding/json"
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

// TestIntegrationHandler_CreatePurchaseOrder
// Test all cases for Create a new purchase order
// Test methods: Create
func TestIntegrationHandler_CreatePurchaseOrder(t *testing.T) {
	// ASSOCIATED USER STORY: CREATE
	// EDGE CASE: create_ok
	// DESCRIPTION: When the data entry is successful, a 201 code will be returned along with the object entered.
	t.Run("it should create a new purchase order", func(t *testing.T) {
		// arrange

		// expected result
		expectedPurchaseOrderID := 1
		buyerID := 1
		productRecordID := 5
		poToSave := domain.PurchaseOrder{
			ID:              expectedPurchaseOrderID,
			OrderNumber:     "123",
			OrderDate:       "2024-01-01",
			TrackingCode:    "abc",
			BuyerID:         buyerID,
			ProductRecordID: productRecordID,
			OrderStatusID:   10,
		}

		expectedStatusCode := http.StatusCreated
		expectedPurchaseOrder := `{"data":{"id":1,"order_number":"123","order_date":"2024-01-01","tracking_code":"abc","buyer_id":1,"product_record_id":5,"order_status_id":10}}`
		purchaseOrderToCreate := RequestBodyPurchaseCreate{
			ID:              expectedPurchaseOrderID,
			OrderNumber:     "123",
			OrderDate:       "2024-01-01",
			TrackingCode:    "abc",
			BuyerID:         buyerID,
			ProductRecordID: productRecordID,
			OrderStatusID:   10,
		}
		// convert purchaseOrderToCreate to json so it can be used as a body request
		jsonPayload, _ := json.Marshal(purchaseOrderToCreate)

		// mock repo
		// create mock for repository
		repoMock := &purchase_order.RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsPurchaseOrder", mock.Anything, 1).Return(false)
		repoMock.On("ExistsBuyer", mock.Anything, buyerID).Return(true)
		repoMock.On("ExistsProductsRecord", mock.Anything, productRecordID).Return(true)
		repoMock.On("Save", mock.Anything, poToSave).Return(expectedPurchaseOrderID, nil)

		// create service with mock of the service
		service := purchase_order.NewService(repoMock)

		// create handler using real service
		handler := NewPurchaseOrder(service)

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
		repoMock.AssertExpectations(t)
	})

}

// TestHandler_ReportPurchaseOrdersByBuyer
// Test all cases for report Purchase Orders By Buyer
// Test methods: ReportPurchaseOrdersByBuyer
func TestIntegrationHandler_ReportPurchaseOrdersByBuyer(t *testing.T) {
	// ASSOCIATED USER STORY: ReportPurchaseOrdersByBuyer
	// EDGE CASE: report_ok
	// DESCRIPTION: When report is successfully created, 200 status code will be returned along with the created report.
	t.Run("should return a report with purchase Orders By Buyer ", func(t *testing.T) {
		// arrange
		buyerID := 1
		expectedStatusCode := http.StatusOK
		expectedPurchaseOrder := `{"data":[{"id":1,"card_number_id":"123","first_name":"John","last_name":"Doe","purchase_orders_count":3},{"id":2,"card_number_id":"456","first_name":"Jane","last_name":"Smith","purchase_orders_count":0}]}`

		expectReport := []domain.PurchaseOrdersByBuyer{
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

		// create mock for repository
		repoMock := &purchase_order.RepositoryMock{}
		// set mock with the expected value
		repoMock.On("ExistsBuyer", mock.Anything, buyerID).Return(true)
		repoMock.On("PurchaseOrdersByBuyers", mock.Anything, buyerID).Return(expectReport, nil)

		// create service with mock of the service
		service := purchase_order.NewService(repoMock)

		// create handler using real service
		handler := NewPurchaseOrder(service)

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
		repoMock.AssertExpectations(t)
	})

}
