package models

import "time"

type Place struct {
	ID        int64     `json:"id" xorm:"id pk autoincr"`
	Name      string    `json:"name" xorm:"name"`
	Latitude  string    `json:"latitude" xorm:"latitude"`
	Longitude string    `json:"longitude" xorm:"longitude"`
	Tag       []string  `json:"tag" xorm:"tag"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}
