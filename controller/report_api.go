package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type ReportController struct {
}

func (r ReportController) Init(g *echo.Group) {
	g.POST("/report", r.Insert)
}

func (ReportController) Insert(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)
	var report models.Report
	if err := e.Bind(&report); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	report.UserID = claims.Id

	if err := e.Validate(&report); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := report.Insert(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, report)
}
