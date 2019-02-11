package models

import (
	"context"
	"strconv"
	"sumwhere/factory"
	"sumwhere/middlewares"
)

const TOPTRIP = "toptrip"
const RECENTTOPTRIP = "recenttoptrip"

type TripRank struct {
	TripPlace
	Rank float64
}

type TripPlace struct {
	Id          int64  `json:"id" xorm:"id pk"`
	Trip        string `json:"trip"`
	Discription string `json:"discription" xorm:"discription"`
	CountryID   int64  `json:"countryId" xorm:"country_id"`
	ImageURL    string `json:"imageURL" xorm:"image_url"`
}

func (TripPlace) Search(ctx context.Context, destination string) (t []TripPlace, err error) {
	err = factory.DB(ctx).Where("trip like ?", "%"+destination+"%").Find(&t)
	return
}

func (TripPlace) GetAll(ctx context.Context, countryID int64) ([]TripPlace, error) {
	var t []TripPlace
	err := factory.DB(ctx).Where("country_id = ?", countryID).Find(&t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (TripPlace) TopTripPlaces(ctx context.Context) ([]TripRank, error) {
	tripIds, err := factory.Redis(ctx, middlewares.ContextGetRedisName).ZRevRangeWithScores("toptrip", 0, 9).Result()
	if err != nil {
		return nil, err
	}

	trips := make([]TripRank, 0)

	for _, z := range tripIds {
		var t TripPlace
		member := z.Member.(string)
		rowId, _ := strconv.Atoi(member)

		result, err := factory.DB(ctx).ID(rowId).Get(&t)
		if err != nil || !result {
			continue
		}

		rank := TripRank{
			TripPlace: t,
			Rank:      z.Score,
		}
		trips = append(trips, rank)
	}

	return trips, nil
}
