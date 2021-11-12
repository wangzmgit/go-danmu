package vo

import "time"

type CommentVo struct {
	ID         uint      `json:"cid"` //评论ID
	CreatedAt  time.Time `json:"created_at"`
	Content    string    `json:"content"` //内容
	Uid        uint      `json:"uid"`
	Name       string    `json:"name"`
	Avatar     string    `json:"avatar"`
	Reply      []ReplyVo `json:"reply"`
	ReplyCount int       `json:"reply_count"`
}

type ReplyVo struct {
	ID        uint      `json:"rid"` //回复id
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"` //内容
	Uid       uint      `json:"uid"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	ReplyUid  uint      `json:"reply_uid"`  //回复的人的uid
	ReplyName string    `json:"reply_name"` //回复的人的昵称
}
