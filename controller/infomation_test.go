package controllers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sumwhere/models"
	"testing"
)

//func TestInfomationController_GetAdvertisement(t *testing.T) {
//	req := httptest.NewRequest(echo.GET, "/advertisement", nil)
//	req.Header.Set(echo.HeaderAuthorization, TOKEN)
//	rec := httptest.NewRecorder()
//	ctx := echoApp.NewContext(req, rec)
//	assert.NoError(t, handleWithFilter(InfomationController{}.GetAdvertisement, ctx))
//	assert.Equal(t, http.StatusOK, rec.Code)
//
//	var v struct {
//		Result  []models.Advertisement `json:"result"`
//		Success bool                   `json:"success"`
//	}
//
//	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
//	require.Equal(t, true, v.Success)
//}

func TestInfomationController_GetEvent(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/event", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	require.NoError(t, handleWithFilter(InfomationController{}.GetEvent, ctx))
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  []models.Event `json:"result"`
		Success bool           `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.NotEqual(t, 0, len(v.Result))
	require.Equal(t, true, v.Success)
}

func TestInfomationController_GetNotice(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/notice", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	require.NoError(t, handleWithFilter(InfomationController{}.GetNotice, ctx))
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  []models.Notice `json:"result"`
		Success bool            `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}
