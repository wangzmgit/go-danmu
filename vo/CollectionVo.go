package vo

import (
	"time"
	"wzm/danmu3.0/model"
)

type CollectionVo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCollectionVo(collection model.Collection) CollectionVo {
	return CollectionVo{
		ID:        collection.ID,
		Title:     collection.Title,
		Cover:     collection.Cover,
		Desc:      collection.Desc,
		CreatedAt: collection.CreatedAt,
	}
}
