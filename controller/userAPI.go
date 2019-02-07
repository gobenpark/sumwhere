package controllers

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
)

type UserController struct {
}

func (u UserController) Init(g *echo.Group) {
	g.GET("/login", u.Login)
	g.GET("/existProfile", u.ExistProfile)
	g.GET("/user", u.GetUser)
	g.GET("/another_user", u.GetById)
	g.GET("/user_with_profile", u.GetUserWithProfile)
	g.POST("/profile", u.CreateProfile)
	g.GET("/characters", u.GetAllCharacter)
	g.DELETE("/signout", u.SignOut)
}

func (UserController) Login(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}
	factory.Logger(e.Request().Context()).WithFields(logrus.Fields{"userInfo": user}).Infoln("userLogin")
	return utils.ReturnApiSucc(e, http.StatusOK, true)
}

func (UserController) ExistProfile(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	if !user.HasProfile {
		return utils.ReturnApiSucc(e, http.StatusOK, false)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, true)
}

func (u UserController) GetUser(e echo.Context) error {

	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	factory.Logger(e.Request().Context()).WithField("user", user).Info("GetUser")

	return utils.ReturnApiSucc(e, http.StatusOK, user)
}

func (UserController) GetUserWithProfile(e echo.Context) error {

	var id int64
	tempId := e.QueryParam("id")

	if tempId == "" {
		users := e.Get("user").(*jwt.Token)
		claims := users.Claims.(*models.JwtCustomClaims)
		id = claims.Id
	} else {
		id, _ = strconv.ParseInt(tempId, 10, 64)
	}

	uwp, err := models.UserWithProfile{}.GetUserWithProfile(e.Request().Context(), id)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, uwp)
}

func (UserController) GetById(e echo.Context) error {
	id := e.QueryParam("id")
	factory.Logger(e.Request().Context()).Info("getby")
	userid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	user, err := models.User{}.GetByUserId(e, userid)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, user)

}

func (UserController) CreateProfile(e echo.Context) error {

	// Get User
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorUserNotExists, err)
	}

	// image1은 필수
	image1, err := e.FormFile("image1")
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, errors.New("첫번째 이미지는 필수입니다."))
	}
	image2, err := e.FormFile("image2")
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, errors.New("이미지 두개는 필수입니다."))
	}

	image3, _ := e.FormFile("image3")
	image4, _ := e.FormFile("image4")

	tripStyleType := e.FormValue("tripStyleType")
	characterType := e.FormValue("characterType")

	var characterModel []models.Character

	err = json.Unmarshal([]byte(characterType), &characterModel)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	job := e.FormValue("job")
	age, err := strconv.Atoi(e.FormValue("age"))
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	profileInput := ProfileInput{
		Age:           age,
		Job:           job,
		TripStyleType: tripStyleType,
		CharacterType: characterModel,
		Image1:        utils.ProfileSaver(image1, user, "image1"),
		Image2:        utils.ProfileSaver(image2, user, "image2"),
		Image3:        utils.ProfileSaver(image3, user, "image3"),
		Image4:        utils.ProfileSaver(image4, user, "image4"),
	}

	if err := e.Validate(profileInput); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	profile, err := profileInput.ToModel()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	profile.UserId = user.Id
	user.Nickname = e.FormValue("nickname")
	user.Gender = e.FormValue("gender")
	user.Point = 50
	user.HasProfile = true
	user.MainProfileImage = profile.Image1
	user.Age = age

	if _, err := profile.Create(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if _, err := user.Update(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, profile)
}

func (UserController) GetAllInterest(e echo.Context) error {
	_, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	interests, err := models.Interest{}.GetAllInterest(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, interests)
}

func (UserController) GetAllCharacter(e echo.Context) error {
	_, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	characters, err := models.Character{}.GetAll(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, characters)
}

func (UserController) SignOut(e echo.Context) error {
	user, err := models.User{}.GetUserByJWT(e)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	result, err := user.Delete(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	if result == 1 {
		return utils.ReturnApiSucc(e, http.StatusOK, true)
	} else {
		return utils.ReturnApiSucc(e, http.StatusOK, false)
	}
}
