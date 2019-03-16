package controllers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatchController_GetTotalCount(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/match/totalcount", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	require.NoError(t, handleWithFilter(MatchController{}.GetTotalCount, ctx))
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  int64 `json:"result"`
		Success bool  `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}

func TestMatchController_GetMatchRequestHistory(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/match/history/request", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	require.NoError(t, handleWithFilter(MatchController{}.GetMatchRequestHistory, ctx))
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  []interface{} `json:"result"`
		Success bool          `json:"success"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}

func TestMatchController_GetMatchReceiveHistory(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/match/history/receive", nil)
	req.Header.Set(echo.HeaderAuthorization, TOKEN)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	require.NoError(t, handleWithFilter(MatchController{}.GetMatchReceiveHistory, ctx))
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  []interface{} `json:"result"`
		Success bool          `json:"success"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}
