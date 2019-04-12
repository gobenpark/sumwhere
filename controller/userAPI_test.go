package controllers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"sumwhere/models"
	"sumwhere/utils"
	"testing"
)

func TestUserController_All(t *testing.T) {

	m := models.GetQuery{
		SortBy:  []string{"email"},
		OrderBy: []string{"desc"},
		Offset:  1,
		Limit:   2,
	}

	bt, err := json.Marshal(&m)
	require.NoError(t, err)
	req := httptest.NewRequest(echo.PATCH, "/user/all", strings.NewReader(string(bt)))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+TokenHelper())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	require.NoError(t, restrictHandleWithFilter(UserController{}.All, ctx))
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  utils.ArrayResult `json:"result"`
		Success bool              `json:"success"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	t.Log(v)
	require.Equal(t, true, v.Success)
}
