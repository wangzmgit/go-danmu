package dto

type UploadVideoRequest struct {
	ID           uint
	Title        string
	Cover        string
	Introduction string
	Original     bool
	Partition    uint //分区ID，3.6.8新增
}

type VideoModifyRequest struct {
	ID           uint
	Title        string
	Cover        string
	Introduction string
	Original     bool
}

type DeleteVideoRequest struct {
	ID uint
}

type UpdateVideoReviewRequest struct {
	ID     uint
	Status int
}

//作者UID
type AuthorUid struct {
	UID uint
}

//管理员导入视频
type ImportVideo struct {
	Title        string
	Cover        string
	Introduction string
	Video        string
}

type GetVideoListDto struct {
	Page      int
	PageSize  int
	Partition int
}
