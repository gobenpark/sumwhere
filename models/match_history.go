package models

import (
	"context"
	"fmt"
	"sumwhere/factory"
	"time"
)

type joinedModel struct {
	MatchHistory `json:"matchHistory"`
	Profile      `json:"profile"`
	User         `json:"user"`
}

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

func (MatchHistory) GetRequest(ctx context.Context, userID int64) (*[]joinedModel, error) {

	var models []joinedModel

	err := factory.DB(ctx).
		Table("match_history").
		Join("INNER", "profile", "profile.user_id = match_history.user_id").
		Join("INNER", "user", "match_history.to_user_id = user.id").
		Where("match_history.user_id = ?", userID).
		Find(&models)
	if err != nil {
		return nil, err
	}
	fmt.Println(models)

	return &models, nil
}
