package controllers

import (
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
)

type TripController struct {
}

func (t TripController) Init(g *echo.Group) {
	g.POST("/trip", t.Create)
	g.DELETE("/trip", t.Delete)
	g.PATCH("/trip/:id", t.Update)
	g.GET("/trip/country", t.GetTripCountry)
	g.GET("/trip", t.GetMyTrip)
	g.GET("/trip/place/:countryid", t.GetTripPlace)
}

/*
 여행 등록
*/
func (TripController) Create(e echo.Context) error {
	var t TripInput
	if err := e.Bind(&t); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	t.UserId = user.Id

	if err := e.Validate(t); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	trip, err := t.ToModel()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if _, err := trip.Create(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, trip)
}

/*
 여행 삭제
*/
func (TripController) Delete(e echo.Context) error {

	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	id, err := strconv.ParseInt(e.QueryParam("tripId"), 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	trip, err := models.Trip{}.Get(e.Request().Context(), id, user.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if _, err := trip.Delete(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, trip)
}

/*
여행 업데이트
*/
func (t TripController) Update(e echo.Context) error {

	paramId := e.Param("id")

	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	var m models.Trip
	if err := e.Bind(&m); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := m.Update(e.Request().Context(), id); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return t.GetMyTrip(e)
}

func (TripController) GetMyTrip(e echo.Context) error {
	factory.Logger(e.Request().Context()).Logger.Info("GetMyTrip")
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, errors.New("유저정보가 없습니다."))
	}

	trips, err := models.TripGroup{}.GetMyTrip(e.Request().Context(), user.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, trips)
}

func (TripController) GetTripCountry(e echo.Context) error {
	countrys, err := models.Country{}.GetAll(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, countrys)
}

func (TripController) GetTripPlace(e echo.Context) error {
	countryID, err := strconv.ParseInt(e.Param("countryid"), 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	ts, err := models.TripPlace{}.GetAll(e.Request().Context(), countryID)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, ts)

}
