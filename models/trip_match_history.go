package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type TripmatchHistory struct {
	UserId    int64     `xorm:"user_id"`
	TripID    int64     `xorm:"trip_id"`
	CreatedAt time.Time `xorm:"created"`
}

func (t *TripmatchHistory) Insert(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(t)
}
