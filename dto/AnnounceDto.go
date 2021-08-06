package dto

import "time"

type AnnounceDto struct {
	ID        uint      `json:"aid"` //公告ID
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"` //内容
	Url       string    `json:"url"`
}
