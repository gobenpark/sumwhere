package models

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	Id    int64  `json:"id"`
	Email string `json:"email" valid:"required"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func (j JwtCustomClaims) Create() (jwtoken string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	jwtoken, err = token.SignedString([]byte("parkbumwoo"))
	return
}

func (j JwtCustomClaims) ToUser() *User {
	return &User{
		Email: j.Email,
		Admin: j.Admin,
	}
}
