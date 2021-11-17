package dto

type UploadVideoRequest struct {
	ID           uint
	Title        string
	Cover        string
	Introduction string
	Original     bool
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
