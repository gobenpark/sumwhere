package models

import (
	"context"
	"errors"
	"github.com/go-xorm/xorm"
	"sumwhere/factory"
	"time"
)

type Match struct {
	Id          int64     `json:"id"`
	UserID      int64     `json:"userId" xorm:"user_id"`
	MatchTypeId int64     `json:"matchTypeId" xorm:"match_type_id"`
	MakerId     int64     `json:"makerId" xorm:"maker_id" valid:"required"`
	CreateAt    time.Time `json:"craeteAt" xorm:"created"`
	UpdateAt    time.Time `json:"updateAt" xorm:"updated"`
}

func (m *Match) Create(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(m)
}

func (m *Match) Update(ctx context.Context) (int64, error) {
	return factory.DB(ctx).ID(m.Id).Update(m)
}

func (Match) TotalCount(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Count(Match{})
}

func (Match) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []Match, err error) {
	queryBuilder := func() xorm.Interface {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)
	go func() {
		v, err := queryBuilder().Count(&Match{})
		if err != nil {
			errc <- err
			return
		}
		totalCount = v
		errc <- nil
	}()

	go func() {
		if err := queryBuilder().Limit(limit, offset).Find(&items); err != nil {
			errc <- err
			return
		}
		errc <- nil
	}()

	if err := <-errc; err != nil {
		return 0, nil, err
	}
	if err := <-errc; err != nil {
		return 0, nil, err
	}
	return
}

func setSortOrder(q xorm.Interface, sortby, order []string) error {
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				if order[i] == "desc" {
					q.Desc(v)
				} else if order[i] == "asc" {
					q.Asc(v)
				} else {
					return errors.New("Invalid order. Must be either [asc|desc]")
				}
			}
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				if order[0] == "desc" {
					q.Desc(v)
				} else if order[0] == "asc" {
					q.Asc(v)
				} else {
					return errors.New("Invalid order. Must be either [asc|desc]")
				}
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return errors.New("'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return errors.New("unused 'order' fields")
		}
	}
	return nil
}
