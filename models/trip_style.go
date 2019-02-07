package models

import (
	"context"
	"sumwhere/factory"
)

type TripStyle struct {
	ID   int64  `json:"id" xorm:"id pk autoincr"`
	Type string `json:"type" xorm:"type not null"`
	Name string `json:"name" xorm:"name notnull"`
}

func (TripStyle) GetAll(ctx context.Context) ([]TripStyle, error) {
	var trips []TripStyle
	err := factory.DB(ctx).Find(&trips)
	if err != nil {
		return nil, err
	}
	return trips, nil
}
