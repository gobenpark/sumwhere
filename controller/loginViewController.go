package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type LoginViewController struct {
}

func (l LoginViewController) Init(e *echo.Group) {
	e.GET("/login", l.LoginView)
}

func (LoginViewController) LoginView(e echo.Context) error {
	return e.Render(http.StatusOK, "login.html", nil)
}
