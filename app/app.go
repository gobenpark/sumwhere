package app

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger"
	"os"
	"os/signal"
	"sumwhere/controller"
	_ "sumwhere/docs"
	"sumwhere/middlewares"
	"sumwhere/models"
	"sumwhere/utils"
	"syscall"
	"time"
)

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

type Sumwhere struct {
	*echo.Echo
}

func NewApp() *Sumwhere {

	return &Sumwhere{
		Echo: echo.New(),
	}
}

// @title Sumwhere API
// @version 2.0
// @description This is a Sumwhere server API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://www.sumwhere.kr
// @contact.email qjadn0914@naver.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host www.sumwhere.kr
// @BasePath /v1
// @schemes https http
func (s Sumwhere) Run() error {
	s.GET("/swagger/*", echoSwagger.WrapHandler)
	v1 := s.Group("/v1")
	api := v1.Group("/api")
	privateApi := api.Group("/restrict")

	privateApi.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("parkbumwoo"),
		Claims: &models.JwtCustomClaims{
			Admin: false,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 8760).Unix(),
			},
		},
	}))

	if err := s.setMiddleWare(); err != nil {
		return err
	}

	s.privateApiController(privateApi)
	s.publicApiController(api)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			s.Logger.Fatal(err)
		}
	}()
	return s.Start(fmt.Sprintf(":%s", "8080"))
}

func (s Sumwhere) setMiddleWare() error {

	switch os.Getenv("RELEASE_SYSTEM") {
	case "kubernetes":
		s.Static("/images", "/images")
		s.Use(middlewares.ContextRedis(middlewares.ContextGetRedisName, initGetRedis()))
	default:
		break
	}

	db, err := initDB()
	if err != nil {
		return err
	}

	fb, err := middlewares.NewFireBaseApp()
	if err != nil {
		return err
	}

	s.Use(middlewares.ContextFireBase(fb))
	s.Use(middlewares.ContextDB("sumwhere", db))
	s.Use(middlewares.ContextRedis(middlewares.ContextSetRedisName, initSetRedis()))
	s.Use(middlewares.Logger())
	s.Pre(middleware.RemoveTrailingSlash())
	s.Use(middleware.CORS())
	s.Use(middleware.RequestID())
	s.Use(middleware.Recover())

	s.Validator = &Validator{}
	s.Static("/", "/public/www.sumwhere.kr")
	return nil
}

func (Sumwhere) privateApiController(e *echo.Group) {
	controllers.UserController{}.Init(e)
	controllers.PurchaseController{}.Init(e)
	controllers.MatchController{}.Init(e)
	controllers.TripController{}.Init(e)
	controllers.TokenController{}.Init(e)
	controllers.ChatRoomController{}.Init(e)
	controllers.InfomationController{}.Init(e)
	controllers.ReportController{}.Init(e)
	controllers.BannerController{}.Init(e)
	controllers.PushController{}.Init(e)
	controllers.MainController{}.Init(e)
	controllers.AdminController{}.Init(e)
}

func (Sumwhere) publicApiController(e *echo.Group) {
	controllers.SignUpController{}.Init(e.Group("/signup"))
	controllers.SignInController{}.Init(e.Group("/signin"))
}

func initDB() (*xorm.Engine, error) {

	var url string
	dbUser := utils.DefaultEnv("DATABASE_USER", "root")
	database := utils.DefaultEnv("DATABASE_DRIVER", "mysql")
	dbPass := utils.DefaultEnv("DATABASE_PASS", "1q2w3e4r")
	dbName := utils.DefaultEnv("DATABASE_NAME", "sumwhere")

	db, err := xorm.NewEngine(database, url)

	if err != nil {
		return nil, err
	}

	switch os.Getenv("RELEASE_SYSTEM") {
	case "kubernetes":
		url = fmt.Sprintf("%s:%s@tcp(mysql-svc.sumwhere:3306)/%s", dbUser, dbPass, dbName)
	default:
		db.ShowSQL(true)
		url = fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", "root", "1q2w3e4r", "sumwhere")
	}

	_ = db.Sync2(
		new(models.User),
		new(models.Profile),
		new(models.Trip),
		new(models.Match),
		new(models.MatchMember),
		new(models.TripStyle),
		new(models.Interest),
		new(models.Character),
		new(models.ChatRoom),
		new(models.ChatMember),
		new(models.Banner),
		new(models.PurchaseProduct),
		new(models.PurchaseHistory),
		new(models.TripPlace),
		new(models.Event),
		new(models.Advertisement),
		new(models.Notice),
		new(models.Report),
		new(models.ReportType),
		new(models.MatchType),
		new(models.Country),
		new(models.MatchHistory),
		new(models.TripMatchHistory),
		new(models.Push),
		new(models.PushHistory),
		new(models.PushType),
	)

	return db, nil
}

func initSetRedis() *redis.Client {
	var client *redis.Client
	if os.Getenv("RELEASE_SYSTEM") == "kubernetes" {
		opt, err := redis.ParseURL("redis://:@redis-redis-ha-announce-0.sumwhere:6379")
		if err != nil {
			panic(err)
		}
		client = redis.NewClient(opt)
		log.Info(client.Ping().Val())
	} else {
		opt, err := redis.ParseURL("redis://:@192.168.1.63:6379")
		if err != nil {
			panic(err)
		}
		client = redis.NewClient(opt)
	}
	return client
}

func initGetRedis() *redis.Client {
	opt, err := redis.ParseURL("redis://:@redis-redis-ha.sumwhere:6379")
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}
