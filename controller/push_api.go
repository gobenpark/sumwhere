package controllers

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"sumwhere/factory"
	"sumwhere/middlewares"
	"sumwhere/models"
	"sumwhere/utils"
)

type PushController struct {
}

func (p PushController) Init(g *echo.Group) {
	g.GET("/push", p.GetAllPush)
	g.PUT("/push", p.UpdatePush)
	g.PUT("/fcmToken", p.FcmTokenUpdate)
	g.GET("/pushHistory", p.GetHistory)
}

func (PushController) GetAllPush(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	push, err := models.Push{}.Get(e.Request().Context(), claims.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, push)
}

func (p PushController) UpdatePush(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	var input models.PushInput
	if err := e.Bind(&input); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	push, err := models.Push{}.Get(e.Request().Context(), claims.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	p.ChangeSubscribe(e.Request().Context(), *push, input)

	push.ChatAlert = input.ChatAlert
	push.FriendAlert = input.FriendAlert
	push.MatchAlert = input.MatchAlert
	push.EventAlert = input.EventAlert

	err = push.Update(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, true)
}

func (PushController) FcmTokenUpdate(e echo.Context) error {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	token := e.Request().PostFormValue("fcmtoken")
	if len(token) == 0 {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, errors.New("토큰을 입력해주세요"))
	}

	push, err := models.Push{}.Get(e.Request().Context(), claims.Id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	push.FcmToken = token
	err = push.Update(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, push)
}

// TODO: 에러 발생시 해결해야될 사항
func (PushController) ChangeSubscribe(ctx context.Context, source models.Push, target models.PushInput) {

	if source.ChatAlert != target.ChatAlert {
		result, err := factory.Firebase(ctx).SetSubscribe(ctx, target.ChatAlert, []string{source.FcmToken}, middlewares.CHATALERT)
		if err != nil {
			log.Error(err)
		}
		log.Info(result)
	}

	if source.MatchAlert != target.MatchAlert {
		result, err := factory.Firebase(ctx).SetSubscribe(ctx, target.MatchAlert, []string{source.FcmToken}, middlewares.MATCHALERT)
		if err != nil {
			log.Error(err)
		}
		log.Info(result)
	}

	if source.FriendAlert != target.FriendAlert {
		result, err := factory.Firebase(ctx).SetSubscribe(ctx, target.FriendAlert, []string{source.FcmToken}, middlewares.FRIENDALERT)
		if err != nil {
			log.Error(err)
		}
		log.Info(result)
	}

	if source.EventAlert != target.EventAlert {
		result, err := factory.Firebase(ctx).SetSubscribe(ctx, target.EventAlert, []string{source.FcmToken}, middlewares.EVENTALERT)
		if err != nil {
			log.Error(err)
		}
		log.Info(result)
	}
}

func (PushController) GetHistory(e echo.Context) error {
	//users := e.Get("user").(*jwt.Token)
	//claims := users.Claims.(*models.JwtCustomClaims)
	return nil
}
