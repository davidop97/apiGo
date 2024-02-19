package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/purchase_order"
	"github.com/gin-gonic/gin"
)

type PurchaseOrder struct {
	poService purchase_order.Service
}

func NewPurchaseOrder(poService purchase_order.Service) *PurchaseOrder {
	return &PurchaseOrder{poService: poService}
}

type RequestBodyPurchaseCreate struct {
	ID              int    `json:"user_id"`
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordID int    `json:"product_record_id"`
	OrderStatusID   int    `json:"order_status_id"`
}

// ShowCreatePurchaseOrders godoc
// Summary Create a new purchase order
// @Description create a new purchase order
// @Tags purchase_orders
// @Produce json
// @Param body body RequestBodyPurchaseCreate true "Purchase orders body"
// @Failure 400 {object} map[string]any
// @Failure 422 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Success 201 {object} map[string]any
// @Router /purchaseOrders [post]
func (po *PurchaseOrder) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request
		var reqBody RequestBodyPurchaseCreate
		// check if body is in the correct json format
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		// process

		// check if all fields come into reqBody
		valid, err := validateEmptys(&reqBody)
		if !valid || err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Missing fields"})
			return
		}
		// validate date
		valid = validateDate(reqBody.OrderDate)
		if !valid {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid date format"})
			return
		}

		// process
		// serialize the purchase order
		purchase := domain.PurchaseOrder{
			ID:              reqBody.ID,
			OrderNumber:     reqBody.OrderNumber,
			OrderDate:       reqBody.OrderDate,
			TrackingCode:    reqBody.TrackingCode,
			BuyerID:         reqBody.BuyerID,
			ProductRecordID: reqBody.ProductRecordID,
			OrderStatusID:   reqBody.OrderStatusID,
		}

		// - save
		// call to purchase order service to save a new purchase
		id, err := po.poService.Save(c, purchase)
		if err != nil {
			switch {
			case errors.Is(err, purchase_order.ErrPurchaseOrderAlreadyExists):
				c.JSON(http.StatusConflict, gin.H{"error": "Purchase order already exists"})
			case errors.Is(err, purchase_order.ErrBuyerIDNotExists):
				c.JSON(http.StatusConflict, gin.H{"error": "Buyer id not found"})
			case errors.Is(err, purchase_order.ErrProductsRecordIDNotExits):
				c.JSON(http.StatusConflict, gin.H{"error": "Product record id not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}
		// return response
		purchase.ID = id
		c.JSON(http.StatusCreated, gin.H{"data": purchase})
	}
}

// GetReport purchase orders per buyer
// Summary Report Purchase Orders By Buyer
// @Description get a buyer
// @Tags purchase_orders
// @Produce json
// @Param id query int false "Purchase Orders By Buyer id"
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Success 204 {object} map[string]any
// @Success 200 {object} map[string]any
// @Router /buyers/reportPurchaseOrders/ [get]
func (po *PurchaseOrder) ReportPurchaseOrdersByBuyer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request
		// check if buyer id come into param
		buyerIDParam := c.Query("id")
		if buyerIDParam == "" {
			// in the case id is not present in the query, it is set to zero
			buyerIDParam = "0"
		}
		buyerID, err := strconv.Atoi(buyerIDParam)
		if err != nil || buyerID < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		// process
		reports, err := po.poService.PurchaseOrdersByBuyer(c, buyerID)
		if err != nil {
			switch {
			case errors.Is(err, purchase_order.ErrBuyerIDNotExists):
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer id not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}
		// check if reports has content
		if len(reports) == 0 {
			c.JSON(http.StatusNoContent, gin.H{})
			return
		}

		// return response
		c.JSON(http.StatusOK, gin.H{"data": reports})
	}

}

// validateEmptys check fields
func validateEmptys(product *RequestBodyPurchaseCreate) (bool, error) {
	switch {
	case product.OrderNumber == "" || product.OrderDate == "" || product.TrackingCode == "":
		return false, errors.New("fields can't be empty")
	case product.BuyerID <= 0 || product.ProductRecordID <= 0 || product.OrderStatusID <= 0:
		if product.BuyerID <= 0 {
			return false, errors.New("buyer id must be greater than 0")
		}
		if product.ProductRecordID <= 0 {
			return false, errors.New("product record id must be greater than 0")
		}
		if product.OrderStatusID <= 0 {
			return false, errors.New("order status id must be greater than 0")
		}
	}
	return true, nil
}

// validateDate check if date has correct format
func validateDate(dateString string) bool {
	_, err := time.Parse("2006-01-02", dateString)
	return err == nil
}
