package vo

import (
	"strconv"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/util"

	"github.com/go-redis/redis"
)

//获取审核状态的video信息
type ReviewVideoVo struct {
	Title        string `json:"title"`
	Cover        string `json:"cover"`
	Introduction string `json:"introduction"`
}

//上传的视频
type UploadVideoVo struct {
	ID        uint      `json:"vid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Review    bool      `json:"review"`
	Clicks    string    `json:"clicks"`
	CreateAt  time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VideoVo struct {
	ID           uint         `json:"vid"`
	Title        string       `json:"title"`
	Cover        string       `json:"cover"`
	Video        string       `json:"video"`
	VideoType    string       `json:"video_type"`
	Introduction string       `json:"introduction"`
	CreateAt     time.Time    `json:"create_at"`
	Original     bool         `json:"original"`
	Author       UserVo       `json:"author"`
	Data         VideoData    `json:"data"`
	Clicks       string       `json:"clicks"`
	SubVideo     []SubVideoVo `json:"sub_video"`
}

//视频数据
type VideoData struct {
	LikeCount    int `json:"like_count"`
	CollectCount int `json:"collect_count"`
}

type SubVideoVo struct {
	ID    uint   `json:"vid"`
	Title string `json:"title"`
	Video string `json:"video"`
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
type InteractiveData struct {
	Collect bool `json:"collect"`
	Like    bool `json:"like"`
	Follow  bool `json:"follow"`
}

type SearchVideoVo struct {
	ID    uint   `json:"vid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

//管理员视频列表
type AdminVideoVo struct {
	ID           uint      `json:"vid"`
	Title        string    `json:"title"`
	Cover        string    `json:"cover"`
	Video        string    `json:"video"`
	VideoType    string    `json:"video_type"`
	Introduction string    `json:"introduction"`
	CreateAt     time.Time `json:"create_at"`
	Original     bool      `json:"original"`
	Uid          uint      `json:"uid"`
}

func ToUploadVideoVo(videos []model.Video) []UploadVideoVo {
	length := len(videos)
	Redis := common.RedisClient
	newVideos := make([]UploadVideoVo, length)
	for i := 0; i < length; i++ {
		newVideos[i].ID = videos[i].ID
		newVideos[i].Title = videos[i].Title
		newVideos[i].Cover = videos[i].Cover
		newVideos[i].Review = videos[i].Review
		newVideos[i].CreateAt = videos[i].CreatedAt
		newVideos[i].UpdatedAt = videos[i].UpdatedAt
		if Redis != nil {
			newVideos[i].Clicks = GetClicksFromRedis(Redis, int(videos[i].ID), strconv.Itoa(videos[i].Clicks))
		}
	}
	return newVideos
}

func ToVideoVo(video model.Video, data VideoData, subVideo []SubVideoVo) VideoVo {
	//通过ID获取视频
	//如果redis可以使用，因为先增加播放量，所以这时的播放量一定存在
	var clicks string
	if common.RedisClient != nil {
		clicks, _ = common.RedisClient.Get(util.VideoClicksKey(int(video.ID))).Result()
	}
	return VideoVo{
		ID:           video.ID,
		Title:        video.Title,
		Cover:        video.Cover,
		Video:        video.Video,
		VideoType:    video.VideoType,
		Introduction: video.Introduction,
		CreateAt:     video.CreatedAt,
		Original:     video.Original,
		Author: UserVo{
			ID:     video.Author.ID,
			Name:   video.Author.Name,
			Sign:   video.Author.Sign,
			Avatar: video.Author.Avatar,
		},
		Data:     data,
		Clicks:   clicks,
		SubVideo: subVideo, //3.3.0新增数据
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

/*********************************************************
** 函数功能: 从Redis获取点击量
** 日    期:2021/8/23
** 参    数:Redis,视频id,数据库中的播放量
**********************************************************/
func GetClicksFromRedis(redis *redis.Client, vid int, dbClicks string) string {
	strClicks, _ := redis.Get(util.VideoClicksKey(vid)).Result()
	if len(strClicks) == 0 {
		//将视频ID存入点击量列表
		redis.RPush(util.ClicksVideoList, vid)
		//将点击量存入redis并设置25小时，防止数据当天过期
		redis.Set(util.VideoClicksKey(vid), dbClicks, time.Hour*25)
		return dbClicks
	}
	return strClicks
}

func ToAdminVideoVo(videos []model.Video) []AdminVideoVo {
	length := len(videos)
	newVideos := make([]AdminVideoVo, length)
	for i := 0; i < length; i++ {
		newVideos[i].ID = videos[i].ID
		newVideos[i].Title = videos[i].Title
		newVideos[i].Cover = videos[i].Cover
		newVideos[i].Video = videos[i].Video
		newVideos[i].VideoType = videos[i].VideoType
		newVideos[i].Introduction = videos[i].Introduction
		newVideos[i].CreateAt = videos[i].CreatedAt
		newVideos[i].Original = videos[i].Original
		newVideos[i].Uid = videos[i].Uid
	}
	return newVideos
}