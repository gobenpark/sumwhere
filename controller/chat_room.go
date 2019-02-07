package controllers

import (
	"net/http"
	"sumwhere/models"
	"sumwhere/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type ChatRoomController struct {
}

var upgrader = websocket.Upgrader{}

func (c ChatRoomController) Init(g *echo.Group) {
	g.GET("/chat/room", c.GetChatRoom)
}

func (ChatRoomController) GetChatRoom(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return err
	}

	result, err := models.ChatRoomJoin{}.GetRoom(e.Request().Context(), user.Id)

	return utils.ReturnApiSucc(e, http.StatusOK, result)
}
