package util

import (
	"fmt"
	"strconv"
)

//储存Redis的key

const (
	ClicksVideoList = "video:clicks:list"
)

func VideoClicksKey(id int) string {
	return fmt.Sprintf("video:clicks:%s", strconv.Itoa(id))
}

func VideoCollectKey(id int) string {
	return fmt.Sprintf("video:collect:%s", strconv.Itoa(id))
}

func VideoLikeKey(id int) string {
	return fmt.Sprintf("video:like:%s", strconv.Itoa(id))
}

func CodeKey(email string) string {
	return fmt.Sprintf("code:%s", email)
}

func LoginCodeKey(email string) string {
	return fmt.Sprintf("login:%s", email)
}

func VideoClicksLimitKey(vid int, ip string) string {
	return fmt.Sprintf("clicks:%d:%s", vid, ip)
}
