package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sumwhere/models"
	"testing"
)

var userJson = `{"concept":" 명소 바꾸기 "}`
var resultJson = `{"result":{"concept":" 명소 바꾸기 "},"success":true,"error":{}}`

func JWT() *jwt.Token {
	token, _ := jwt.ParseWithClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OCwiZW1haWwiOiIiLCJhZG1pbiI6ZmFsc2UsImV4cCI6MTU3OTE1MTAxOX0.huD7yQUMvbTAcRyh9oKvayPGDsN4lzLWuiST4S-IJe4", models.JwtCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte("parkbumwoo"), nil
	})
	return token
}

func TestTripController_Update(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", JWT())
	c.SetPath("/trip/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	e.ServeHTTP(rec, req)

	// Assertions
	if assert.NoError(t, TripController{}.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, resultJson, rec.Body.String())
	}
}
