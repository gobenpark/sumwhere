package controllers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sumwhere/models"
	"testing"
)

func TestMainController_TopTrip(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/main/toptrip", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OCwiZW1haWwiOiIiLCJhZG1pbiI6ZmFsc2UsImV4cCI6MTU3OTE1MTAxOX0.huD7yQUMvbTAcRyh9oKvayPGDsN4lzLWuiST4S-IJe4")
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	assert.NoError(t, handleWithFilter(MainController{}.TopTrip, ctx))
	assert.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  models.TripPlaceType `json:"result"`
		Success bool                 `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}
