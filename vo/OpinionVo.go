package vo

import "time"

type OpinionVo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Telephone string    `json:"telephone"`
	Gender    int       `json:"gender"`
	Desc      string    `json:"desc"`
	Uid       uint      `json:"uid"`
	CreatedAt time.Time `json:"created_at"`
}
