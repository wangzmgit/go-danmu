package vo

import (
	"time"

	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/util"
)

type UserVo struct {
	ID       uint      `json:"uid"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Sign     string    `json:"sign"`
	Avatar   string    `json:"avatar"`
	Gender   int       `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

type AdminUserVo struct {
	ID        uint      `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Sign      string    `json:"sign"`
	Avatar    string    `json:"avatar"`
	Gender    int       `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserVo(user model.User) UserVo {
	return UserVo{
		ID:       user.ID,
		Name:     user.Name,
		Email:    util.HideEmail(user.Email),
		Sign:     user.Sign,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Birthday: user.Birthday,
	}
}

func ToAuthorVo(user model.User) UserVo {
	return UserVo{
		ID:     user.ID,
		Name:   user.Name,
		Sign:   user.Sign,
		Avatar: user.Avatar,
	}
}
