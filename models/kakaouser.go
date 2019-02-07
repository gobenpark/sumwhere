package models

import (
	"context"
	"fmt"
	"sumwhere/factory"
	"time"
)

type KakaoUser struct {
	ID         int `json:"id"`
	Token      string
	Properties struct {
		Nickname       string `json:"nickname"`
		ProfileImage   string `json:"profile_image"`
		ThumbnailImage string `json:"thumbnail_image"`
	} `json:"properties"`
	KakaoAccount struct {
		HasEmail        bool   `json:"has_email"`
		IsEmailValid    bool   `json:"is_email_valid"`
		IsEmailVerified bool   `json:"is_email_verified"`
		Email           string `json:"email"`
		HasAgeRange     bool   `json:"has_age_range"`
		HasBirthday     bool   `json:"has_birthday"`
		HasGender       bool   `json:"has_gender"`
	} `json:"kakao_account"`
}

func (ku *KakaoUser) SearchAndCreate(ctx context.Context) (*User, error) {
	var u User
	result, err := factory.DB(ctx).Where("join_type = ? AND sns_id = ?", "KAKAO", ku.ID).Get(&u)
	if err != nil {
		return nil, err
	}

	if !result {
		user := &User{
			JoinType:  "KAKAO",
			Token:     ku.Token,
			SNSID:     fmt.Sprintf("%d", ku.ID),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if _, err := user.Create(ctx); err != nil {
			return nil, err
		}
		return user, nil
	} else {
		if u.Token != ku.Token {
			u.Token = ku.Token
			if _, err := u.Update(ctx); err != nil {
				return nil, err
			}
		}
	}

	return &u, nil
}
