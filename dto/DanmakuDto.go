package dto

type DanmakuDto struct {
	Time  uint   `json:"time"`
	Type  int    `json:"type"`
	Color string `json:"color"`
	Text  string `json:"text"`
}
