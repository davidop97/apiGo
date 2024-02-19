package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/seller"
	"github.com/davidop97/apiGo/pkg/web"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidID           = "invalid id"
	ErrInvalidSeller       = "invalid seller"
	ErrFieldsCannotBeEmpty = "fields can not be empty. CID must be 1 or greater"
	ErrSellerNotFound      = "seller not found"
	ErrGreaterID           = "id must be 1 or greater"
	ErrSellerAlreadyExists = "seller already exists"
	ErrIdLocalityNotExists = "id locality not exists"
	ErrLocality            = "invalid or missing locality. Locality_id must be 1 or greater"
	ErrCID                 = "invalid or missing cid. CID must be 1 or greater"
	ErrLocalityIDGreater   = "locality_id must be 1 or greater"
	ErrInvalidLocalityID   = "invalid locality_id"
)

type Seller struct {
	sellerService seller.Service
}

func NewSeller(s seller.Service) *Seller {
	return &Seller{
		sellerService: s,
	}
}

// @Summary Get all the sellers available.
// @Description Get all the sellers available or an error if the list is empty or an internal error occurs.
// @Tags domain.Seller
// @Produce json
// @Success 200 {array} domain.Seller "List of all sellers"
// @Failure 404 {object} string "Sellers not found"
// @Failure 500 {object} string "Server Internal error"
// @Router /seller [get]
func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		allSeller, err := s.sellerService.GetAllSellers(c)
		//Check if an ErrNotFound error occurs and return a 404 status code.
		if err != nil {
			if errors.Is(err, seller.ErrNotFound) {
				web.Error(c, http.StatusNotFound, "Sellers not found")
			} else {
				//Or return a 500 status code error by default if an internal error occurs.
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}
		//If no errors occurs, return a 200 status code and the list of sellers.
		web.Success(c, http.StatusOK, allSeller)
	}
}

// @Summary Get a seller by id or an error if that id not exists.
// @Description Get a seller by id or an error if that id not exists or an internal error occurs.
// @Tags domain.Seller
// @Produce json
// @Success 200 {object} domain.Seller "Seller requested"
// @Failure 400 {object} string "invalid id"
// @Failure 404 {object} string "Seller not found"
// @Failure 500 {object} string "Server Internal error"
// @Param id path int true "id from the seller"
// @Router /seller/{id} [get]
func (s *Seller) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		//Check if the id is a valid integer.
		id, err := strconv.Atoi(idParam)
		if err != nil {
			//In case to be invalid, return a 400 status code and an error message.
			web.Error(c, http.StatusBadRequest, ErrInvalidID)
			return
		}
		sellerById, err := s.sellerService.GetSellerByID(c, id)
		//Check if an ErrNotFound error occurs and return a 404 status code.
		if err != nil {
			if errors.Is(err, seller.ErrNotFound) {
				web.Error(c, http.StatusNotFound, ErrSellerNotFound)
			} else {
				//Or return a 500 status code error by default if an internal error occurs.
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}
		//If no errors occurs, return a 200 status code and the requested seller.
		web.Success(c, http.StatusOK, sellerById)
	}
}

