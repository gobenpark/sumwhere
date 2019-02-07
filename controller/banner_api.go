package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"
)

type BannerController struct {
}

func (b BannerController) Init(g *echo.Group) {
	g.GET("/banner", b.GetBanner)
}

func (BannerController) GetBanner(e echo.Context) error {
	banner, err := models.Banner{}.GetAll(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, banner)
}
