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

func TestBannerController_GetBanner(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/banner", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	assert.NoError(t, handleWithFilter(BannerController{}.GetBanner, ctx))
	assert.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  []models.Banner `json:"result"`
		Success bool            `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}
