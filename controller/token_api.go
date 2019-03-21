package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
)

type TokenController struct {
}

func (t TokenController) Init(g *echo.Group) {
	g.GET("/token/vaild", t.LoginToken)
	g.PUT("/kakao/token_update", t.KakaoTokenUpdate)
}

func (TokenController) LoginToken(e echo.Context) error {
	factory.Logger(e.Request().Context()).Infoln("login token")

	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusUnauthorized, utils.ApiErrorUserNotExists, err)
	}

	if user.Id == 0 {
		return utils.ReturnApiSucc(e, http.StatusOK, false)
	} else {
		return utils.ReturnApiSucc(e, http.StatusOK, true)
	}
}

func (TokenController) KakaoTokenUpdate(e echo.Context) error {
	token := e.QueryParam("token")
	//id := e.QueryParam("id")

	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, err)
	}
	tokenHeader := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", tokenHeader)

	var body []byte
	_, err = req.Body.Read(body)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, string(body))
}
