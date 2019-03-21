package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sumwhere/factory"
	"sumwhere/middlewares"
	"sumwhere/models"
	"sumwhere/utils"
	"time"
)

type MatchController struct {
}

func (m MatchController) Init(g *echo.Group) {
	g.POST("/match/member", m.JoinMember)
	g.POST("/match/request", m.MatchRequest)
	g.GET("/match/list", m.GetMatchList)
	g.GET("/match/new", m.NewMatchList)
	g.GET("/match/check", m.MatchRequestCheck)
	g.GET("/match/type", m.GetMatchTypes)
	g.GET("/match/totalcount", m.GetTotalCount)

	g.GET("/match/history/request", m.GetMatchRequestHistory)
	g.GET("/match/history/receive", m.GetMatchReceiveHistory)

}

func (MatchController) MatchListFromMysql(e echo.Context, userID int64, trip *models.Trip, count int) ([]models.TripUserGroup, error) {
	trips, err := models.TripUserGroup{}.Join(e.Request().Context(), trip, count)
	if err != nil {
		return nil, utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	return trips, nil
}

/*
- tripMatch_history 테이블 의 내용을 제외한 4명의 사용자 리턴
*/
func (m MatchController) GetMatchList(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	tid, err := strconv.ParseInt(e.QueryParam("tripId"), 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	trip, err := models.Trip{}.Get(e.Request().Context(), tid, claims.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	trips, err := m.MatchListFromMysql(e, claims.Id, trip, 4)
	if err != nil {
		return err
	}

	return utils.ReturnApiSucc(e, http.StatusOK, trips)
}

func (MatchController) MatchRequestCheck(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	result := factory.Redis(e.Request().Context(), middlewares.ContextGetRedisName).ZScore(middlewares.FREEMATCH_COUNT, fmt.Sprintf("%d", claims.Id)).Val()

	possibleCount := 2 - result

	return utils.ReturnApiSucc(e, http.StatusOK, possibleCount)
}

func (MatchController) MatchRequest(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)
	factory.Logger(e.Request().Context()).Info("MatchRequest")
	var m models.MatchRequestDTO
	if err := e.Bind(&m); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorTokenInvaild, err)
	}

	if err := e.Validate(m); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorTokenInvaild, err)
	}

	historyModel := m.ToModel(claims.Id)
	if err := historyModel.Insert(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	push, err := models.Push{}.Get(e.Request().Context(), m.ToUserID)
	if err == nil {
		factory.Firebase(e.Request().Context()).SendMessage("", "매칭 요청이 도착했어요!", push.FcmToken)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, historyModel)
}

/*
TODO: 4명이 안될경우 포인트는?

*/
func (MatchController) NewMatchList(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	queryTID := e.QueryParam("tripId")
	if len(queryTID) == 0 {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, errors.New("tripId는 필수입니다."))
	}

	tid, err := strconv.ParseInt(e.QueryParam("tripId"), 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	trip, err := models.Trip{}.Get(e.Request().Context(), tid, user.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	trips, _ := models.TripUserGroup{}.Join(e.Request().Context(), trip, 4)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	//TODO: 4명이 안될경우 해야될정책은 ?

	// 포인트 체킹
	if user.Point < 5 {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorNotEnoughPoint, errors.New("포인트가 부족합니다."))
	}

	user.Point -= 5
	_, err = user.Update(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorUpdate, err)
	}

	var m []models.TripUserGroup
	if err := json.Unmarshal([]byte(factory.Redis(e.Request().Context(), middlewares.ContextGetRedisName).Get(fmt.Sprintf("trip:%d", tid)).Val()), &m); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorSystem, err)
	}

	m = append(m, trips...)

	bytes, err := json.Marshal(m)

	if err := factory.Redis(e.Request().Context(), middlewares.ContextSetRedisName).Set(fmt.Sprintf("trip:%d", tid), bytes, 24*time.Hour).Err(); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorRedis, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, m)
}

func (MatchController) GetAll(e echo.Context) error {
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

	totalCount, items, err := models.Match{}.GetAll(e.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, utils.ArrayResult{
		TotalCount: totalCount,
		Items:      items,
	})
}

// Member API

func (MatchController) JoinMember(e echo.Context) error {

	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	var mm models.MatchMember

	if err := e.Bind(&mm); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := e.Validate(mm); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	//TODO: 등록시 푸시 알림 및 알림 전송
	mm.UserId = user.Id
	if _, err := mm.Create(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, mm)
}

func (MatchController) GetMatchTypes(e echo.Context) error {
	types, err := models.MatchType{}.GetAll(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, types)
}

func (MatchController) GetTotalCount(e echo.Context) error {
	count, err := models.Match{}.TotalCount(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	if count < 30 {
		return utils.ReturnApiSucc(e, http.StatusOK, 30)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, count)
}

func (MatchController) GetMatchRequestHistory(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	model, err := models.MatchHistory{}.GetRequest(e.Request().Context(), claims.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, model)
}

func (MatchController) GetMatchReceiveHistory(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	model, err := models.MatchHistory{}.GetReceive(e.Request().Context(), claims.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, model)
}
