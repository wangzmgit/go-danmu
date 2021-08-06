package dto

import (
	"time"
	"wzm/danmu3.0/model"
)

type UserDto struct {
	ID      uint      `json:"uid"`
	Name     string    `json:"name"`
	Sign     string    `json:"sign"`
	Avatar   string    `json:"avatar"`
	Gender   int       `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

type AdminUserDto struct {
	ID      uint      `json:"uid"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Sign     string    `json:"sign"`
	Avatar   string    `json:"avatar"`
	Gender   int       `json:"gender"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		ID:      user.ID,
		Name:     user.Name,
		Sign:     user.Sign,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Birthday: user.Birthday,
	}
}
