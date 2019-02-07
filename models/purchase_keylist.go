package models

import (
	"context"
	"errors"
	"sumwhere/factory"
)

type PurchaseProduct struct {
	ShowName    string  `json:"showName" xorm:"show_name"`
	ProductName string  `json:"productName" xorm:"product_name"`
	Increase    float32 `json:"increase" xorm:"increase"`
	Price       float32 `json:"price" xorm:"price"`
}

func (PurchaseProduct) GetByIdentifier(ctx context.Context, productName string) (*PurchaseProduct, error) {
	var model PurchaseProduct
	result, err := factory.DB(ctx).Where("product_name=?", productName).Get(&model)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("물품이 존재하지 않습니다.")
	}

	return &model, nil

}
