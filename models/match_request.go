package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type MatchRequestJoin struct {
	FromMatchModel Trip `json:"fromMatchModel" xorm:"extends"`
	ToMatchModel   Trip `json:"toMatchModel" xorm:"extends"`
}

type MatchRequest struct {
	FromMatchId int64     `json:"fromMatchId" xorm:"from_match_id"`
	ToMatchId   int64     `json:"toMatchId" xorm:"to_match_id"`
	Accepted    bool      `json:"accepted" xorm:"accepted default 0"`
	CreateAt    time.Time `json:"createAt" xorm:"create_at created"`
}

func (m *MatchRequest) Insert(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(m)
}

func (MatchRequestJoin) FindReceiveRequest(ctx context.Context, userId int64) (m []MatchRequestJoin, err error) {
	err = factory.DB(ctx).Table("match_request").
		Join("INNER", "trip t1", "to_match_id = t1.id AND t1.user_id = ?", userId).
		Join("INNER", "trip t2", "from_match_id = t2.id").Find(&m)
	return
}

func (MatchRequestJoin) FindSendRequest(ctx context.Context, userID int64) (m []MatchRequestJoin, err error) {
	err = factory.DB(ctx).Table("match_request").
		Join("INNER", "trip t1", "from_match_id = t1.id AND t1.user_id = ?", userID).
		Join("INNER", "trip t2", "to_match_id = t2.id").Find(&m)
	return
}

func (MatchRequest) Get(ctx context.Context, from, to int64) (*MatchRequest, error) {
	var m MatchRequest
	result, err := factory.DB(ctx).Where("from_match_id = ?", from).And("to_match_id = ?", to).Get(&m)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, nil
	}

	return &m, nil
}
