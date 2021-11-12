package dto

//删除评论回复
type CommentDeleteRequest struct {
	ID uint
}

type CommentRequest struct {
	Content string
	Vid     uint
}

type ReplyRequest struct {
	Cid       uint
	Content   string
	ReplyUid  uint
	ReplyName string
}
