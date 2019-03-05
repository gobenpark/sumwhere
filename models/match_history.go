package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type MatchRequestDTO struct {
	TripID   int64 `json:"tripId" valid:"required"`
	ToTripID int64 `json:"toTripId" valid:"required"`
	ToUserID int64 `json:"toUserId" valid:"required"`
}

func (m MatchRequestDTO) ToModel(userID int64) *MatchHistory {
	return &MatchHistory{
		UserID:   userID,
		TripID:   m.TripID,
		ToUserID: m.ToUserID,
		ToTripID: m.ToTripID,
	}
}

type MatchHistory struct {
	ID       int64     `json:"id" xorm:"id pk autoincr"`
	UserID   int64     `json:"userId" xorm:"user_id"`
	TripID   int64     `json:"tripId" xorm:"trip_id"`
	ToUserID int64     `json:"toUserId" xorm:"to_user_id"`
	ToTripID int64     `json:"toTripId" xorm:"to_trip_id"`
	CreateAt time.Time `json:"createAt" xorm:"created"`
}

func (m *MatchHistory) Insert(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(m)
	if err != nil {
		return err
	}
	return nil
}

func (MatchHistory) GetRequest(ctx context.Context) error {
	return nil
}