// @Summary Create a new seller or return an error if the new seller cannot be created.
// @Description Create a new seller or return an error if the new seller has invalid
// format, empty fields or if its already exists.
// @Tags domain.Seller
// @Produce json
// @Accept json
// @Success 201 {object} domain.Seller "New created seller"
// @Failure 409 {object} string "seller already exists"
// @Failure 422 {object} string "invalid JSON"
// @Failure 422 {object} string "invalid seller"
// @Failure 422 {object} string "id locality not exists"
// @Failure 422 {object} string "invalid or missing locality. Locality_id must be 1 or greater"
// @Failure 500 {object} string "Server Internal error"
// @Param Seller body domain.Seller true "Struct of Seller domain"
// @Router /seller [post]
func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get a copy of the body request for using it later in the process of creating the new seller with a map of strings and interfaces
		//to check the fields of the request and a struct of the domain for create the new seller.
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			//if cannot read the body, return a 500 status code and an error message.
			web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			return
		}

		var bodyJson map[string]interface{}
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			//If an error occurs when unmarshal the body, return a 400 Bad Request status code and an error message.
			web.Error(c, http.StatusBadRequest, ErrInvalidJSON)
			return
		}

		//Check if the cid is a valid integer greater than 0.
		cidSeller, ok := bodyJson["cid"].(float64)
		if !ok || cidSeller <= 0 || cidSeller != float64(int(cidSeller)) {
			//If not, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrCID)
			return
		}

		//Check if the locality_id is a valid integer greater than 0.
		locality_id, ok := bodyJson["locality_id"].(float64)
		if !ok || locality_id <= 0 || locality_id != float64(int(locality_id)) {
			//If not, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrLocality)
			return
		}

		//Check if locality id exists.
		existLocality := s.sellerService.GetLocalityIdFromSeller(c, int(locality_id))
		if !existLocality {
			//If not, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrIdLocalityNotExists)
			return
		}

		//Only string fields of the structure.
		requiredStringFields := []string{"company_name", "address", "telephone"}
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
		fields := []string{"cid", "company_name", "address", "telephone", "locality_id"}

		// Check if all the fields are in the request using the function validExistenceOfField and the fields of the structure.
		err = validExistenceOfField(bodyJson, fields...)
		if err != nil {
			//If an error occurs, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		//Process of creating the new data.
		//Create a new structure with the data of the request to avoid expose sensitive data of the domain.
		sellerRequest := domain.Seller{
			CID:         int(bodyJson["cid"].(float64)),
			CompanyName: bodyJson["company_name"].(string),
			Address:     bodyJson["address"].(string),
			Telephone:   bodyJson["telephone"].(string),
			IDLocality:  int(bodyJson["locality_id"].(float64)),
		}
		//var sellerRequest domain.Seller
		err = json.Unmarshal(body, &sellerRequest)
		if err != nil {
			//If an error occurs when unmarshal the struct, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrInvalidJSON)
			return
		}

		// Check if there are fields empty or incorrect type of data.
		validEmpty := validateSellerEmpty(&sellerRequest)
		if !validEmpty {
			//If an error occurs, return a 422 status code and an error message.
			web.Error(c, http.StatusUnprocessableEntity, ErrInvalidSeller)
			return
		}

		// Save the new seller
		newSellerID, err := s.sellerService.Save(c, sellerRequest)
		if err != nil {
			if errors.Is(err, seller.ErrSellerAlreadyExists) {
				web.Error(c, http.StatusConflict, ErrSellerAlreadyExists)
			} else {
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}

		//Create a new structure with the data of the new seller created to avoid expose sensitive data of the domain
		//like the sellerRequest does too.
		newSellerCreated := domain.Seller{
			ID:          newSellerID,
			CID:         int(bodyJson["cid"].(float64)),
			CompanyName: bodyJson["company_name"].(string),
			Address:     bodyJson["address"].(string),
			Telephone:   bodyJson["telephone"].(string),
			IDLocality:  int(bodyJson["locality_id"].(float64)),
		}
		//Return a 201 status code and the new seller created.
		web.Success(c, http.StatusCreated, newSellerCreated)
	}
}

