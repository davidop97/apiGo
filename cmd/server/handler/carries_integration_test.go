package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidop97/apiGo/internal/carries"
	"github.com/davidop97/apiGo/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerIntegration_CarriesGetAllCarriesByLocalityID(t *testing.T) {
	t.Run("it should return 200 when the carries are found", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		localityCarries := domain.LocalityCarries{
			LocalityID:   "1",
			LocalityName: "Test Locality Name",
			CarriesCount: 1,
		}
		localityId := 1
		localityCarriesJson, err := json.Marshal(localityCarries)
		if err != nil {
			panic(err)
		}
		expectedBody := `{"data":` + string(localityCarriesJson) + `}`

		repo := &carries.RepositoryMock{}
		repo.On("GetAllCarriesByLocalityID", mock.Anything, localityId).Return(localityCarries, nil)

		service := carries.NewService(repo)

		handler := NewCarry(service)
		server.GET("/api/v1/localities/reportCarries", handler.GetCarriesByLocality())

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/localities/reportCarries?id=%d", localityId), nil)
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		repo.AssertExpectations(t)
	})
	t.Run("it should return 404 when the locality does not exists", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		expectedBody := `{"message":"Locality not found"}`

		repo := &carries.RepositoryMock{}
		repo.On("GetAllCarriesByLocalityID", mock.Anything, 1).Return(domain.LocalityCarries{}, carries.ErrLocalityCarriesNotFound)

		service := carries.NewService(repo)

		handler := NewCarry(service)
		server.GET("/api/v1/localities/reportCarries", handler.GetCarriesByLocality())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries?id=1", nil)
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		repo.AssertExpectations(t)
	})
}
