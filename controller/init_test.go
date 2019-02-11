package controllers

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"runtime"
	"sumwhere/middlewares"
	"sumwhere/models"
	"time"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func init() {
	runtime.GOMAXPROCS(1)
	xormEngine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(210.100.177.146:33060)/%s", "root", "1q2w3e4r", "sumwhere"))
	if err != nil {
		panic(err)
	}
	_ = xormEngine.Sync(new(models.Banner))

	echoApp = echo.New()
	echoApp.Validator = &Validator{}

	db := middlewares.ContextDB("test", xormEngine)
	token := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("parkbumwoo"),
		Claims: &models.JwtCustomClaims{
			Admin: false,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 8760).Unix(),
			},
		},
	})

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return token(db(handlerFunc))(c)
	}
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
