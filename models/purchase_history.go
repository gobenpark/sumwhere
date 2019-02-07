package models

import (
	"context"
	"sumwhere/factory"
	"time"
)

type PurchaseHistory struct {
	ID            int64     `json:"id" xorm:"id pk autoincr"`
	UserID        int64     `json:"user_id" xorm:"user_id"`
	Message       string    `json:"message" xorm:"message"`
	PositiveValue bool      `json:"positive_value"`
	Key           int32     `json:"key"`
	CreateAt      time.Time `json:"create_at" xorm:"created"`
}

func (p PurchaseHistory) GetUserHistory(ctx context.Context, userID int64, pageNum int) ([]PurchaseHistory, error) {

	var historys []PurchaseHistory
	if err := factory.DB(ctx).Where("user_id = ?", userID).Desc("create_at").Limit(10, 10*pageNum).Find(&historys); err != nil {
		return nil, err
	}

	return historys, nil
}

func (p PurchaseHistory) AddBuyKey(ctx context.Context, userID int64, key int64) error {
	history := &PurchaseHistory{
		UserID:        userID,
		Message:       "키 구매",
		PositiveValue: true,
		Key:           int32(key),
	}

	_, err := factory.DB(ctx).Insert(history)
	if err != nil {
		return err
	}
	return nil
}
