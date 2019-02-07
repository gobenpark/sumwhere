package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type Banner struct {
	ID       int64     `json:"id" xorm:"id pk autoincr"`
	ImageURL string    `json:"image_url" xorm:"image_url"`
	CreateAt time.Time `json:"create_at" xorm:"created"`
	UpdateAt time.Time `json:"update_at" xorm:"updated"`
}

func (Banner) GetAll(ctx context.Context) ([]Banner, error) {
	var b []Banner
	err := factory.DB(ctx).Find(&b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
