package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type MatchMember struct {
	Id       int64     `json:"id"`
	MatchId  int64     `json:"matchId" xorm:"match_id" valid:"required"`
	UserId   int64     `json:"userId" xorm:"user_id"`
	JoinDate time.Time `xorm:"created"`
}

func (m *MatchMember) Create(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(m)
}
