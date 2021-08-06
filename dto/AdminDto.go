package dto

import (
	"time"
	"wzm/danmu3.0/model"
)

type AdminListDto struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Authority int    `json:"authority"`
}

type AdminVideoDto struct {
	ID           uint      `json:"vid"`
	Title        string    `json:"title"`
	Cover        string    `json:"cover"`
	Video        string    `json:"video"`
	Introduction string    `json:"introduction"`
	CreateAt     time.Time `json:"create_at"`
	Original     bool      `json:"original"`
	Uid          uint      `json:"uid"`
}

type AdminAnnounceDto struct {
	ID        uint      `json:"aid"` //公告ID
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"` //内容
	Url       string    `json:"url"`
}

type AdminCarouselDto struct {
	ID        uint      `json:"id"`
	Img       string    `json:"img"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

func ToAdminVideoDto(videos []model.Video) []AdminVideoDto {
	length := len(videos)
	newVideos := make([]AdminVideoDto, length)
	for i := 0; i < length; i++ {
		newVideos[i].ID = videos[i].ID
		newVideos[i].Title = videos[i].Title
		newVideos[i].Cover = videos[i].Cover
		newVideos[i].Video = videos[i].Video
		newVideos[i].Introduction = videos[i].Introduction
		newVideos[i].CreateAt = videos[i].CreatedAt
		newVideos[i].Original = videos[i].Original
		newVideos[i].Uid = videos[i].Uid
	}
	return newVideos
}
