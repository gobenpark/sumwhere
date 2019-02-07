package models

import (
	"context"
	"sumwhere/factory"
)

type Character struct {
	Id       int64  `json:"id" xorm:"id pk autoincr"`
	TypeName string `json:"typeName" xorm:"type_name"`
}

func (Character) GetAll(ctx context.Context) (characters []Character, err error) {
	err = factory.DB(ctx).Find(&characters)
	return
}
