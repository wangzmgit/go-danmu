package dto

type CreateCollectionDto struct {
	Title string
	Cover string
	Desc  string
}

type DeleteCollectionDto struct {
	ID uint
}

type AddVideoDto struct {
	Vid uint //视频ID
	Cid uint //合集ID
}

type DeleteVideoDto struct {
	Vid uint //视频ID
	Cid uint //合集ID
}
