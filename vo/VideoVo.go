package vo

import (
	"time"

	"kuukaa.fun/danmu-v4/model"
)

//上传的视频
type UploadVideoVo struct {
	ID        uint      `json:"vid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Review    bool      `json:"review"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VideoVo struct {
	ID           uint         `json:"vid"`
	Title        string       `json:"title"`
	Cover        string       `json:"cover"`
	VideoType    string       `json:"video_type"`
	Desc         string       `json:"desc"`
	CreatedAt    time.Time    `json:"created_at"`
	Copyright    bool         `json:"copyright"`
	Author       UserVo       `json:"author"`
	Resource     []ResourceVo `json:"resource"`
	LikeCount    int          `json:"like_count"`
	CollectCount int          `json:"collect_count"`
	Clicks       int          `json:"clicks"`
}

//收藏的视频
type CollectVideoVo struct {
	ID    uint   `json:"vid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

//推荐视频
type RecommendVideoVo struct {
	ID     uint   `json:"vid"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	Author string `json:"author"`
	Clicks string `json:"clicks"`
}

//视频交互数据
type InteractiveVo struct {
	Collect bool `json:"collect"`
	Like    bool `json:"like"`
	Follow  bool `json:"follow"`
}

//搜索的视频
type SearchVideoVo struct {
	ID    uint   `json:"vid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

//合集里的视频
type CollectionVideoVo struct {
	ID        uint      `json:"vid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	CreatedAt time.Time `json:"created_at"`
	Desc      string    `json:"desc"`
}

//管理员视频列表
type AdminVideoListVo struct {
	ID          uint      `json:"vid"`
	Title       string    `json:"title"`
	Cover       string    `json:"cover"`
	Desc        string    `json:"desc"`
	CreatedAt   time.Time `json:"created_at"`
	Copyright   bool      `json:"copyright"`
	Uid         uint      `json:"uid"`
	VideoType   string    `json:"video_type"`
	Partition   string    `json:"partition"` //分区
	PartitionID uint      `json:"-"`
}

func ToUploadVideoVo(videos []model.Video) []UploadVideoVo {
	length := len(videos)
	newVideos := make([]UploadVideoVo, length)
	for i := 0; i < length; i++ {
		newVideos[i].ID = videos[i].ID
		newVideos[i].Title = videos[i].Title
		newVideos[i].Cover = videos[i].Cover
		newVideos[i].Review = videos[i].Review
		newVideos[i].CreatedAt = videos[i].CreatedAt
		newVideos[i].UpdatedAt = videos[i].UpdatedAt
		newVideos[i].Clicks = videos[i].Clicks
	}
	return newVideos
}

func ToVideoVo(video model.Video, like, collect int, resource []model.Resource) VideoVo {
	length := len(resource)
	newResource := make([]ResourceVo, length)
	for i := 0; i < length; i++ {
		newResource[i].ID = resource[i].UUID
		newResource[i].Title = resource[i].Title
		newResource[i].Res360 = resource[i].Res360
		newResource[i].Res480 = resource[i].Res480
		newResource[i].Res720 = resource[i].Res720
		newResource[i].Res1080 = resource[i].Res1080
		newResource[i].Original = resource[i].Original
	}
	return VideoVo{
		ID:        video.ID,
		Title:     video.Title,
		Cover:     video.Cover,
		VideoType: video.VideoType,
		Desc:      video.Desc,
		CreatedAt: video.CreatedAt,
		Copyright: video.Copyright,
		Author: UserVo{
			ID:     video.Author.ID,
			Name:   video.Author.Name,
			Sign:   video.Author.Sign,
			Avatar: video.Author.Avatar,
		},
		Resource:     newResource,
		LikeCount:    like,
		CollectCount: collect,
		Clicks:       video.Clicks,
	}
}

func ToCollectVideoVo(videos []model.Interactive) []CollectVideoVo {
	length := len(videos)
	newVideos := make([]CollectVideoVo, length)
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
