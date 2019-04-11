package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type AdminController struct {
}

func (a AdminController) Init(g *echo.Group) {
	g.PATCH("/admin/assgin", a.Assgin)
}

// AdminController godoc
// @Summary 어드민 권한 부여
// @tags Admin
// @Description 기존 어드민사용자가 추가로 어드민 권한을 유저에게 부여
// @Param user body models.User true "email 과 passward만 쓸것"
// @Accept  json
// @Produce  json
// @Success 200 {object} controllers.Token "access 토큰을 반환"
// @Router /admin/assgin [patch]
func (AdminController) Assgin(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	if !user.Admin {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorNotEnoughPoint, errors.New("unauthorized user"))
	}

	return utils.ReturnApiSucc(e, http.StatusOK, nil)
}
