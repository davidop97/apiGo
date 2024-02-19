package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/section"
	"github.com/gin-gonic/gin"
)

// Body request parameter names for Sections
var (
	sectionNumber      = "section_number"
	currentTemperature = "current_temperature"
	minTemperature     = "minimum_temperature"
	currentCapacity    = "current_capacity"
	minCapacity        = "minimum_capacity"
	maxCapacity        = "maximum_capacity"
	warehouseID        = "warehouse_id"
	productTypeID      = "product_type_id"
)

// SectionRequest is a struct that represent a body request for a section
type SectionRequest struct {
	ID                 int `json:"id"`
	SectionNumber      int `json:"section_number"`
	CurrentTemperature int `json:"current_temperature"`
	MinimumTemperature int `json:"minimum_temperature"`
	CurrentCapacity    int `json:"current_capacity"`
	MinimumCapacity    int `json:"minimum_capacity"`
	MaximumCapacity    int `json:"maximum_capacity"`
	WarehouseID        int `json:"warehouse_id"`
	ProductTypeID      int `json:"product_type_id"`
}

// Section is a struct that contains handlers for section
type Section struct {
	sectionService section.Service
}

// NewSection is a function that returns a new instance of Section
func NewSection(s section.Service) *Section {
	return &Section{
		sectionService: s,
	}
}

// GetAll godoc
// @Summary Retrieves all sections.
// @Description Gets a list of all the sections.
// @Tags sections
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections [get]
func (s *Section) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		l, err := s.sectionService.GetAll(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": l})
	}
}

