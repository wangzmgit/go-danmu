package vo

import "time"

//获取审核状态的video信息
type ReviewVideoVo struct {
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Desc      string `json:"desc"`
	Partition string `json:"partition"` //分区名，3.6.8新增
}

type ReviewVideoListVo struct {
	ID        uint      `json:"vid"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	VideoType string    `json:"video_type"`
	Desc      string    `json:"desc"`
	Uid       uint      `json:"uid"`
	Copyright bool      `json:"copyright"`
	Partition string    `json:"partition"`
}
