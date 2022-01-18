package vo

import (
	"time"

	"github.com/google/uuid"
	"kuukaa.fun/danmu-v4/model"
)

type ResourceVo struct {
	ID uuid.UUID `json:"id"`
	//分P使用的标题
	Title string `json:"title"`
	//不同分辨率
	Res360  string `json:"res360"`
	Res480  string `json:"res480"`
	Res720  string `json:"res720"`
	Res1080 string `json:"res1080"`
	//原始分辨率，适用于早期版本未指定分辨率的视频
	//或者不进行转码处理的情况
	Original string `json:"original"`
}

type ResourceInfoVo struct {
	UUID      uuid.UUID `json:"uuid"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func ToReviewResourceVo(resource model.Resource) ResourceVo {
	return ResourceVo{
		ID:       resource.UUID,
		Title:    resource.Title,
		Res360:   resource.Res360,
		Res480:   resource.Res480,
		Res720:   resource.Res720,
		Res1080:  resource.Res1080,
		Original: resource.Original,
	}
}
