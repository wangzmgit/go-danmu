package dto

//删除评论回复
type CommentIdDto struct {
	ID uint
}

type CommentDto struct {
	Content string
	Vid     uint
}

type ReplyDto struct {
	Cid       uint
	Content   string
	ReplyUid  uint
	ReplyName string
}
