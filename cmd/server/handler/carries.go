package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/davidop97/apiGo/internal/carries"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/gin-gonic/gin"
)

type Carry struct {
	carriesService carries.Service
}

func NewCarry(c carries.Service) *Carry {
	return &Carry{
		carriesService: c,
	}
}

// ShowGetAll godoc
// @Summary Get all carries, returns empty list if there are no carries.
// @Tags carries
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carries [get]
func (c *Carry) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		carriesList, err := c.carriesService.GetAll(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"data": carriesList,
		})
	}
}

// ShowSave godoc
// @Summary Save a carry, returns error if the carry already exists or if the data is incorrect.
// @Tags carries
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carries [post]
func (c *Carry) Save() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Convert the body request into a map for checking
		var body map[string]interface{}
		err := json.NewDecoder(ctx.Request.Body).Decode(&body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		//Validate that all fields are in the request
		err = ValidateBody(body, "cid", "company_name", "address", "telephone", "locality_id")
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Invalid request body"})
			return
		}

		localityId, err := strconv.Atoi(body["locality_id"].(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}
		body["locality_id"] = localityId

		//Convert into bytes
		fmtBody, err := json.Marshal(body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		//Convert into the domain.Carry struct
		var carry domain.Carries
		json.Unmarshal(fmtBody, &carry)
		id, err := c.carriesService.Save(ctx, carry)
		if err != nil {
			if errors.Is(err, carries.ErrIncorrectData) {
				ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Incorrect data"})
				return
			} else if errors.Is(err, carries.ErrDuplicateCarry) {
				ctx.JSON(http.StatusConflict, map[string]interface{}{"message": "Carry already exists"})
				return
			} else if errors.Is(err, carries.ErrLocalityCarriesNotFound) {
				ctx.JSON(http.StatusNotFound, map[string]interface{}{"message": "Locality not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		carry.ID = id
		ctx.JSON(http.StatusCreated, map[string]interface{}{
			"data": carry,
		})
	}
}

// ShowGetAllByLocality godoc
// @Summary Get all carries by locality, returns empty list if there are no carries.
// @Tags carries
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carries/locality [get]
func (c *Carry) GetCarriesByLocality() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Query("id") == "" {
			carriesList, err := c.carriesService.GetAllCarriesByLocality(ctx)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
				return
			}

			ctx.JSON(http.StatusOK, map[string]interface{}{
				"data": carriesList,
			})
			return
		}

		id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		localityCarries, err := c.carriesService.GetAllCarriesByLocalityID(ctx, id)
		if err != nil {
			if errors.Is(err, carries.ErrLocalityCarriesNotFound) {
				ctx.JSON(http.StatusNotFound, map[string]interface{}{"message": "Locality not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"data": localityCarries,
		})

	}
}
