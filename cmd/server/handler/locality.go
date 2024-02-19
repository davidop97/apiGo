package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/locality"
	"github.com/davidop97/apiGo/pkg/web"
	"github.com/gin-gonic/gin"
)

var (
	ErrLocalityNotFound      = "locality not found"
	ErrLocalityAlreadyExists = "locality already exists"
	ErrInvalidPostalCode     = "Invalid or missing 'postal_code'"
	ErrReportEmpty           = "Sellers Not found for the requested ID"
	ErrLocalityReport        = "Error getting the report for the requested ID. Id must be greater than 0"
	ErrInvalidJSON           = "invalid json"
)

type Locality struct {
	localityService locality.Service
}

func NewLocality(l locality.Service) *Locality {
	return &Locality{
		localityService: l,
	}
}

// @Summary Get a locality by id or an error if that id not exists.
// @Description Get a locality by id or an error if that id not exists or an internal error occurs.
// @Tags domain.Locality
// @Produce json
// @Success 200 {object} domain.Locality "Locality requested"
// @Failure 400 {object} string "invalid id"
// @Failure 404 {object} string "Locality not found"
// @Failure 500 {object} string "Server Internal error"
// @Param id path int true "id from the locality"
// @Router /locality/{id} [get]
func (l *Locality) GetLocalityById() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		//Check if the id is a valid integer.
		id, err := strconv.Atoi(idParam)
		if err != nil {
			//In case to be invalid, return a 400 status code and an error message.
			web.Error(c, http.StatusBadRequest, ErrInvalidID)
			return
		}
		//Check if the id is greater than 0.
		if id < 1 {
			//In case to be lower, return a 400 status code and an error message.
			web.Error(c, http.StatusBadRequest, ErrGreaterID)
			return
		}

		//Get the locality by id using the locality service.
		localityById, err := l.localityService.GetLocalityByID(c, id)
		//Check if an ErrNotFound error occurs and return a 404 status code.
		if err != nil {
			if errors.Is(err, locality.ErrLocalityNotFound) {
				web.Error(c, http.StatusNotFound, ErrLocalityNotFound)
			} else {
				//Or return a 500 status code error by default if an internal error occurs.
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}
		//If no errors occurs, return a 200 status code and the requested locality.
		web.Success(c, http.StatusOK, localityById)
	}
}

// @Summary Get all the localities available.
// @Description Get all the localities available or an error if the list is empty or an internal error occurs.
// @Tags domain.Locality
// @Produce json
// @Success 200 {array} domain.Locality "List of all localities"
// @Failure 404 {object} string "Localities not found"
// @Failure 500 {object} string "Server Internal error"
// @Router /localities [get]
func (l *Locality) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		allLocalities, err := l.localityService.GetAll(c)
		//Check if an ErrNotFound error occurs and return a 404 status code.
		if err != nil {
			if errors.Is(err, locality.ErrNoRows) {
				web.Error(c, http.StatusNotFound, "Localities not found")
			} else {
				//Or return a 500 status code error by default if an internal error occurs.
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}
		//If no errors occurs, return a 200 status code and the list of localities.
		web.Success(c, http.StatusOK, allLocalities)
	}
}

