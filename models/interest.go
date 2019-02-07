package models

import (
	"context"
	"sumwhere/factory"
)

type Interest struct {
	Id       int64  `json:"id" xorm:"id pk autoincr"`
	TypeName string `json:"typeName" xorm:"type_name"`
}

func (Interest) GetAllInterest(ctx context.Context) (c []Interest, err error) {
	err = factory.DB(ctx).Find(&c)
	return
}
