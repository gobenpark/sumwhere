package controllers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"sumwhere/models"
	"testing"
)

func TestAdminController_Assgin(t *testing.T) {

	useridJSON := `{"id": 6}`
	req := httptest.NewRequest(echo.PATCH, "/admin/assgin", strings.NewReader(useridJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+TokenHelper())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	require.NoError(t, restrictHandleWithFilter(AdminController{}.Assgin, ctx))
	t.Log(rec.Body)
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  models.User `json:"result"`
		Success bool        `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}
