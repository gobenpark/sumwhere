package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type MainController struct {
}

func (m MainController) Init(g *echo.Group) {
	g.GET("/main/toptrip", m.TopTrip)
}

func (MainController) TopTrip(e echo.Context) error {
	trips, err := models.TripPlaceType{}.TopTripPlaces(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, trips)
}
