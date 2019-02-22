package controllers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"sumwhere/models"
	"testing"
)

var userJson = `{"concept":" 명소 바꾸기 "}`
var resultJson = `{"result":{"concept":" 명소 바꾸기 "},"success":true,"error":{}}`

func TestTripController_Update(t *testing.T) {
	req := httptest.NewRequest(echo.PATCH, "/trip", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("57")
	assert.NoError(t, handleWithFilter(TripController{}.Update, ctx))
	assert.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  models.Trip `json:"result"`
		Success bool        `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}

func TestTripController_GetTripPlace(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/trip", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("countryid")
	ctx.SetParamValues("1")
	assert.NoError(t, handleWithFilter(TripController{}.GetTripPlace, ctx))
	assert.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  []models.TripPlace `json:"result"`
		Success bool               `json:"success"`
	}

	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	assert.Equal(t, true, v.Success)
	t.Log(v.Result)
}

func TestTripController_GetTripCountry(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/trip/country", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	assert.NoError(t, handleWithFilter(TripController{}.GetTripCountry, ctx))
	assert.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  []models.Country `json:"result"`
		Success bool             `json:"success"`
	}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	assert.Equal(t, true, v.Success)
}
