package handler

import (
	"bytes"
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

func TestHandler_CarriesRead(t *testing.T) {
	t.Run("it should return all carries", func(t *testing.T) {
		server := gin.New()

		// Arrange.
		expectedCarries := []domain.Carries{
			{
				ID:          1,
				CID:         "Test Name",
				CompanyName: "Test Company Name",
				Address:     "Test Address",
				Telephone:   "Test Telephone",
				LocalityID:  1,
			},
			{
				ID:          2,
				CID:         "Test Name 2",
				CompanyName: "Test Company Name 2",
				Address:     "Test Address 2",
				Telephone:   "Test Telephone 2",
				LocalityID:  2,
			},
		}

		carriesJson, err := json.Marshal(expectedCarries)
		if err != nil {
			panic(err)
		}
		expectedBody := `{"data":` + string(carriesJson) + `}`

		service := &carries.ServiceMock{}
		service.On("GetAll", mock.Anything).Return(expectedCarries, nil)

		handler := NewCarry(service)
		server.GET("/api/v1/carries", handler.GetAll())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/carries", nil)
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
}

func TestHandler_CarriesCreate(t *testing.T) {
	t.Run("it should create a carry", func(t *testing.T) {
		server := gin.New()

		// Arrange.
		carry := domain.Carries{
			ID:          1,
			CID:         "Test Name",
			CompanyName: "Test Company Name",
			Address:     "Test Address",
			Telephone:   "Test Telephone",
			LocalityID:  1,
		}
		carryJson, err := json.Marshal(carry)
		if err != nil {
			panic(err)
		}
		expectedBody := `{"data":` + string(carryJson) + `}`

		reqBody := `{
			"cid": "Test Name",
			"company_name": "Test Company Name",
			"address": "Test Address",
			"telephone": "Test Telephone",
			"locality_id": "1"
		}`

		service := &carries.ServiceMock{}
		service.On("Save", mock.Anything, mock.Anything).Return(carry.ID, nil)

		handler := NewCarry(service)
		server.POST("/api/v1/carries", handler.Save())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer([]byte(reqBody)))
		res := httptest.NewRecorder()

		t.Log(string(carryJson))

		// Act.
		server.ServeHTTP(res, req)

		t.Log(carryJson)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return 409 when the carry already exists", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		reqBody := `{
			"cid": "Test Name",
			"company_name": "Test Company Name",
			"address": "Test Address",
			"telephone": "Test Telephone",
			"locality_id": "1"
		}`

		expectedBody := `{"message":"Carry already exists"}`

		service := &carries.ServiceMock{}
		service.On("Save", mock.Anything, mock.Anything).Return(0, carries.ErrDuplicateCarry)

		handler := NewCarry(service)
		server.POST("/api/v1/carries", handler.Save())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer([]byte(reqBody)))
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return 422 when the request body is invalid", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		reqBody := `{
			"cid": "Test Name",
			"company_name": "Test Company Name",
			"telephone": "Test Telephone",
			"locality_id": "1"
		}`

		expectedBody := `{"message":"Invalid request body"}`

		service := &carries.ServiceMock{}

		handler := NewCarry(service)
		server.POST("/api/v1/carries", handler.Save())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer([]byte(reqBody)))
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return 500 when the body is invalid", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		reqBody := `{
			"cid": "Test Name",
			"company_name": "Test Company Name",
			"address": "Test Address",
			"telephone": "Test Telephone",
			"locality_id": "1"
		`

		expectedBody := `{"message":"Internal server error"}`

		service := &carries.ServiceMock{}

		handler := NewCarry(service)
		server.POST("/api/v1/carries", handler.Save())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer([]byte(reqBody)))
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return 422 if some field has invalid data", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		reqBody := `{
			"cid": "Test Name",
			"company_name": "Test Company Name",
			"address": "",
			"telephone": "",
			"locality_id": "1"
		}`

		expectedBody := `{"message": "Incorrect data"}`

		service := &carries.ServiceMock{}
		service.On("Save", mock.Anything, mock.Anything).Return(0, carries.ErrIncorrectData)

		handler := NewCarry(service)
		server.POST("/api/v1/carries", handler.Save())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer([]byte(reqBody)))
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return 404 if the locality does not exists", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		reqBody := `{
			"cid": "Test Name",
			"company_name": "Test Company Name",
			"address": "Test Address",
			"telephone": "Test Telephone",
			"locality_id": "1"
		}`

		expectedBody := `{"message": "Locality not found"}`

		service := &carries.ServiceMock{}
		service.On("Save", mock.Anything, mock.Anything).Return(0, carries.ErrLocalityCarriesNotFound)

		handler := NewCarry(service)
		server.POST("/api/v1/carries", handler.Save())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer([]byte(reqBody)))
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return 500 when the locality_id is not a number", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		reqBody := `{
			"cid": "Test Name",
			"company_name": "Test Company Name",
			"address": "Test Address",
			"telephone": "Test Telephone",
			"locality_id": "a"
		}`

		expectedBody := `{"message":"Internal server error"}`

		service := &carries.ServiceMock{}

		handler := NewCarry(service)
		server.POST("/api/v1/carries", handler.Save())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/carries", bytes.NewBuffer([]byte(reqBody)))
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
}

func TestHandler_CarriesGetAllCarriesByLocality(t *testing.T) {
	t.Run("it should return 200 when the carries are found", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		localityCarries := []domain.LocalityCarries{
			{
				LocalityID:   "1",
				LocalityName: "Test Locality Name",
				CarriesCount: 1,
			},
			{
				LocalityID:   "2",
				LocalityName: "Test Locality Name 2",
				CarriesCount: 2,
			},
		}
		localityCarriesJson, err := json.Marshal(localityCarries)
		if err != nil {
			panic(err)
		}
		expectedBody := `{"data":` + string(localityCarriesJson) + `}`

		service := &carries.ServiceMock{}
		service.On("GetAllCarriesByLocality", mock.Anything).Return(localityCarries, nil)

		handler := NewCarry(service)
		server.GET("/api/v1/localities/reportCarries", handler.GetCarriesByLocality())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries", nil)
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
}

func TestHandler_CarriesGetAllCarriesByLocalityID(t *testing.T) {
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

		service := &carries.ServiceMock{}
		service.On("GetAllCarriesByLocalityID", mock.Anything, localityId).Return(localityCarries, nil)

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
		service.AssertExpectations(t)
	})
	t.Run("it should return 500 when the locality_id is not a number", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		expectedBody := `{"message":"Internal server error"}`

		service := &carries.ServiceMock{}

		handler := NewCarry(service)
		server.GET("/api/v1/localities/reportCarries", handler.GetCarriesByLocality())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries?id=a", nil)
		res := httptest.NewRecorder()

		// Act.
		server.ServeHTTP(res, req)

		// Assert.
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, expectedBody, res.Body.String())
		service.AssertExpectations(t)
	})
	t.Run("it should return 404 when the locality does not exists", func(t *testing.T) {
		// Arrange.
		server := gin.New()

		expectedBody := `{"message":"Locality not found"}`

		service := &carries.ServiceMock{}
		service.On("GetAllCarriesByLocalityID", mock.Anything, 1).Return(domain.LocalityCarries{}, carries.ErrLocalityCarriesNotFound)

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
		service.AssertExpectations(t)
	})
}
