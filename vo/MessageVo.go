package vo

import "time"

type MessagesListVo struct {
	ID        uint      `json:"id"` //消息ID
	CreatedAt time.Time `json:"created_at"`
	Uid       uint      `json:"uid"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Status    bool      `json:"status"` //已读状态
}

type MessageDetailsVo struct {
	Fid       uint      `json:"fid"`
	FromId    uint      `json:"from_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

//单个消息
type MessageVo struct {
	Fid     uint   `json:"fid"`
	Content string `json:"content"`
}
