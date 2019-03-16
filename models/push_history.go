package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type PushHistoryJoin struct {
	PushHistory `json:"pushHistory" xorm:"extends"`
	PushType    `json:"pushType" xorm:"extends"`
}

type PushHistory struct {
	ID        int64     `json:"id" xorm:"id pk autoincr"`
	UserID    int64     `json:"userId" xorm:"user_id"`
	TypeID    int64     `json:"typeId" xorm:"type_id"`
	Title     string    `json:"title" xorm:"title"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
}

func (p *PushHistory) Insert(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(p)
}

func (PushHistory) Get(ctx context.Context, userID int64) ([]PushHistoryJoin, error) {
	var historys []PushHistoryJoin
	err := factory.DB(ctx).Table("push_history").Join("INNER", "push_type", "push_history.type_id = push_type.id").Where("user_id = ?", userID).Desc("created_at").Find(&historys)
	if err != nil {
		return nil, err
	}
	return historys, nil
}
