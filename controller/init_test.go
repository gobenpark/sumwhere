package controllers

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"runtime"
	"sumwhere/middlewares"
	"sumwhere/models"
	"time"
)

const TOKEN string = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OCwiZW1haWwiOiIiLCJhZG1pbiI6ZmFsc2UsImV4cCI6MTU3OTE1MTAxOX0.huD7yQUMvbTAcRyh9oKvayPGDsN4lzLWuiST4S-IJe4"

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

	opt, err := redis.ParseURL("redis://:@1.215.236.26:53379")
	if err != nil {
		panic(err)
	}

	rclient := redis.NewClient(opt)

	_ = xormEngine.Sync2(new(models.Banner),
		new(models.Country),
		new(models.Advertisement),
		new(models.Notice),
		new(models.Event),
		new(models.TripPlace))

	fmt.Println("start")
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

	redisClient := middlewares.ContextRedis(middlewares.ContextSetRedisName, rclient)

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return redisClient(token(db(handlerFunc)))(c)
	}
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
