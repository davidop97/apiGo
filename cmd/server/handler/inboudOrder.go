package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/davidop97/apiGo/internal/domain"
	inboudorder "github.com/davidop97/apiGo/internal/inboudOrder"
	"github.com/gin-gonic/gin"
)

type InboudOrder struct {
	inboudOrderService inboudorder.Service
}

func NewInboudOrder(inboudOrderService inboudorder.Service) *InboudOrder {
	return &InboudOrder{
		inboudOrderService: inboudOrderService,
	}
}

// Body request parameters
var (
	// id             = "id"
	OrderDate      = "order_date"
	OrderNumber    = "order_number"
	EmployeeID     = "employee_id"
	ProductBatchID = "product_batch_id"
	WarehouseId    = "warehouse_id"
)

// EmployeeRequest is a struct that represents a body request for a employee
type InboudOrderRequest struct {
	ID             int    `json:"id"`
	OrderDate      string `json:"order_date"`
	OrderNumber    string `json:"order_number"`
	EmployeeID     int    `json:"employee_id"`
	ProductBatchID int    `json:"product_batch_id"`
	WarehouseID    int    `json:"warehouse_id"`
}

// @Summary Get all reports with inboudOrders
// @Tags inboundOrders
// @Produce json
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Success 500 {object} map[string]any
// @Success 200 {object} map[string]any
// @Router /employees/reportInboundOrder [get]
func (i *InboudOrder) GetAllReports() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener todos los informes
		reports, err := i.inboudOrderService.GetAllReports(c)
		if err != nil {
			// Manejar el error, por ejemplo, devolver un JSON con un mensaje de error y un código HTTP 500
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		// Construir la respuesta JSON
		response := gin.H{"data": reports}

		// Enviar la respuesta JSON
		c.JSON(http.StatusOK, response)
	}
}

// GetReport inbound orders per employee
// Summary Report inbound By employee
// @Description get inbound orders
// @Tags inboundOrders
// @Produce json
// @Param id query int false "Inbound Orders By Employee id"
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Success 500 {object} map[string]any
// @Success 200 {object} map[string]any
// @Router /employees/reportInboundOrders [get]
func (i *InboudOrder) GenerateReport() gin.HandlerFunc {
	return func(c *gin.Context) {

		employeeIDStr := c.Query("id")
		if employeeIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID is required"})
			return
		}
		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
			return
		}

		report, err := i.inboudOrderService.GenerateReport(c.Request.Context(), employeeID)
		if err != nil {
			switch {
			case errors.Is(err, inboudorder.ErrEmployeeNotFound):
				c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"data": report})

	}
}

// @Summary Create new inboundOrder
// @Tags inboundOrders
// @Produce json
// @Param body body InboudOrderRequest true "Inbound orders body"
// @Failure 400 {object} map[string]any
// @Failure 422 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Success 201 {object} map[string]any
// @Router /inboundOrders [post]
func (i *InboudOrder) CreateInboundOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error, imposible read all"})
			return
		}

		// - check if empty fields
		var bodyJson map[string]interface{}
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		err = checkEmptyValues(bodyJson, OrderNumber, EmployeeID, ProductBatchID, WarehouseID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		// - Create request struct
		req := InboudOrderRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}

		// - Set current date as the order date in the timezone "America/Bogota"
		loc, err := time.LoadLocation("America/Bogota")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error, unable to load timezone"})
			return
		}
		req.OrderDate = time.Now().In(loc).Format("2006-01-02")

		// - Check if negative IDs
		err = checkNegativeInbound(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// - Create a new section in the database
		NewInboudOrder := requestToInboundOrder(req)
		id, err := i.inboudOrderService.CreateInboundOrder(c, NewInboudOrder)
		if err != nil {
			switch {
			case errors.Is(err, inboudorder.ErrInboundOrderAlreadyExists):
				c.JSON(http.StatusConflict, gin.H{"message": "duplicate inbound order number"})
				return
			case errors.Is(err, inboudorder.ErrEmployeeDoesNotExists):
				c.JSON(http.StatusConflict, gin.H{"message": "Employee does not exists"})
				return
			case errors.Is(err, inboudorder.ErrWarehouseDoesNotExists):
				c.JSON(http.StatusConflict, gin.H{"message": "Warehouse does not exists"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error, impossible to creat a new inbound order"})
			}
		}

		// Response
		NewInboudOrder.ID = id
		c.JSON(http.StatusCreated, gin.H{"data": NewInboudOrder})
	}
}

func requestToInboundOrder(request InboudOrderRequest) (inboundOrder domain.InboudOrder) {
	inboundOrder.ID = request.ID
	inboundOrder.OrderDate = request.OrderDate
	inboundOrder.OrderNumber = request.OrderNumber
	inboundOrder.EmployeeID = request.EmployeeID
	inboundOrder.ProductBatchID = request.ProductBatchID
	inboundOrder.WarehouseID = request.WarehouseID

	return
}

// employeeToRequest creates a Request struct from a Section struct
func InboundOrderToRequest(inbound domain.InboudOrder) (request InboudOrderRequest) {
	request.ID = inbound.ID
	request.OrderDate = inbound.OrderDate
	request.OrderNumber = inbound.OrderNumber
	request.EmployeeID = inbound.EmployeeID
	request.ProductBatchID = inbound.ProductBatchID
	request.WarehouseID = inbound.WarehouseID
	return
}

// Checks wether negative IDS on struct
func checkNegativeInbound(req InboudOrderRequest) error {
	if req.WarehouseID <= 0 {
		return fmt.Errorf("negative warehouse_id")
	}
	return nil
}
