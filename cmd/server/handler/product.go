package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/product"
	"github.com/davidop97/apiGo/pkg/web"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidIDD      = "invalid id"
	ErrProductNotFound = "product not found"
	ErrInternalServer  = "internal server error"
	ErrInvalidDate     = "invalid date format"
	ProductDeleted     = "product deleted"
)

// Product struct represents a product handler.
type Product struct {
	// productService product.Service
	// productGroup *gin.RouterGroup
	service product.Service
}

// NewProduct creates a new instance of the Product struct.
func NewProduct(service product.Service) *Product {
	return &Product{
		service: service,
	}
}

// Ping handles the ping endpoint.
// @Summary Responds with a pong message.
// @Tags domain.Product
// @Tags ping
// @Tags domain.Product
// @Produce json
// @Success 200
// @Router /ping [get]
func (p *Product) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		web.Response(c, http.StatusOK, "pong")
	}
}

// GetAll handles the endpoint to retrieve all products.
// @Summary Retrieves a list of all products.
// @Tags products
// @Produce json
// @Success 200 {object} domain.Product "List of all products"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products [get]
func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := p.service.GetAll(c)
		if err != nil {
			web.Response(c, http.StatusInternalServerError, ErrInternalServer)
			return
		}

		web.Success(c, http.StatusOK, products)
	}
}

// Get handles the endpoint to retrieve a specific product by ID.
// @Summary Retrieves a product by ID.
// @Tags products
// @Produce json
// @Tags domain.Product
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product "Product data"
// @Failure 400 {object} string "Invalid ID"
// @Failure 404 {object} string "Product Not Found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products/{id} [get]
func (p *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Response(c, http.StatusBadRequest, ErrInvalidID)
			return
		}

		products, err := p.service.Get(c, id)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrNotFound):
				web.Response(c, http.StatusNotFound, ErrProductNotFound)
				return
			default:
				web.Response(c, http.StatusInternalServerError, ErrInternalServer)
				return
			}
		}

		web.Success(c, http.StatusOK, products)
	}
}

// Create handles the endpoint to create a new product.
// @Summary Creates a new product.
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product to be created"
// @Success 201 {object} domain.Product "Created product data"
// @Failure 400 {object} string "Invalid JSON"
// @Failure 422 {object} string "Invalid JSON"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products [post]
func (p *Product) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.Product

		// Check if JSON is valid
		err := c.ShouldBindJSON(&req)
		if err != nil {
			web.Response(c, http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}

		if !checkStruct(&req) {
			web.Response(c, http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}

		// Create product
		prodCreate, err := p.service.Save(c, req)
		//switch err and return error
		if err != nil {
			switch {
			case errors.Is(err, product.ErrProductCodeExists):
				web.Response(c, http.StatusConflict, err.Error())
				return
			default:
				web.Response(c, http.StatusInternalServerError, ErrInternalServer)
				return
			}
		}

		products := domain.Product{
			ID:             prodCreate,
			Description:    req.Description,
			ExpirationRate: req.ExpirationRate,
			FreezingRate:   req.FreezingRate,
			Height:         req.Height,
			Length:         req.Length,
			Netweight:      req.Netweight,
			ProductCode:    req.ProductCode,
			RecomFreezTemp: req.RecomFreezTemp,
			Width:          req.Width,
			ProductTypeID:  req.ProductTypeID,
			SellerID:       req.SellerID,
		}

		web.Success(c, http.StatusCreated, products)
	}
}

func checkStruct(product *domain.Product) bool {
	if product.Description == "" {
		return false
	}

	if product.ExpirationRate == 0 || product.ExpirationRate > 100 {
		return false
	}

	if product.FreezingRate == 0 || product.FreezingRate > 100 {
		return false
	}

	if product.Height == 0 {
		return false
	}

	if product.Length == 0 {
		return false
	}

	if product.Netweight == 0 {
		return false
	}

	if product.ProductCode == "" || len(product.ProductCode) > 100 {
		return false
	}

	if product.RecomFreezTemp == 0 || product.RecomFreezTemp > 100 {
		return false
	}

	if product.Width == 0 {
		return false
	}

	if product.ProductTypeID == 0 {
		return false
	}

	return true
}

