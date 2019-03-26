package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type joinedModel struct {
	MatchHistory `json:"matchHistory" xorm:"extends"`
	Trip         `json:"trip" xorm:"extends"`
	TripPlace    `json:"tripPlace" xorm:"extends"`
	User         `json:"user" xorm:"extends"`
	Profile      `json:"profile" xorm:"extends"`
}

type MatchRequestDTO struct {
	TripPlaceID int64 `json:"tripPlaceId" valid:"required"`
	TripID      int64 `json:"tripId" valid:"required"`
	ToTripID    int64 `json:"toTripId" valid:"required"`
	ToUserID    int64 `json:"toUserId" valid:"required"`
}

func (m MatchRequestDTO) ToModel(userID int64) *MatchHistory {
	return &MatchHistory{
		UserID:      userID,
		TripPlaceID: m.TripPlaceID,
		TripID:      m.TripID,
		ToUserID:    m.ToUserID,
		ToTripID:    m.ToTripID,
	}
}

type MatchHistory struct {
	ID          int64     `json:"id" xorm:"id pk autoincr"`
	UserID      int64     `json:"userId" xorm:"user_id"`
	TripID      int64     `json:"tripId" xorm:"trip_id"`
	TripPlaceID int64     `json:"tripPlaceId" xorm:"trip_place_id"`
	ToUserID    int64     `json:"toUserId" xorm:"to_user_id"`
	ToTripID    int64     `json:"toTripId" xorm:"to_trip_id"`
	State       string    `json:"state" xorm:"state VARCHAR(255) default NONE"`
	DeletedAt   time.Time `xorm:"deleted"`
	CreatedAt   time.Time `json:"createdAt" xorm:"created"`
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
		Join("INNER", "user", "match_history.user_id = user.id").
		Join("INNER", "profile", "profile.user_id = match_history.user_id").
		Where("match_history.to_user_id = ?", userID).
		And("match_history.state = ?", "NONE").
		Find(&models)
	if err != nil {
		return nil, err
	}
	return &models, nil
}

func (MatchHistory) StateUpdate(ctx context.Context, historyID int64) (*MatchHistory, error) {
	var history MatchHistory
	result, err := factory.DB(ctx).ID(historyID).Get(&history)
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, nil
	}

	history.State = "ACCEPT"

	_, err = factory.DB(ctx).Update(&history)
	if err != nil {
		return nil, err
	}

	return &history, nil
}
