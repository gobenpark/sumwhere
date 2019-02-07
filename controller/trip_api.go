package controllers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
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
	g.GET("/trip/:name/validate", t.Validation)
	g.DELETE("/trip", t.Delete)
	g.PATCH("/trip/:id", t.Update)
	g.GET("/trip", t.GetMyTrip)
	g.GET("/alltriplist", t.GetAllTrip)
	g.GET("/tripplaces", t.DestinationSearch)
	g.GET("/tripstyle", t.GetAllTripStyle)
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

func (TripController) Validation(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}
	name := e.Param("name")
	var input ValidateInput
	if err := e.Bind(&input); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	factory.Logger(e.Request().Context()).WithFields(logrus.Fields{
		"user_id": user.Id,
		"id":      input.TripId,
		"start":   input.StartAt,
		"end":     input.EndAt,
	}).Info("Insert Validate")

	switch name {
	case "destination":
		result, err := models.Trip{}.Exist(e.Request().Context(), user.Id, fmt.Sprintf("triptype_id = %d", input.TripId))
		if err != nil {
			return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
		}
		return utils.ReturnApiSucc(e, http.StatusOK, result)
	case "date":
		query := fmt.Sprintf("DATE(start_date) BETWEEN '%s' AND '%s' OR DATE(end_date) BETWEEN '%s' AND '%s'", input.StartAt, input.EndAt, input.StartAt, input.EndAt)
		result, err := models.Trip{}.Exist(e.Request().Context(), user.Id, query)
		if err != nil {
			return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
		}
		return utils.ReturnApiSucc(e, http.StatusOK, result)
	default:
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, errors.New("올바른 값을 입력해주세요"))
	}

	return nil
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

func (TripController) DestinationSearch(e echo.Context) error {
	destination := e.QueryParam("name")
	factory.Logger(e.Request().Context()).WithField("trip_place", destination).Infoln("DestinationSearch")

	result, err := models.TripPlaceType{}.Search(e.Request().Context(), destination)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, result)
}

func (TripController) GetAllTrip(e echo.Context) error {
	_, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	var v SearchInput
	if err := e.Bind(&v); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}

	factory.Logger(e.Request().Context()).WithFields(logrus.Fields{
		"sortby":         v.Sortby,
		"order":          v.Order,
		"maxResultCount": v.MaxResultCount,
		"skipCount":      v.SkipCount,
	}).Info("SearchStart")

	types, err := models.TripPlaceType{}.GetAll(e.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, types)
}

func (TripController) GetAllTripStyle(e echo.Context) error {
	_, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	tripStyles, err := models.TripStyle{}.GetAll(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, tripStyles)
}
