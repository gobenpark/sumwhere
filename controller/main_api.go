package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type MainController struct {
}

func (m MainController) Init(g *echo.Group) {
	g.GET("/main/list", m.MainList)
	g.GET("/main/toptrip", m.TopTrip)
}

func (MainController) TopTrip(e echo.Context) error {

	trips, err := models.TripPlace{}.TopTripPlaces(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, trips)
}

func (MainController) MainList(e echo.Context) error {
	m, err := models.TripPlace{}.GetCountryJoind(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, m)
}
