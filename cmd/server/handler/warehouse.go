package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/warehouse"
	"github.com/gin-gonic/gin"
)

type Warehouse struct {
	warehouseService warehouse.Service
}

func NewWarehouse(w warehouse.Service) *Warehouse {
	return &Warehouse{
		warehouseService: w,
	}
}

// ShowGet godoc
// @Summary Get the warehouses by ID, returns error if the warehouse doesn't exists.
// @Tags warehouses
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Param id path int true "Warehouse ID"
// @Router /warehouses/{id} [get]
func (w *Warehouse) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Check the id in the parameters
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		wh, err := w.warehouseService.Get(c, id)
		if err != nil {
			if errors.Is(err, warehouse.ErrNotFound) {
				c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not found"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
				return
			}
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"data": wh,
		})
	}
}

// ShowGetAll godoc
// @Summary Get all the warehouses available.
// @Tags warehouses
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /warehouses [get]
func (w *Warehouse) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouses, err := w.warehouseService.GetAll(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"data": warehouses,
		})
	}
}

// ShowCreate godoc
// @Summary Creates a warehouse, returns error if the warehouse doesn't match the standard input.
// @Tags warehouses
// @Produce json
// @Accept json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Param body body Warehouse true "Warehouse struct"
// @Router /warehouses [post]
func (w *Warehouse) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Convert the body request into a map for checking
		var body map[string]interface{}
		err := json.NewDecoder(c.Request.Body).Decode(&body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		//Validate that all fields are in the request
		err = ValidateBody(body, "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature")
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Invalid request body"})
			return
		}

		//Convert into bytes
		fmtBody, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		//Convert into the Warehouse struct
		var whouse domain.Warehouse
		json.Unmarshal(fmtBody, &whouse)
		id, err := w.warehouseService.Save(c, whouse)
		if err != nil {
			if errors.Is(err, warehouse.ErrDuplicateWarehouse) {
				c.JSON(http.StatusConflict, map[string]interface{}{"message": "Warehouse already exists"})
				return
			} else if errors.Is(err, warehouse.ErrIncorrectData) {
				c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Incorrect data"})
				return
			}

			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		whouse.ID = id
		c.JSON(http.StatusCreated, map[string]interface{}{
			"data": whouse,
		})
	}
}

// ShowUpdate godoc
// @Summary Update the warehouses by ID, returns error if the warehouse doesn't exists.
// @Tags warehouses
// @Produce json
// @Accept json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Param id path int true "Warehouse ID"
// @Router /warehouses/{id} [patch]
func (w *Warehouse) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		//Get the actual data of the warehouse
		warehouse, err := w.warehouseService.Get(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Warehouse not found"})
			return
		}

		//Bind only the new data in the body to the warehouse
		err = c.ShouldBindJSON(&warehouse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		//Update the warehouse in the database
		err = w.warehouseService.Update(c, warehouse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{"data": warehouse})
	}
}

// ShowUpdate godoc
// @Summary Delete the warehouses by ID, returns error if the warehouse doesn't exists.
// @Tags warehouses
// @Success 204 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Param id path int true "Warehouse ID"
// @Router /warehouses/{id} [delete]
func (w *Warehouse) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		err = w.warehouseService.Delete(c, id)
		if err != nil {
			if errors.Is(err, warehouse.ErrNotFound) {
				c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Warehouse not found"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
				return
			}
		}

		c.JSON(http.StatusNoContent, "")
	}
}

func ValidateBody(body map[string]interface{}, fields ...string) error {
	err := errors.New("invalid request body") // fix: Changed error string to lowercase
	//Check if all the fields exist
	for _, v := range fields {
		_, ok := body[v]
		if !ok {
			return err
		}
	}
	return nil
}
