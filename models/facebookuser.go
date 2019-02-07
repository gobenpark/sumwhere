package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type FaceBookUser struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string
}

func (fu *FaceBookUser) SearchAndCreate(ctx context.Context) (*User, error) {
	var u User
	result, err := factory.DB(ctx).Where("join_type = ? AND sns_id = ?", "FACEBOOK", fu.ID).Get(&u)
	if err != nil {
		return nil, err
	}

	if !result {
		user := &User{
			Email:     fu.Email,
			Username:  fu.Name,
			SNSID:     fu.ID,
			JoinType:  "FACEBOOK",
			Token:     fu.Token,
			Point:     50,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if _, err := user.Create(ctx); err != nil {
			return nil, err
		}
		return user, nil
	} else {
		if u.Token != fu.Token {
			u.Token = fu.Token
			if _, err := u.Update(ctx); err != nil {
				return nil, err
			}
		}
	}

	return &u, nil
}
