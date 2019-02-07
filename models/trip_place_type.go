package models

import (
	"context"
	"github.com/go-xorm/xorm"
	"strconv"
	"sumwhere/factory"
	"sumwhere/middlewares"
)

const TOPTRIP = "toptrip"
const RECENTTOPTRIP = "recenttoptrip"

type TripRank struct {
	TripPlaceType
	Rank float64
}

type TripPlaceType struct {
	Id       int64  `json:"id" xorm:"id pk"`
	Trip     string `json:"trip"`
	Country  string `json:"country" xorm:"country"`
	ImageURL string `json:"imageURL" xorm:"image_url"`
}

func (TripPlaceType) Search(ctx context.Context, destination string) (t []TripPlaceType, err error) {
	err = factory.DB(ctx).Where("trip like ?", "%"+destination+"%").Find(&t)
	return
}

func (TripPlaceType) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (items []TripPlaceType, err error) {

	queryBuilder := func() xorm.Interface {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)

	go func() {
		if err := queryBuilder().Limit(limit, offset).Find(&items); err != nil {
			errc <- err
			return
		}
		errc <- nil
	}()

	//if err := <-errc; err != nil {
	//	return 0, nil, err
	//}
	if err := <-errc; err != nil {
		return nil, err
	}
	return
}

func (TripPlaceType) TopTripPlaces(ctx context.Context) ([]TripRank, error) {
	tripIds, err := factory.Redis(ctx, middlewares.ContextGetRedisName).ZRevRangeWithScores("toptrip", 0, 9).Result()
	if err != nil {
		return nil, err
	}

	trips := make([]TripRank, 0)

	for _, z := range tripIds {
		var t TripPlaceType
		member := z.Member.(string)
		rowId, _ := strconv.Atoi(member)

		result, err := factory.DB(ctx).ID(rowId).Get(&t)
		if err != nil || !result {
			continue
		}

		rank := TripRank{
			TripPlaceType: t,
			Rank:          z.Score,
		}
		trips = append(trips, rank)
	}

	return trips, nil
}
