package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type SignUpController struct {
}

func (c SignUpController) Init(g *echo.Group) {
	g.POST("/email", c.Email)
	g.GET("/nickname/:nickname", c.NickNameExist)

}

// signupkakao godoc
// @Summary 이메일 가입
// @tags signup
// @Param user body models.User true "email 과 passward만 쓸것"
// @Description 이메일을 이용한 가입
// @Accept  json
// @Produce  json
// @Success 200 {object} controllers.Token "access 토큰을 반환"
// @Router /signup/email [post]
func (SignUpController) Email(e echo.Context) error {

	var u models.User

	if err := e.Bind(&u); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	u.JoinType = "EMAIL"

	if err := e.Validate(&u); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	_, err := u.Create(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorDB, err)
	}

	t, err := u.JwtTokenCreate()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, Token{t})
}

// signupkakao godoc
// @Summary 닉네임 확인
// @tags signup
// @Description 닉네임이 존재하는지 확인
// @Accept  json
// @Produce  json
// @Success 200 {boolean} true "bool 타입의 결과 있으면 true 없으면 false"
// @Router /signup/nickname [get]
func (SignUpController) NickNameExist(e echo.Context) error {
	nickname := e.Param("nickname")
	result, err := models.User{}.SearchUserByNickname(e.Request().Context(), nickname)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if result != 0 {
		return utils.ReturnApiSucc(e, http.StatusOK, false)
	} else {
		return utils.ReturnApiSucc(e, http.StatusOK, true)
	}
}
