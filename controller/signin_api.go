package controllers

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
)

type Token struct {
	AccessToken string `json:"access_token" valid:"required" example:"ldifgj1lij31t9gsegl"`
}

type SignInController struct {
}

func (c SignInController) Init(g *echo.Group) {
	g.POST("/facebook", c.FaceBook)
	g.POST("/kakao", c.Kakao)
	g.POST("/email", c.Email)
}

// ShowAccount godoc
// @Summary 이메일 로그인
// @tags signin
// @Description 이메일을 이용한 로그인
// @Param user body models.User true "email 과 passward만 쓸것"
// @Accept  json
// @Produce  json
// @Success 200 {object} controllers.Token "access 토큰을 반환"
// @Router /signin/email [post]
func (SignInController) Email(e echo.Context) error {
	var u struct {
		Email    string `json:"email" valid:"required"`
		Password string `json:"password" valid:"required"`
	}

	if err := e.Bind(&u); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := e.Validate(&u); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	user, err := models.User{}.ValidateEmailLogin(e.Request().Context(), u.Email, u.Password)
	if user == nil && err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorDB, err)
	} else if user == nil && err == nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, errors.New("존재하지 않는 계정입니다."))
	} else if user != nil && err != nil {
		// 비밀번호 틀릴경우
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorPassword, err)
	}

	t, err := user.JwtTokenCreate()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, Token{t})
}

// signinfacebook godoc
// @Summary 페이스북 로그인
// @tags signin
// @Param token body controllers.Token true "토큰 전송 "
// @Description 페이스북을 이용한 로그인
// @Accept  json
// @Produce  json
// @Success 200 {object} controllers.Token "access 토큰을 반환"
// @Router /signin/facebook [post]
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
	return utils.ReturnApiSucc(e, http.StatusOK, Token{t})
}

// signinkakao godoc
// @Summary kakao 로그인
// @tags signin
// @Param token body controllers.Token true "토큰 전송 "
// @Description 카카오를 이용한 로그인
// @Accept  json
// @Param token body controllers.Token true "Bottle ID"
// @Produce  json
// @Success 200 {object} controllers.Token "access 토큰을 반환"
// @Router /signin/kakao [post]
func (SignInController) Kakao(e echo.Context) error {
	factory.Logger(e.Request().Context()).Info("KakaoLogin")

	// 카카오 유저 가져오기
	model, err := KakaoUtil(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorKakaoAuth, err)
	}

	// 유저 존재하면 유저반환 없으면 생성후 반환
	user, err := model.SearchAndCreate(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorKakaoAuth, err)
	}

	factory.Logger(e.Request().Context()).
		WithFields(logrus.Fields{"userInfo": user}).
		Infoln("Kakao Login")

	// 토큰반환
	t, err := user.JwtTokenCreate()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorKakaoAuth, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, Token{t})
}

func FacebookUtil(c echo.Context) (*models.FaceBookUser, error) {
	var token Token

	if err := c.Bind(&token); err != nil {
		return nil, err
	}

	if err := c.Validate(&token); err != nil {
		return nil, err
	}

	url := "https://graph.facebook.com/v3.0/me?fields=id,email,name&access_token="
	res, err := http.Get(url + token.AccessToken)
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
	facebookuser.Token = token.AccessToken

	return &facebookuser, nil
}

func KakaoUtil(c echo.Context) (*models.KakaoUser, error) {

	var token Token

	if err := c.Bind(&token); err != nil {
		return nil, err
	}

	if err := c.Validate(&token); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer "+token.AccessToken)

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

	kakaoUser.Token = token.AccessToken
	if reflect.DeepEqual(kakaoUser, models.KakaoUser{}) {
		return nil, errors.New("empty kakao model")
	}

	return &kakaoUser, nil
}
