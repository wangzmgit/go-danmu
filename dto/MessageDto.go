package dto

import "time"

type MessagesListDto struct {
	ID        uint      `json:"id"` //消息ID
	CreatedAt time.Time `json:"created_at"`
	Uid       uint      `json:"uid"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
}

type MessageDetailsDto struct {
	Fid       uint      `json:"fid"`
	FromId    uint      `json:"from_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
