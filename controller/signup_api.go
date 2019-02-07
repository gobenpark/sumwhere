package controllers

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type SignUpController struct {
}

func (c SignUpController) Init(g *echo.Group) {
	g.POST("", c.Email)
	g.GET("/nickname/:nickname", c.NickNameExist)

}

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

	log.Info(u.Id)
	t, err := u.JwtTokenCreate()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, echo.Map{"token": t})
}

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
