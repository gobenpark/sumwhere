package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"sumwhere/factory"
	"time"
)

type UserWithProfile struct {
	User    `json:"user" xorm:"extends"`
	Profile `json:"profile" xorm:"extends"`
}

func (UserWithProfile) TableName() string {
	return "user"
}

type (
	User struct {
		Id               int64     `json:"id" xorm:"'id' pk autoincr"`
		Email            string    `json:"email" xorm:"varchar(50)" valid:"email"`
		Password         string    `json:"password" valid:"required"`
		Username         string    `json:"username" xorm:"username"`
		Gender           string    `json:"gender"`
		Age              int       `json:"age" xorm:"age"`
		Nickname         string    `json:"nickname" xorm:"nickname VARCHAR(10)"`
		HasProfile       bool      `json:"hasProfile" xorm:"has_profile default 0"`
		JoinType         string    `json:"joinType" xorm:"join_type VARCHAR(50) notnull" valid:"required"`
		Token            string    `json:"token" xorm:"token"`
		SNSID            string    `json:"snsId" xorm:"sns_id"`
		Point            int64     `json:"point" xorm:"point default 1"`
		MainProfileImage string    `json:"mainProfileImage" xorm:"main_profile_image"`
		Admin            bool      `json:"is_admin"`
		CreatedAt        time.Time `xorm:"created"`
		UpdatedAt        time.Time `xorm:"updated"`
		DeletedAt        time.Time `xorm:"deleted"`
	}
)

func (u *User) CheckDuplicate(ctx context.Context) error {
	result, err := factory.DB(ctx).Where("email = ?", u.Email).And("join_type = ?", u.JoinType).Count(User{})
	if err != nil {
		return err
	}

	if result > 0 {
		return errors.New("exist user")
	}
	return nil
}

func (u *User) Create(ctx context.Context) (int64, error) {

	if err := u.CheckDuplicate(ctx); err != nil {
		return 0, err
	}

	u.Password = GenerateHash(u.Password)

	result, err := factory.DB(ctx).Insert(u)
	if err != nil {
		return result, err
	}

	push := &Push{
		UserID: u.Id,
	}

	if err := push.Insert(ctx); err != nil {
		return result, err
	}

	return result, err
}

func (u *User) Update(ctx context.Context) (int64, error) {
	return factory.DB(ctx).ID(u.Id).UseBool("has_profile").Update(u)
}

func (u *User) Delete(ctx context.Context) (int64, error) {
	return factory.DB(ctx).ID(u.Id).Delete(User{})
}

func (User) GetUserByJWT(e echo.Context) (*User, error) {
	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*JwtCustomClaims)
	var u User
	result, err := factory.DB(e.Request().Context()).ID(claims.Id).Get(&u)
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, errors.New("유저가 존재하지 않습니다.")
	}

	return &u, nil
}

func (u User) UpdateKeyWithHistory(e echo.Context, addKey string) error {
	user, err := u.GetUserByJWT(e)
	if err != nil {
		return err
	}
	key := user.Point

	product, err := PurchaseProduct{}.GetByIdentifier(e.Request().Context(), addKey)
	if err != nil {
		return err
	}

	convertedKey, err := strconv.Atoi(addKey)
	if err != nil {
		return err
	}

	resultKey := int64(float32(convertedKey) + (float32(convertedKey) * product.Increase))

	key += resultKey
	user.Point = key

	_, err = user.Update(e.Request().Context())
	if err != nil {
		return err
	}

	err = PurchaseHistory{}.AddBuyKey(e.Request().Context(), user.Id, resultKey)
	if err != nil {
		return err
	}

	return nil
}

func (UserWithProfile) GetUserWithProfile(ctx context.Context, id int64) (*UserWithProfile, error) {
	var userWithProfile UserWithProfile
	result, err := factory.DB(ctx).Where("user_id = ?", id).Join("INNER", "profile", "profile.user_id = user.id").Get(&userWithProfile)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("프로필이 존재하지 않습니다.")
	}

	return &userWithProfile, nil
}

func (User) GetByEmailWithPassword(ctx context.Context, email, password string) (*User, error) {
	var u User
	has, err := factory.DB(ctx).Where("email = ? AND join_type = ?", email, "EMAIL").Get(&u)
	if err != nil {
		return nil, err
	}

	if !has {
		fmt.Println("this")
		return nil, nil
	}

	if !CheckPasswordHash(password, u.Password) {
		return &u, errors.New("비밀번호가 틀렸습니다.")
	}

	return &u, nil
}

func (User) GetByUserId(e echo.Context, id int64) (*User, error) {
	var u User

	if has, err := factory.DB(e.Request().Context()).Where("id=?", id).Get(&u); err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("does not exists user")
	}
	return &u, nil
}

func (User) GetByEmail(e echo.Context, email string) (*User, error) {
	var u User
	if has, err := factory.DB(e.Request().Context()).Where("email = ?", email).Get(&u); err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New(email + " does not exists")
	}
	return &u, nil
}

func (User) SearchUserByNickname(ctx context.Context, nickname string) (result int64, err error) {
	result, err = factory.DB(ctx).Table(User{}).Where("nickname = ?", nickname).Count()
	return
}

func (u *User) JwtTokenCreate() (t string, err error) {
	claims := &JwtCustomClaims{
		Id:    u.Id,
		Email: u.Email,
		Admin: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 8760).Unix(),
		},
	}
	t, err = claims.Create()
	return
}

func GenerateHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return password
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
