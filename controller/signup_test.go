package controllers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignUpController_Email(t *testing.T) {
	loginJSON := `{"email":"qbne9dfhdfg2@naver.com","password":"1q2w3e4r"}`
	req := httptest.NewRequest(echo.PATCH, "/signin/email", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	require.NoError(t, nonrestrictHandleWithFilter(SignUpController{}.Email, ctx))
	require.Equal(t, http.StatusOK, rec.Code)

	var v struct {
		Result  Token `json:"result"`
		Success bool  `json:"success"`
	}

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, true, v.Success)
}
