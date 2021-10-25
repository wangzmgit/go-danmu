package dto

import (
	"time"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/util"
)

type UserDto struct {
	ID       uint      `json:"uid"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Sign     string    `json:"sign"`
	Avatar   string    `json:"avatar"`
	Gender   int       `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		ID:       user.ID,
		Name:     user.Name,
		Email:    util.HideEmail(user.Email),
		Sign:     user.Sign,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Birthday: user.Birthday,
	}
}
