package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/davidop97/apiGo/internal/batch"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/gin-gonic/gin"
)

// Product batch request field names
var (
	batchNumber         = "batch_number"
	batchCurQuantity    = "current_quantity"
	batchCurTemperature = "current_temperature"
	batchDueDate        = "due_date"
	batchInitQuantity   = "initial_quantity"
	batchManufDate      = "manufacturing_date"
	batchManufHour      = "manufacturing_hour"
	batchMinTemperature = "minimum_temperature"
	batchProdID         = "product_id"
	batchSectID         = "section_id"
)

// BatchRequests is a struct that represents body request for ProductBatch
type BatchRequest struct {
	BatchNumber        int    `json:"batch_number"`
	CurrentQuantity    int    `json:"current_quantity"`
	CurrentTemperature int    `json:"current_temperature"`
	DueDate            string `json:"due_date"`
	InitialQuantity    int    `json:"initial_quantity"`
	ManufacturingDate  string `json:"manufacturing_date"`
	ManufacturingHour  int    `json:"manufacturing_hour"`
	MinimumTemperature int    `json:"minimum_temperature"`
	ProductID          int    `json:"product_id"`
	SectionID          int    `json:"section_id"`
}

// productBatch is a struct that contains handler functions for ProductBatches
type ProductBatch struct {
	batchService batch.Service
}

// NewProductBatch returns a new instance of ProductBatch
func NewProductBatch(b batch.Service) *ProductBatch {
	return &ProductBatch{
		batchService: b,
	}
}

// GetAll godoc
// @Summary Retrieves all product batches.
// @Description Gets a list of all product batches.
// @Tags productBatches
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /productBatches [get]
func (b *ProductBatch) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		l, err := b.batchService.GetAll(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": l})
	}
}

// Create godoc
// @Summary Creates a new product batch.
// @Description Creates a new product batch based on the provided data.
// @Tags productBatches
// @Accept json
// @Produce json
// @Param sectionData body object true "Product batch data to create" format(json)
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /productBatches [post]
func (b *ProductBatch) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get a copy of the body request
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		// - check if missing fields
		var bodyJson map[string]interface{}
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			switch e := err.(type) {
			case *json.SyntaxError:
				message := fmt.Sprintf("bad request: invalid syntax at position %v", e.Offset)
				c.JSON(http.StatusBadRequest, gin.H{"message": message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}
		ok, field := batchComplete(bodyJson, batchNumber, batchCurQuantity, batchCurTemperature, batchDueDate, batchInitQuantity, batchManufDate, batchManufHour, batchMinTemperature, batchProdID, batchSectID)
		if !ok {
			message := fmt.Sprintf("bad request: field %s is missing", field)
			c.JSON(http.StatusBadRequest, gin.H{"message": message})
			return
		}

		// - create request struct
		req := BatchRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			switch e := err.(type) {
			case *json.UnmarshalTypeError:
				message := fmt.Sprintf("bad request: type %s was provided at %s field, %s was expected.", e.Value, e.Field, e.Type)
				c.JSON(http.StatusBadRequest, gin.H{"message": message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}

		// - check if fields are valid
		ok, errorMessage := validBatchRequest(req)
		if !ok {
			message := fmt.Sprintf("bad request: %s", errorMessage)
			c.JSON(http.StatusBadRequest, gin.H{"message": message})
			return
		}

		// Process
		// - create new product batch in the database
		productBatch := requestToBatch(req)
		id, err := b.batchService.Save(c, productBatch)
		if err != nil {
			switch {
			case errors.Is(err, batch.ErrDuplicateBatchNumber):
				c.JSON(http.StatusConflict, gin.H{"message": "batch number must be unique, provided already exists"})
			case errors.Is(err, batch.ErrProductNotFound):
				c.JSON(http.StatusConflict, gin.H{"message": "can't create batch, provided product id was not found"})
			case errors.Is(err, batch.ErrSectionNotFound):
				c.JSON(http.StatusConflict, gin.H{"message": "can't create batch, provided section id was not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}
		productBatch.ID = id
		c.JSON(http.StatusCreated, gin.H{"data": productBatch})
	}
}

// batchComplete is an auxiliary function that checks if a body request (as a map) include required fields
func batchComplete(m map[string]interface{}, fields ...string) (complete bool, missing string) {
	for _, field := range fields {
		_, ok := m[field]
		if !ok {
			missing = field
			return
		}
	}
	complete = true
	return
}

// validBatchRequest checks if a body request's fields are valid
func validBatchRequest(b BatchRequest) (ok bool, message string) {
	// Check if values are positive
	switch {
	case b.CurrentQuantity < 0:
		message = "current_quantity must be equal or greater than 0"
		return
	case b.InitialQuantity < 0:
		message = "initial_quantity must be equal or greater than 0"
		return
	case b.ManufacturingHour < 0 || b.ManufacturingHour > 23:
		message = "manufacturing_hour value must be within range [0 - 23]"
		return
	case b.ProductID < 0:
		message = "product_id must be equal or greater than 0"
		return
	case b.SectionID < 0:
		message = "section_id must be equal or greater than 0"
		return
	}

	// Check if date format is correct
	dateFormat := "2006-01-02"
	_, err := time.Parse(dateFormat, b.DueDate)
	if err != nil {
		message = "due_date should match format YYYY-MM-DD"
		return
	}
	_, err = time.Parse(dateFormat, b.ManufacturingDate)
	if err != nil {
		message = "manufacturing_date should match format YYYY-MM-DD"
		return
	}

	ok = true
	return
}

// requestToBatch is an auxiliary function that takes a Request and returns its corresponding ProductBatch domain struct
func requestToBatch(req BatchRequest) (b domain.ProductBatch) {
	b.BatchNumber = req.BatchNumber
	b.CurrentQuantity = req.CurrentQuantity
	b.CurrentTemperature = req.CurrentTemperature
	b.DueDate = req.DueDate
	b.InitialQuantity = req.InitialQuantity
	b.ManufacturingDate = req.ManufacturingDate
	b.ManufacturingHour = req.ManufacturingHour
	b.MinimumTemperature = req.MinimumTemperature
	b.ProductID = req.ProductID
	b.SectionID = req.SectionID
	return
}
