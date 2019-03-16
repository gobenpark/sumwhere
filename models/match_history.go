package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type joinedModel struct {
	Trip      `json:"trip" xorm:"extends"`
	TripPlace `json:"tripPlace" xorm:"extends"`
	User      `json:"user" xorm:"extends"`
	Profile   `json:"profile" xorm:"extends"`
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
	Accept   bool      `json:"accept" xorm:"accept default 0"`
	DeleteAt time.Time `xorm:"deleted"`
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
		Join("INNER", "trip", "match_history.trip_id = trip.id").
		Join("INNER", "trip_place", "match_history.trip_place_id = trip_place.id").
		Join("INNER", "user", "match_history.to_user_id = user.id").
		Join("INNER", "profile", "profile.user_id = match_history.to_user_id").
		Where("match_history.user_id = ?", userID).
		Find(&models)
	if err != nil {
		return nil, err
	}

	return &models, nil
}

func (MatchHistory) GetReceive(ctx context.Context, userID int64) (*[]joinedModel, error) {
	var models []joinedModel

	err := factory.DB(ctx).
		Table("match_history").
		Join("INNER", "trip", "match_history.trip_id = trip.id").
		Join("INNER", "trip_place", "match_history.trip_place_id = trip_place.id").
		Join("INNER", "user", "match_history.to_user_id = user.id").
		Join("INNER", "profile", "profile.user_id = match_history.to_user_id").
		Where("match_history.to_user_id = ?", userID).
		Find(&models)
	if err != nil {
		return nil, err
	}
	return &models, nil
}
