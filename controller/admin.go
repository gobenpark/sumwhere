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
// @tags admin
// @Description 기존 어드민사용자가 추가로 어드민 권한을 유저에게 ≠부여
// @Param id body controllers.AdminController false "{'id':1}"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User "해당 유저정보"
// @Router /admin/assgin [patch]
func (AdminController) Assgin(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	if !user.Admin {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorNotEnoughPoint, errors.New("unauthorized user"))
	}

	var i struct {
		ID int64 `json:"id" valid:"required"`
	}

	if err := e.Bind(&i); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := e.Validate(&i); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	targetUser, err := models.User{}.GetByUserId(e, i.ID)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorUserNotExists, err)
	}

	targetUser.Admin = true
	if _, err := targetUser.Update(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, targetUser)
}
