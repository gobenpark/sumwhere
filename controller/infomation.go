package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type InfomationController struct{}

func (i InfomationController) Init(g *echo.Group) {
	g.GET("/advertisement", i.GetNotice)
	g.GET("/notice", i.GetNotice)
	g.GET("/event", i.GetEvent)
}

func (InfomationController) GetNotice(e echo.Context) error {
	model, err := models.Notice{}.GetAll(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, model)
}

func (InfomationController) GetEvent(e echo.Context) error {
	model, err := models.Event{}.GetAll(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, model)
}

func (InfomationController) GetAdvertisement(e echo.Context) error {
	return utils.ReturnApiSucc(e, http.StatusOK, []models.Advertisement{})
}
