package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type Event struct {
	ID       int64     `json:"id" xorm:"id pk autoincr"`
	ImageURL string    `json:"image_url" xorm:"image_url"`
	Title    string    `json:"title" xorm:"VARCHAR(255) title"`
	Text     string    `json:"text" xorm:"TEXT"`
	StartAt  time.Time `json:"startAt" xorm:"DATETIME start_at"`
	EndAt    time.Time `json:"endAt" xorm:"DATETIME end_at"`
}

func (Event) GetAll(ctx context.Context) ([]Event, error) {
	var e []Event
	err := factory.DB(ctx).Find(&e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
