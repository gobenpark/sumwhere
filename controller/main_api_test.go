package controllers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"sumwhere/models"
	"testing"
)

func TestMainController_TopTrip(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/main/toptrip", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	assert.NoError(t, handleWithFilter(MainController{}.TopTrip, ctx))

	var v struct {
		Result  models.TripPlace `json:"result"`
		Success bool             `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
}
