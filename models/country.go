package models

import (
	"context"
	"sumwhere/factory"
)

type Country struct {
	ID       int64  `json:"id" xorm:"id pk autoincr"`
	Name     string `json:"name" xorm:"name"`
	ImageURL string `json:"imageUrl" xorm:"image_url"`
}

func (Country) GetAll(ctx context.Context) ([]Country, error) {
	var c []Country
	err := factory.DB(ctx).Find(&c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
