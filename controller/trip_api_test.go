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
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OCwiZW1haWwiOiIiLCJhZG1pbiI6ZmFsc2UsImV4cCI6MTU3OTE1MTAxOX0.huD7yQUMvbTAcRyh9oKvayPGDsN4lzLWuiST4S-IJe4")
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
