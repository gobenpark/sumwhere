package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type Profile struct {
	Id            int64       `json:"id" xorm:"id pk autoincr"`
	UserId        int64       `json:"userId" xorm:"user_id"`
	Age           int         `json:"age" xorm:"age"`
	Job           string      `json:"job" xorm:"job"`
	CharacterType []Character `json:"characterType"`
	TripStyleType string      `json:"tripStyleType"`
	Image1        string      `json:"image1"`
	Image2        string      `json:"image2"`
	Image3        string      `json:"image3"`
	Image4        string      `json:"image4"`
	CreateAt      time.Time   `json:"createAt" xorm:"created"`
	UpdateAt      time.Time   `json:"updateAt" xorm:"updated"`
}

func (p *Profile) Create(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(p)
}
