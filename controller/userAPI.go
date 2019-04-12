package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
)

type UserController struct {
}

func (u UserController) Init(g *echo.Group) {
	g.GET("/existProfile", u.ExistProfile)
	g.GET("/user", u.GetUser)
	g.GET("/user/all", u.All)
	g.GET("/another_user", u.GetById)
	g.GET("/user_with_profile", u.GetUserWithProfile)
	g.POST("/profile", u.CreateProfile)
	g.GET("/characters", u.GetAllCharacter)
	g.DELETE("/signout", u.SignOut)
}

// UserController godoc
// @Summary 유저 리스트 반환
// @tags user
// @Description 가입자 리스트 반환
// @Param query body models.GetQuery true "순서,순서마다의 정렬순서 (desc|asc), 어디서부터, 몇개를 가져올지"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ArrayResult "쿼리의 결과"
// @Router /user/all [get]
func (UserController) All(e echo.Context) error {

	var i models.GetQuery

	if err := e.Bind(&i); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := e.Validate(&i); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	count, item, err := models.User{}.GetAll(e.Request().Context(), i.SortBy, i.OrderBy, i.Offset, i.Limit)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, utils.ArrayResult{
		Items:      item,
		TotalCount: count,
	})
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
		Image1:        ProfileSaver(image1, user, "image1"),
		Image2:        ProfileSaver(image2, user, "image2"),
		Image3:        ProfileSaver(image3, user, "image3"),
		Image4:        ProfileSaver(image4, user, "image4"),
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

func ProfileSaver(header *multipart.FileHeader, user *models.User, imageName string) string {

	if header == nil {
		return ""
	}

	file, err := header.Open()
	if err != nil {
		return ""
	}

	defer file.Close()
	path := fmt.Sprintf("/images/%d/profile/", user.Id)
	CreateDirIfNotExist(path)

	dst, err := os.Create(path + fmt.Sprintf("%s.jpg", imageName))
	if err != nil {
		return ""
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return ""
	}
	return fmt.Sprintf("/%d/profile/%s.jpg", user.Id, imageName)
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
