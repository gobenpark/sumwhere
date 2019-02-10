package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"sumwhere/factory"
	"time"
)

type TripUserGroup struct {
	User `json:"user" xorm:"extends"`
	Trip `json:"trip" xorm:"extends"`
}

func (TripUserGroup) TableName() string {
	return "user"
}

type TripGroup struct {
	Trip      `json:"trip" xorm:"extends"`
	TripPlace `json:"tripType" xorm:"extends"`
}

func (TripGroup) TableName() string {
	return "trip"
}

type Trip struct {
	Id          int64     `json:"id" xorm:"id pk autoincr"`
	UserId      int64     `json:"userId" xorm:"user_id" valid:"required"`
	MatchTypeId int64     `json:"matchTypeId" xorm:"match_type_id"`
	Concept     string    `json:"concept" xorm:"concept not null"`
	TripTypeId  int64     `json:"tripTypeId" xorm:"triptype_id"`
	GenderType  string    `json:"genderType" xorm:"gender_type VARCHAR(20)"`
	StartDate   time.Time `json:"startDate" xorm:"start_date"`
	EndDate     time.Time `json:"endDate" xorm:"end_date"`
	CreateAt    time.Time `json:"createAt" xorm:"created"`
	UpdateAt    time.Time `json:"updateAt" xorm:"updated"`
	DeleteAt    time.Time `xorm:"deleted"`
}

func (t *Trip) Create(ctx context.Context) (int64, error) {
	//tripId := strconv.FormatInt(t.TripTypeId, 10)
	//if factory.Redis(ctx) != nil {
	//	factory.Redis(ctx).ZIncrBy(TOPTRIP, 1, tripId)
	//}
	return factory.DB(ctx).Insert(t)
}

func (t *Trip) Delete(ctx context.Context) (int64, error) {
	return factory.DB(ctx).ID(t.Id).Delete(t)
}

func (t *Trip) Update(ctx context.Context, id int64) error {
	_, err := factory.DB(ctx).ID(id).Update(t)
	if err != nil {
		return err
	}
	return nil
}

func (Trip) Get(ctx context.Context, tripId, userId int64) (*Trip, error) {
	var t Trip

	result, err := factory.DB(ctx).ID(tripId).Where("user_id = ?", userId).Get(&t)
	if err != nil {
		return nil, err
	} else if !result {
		return nil, errors.New("해당 데이터가 존재하지 않습니다.")
	} else {
		return &t, nil
	}
}

func (TripGroup) GetMyTrip(ctx context.Context, id int64) (*TripGroup, error) {
	var item TripGroup
	result, err := factory.DB(ctx).Where("user_id = ?", id).Join("INNER", "trip_place", "trip_place.id = trip.triptype_id").Get(&item)
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, nil
	}
	return &item, nil
}

func (TripGroup) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []TripGroup, err error) {

	queryBuilder := func() xorm.Interface {
		q := factory.DB(ctx).Join("INNER", "trip_place", "trip_place.id = trip.triptype_id")
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)

	//go func() {
	//	v, err := queryBuilder().Count(&TravelGroup{})
	//	if err != nil {
	//		errc <- err
	//		return
	//	}
	//	totalCount = v
	//	errc <- nil
	//}()

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
		return 0, nil, err
	}
	return
}

func (Trip) GetMyTrips(ctx context.Context, id int64) (t []Trip, err error) {
	err = factory.DB(ctx).Where("user_id = ?", id).Find(&t)
	return
}

func (Trip) Exist(ctx context.Context, id int64, query string) (int64, error) {
	result, err := factory.DB(ctx).Where("user_id = ?", id).And(query).Count(&Trip{})
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (TripUserGroup) Join(ctx context.Context, trip *Trip, count int) (tripGroup []TripUserGroup, err error) {

	startDate := trip.StartDate.Format("2006-01-02")
	endDate := trip.EndDate.Format("2006-01-02")

	query := fmt.Sprintf("SELECT user.*, trip.* "+
		"FROM user LEFT JOIN (trip LEFT OUTER JOIN tripmatch_history "+
		"ON trip.id = tripmatch_history.trip_id) on user.id = trip.user_id "+
		"WHERE (tripmatch_history.trip_id is null) "+
		"AND (trip.user_id != %d) "+
		"AND (user.gender = 'male') "+
		"AND (DATE(start_date) BETWEEN '%s' AND '%s' OR DATE(end_date) BETWEEN '%s' AND '%s') "+
		"limit 0,%d", trip.UserId, startDate, endDate, startDate, endDate, count)

	err = factory.DB(ctx).SQL(query).Find(&tripGroup)
	return
}

func (t *TripUserGroup) InsertHistory(ctx context.Context, userID int64) (int64, error) {
	h := TripmatchHistory{
		UserId: userID,
		TripID: t.Trip.Id,
	}

	return h.Insert(ctx)
}
