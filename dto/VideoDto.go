package dto

type UploadVideoDto struct {
	ID        uint
	Title     string
	Cover     string
	Desc      string
	Copyright bool
	Partition uint //分区ID，3.6.8新增
}

//修改视频
type ModifyVideoDto struct {
	ID        uint
	Title     string
	Cover     string
	Desc      string
	Copyright bool
}

//只有视频id为参数
type VideoIdDto struct {
	ID uint
}

//管理员导入视频
type ImportVideo struct {
	Title string
	Cover string
	Desc  string
	Type  string
}

type ImportResourceDto struct {
	Vid      uint
	Title    string
	Res360   string
	Res480   string
	Res720   string
	Res1080  string
	Original string
}

//获取视频列表
type GetVideoListDto struct {
	Page      int
	PageSize  int
	Partition int
}
