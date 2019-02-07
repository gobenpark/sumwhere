package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type Report struct {
	ID           int64     `json:"id" xorm:"id pk autoincr" valid:"-"`
	UserID       int64     `json:"user_id" xorm:"user_id" valid:"-"`
	TargetUserID int64     `json:"target_user_id" xorm:"target_user_id"`
	ReportType   int64     `json:"report_type" xorm:"report_type"`
	Comment      string    `json:"comment" xorm:"TEXT comment"`
	CreateAt     time.Time `json:"create_at" xorm:"created"`
	UpdateAt     time.Time `json:"update_at" xorm:"updated"`
}

func (r *Report) Insert(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(r)
	if err != nil {
		return err
	}
	return nil
}
