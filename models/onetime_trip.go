package models

import "time"

type OneTimeTrip struct {
	ID        int64     `json:"id" xorm:"id pk autoincr"`
	UserID    int64     `json:"user_id" xorm:"user_id"`
	Lat       float32   `json:"lat" xorm:"lat"`
	Lon       float32   `json:"lon" xorm:"lon"`
	StartTime time.Time `json:"start_time" xorm:"start_time"`
	CreateAt  time.Time `json:"create_at" xorm:"created"`
}
