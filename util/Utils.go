package util

import (
	"math/rand"
	"regexp"
	"time"
)

/*********************************************************
** 函数功能: 随机字符生成
** 参    数：随机数位数
** 日    期:2021/7/10
** 修 改 人:
** 日    期:
** 描    述:
**********************************************************/
func RandomString(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

/*********************************************************
** 函数功能: 邮箱格式匹配
** 参    数：邮箱字符串
** 日    期:2021/7/10
** 修 改 人:
** 日    期:
** 描    述:
**********************************************************/
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

/*********************************************************
** 函数功能: 手机号格式匹配
** 参    数：手机号字符串
** 日    期:2021/7/10
** 修 改 人:
** 日    期:
** 描    述:
**********************************************************/
func VerifyTelephoneFormat(telephone string) bool {
	pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(telephone)
}

/*********************************************************
** 函数功能: 随机数字
** 日    期:2021/7/23
**********************************************************/
func RandomCode(n int) string {
	var letters = []byte("1234567890")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

/*********************************************************
** 函数功能: 是否存在SQL注入
** 日    期: 2022/1/4
**********************************************************/
func ExistSQLInject(sql string) bool {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		return false
	}
	return re.MatchString(sql)
}

/*********************************************************
** 函数功能: 隐藏邮箱
** 日    期:2021/10/21
**********************************************************/
func HideEmail(email string) string {
	pattern := `(\w?)(\w+)(\w)(@\w+\.[a-z]+(\.[a-z]+)?)`
	reg := regexp.MustCompile(pattern)
	return reg.ReplaceAllString(email, "$1****$3$4")
}