// @Summary Update seller or return an error otherwise.
// @Description Update a seller or return an error if that seller not exist,
// or has invalid format.
// @Tags domain.Seller
// @Produce json
// @Accept json
// @Success 200 {object} domain.Seller "Seller updated"
// @Failure 400 {object} string "invalid id"
// @Failure 404 {object} string "Seller not found"
// @Failure 500 {object} string "Server Internal error"
// @Param id path int true "id from the seller"
// @Param Seller body domain.Seller true "Struct of Seller domain"
// @Router /seller/{id}  [patch]
func (s *Seller) Update() gin.HandlerFunc {

	return func(c *gin.Context) {
		idParam := c.Param("id")
		//Check if the id is a valid integer.
		id, err := strconv.Atoi(idParam)

		//Check if the id is greater than 0.
		if id < 1 {
			//In case to be lower, return a 400 status code and an error message.
			web.Error(c, http.StatusBadRequest, ErrGreaterID)
			return
		}

		//Check if occurs an error when convert the string to int.
		if err != nil {
			//In case to be invalid, return a 400 status code and an error message.
			web.Error(c, http.StatusBadRequest, ErrInvalidID)
			return
		}
		//Check if the seller exists.
		sellerID, err := s.sellerService.GetSellerByID(c, id)
		if err != nil {
			//In case to not exist, return a 404 status code and an error message.
			if errors.Is(err, seller.ErrNotFound) {
				web.Error(c, http.StatusNotFound, ErrSellerNotFound)
			} else {
				//Otherwise, return a 500 status code and an error message.
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}

		//Check if the request has a valid format.
		if err := c.ShouldBindJSON(&sellerID); err != nil {
			//In case to be invalid, return a 400 status code and an error message.
			web.Error(c, http.StatusBadRequest, ErrInvalidJSON)
			return
		}

		// It checks if locality_id is present in the request by checking if it is greater than zero.
		if sellerID.IDLocality > 0 {

			// Check if locality_id exists in the database.
			existLocality := s.sellerService.GetLocalityIdFromSeller(c, sellerID.IDLocality)
			if !existLocality {
				// If not, return a 422 status code and an error message.
				web.Error(c, http.StatusUnprocessableEntity, ErrIdLocalityNotExists)
				return
			}
		}

		//Update the seller using a local structure to avoid expose sensitive data of the domain.
		update := domain.Seller{
			ID:          id,
			CID:         sellerID.CID,
			CompanyName: sellerID.CompanyName,
			Address:     sellerID.Address,
			Telephone:   sellerID.Telephone,
			IDLocality:  sellerID.IDLocality,
		}

		err2 := s.sellerService.Update(c, update, id)
		//Check if occurs an error when update the seller. It can be a 500 status code because the existence of the seller is checked before.
		if err2 != nil {
			//In case of error, return a 500 status code and an error message.
			web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			return
		}
		//Return a 200 status code and the seller updated.
		web.Success(c, http.StatusOK, update)
	}
}

// @Summary Delete seller or return an error otherwise.
// @Description Delete a seller using its id or return an error if that
// seller not exist.
// @Tags domain.Seller
// @Success 204 {object} string "Seller deleted"
// @Failure 400 {object} string "invalid id"
// @Failure 400 {object} string "id must be 1 or greater"
// @Failure 404 {object} string "Seller not found"
// @Failure 500 {object} string "Server Internal error"
// @Param id path int true "id from the seller"
// @Router /seller/{id} [delete]
func (s *Seller) Delete() gin.HandlerFunc {
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
		//Delete the seller using its id.
		err = s.sellerService.Delete(c, id)
		//Check if occurs an error when delete the seller.
		if err != nil {
			//In case to not exist, return a 404 status code and an error message.
			if errors.Is(err, seller.ErrNotFound) {
				web.Error(c, http.StatusNotFound, ErrSellerNotFound)
				//Otherwise, return a 500 status code and an error message.
			} else {
				web.Error(c, http.StatusInternalServerError, ErrInternalServer)
			}
			return
		}
		//Return a 204 status code if the seller was deleted.
		web.Success(c, http.StatusNoContent, "Seller deleted")
	}
}

// Check if the fields are empty or seller.CID has the correct number.
func validateSellerEmpty(seller *domain.Seller) bool {
	expectedType := "string"
	expectedTypeCID := "int"
	switch {
	//Checkk if CID is greater than 0, address, company name and telephone are not empty.
	case seller.CID <= 0 || seller.Address == "" || seller.CompanyName == "" || seller.Telephone == "":
		return false
	//Check if CID has the correct type (int).
	case fmt.Sprintf("%T", seller.CID) != expectedTypeCID:
		return false
	//Check if address, company name and telephone have the correct type (String).
	case fmt.Sprintf("%T", seller.Address) != expectedType || fmt.Sprintf("%T", seller.CompanyName) != expectedType || fmt.Sprintf("%T", seller.Telephone) != expectedType:
		return false
	}
	//If all the fields are correct, return true.
	return true
}

// Check if all the fields of the request exist using a map. Return an error otherwise.
func validExistenceOfField(body map[string]interface{}, fields ...string) error {

	for _, field := range fields {
		//Check if the field exist in the map.
		if _, ok := body[field]; !ok {
			//If the field does not exist, return an error describing which field is required.
			err := field + " is required"
			return errors.New(err)
		}
	}
	//If all the fields exist, return nil.
	return nil
}