// Get godoc
// @Summary Retrieves a specific section.
// @Description Gets details of a specific section using its ID.
// @Tags sections
// @Produce json
// @Param id path int true "ID of the section item"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /sections/{id} [get]
func (s *Section) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad ID"})
			return
		}
		i, err := s.sectionService.Get(c, id)
		if err != nil {
			switch {
			case errors.Is(err, section.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"message": "section not found"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": i})
	}
}

// Create godoc
// @Summary Creates a new section.
// @Description Creates a new section based on the provided data.
// @Tags sections
// @Accept json
// @Produce json
// @Param sectionData body object true "Section data to create" format(json)
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /sections [post]
func (s *Section) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get a copy of the body request
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}

		// - check if empty fields
		var bodyJson map[string]interface{}
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			switch err.(type) {
			case *json.SyntaxError:
				c.JSON(http.StatusBadRequest, gin.H{"message": "bad request: invalid JSON syntax"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}
		err = checkCompleteSection(bodyJson, sectionNumber, currentTemperature, minTemperature, currentCapacity, minCapacity, maxCapacity, warehouseID, productTypeID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		// - create Request struct
		req := SectionRequest{}
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

		// - check if negative IDs
		err = checkSectionNegative(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Process
		// - create new section in the database
		sect := requestToSection(req)
		id, err := s.sectionService.Save(c, sect)
		if err != nil {
			switch {
			case errors.Is(err, section.ErrDuplicateSectNumber):
				c.JSON(http.StatusConflict, gin.H{"message": "duplicate section number"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
			return
		}

		//Response
		sect.ID = id
		c.JSON(http.StatusCreated, gin.H{"data": sect})
	}
}

// Update godoc
// @Summary Updates an existing section.
// @Description Updates an existing section based on the provided data.
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "ID of the section to update"
// @Param sectionData body object true "Updated section data" format(json)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections/{id} [patch]
func (s *Section) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get ID
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
			return
		}
		// - get original section
		res, err := s.sectionService.Get(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "section not found"})
			return
		}
		// - create Request strutct
		req := sectionToRequest(res)
		// - apply changes
		err = c.ShouldBindJSON(&req)
		if err != nil {
			switch e := err.(type) {
			case *json.SyntaxError:
				c.JSON(http.StatusBadRequest, gin.H{"message": "bad request: invalid JSON syntax"})
			case *json.UnmarshalTypeError:
				message := fmt.Sprintf("bad request: type %s was provided at %s field, %s was expected.", e.Value, e.Field, e.Type)
				c.JSON(http.StatusBadRequest, gin.H{"message": message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}
		req.ID = id
		// - check if negative IDs
		err = checkSectionNegative(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Process
		// - save changes
		sect := requestToSection(req)
		err = s.sectionService.Update(c, sect)
		if err != nil {
			switch {
			case errors.Is(err, section.ErrDuplicateSectNumber):
				c.JSON(http.StatusConflict, gin.H{"message": "duplicate section number"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
			return
		}

		// Response
		c.JSON(http.StatusOK, gin.H{"data": sect})
	}
}

// Delete godoc
// @Summary Deletes a section.
// @Description Deletes the section corresponding to the provided ID.
// @Tags sections
// @Produce json
// @Param id path int true "ID of the section to delete"
// @Success 204 "No content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections/{id} [delete]
func (s *Section) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get id from params
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
			return
		}

		// Process
		// - Delete item that matches given id
		err = s.sectionService.Delete(c, id)

		// Response
		if err != nil {
			switch {
			case errors.Is(err, section.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"message": "section not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// ProductCount godoc
// @Summary Gets number of products stored in a section.
// @Description Gets the total ammount of products stored in a given section. If section is not specified, gets ammount of products for each section.
// @Tags sections
// @Produce json
// @Param id query int false "ID of the specific section"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections/reportProducts [get]
func (s *Section) ProductCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get query parameter
		var id int
		idParam := c.Query("id")
		// - check if a section id was provided as a query parameter
		// -- if an id was provided, set id to the provided section id
		if idParam != "" {
			idValue, err := strconv.Atoi(idParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
				return
			}
			id = idValue
		}

		// Process
		l, err := s.sectionService.ProductCount(c, id)
		if err != nil {
			switch {
			case errors.Is(err, section.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"message": "section not found"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
				return
			}

		}

		// Response
		c.JSON(http.StatusOK, gin.H{"data": l})
	}
}

// checkEmpty checks wether a map cointains given keys.
// It returns an error if any of the specified keys are missing in the map.
func checkCompleteSection(m map[string]interface{}, params ...string) error {
	for _, param := range params {
		_, ok := m[param]
		if !ok {
			return fmt.Errorf("field %s is missing", param)
		}
	}
	return nil
}

// checkSectionNegative validates values that should be positive
func checkSectionNegative(req SectionRequest) error {
	switch {
	case req.ProductTypeID <= 0:
		return fmt.Errorf("negative product_type_id")
	case req.WarehouseID <= 0:
		return fmt.Errorf("negative warehouse_id")
	}
	return nil
}

// requestToSection creates a Section struct from a Request struct
func requestToSection(request SectionRequest) (section domain.Section) {
	section.ID = request.ID
	section.SectionNumber = request.SectionNumber
	section.CurrentTemperature = request.CurrentTemperature
	section.MinimumTemperature = request.MinimumTemperature
	section.CurrentCapacity = request.CurrentCapacity
	section.MinimumCapacity = request.MinimumCapacity
	section.MaximumCapacity = request.MaximumCapacity
	section.WarehouseID = request.WarehouseID
	section.ProductTypeID = request.ProductTypeID
	return
}

// sectionToRequest creates a Request struct from a Section struct
func sectionToRequest(section domain.Section) (request SectionRequest) {
	request.ID = section.ID
	request.SectionNumber = section.SectionNumber
	request.CurrentTemperature = section.CurrentTemperature
	request.MinimumTemperature = section.MinimumTemperature
	request.CurrentCapacity = section.CurrentCapacity
	request.MinimumCapacity = section.MinimumCapacity
	request.MaximumCapacity = section.MaximumCapacity
	request.WarehouseID = section.WarehouseID
	request.ProductTypeID = section.ProductTypeID
	return
}
