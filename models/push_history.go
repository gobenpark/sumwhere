package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type PushHistory struct {
	ID       int64     `json:"id" xorm:"id pk autoincr"`
	UserID   int64     `json:"userId" xorm:"user_id"`
	Title    string    `json:"title" xorm:"title"`
	CreateAt time.Time `json:"createAt" xorm:"created"`
}

func (p *PushHistory) Insert(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(p)
}

func (PushHistory) Get(ctx context.Context, userID int64) ([]PushHistory, error) {
	var historys []PushHistory
	err := factory.DB(ctx).Where("user_id = ?", userID).Desc("created_at").Find(&historys)
	if err != nil {
		return nil, err
	}
	return historys, nil
}
