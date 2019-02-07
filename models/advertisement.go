package models

import "time"

type Advertisement struct {
	ID       int64     `json:"id" xorm:"id pk autoincr"`
	Image    string    `json:"image" xorm:"image"`
	CreateAt time.Time `json:"create_at" xorm:"created"`
	UpdateAt time.Time `json:"update_at" xorm:"updated"`
}
