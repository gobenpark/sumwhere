package models

import (
	"context"
	"sumwhere/factory"
)

type MatchType struct {
	ID       int64  `json:"id" xorm:"id pk autoincr"`
	Title    string `json:"title" xorm:"title"`
	SubTitle string `json:"subTitle" xorm:"sub_title"`
	ImageURL string `json:"imageUrl" xorm:"image_url"`
	IsEnable bool   `json:"isEnable" xorm:"is_enable"`
}

func (MatchType) GetAll(ctx context.Context) ([]MatchType, error) {
	var m []MatchType
	err := factory.DB(ctx).Find(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