// Update handles the endpoint to update an existing product by ID.
// @Summary Updates an existing product by ID.
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body domain.Product true "Updated product object"
// @Success 200 {object} domain.Product "Updated product data"
// @Failure 400 {object} string "Invalid ID"
// @Failure 404 {object} string "Product Not Found"
// @Failure 422 {object} string "Invalid JSON"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products/{id} [put]
func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			web.Response(c, http.StatusBadRequest, ErrInvalidID)
			return
		}

		//check if product exists
		products, err := p.service.Get(c, id)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrNotFound):
				web.Response(c, http.StatusNotFound, ErrProductNotFound)
				return
			default:
				web.Response(c, http.StatusInternalServerError, ErrInternalServer)
				return
			}
		}

		err = c.ShouldBindJSON(&products)
		if err != nil {
			web.Response(c, http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}

		update := domain.Product{
			ID:             id,
			Description:    products.Description,
			ExpirationRate: products.ExpirationRate,
			FreezingRate:   products.FreezingRate,
			Height:         products.Height,
			Length:         products.Length,
			Netweight:      products.Netweight,
			ProductCode:    products.ProductCode,
			RecomFreezTemp: products.RecomFreezTemp,
			Width:          products.Width,
			ProductTypeID:  products.ProductTypeID,
			SellerID:       products.SellerID,
		}

		err = p.service.Update(c, update)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrProductCodeExists):
				web.Response(c, http.StatusConflict, err.Error())
				return
			default:
				web.Response(c, http.StatusInternalServerError, ErrInternalServer)
				return
			}
		}

		web.Success(c, http.StatusOK, update)

	}
}

// Delete handles the endpoint to delete a product by ID.
// @Summary Deletes a product by ID.
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Product Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/{id} [delete]
func (p *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Response(c, http.StatusBadRequest, ErrInvalidID)
			return
		}

		err = p.service.Delete(c, id)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrNotFound):
				web.Response(c, http.StatusNotFound, ErrProductNotFound)
				return
			default:
				web.Response(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		web.Success(c, http.StatusNoContent, ProductDeleted)
	}
}

// CreateProductRecord handles the endpoint to create a new product record.
// @Summary Creates a new product record.
// @Tags productrecords
// @Accept json
// @Produce json
// @Param product body domain.ProductRecordCreate true "Product Record to be created"
// @Success 201 {object} domain.ProductRecord
// @Failure 422 {string} string "Invalid JSON"
// @Failure 500 {string} string "Internal Server Error"
// @Router /productrecord [post]
func (p *Product) CreateProductRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.ProductRecordCreate

		// Check if JSON is valid
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}

		// Check if all fields are valid
		if req.LastUpdate == "" || req.PurchasePrice <= 0 || req.SalePrice <= 0 || req.ProductID <= 0 {
			c.JSON(http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}

		//check if date is in format yyyy-mm-dd
		if req.LastUpdate != "" {
			_, err := time.Parse("2006-01-02", req.LastUpdate)
			if err != nil {
				web.Response(c, http.StatusUnprocessableEntity, ErrInvalidDate)
				return
			}
		}

		products, err := p.service.CreateProductRecord(c, req)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrNotFound):
				web.Response(c, http.StatusConflict, err.Error())
				return
			default:
				web.Response(c, http.StatusInternalServerError, ErrInternalServer)
				return
			}
		}

		// Create product updated with the id of the product record created
		update := domain.ProductRecord{
			ID:            products,
			LastUpdate:    req.LastUpdate,
			PurchasePrice: req.PurchasePrice,
			SalePrice:     req.SalePrice,
			ProductID:     req.ProductID,
		}

		web.Success(c, http.StatusCreated, update)
	}
}

// GetProductRecord handles the endpoint to retrieve product records by product ID.
// @Summary Retrieves product records by product ID or all product records if idProduct is 0.
// @Tags productrecords
// @Produce json
// @Param id query int false "Product ID"
// @Success 200 {array} domain.ProductRecordGet
// @Failure 404 {string} string "Product Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /productrecord [get]
func (p *Product) GetProductRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int // id = 0
		idParam := c.Query("id")

		// Check if id is in the query param
		if idParam != "" {
			var err error
			id, err = strconv.Atoi(idParam)
			if err != nil {
				web.Response(c, http.StatusBadRequest, ErrInvalidID)
				return
			}
		}

		products, err := p.service.GetProductRecord(c, id)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrNotFound):
				web.Response(c, http.StatusNotFound, err.Error())
				return
			default:
				web.Response(c, http.StatusInternalServerError, ErrInternalServer)
				return
			}
		}

		web.Success(c, http.StatusOK, products)
	}
}
