package vo

import "time"

type CarouselVo struct {
	Img string `json:"img"`
	Url string `json:"url"`
}

//admin 轮播图
type AdminCarouselVo struct {
	ID        uint      `json:"id"`
	Img       string    `json:"img"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}
