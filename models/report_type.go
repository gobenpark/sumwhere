package models

type ReportType struct {
	ID   int64  `json:"id" xorm:"id pk autoincr"`
	Name string `json:"name" xorm:"name not null"`
}
