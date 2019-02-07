package models

import (
	"context"
	"errors"
	"sumwhere/factory"
)

type PushInput struct {
	MatchAlert  bool `json:"matchAlert"`
	FriendAlert bool `json:"friendAlert"`
	ChatAlert   bool `json:"chatAlert"`
	EventAlert  bool `json:"eventAlert"`
}

type Push struct {
	ID          int64  `json:"id" xorm:"id pk autoincr"`
	UserID      int64  `json:"userId" xorm:"user_id unique"`
	FcmToken    string `json:"fcmToken" xorm:"fcm_token"`
	MatchAlert  bool   `json:"matchAlert" xorm:"match_alert default 1"`
	FriendAlert bool   `json:"friendAlert" xorm:"friend_alert default 1"`
	ChatAlert   bool   `json:"chatAlert" xorm:"chat_alert default 1"`
	EventAlert  bool   `json:"eventAlert" xorm:"event_alert default 1"`
}

func (Push) Get(ctx context.Context, userID int64) (*Push, error) {
	var p Push
	result, err := factory.DB(ctx).Where("user_id = ?", userID).Get(&p)
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, errors.New("Not Exist Push")
	}

	return &p, nil
}

func (p *Push) Update(ctx context.Context) error {
	_, err := factory.DB(ctx).ID(p.ID).UseBool("match_alert", "friend_alert", "chat_alert", "event_alert").Update(p)

	if err != nil {
		return err
	}
	return nil
}

func (p *Push) Insert(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(p)

	if err != nil {
		return err
	}
	return nil
}
