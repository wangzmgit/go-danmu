package dto

import (
	"strconv"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/util"
)

//视频数据
type VideoData struct {
	LikeCount    int `json:"like_count"`
	CollectCount int `json:"collect_count"`
}

type UploadVideoDto struct {
	ID        uint      `json:"vid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Review    bool      `json:"review"`
	Clicks    int       `json:"clicks"`
	CreateAt  time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VideoDto struct {
	ID           uint      `json:"vid"`
	Title        string    `json:"title"`
	Cover        string    `json:"cover"`
	Video        string    `json:"video"`
	Introduction string    `json:"introduction"`
	CreateAt     time.Time `json:"create_at"`
	Original     bool      `json:"original"`
	Author       UserDto   `json:"author"`
	Data         VideoData `json:"data"`
	Clicks       string    `json:"clicks"`
}

type CollectVideoDto struct {
	ID    uint   `json:"vid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

type SearchVideoDto struct {
	ID    uint   `json:"vid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

func ToUploadVideoDto(videos []model.Video) []UploadVideoDto {
	length := len(videos)
	newVideos := make([]UploadVideoDto, length)
	for i := 0; i < length; i++ {
		newVideos[i].ID = videos[i].ID
		newVideos[i].Title = videos[i].Title
		newVideos[i].Cover = videos[i].Cover
		newVideos[i].Review = videos[i].Review
		newVideos[i].CreateAt = videos[i].CreatedAt
		newVideos[i].UpdatedAt = videos[i].UpdatedAt
		newVideos[i].Clicks = videos[i].Clicks
		//从redis中拉取数据
		strClicks, _ := common.RedisClient.Get(util.VideoClicksKey(int(videos[i].ID))).Result()
		if strClicks == "" {
			//将视频ID存入点击量列表
			common.RedisClient.RPush(util.ClicksVideoList, videos[i].ID)
			//将点击量存入redis并设置25小时，防止数据当天过期
			common.RedisClient.Set(util.VideoClicksKey(int(videos[i].ID)), videos[i].Clicks, time.Hour*25)
		} else {
			newVideos[i].Clicks, _ = strconv.Atoi(strClicks)
		}
	}
	return newVideos
}

func ToVideoDto(video model.Video, data VideoData) VideoDto {
	//通过ID获取视频
	//因为先增加播放量，所以这时的播放量一定存在
	clicks, _ := common.RedisClient.Get(util.VideoClicksKey(int(video.ID))).Result()
	return VideoDto{
		ID:           video.ID,
		Title:        video.Title,
		Cover:        video.Cover,
		Video:        video.Video,
		Introduction: video.Introduction,
		CreateAt:     video.CreatedAt,
		Original:     video.Original,
		Author: UserDto{
			ID:     video.Author.ID,
			Name:   video.Author.Name,
			Sign:   video.Author.Sign,
			Avatar: video.Author.Avatar,
		},
		Data:   data,
		Clicks: clicks,
	}
}

func ToCollectVideoDto(videos []model.Interactive) []CollectVideoDto {
	length := len(videos)
	newVideos := make([]CollectVideoDto, length)
	for i := 0; i < length; i++ {
		if videos[i].Video.Review {
			newVideos[i].ID = videos[i].Video.ID
			newVideos[i].Title = videos[i].Video.Title
			newVideos[i].Cover = videos[i].Video.Cover
		} else {
			newVideos[i].ID = videos[i].Video.ID
		}
	}
	return newVideos
}
