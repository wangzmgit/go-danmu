package vo

type TotalData struct {
	User    int    `json:"user"`
	Video   int    `json:"video"`
	Review  int    `json:"review"`
	Message int    `json:"message"`
	Version string `json:"version"`
	Redis   bool   `json:"redis"`
	FFmpeg  bool   `json:"ffmpeg"`
}

type OneDayData struct {
	User  int    `json:"user"`
	Video int    `json:"video"`
	Date  string `json:"date"`
}
