package models

import (
	"context"
	"sumwhere/factory"
)

type TripStyle struct {
	ID       int64  `json:"id" xorm:"id pk autoincr"`
	Type     string `json:"type" xorm:"type not null"`
	Name     string `json:"name" xorm:"name notnull"`
	ImageURL string `json:"imageUrl" xorm:"image_url"`
}

func (TripStyle) GetAll(ctx context.Context) ([]TripStyle, error) {
	var trips []TripStyle
	err := factory.DB(ctx).Find(&trips)
	if err != nil {
		return nil, err
	}
	return trips, nil
}

func (TripStyle) GetFromIDS(ctx context.Context, ids []string) ([]TripStyle, error) {

	var styles []TripStyle
	query := factory.DB(ctx).Where("id = ?", ids[0])

	for _, id := range ids[1:] {
		query.Or("id = ?", id)
	}

	if err := query.Find(&styles); err != nil {
		return nil, err
	}

	return styles, nil
}
