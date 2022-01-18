package vo

import "time"

//获取审核状态的video信息
type ReviewVideoVo struct {
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Desc      string `json:"desc"`
	Partition string `json:"partition"` //分区名，3.6.8新增
}

//管理员视频列表
type AdminVideoVo struct {
	ID        uint      `json:"vid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Desc      string    `json:"desc"`
	CreateAt  time.Time `json:"create_at"`
	Copyright bool      `json:"copyright"`
	Uid       uint      `json:"uid"`
}
