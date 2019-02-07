package models

import (
	"context"
	"fmt"
	"sumwhere/factory"
	"time"
)

type ChatRoomJoin struct {
	ChatRoom   `json:"chatRoom" xorm:"extends"`
	ChatMember `json:"chatMember" xorm:"extends"`
}

type ChatRoom struct {
	Id       int64     `json:"id" xorm:"id pk autoincr"`
	CreateAt time.Time `json:"createAt" xorm:"created"`
	UpdateAt time.Time `json:"updateAt" xorm:"updated"`
}

func (ChatRoomJoin) TableName() string {
	return "chat_room"
}

func (c *ChatRoom) Insert(ctx context.Context) (int64, error) {
	result, err := factory.DB(ctx).Insert(c)
	return result, err
}

func (ChatRoomJoin) GetRoom(ctx context.Context, userId int64) ([]ChatRoomJoin, error) {
	fmt.Println(userId)
	var c []ChatRoomJoin
	err := factory.DB(ctx).
		Join("LEFT", "chat_member", "chat_room.id = chat_member.room_id").
		Where("chat_member.user_id = ?", userId).Find(&c)

	if err != nil {
		return nil, err
	}
	return c, nil
}

//select * from chat_room left join chat_member m on chat_room.id = m.room_id where m.user_id = 91
