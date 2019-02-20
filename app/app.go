package app

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"sumwhere/controller"
	"sumwhere/middlewares"
	"sumwhere/models"
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

func (s Sumwhere) Run() error {
	v1 := s.Group("/v1")
	privateV1 := v1.Group("/restrict")

	privateV1.Use(middleware.JWTWithConfig(middleware.JWTConfig{
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

	s.privateController(privateV1)
	s.publicController(v1)

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

	s.Use(middlewares.ContextFireBase(&fb))
	s.Use(middlewares.ContextDB("sumwhere", db))
	s.Use(middlewares.ContextRedis(middlewares.ContextSetRedisName, initSetRedis()))
	s.Use(middlewares.Logger())
	s.Pre(middleware.RemoveTrailingSlash())
	s.Use(middleware.CORS())
	s.Use(middleware.RequestID())
	s.Use(middleware.Recover())
	s.Validator = &Validator{}
	s.Static("/", "/static/www.sumwhere.kr")
	return nil
}

func (Sumwhere) privateController(e *echo.Group) {
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
}

func (Sumwhere) publicController(e *echo.Group) {
	controllers.SignUpController{}.Init(e.Group("/signup"))
	controllers.SignInController{}.Init(e.Group("/signin"))
}

func initDB() (*xorm.Engine, error) {

	var url string
	dbUser := os.Getenv("DATABASE_USER")
	database := os.Getenv("DATABASE_DRIVER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	switch os.Getenv("RELEASE_SYSTEM") {
	case "kubernetes":
		url = fmt.Sprintf("%s:%s@tcp(mysql-svc.sumwhere:3306)/%s", dbUser, dbPass, dbName)
	default:
		database = "mysql"
		url = fmt.Sprintf("%s:%s@tcp(192.168.0.192:3306)/%s", "root", "1q2w3e4r", "sumwhere")
	}

	db, err := xorm.NewEngine(database, url)
	if err != nil {
		return nil, err
	}
	db.ShowSQL(true)

	_ = db.Sync2(
		new(models.User),
		new(models.Profile),
		new(models.Trip),
		new(models.Match),
		new(models.MatchMember),
		new(models.TripStyle),
		new(models.Interest),
		new(models.Character),
		new(models.TripmatchHistory),
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
		new(models.Push),
		new(models.PushHistory),
		new(models.MatchType),
		new(models.Country),
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
		opt, err := redis.ParseURL("redis://:@1.215.236.26:53379")
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
