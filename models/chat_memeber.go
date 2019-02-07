package models

type ChatMember struct {
	Id         int64 `json:"id" xorm:"id pk autoincr"`
	ChatRoomId int64 `json:"chatRoomId" xorm:"room_id"`
	UserId     int64 `json:"userId" xorm:"user_id"`
}
