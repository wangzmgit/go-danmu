package dto

type FollowDto struct {
	ID     uint   `json:"uid"`
	Name   string `json:"name"`
	Sign   string `json:"sign"`
	Avatar string `json:"avatar"`
}
