package models

type PushType struct {
	ID       int64  `json:"id" xorm:"id pk autoincr"`
	Name     string `json:"name" xorm:"name"`
	ImageURL string `json:"imageURL" xorm:"image_url"`
}
