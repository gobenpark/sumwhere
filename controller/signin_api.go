package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
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

	model, err := FacebookUtil(e)

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
	model, err := KakaoUtil(e)
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

func FacebookUtil(c echo.Context) (*models.FaceBookUser, error) {
	req := c.Request()
	token := req.PostFormValue("access_token")

	url := "https://graph.facebook.com/v3.0/me?fields=id,email,name&access_token="
	res, err := http.Get(url + token)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var facebookuser models.FaceBookUser
	if err := json.Unmarshal(data, &facebookuser); err != nil {
		return nil, err
	}

	if reflect.DeepEqual(facebookuser, models.FaceBookUser{}) {
		return nil, errors.New("empty facebook model")
	}
	facebookuser.Token = token

	return &facebookuser, nil
}

func KakaoUtil(c echo.Context) (*models.KakaoUser, error) {
	req := c.Request()
	token := req.PostFormValue("access_token")
	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	user, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var kakaoUser models.KakaoUser
	json.Unmarshal(user, &kakaoUser)

	kakaoUser.Token = token
	if reflect.DeepEqual(kakaoUser, models.KakaoUser{}) {
		return nil, errors.New("empty kakao model")
	}

	return &kakaoUser, nil
}
