package controllers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
	gutil "sumwhere/utils"
)

type SignInController struct {
}

func (c SignInController) Init(g *echo.Group) {
	g.GET("/email", c.Email)
	g.POST("/facebook", c.FaceBook)
	g.POST("/kakao", c.Kakao)
}

func (SignInController) Email(e echo.Context) error {

	email := e.FormValue("email")
	password := e.FormValue("password")

	user, err := models.User{}.GetByEmailWithPassword(e.Request().Context(), email, password)
	if user == nil && err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorDB, err)
	} else if user == nil && err == nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, errors.New("존재하지 않는 계정입니다."))
	} else if user != nil && err != nil {
		// 비밀번호 틀릴경우
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorPassword, err)
	}

	factory.Logger(e.Request().Context()).
		WithFields(logrus.Fields{"userInfo": user}).
		Infoln("Email Login")

	t, err := user.JwtTokenCreate()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, echo.Map{"token": t})
}

func (SignInController) FaceBook(e echo.Context) error {
	factory.Logger(e.Request().Context()).Info("FacebookLogin")

	model, err := gutil.FacebookUtil(e)

	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorFacebookAuth, err)
	}

	user, err := model.SearchAndCreate(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorFacebookAuth, err)
	}

	factory.Logger(e.Request().Context()).
		WithFields(logrus.Fields{"userInfo": user}).
		Infoln("FaceBook Login")

	t, err := user.JwtTokenCreate()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, echo.Map{"token": t})
}

func (SignInController) Kakao(e echo.Context) error {
	factory.Logger(e.Request().Context()).Info("KakaoLogin")

	// 카카오 유저 가져오기
	model, err := gutil.KakaoUtil(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorKakaoAuth, err)
	}

	fmt.Printf("%#v", model)

	// 유저 존재하면 유저반환 없으면 생성후 반환
	user, err := model.SearchAndCreate(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorKakaoAuth, err)
	}

	factory.Logger(e.Request().Context()).
		WithFields(logrus.Fields{"userInfo": user}).
		Infoln("Kakao Login")

	// 토큰반환
	token, err := user.JwtTokenCreate()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorKakaoAuth, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, echo.Map{"token": token})
}
