package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/davidop97/apiGo/internal/buyer"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/gin-gonic/gin"
)

type Buyer struct {
	buyerService buyer.Service
}

func NewBuyer(b buyer.Service) *Buyer {
	return &Buyer{
		buyerService: b,
	}
}

type Request struct {
	CardNumberID string `json:"card_number_id,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

type RequestBodyBuyerCreate struct {
	CardNumberID string `json:"card_number_id,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

// ShowGetBuyer godoc
// Summary Get a buyer
// @Description get a buyer
// @Tags domain.Buyer
// @Tags buyers
// @Produce json
// @Param id path int true "Buyer id"
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Success 200 {object} map[string]any
// @Router /buyers/{id} [get]
func (b *Buyer) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request
		// obtain id from Param
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		// process
		// obtain a buyer by id
		b, err := b.buyerService.Get(c, id)
		// check if any error appears
		if err != nil {
			switch {
			case errors.Is(err, buyer.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": " Internal error"})
			}
			return
		}

		// response
		c.JSON(http.StatusOK, gin.H{"data": b})
	}
}

// ShowGetAll returns all buyers
// Summary Show all buyers
// @Description get all buyers
// @Tags domain.Buyer
// @Tags buyers
// @Produce json
// @Failure 500 {object} map[string]any
// @Success 200 {object} map[string]any
// @Router /buyers [get]
func (b *Buyer) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// obtain all buyer
		buyers, err := b.buyerService.GetAll(c)
		// check for errors
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error", nil)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": buyers,
		})
	}
}

// PostBuyer save the given buyer
// Summary Create a buyer
// @Description create a buyer
// @Tags domain.Buyer
// @Tags buyers
// @Accept json
// @Produce json
// @Param body body domain.Buyer true "Buyer body"
// @Failure 400 {object} map[string]any
// @Failure 422 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Success 201 {object} map[string]any
// @Router /buyers [post]
func (b *Buyer) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestBodyBuyerCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		// check if all the fields come into request
		if req.CardNumberID == "" || req.FirstName == "" || req.LastName == "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Missing fields"})
			return
		}

		// process
		buyerCreate := domain.Buyer{
			CardNumberID: req.CardNumberID,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}

		id, err := b.buyerService.Save(c, buyerCreate)
		// check for error
		if err != nil {
			switch {
			case errors.Is(err, buyer.ErrAlreadyExists):
				c.JSON(http.StatusConflict, gin.H{"error": "Buyer already exists"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}
		// response
		newData := domain.Buyer{
			CardNumberID: req.CardNumberID,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}

		newData.ID = id
		c.JSON(http.StatusCreated, gin.H{"data": newData})
	}
}

// UpdateBuyer updates a buyer
// @Summary Update a buyer
// @Tags domain.Buyer
// @Description Update a buyer
// @Tags buyers
// @Param id path int true "Update buyer ID"
// @Accept json
// @Produce json
// @Param request body Request true "Buyer update request"
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Success 200 {object} map[string]any
// @Router /buyers/{id} [patch]
func (b *Buyer) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request
		idParam, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		// check if id exists
		bs, err := b.buyerService.Get(c, idParam)
		// check for errors
		if err != nil {
			switch {
			case errors.Is(err, buyer.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		// prepare request
		toUpdate := domain.Buyer{
			ID:           bs.ID,
			CardNumberID: req.CardNumberID,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}

		// process
		err = b.buyerService.Update(c, idParam, toUpdate, &bs)
		if err != nil {
			switch {
			case errors.Is(err, buyer.ErrAlreadyExists):
				c.JSON(http.StatusConflict, gin.H{"error": "Buyer already exists"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}
		// response
		c.JSON(http.StatusOK, gin.H{"data": bs})
	}
}

// DeleteBuyer deletes a buyer
// Summary Delete a buyer
// @Description delete a buyer
// @Tags domain.Buyer
// @Tags buyers
// @Param id path int true "Delete buyer ID"
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Success 204 {object} map[string]any
// @Router /buyers/{id} [delete]
func (b *Buyer) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request
		// get id from Param
		id, err := strconv.Atoi(c.Param("id"))
		// check for error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		// process
		// call to buyer service to delete a buyer by id
		err = b.buyerService.Delete(c, id)
		if err != nil {
			switch {
			case errors.Is(err, buyer.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "Buyer not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": " Internal error"})
			}
			return
		}

		// response
		c.JSON(http.StatusNoContent, nil)
	}
}
