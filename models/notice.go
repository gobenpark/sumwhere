package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type Notice struct {
	ID       int64     `json:"id" xorm:"id pk autoincr"`
	Title    string    `json:"title" xorm:"VARCHAR(255) title"`
	Text     string    `json:"text" xorm:"TEXT"`
	CreateAt time.Time `json:"create_at" xorm:"created"`
	UpdateAt time.Time `json:"update_at" xorm:"updated"`
}

func (Notice) GetAll(ctx context.Context) ([]Notice, error) {
	var n []Notice
	err := factory.DB(ctx).Find(&n)
	if err != nil {
		return nil, err
	}

	return n, nil
}