// @Summary Create a new locality or an error if that locality cannot be created.
// @Description Create a new locality or an error if that locality cannot be created or an internal error occurs.
// @Tags domain.Locality
// @Produce json
// @Success 201 {object} domain.Locality "New Locality created"
// @Failure 409 {object} string "locality already exists"
// @Failure 422 {object} string "invalid JSON"
// @Failure 422 {object} string "invalid or missing field"
// @Failure 500 {object} string "Server Internal error"
// @Param Locality body domain.Locality true "Struct of Locality domain"
// @Router /localities [post]
func (l *Locality) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get a copy of the body request for using it later in the process of creating the new locality with a map of strings and interfaces
		//to check the fields of the request and a struct of the domain for create the new locality.
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			//if cannot read the body, return a 500 status code and an error message.
			web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			return
		}

		var bodyJson map[string]interface{}
		//Unmarshal the body to a map of strings and interfaces.
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			//If an error occurs when unmarshal the body, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}
		//Check if the postal code is a valid integer greater than 0.
		postalCode, ok := bodyJson["postal_code"].(float64)
		if !ok || postalCode <= 0 || postalCode != float64(int(postalCode)) {
			//If not, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrInvalidPostalCode)
			return
		}
		//Only string fields of the structure.
		requiredStringFields := []string{
			"locality_name", "province_name", "country_name",
		}
		//Check if the requested string fields are string type data and not empty.
		for _, field := range requiredStringFields {
			value, ok := bodyJson[field].(string)
			//If not, return a 400 status code and an error message detailing which field is invalid or missing.
			if !ok || value == "" {
				web.Error(c, http.StatusUnprocessableEntity, fmt.Sprintf("Invalid or missing '%s'", field))
				return
			}
		}

		//Fields of the structure
		fields := []string{"postal_code", "locality_name", "province_name", "country_name"}

		// Check if all the fields are in the request using the fields of the structure.
		err = validExistenceOfField(bodyJson, fields...)
		if err != nil {
			//If an error occurs, return a 422 status code and an error message detailing which field is missing.
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		//Process of creating the new data.
		//Create a new structure with the data of the request to avoid expose sensitive data of the domain.
		localityRequest := domain.Locality{
			PostalCode:   int(bodyJson["postal_code"].(float64)),
			LocalityName: bodyJson["locality_name"].(string),
			ProvinceName: bodyJson["province_name"].(string),
			CountryName:  bodyJson["country_name"].(string),
		}

		err = json.Unmarshal(body, &localityRequest)
		if err != nil {
			//If an error occurs when unmarshal the struct, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}

		// Save the new locality
		newLocalityID, err := l.localityService.Save(c, localityRequest)
		if err != nil {
			if errors.Is(err, locality.ErrLocalityAlreadyExists) {
				web.Error(c, http.StatusConflict, ErrLocalityAlreadyExists)
			} else {
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}

		//Create a new structure with the data of the new locality created to avoid expose sensitive data of the domain
		//like the localityRequest does too.
		newLocalityCreated := domain.Locality{
			ID:           newLocalityID,
			PostalCode:   int(bodyJson["postal_code"].(float64)),
			LocalityName: bodyJson["locality_name"].(string),
			ProvinceName: bodyJson["province_name"].(string),
			CountryName:  bodyJson["country_name"].(string),
		}
		//Return a 201 status code and the new locality created.
		web.Success(c, http.StatusCreated, newLocalityCreated)
	}
}

// @Summary Get a report of sellers by locality.
// @Description Get a report of sellers by locality using its id. Or show all the sellers of all the localities if the id is not provided.
// @Tags domain.Locality
// @Produce json
// @Param id query int false "locality id"
// @Success 200 {object} domain.ReportSellers "Report of sellers by locality"
// @Failure 400 {object} string "invalid id"
// @Failure 400 {object} string "error getting the report for the requested ID. Id must be greater than 0"
// @Failure 404 {object} string "locality not found"
// @Failure 404 {object} string "sellers not found for the requested ID"
// @Failure 500 {object} string "server internal error"
// @Router /localities/reportSellers [get]
func (s *Locality) GetReportSellers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		idParam := c.Query("id")
		//Check if the id is a valid integer grater than 0.
		if idParam != "" {
			var err error
			id, err = strconv.Atoi(idParam)
			if err != nil {
				//In case to be invalid, return a 400 status code and an error message.
				web.Error(c, http.StatusBadRequest, ErrInvalidID)
				return
			}
			if id < 1 {
				//In case to lower than 1, return a 400 status code and an error message.
				web.Error(c, http.StatusBadRequest, ErrLocalityReport)
				return
			}
			//Check if the locality exists.
			_, err2 := s.localityService.GetLocalityByID(c, id)
			if err != nil {
				if errors.Is(err2, locality.ErrLocalityNotFound) {
					//If not, return a 404 status code and an error message.
					web.Error(c, http.StatusNotFound, ErrLocalityNotFound)
				} else {
					//Or return a 500 status code error by default if an internal error occurs.
					web.Error(c, http.StatusInternalServerError, ErrInternalServer)
				}
				return
			}
		}

		//Check if exists sellers in the locality.
		report, err := s.localityService.GetReportSellers(c, id)
		if err != nil {
			if errors.Is(err, locality.ErrNoRows) {
				//If there are no rows, return a 404 status code and an error message.
				web.Error(c, http.StatusNotFound, ErrReportEmpty)
			} else {
				//Or return a 500 status code error by default if an internal error occurs.
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}
		//If no errors occurs, return a 200 status code and the requested report.
		web.Success(c, http.StatusOK, report)
	}
}
