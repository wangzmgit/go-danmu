package common

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// token加密秘钥
var adminJwtKey = []byte(viper.GetString("server.admin_jwt_secret"))

type AdminClaims struct {
	AdminID uint
	jwt.StandardClaims
}

/*********************************************************
** 函数功能: 发放管理员Token
** 日    期:2021/8/1
**********************************************************/
func ReleaseAdminToken(AdminID uint) (string, error) {
	//token过期时间
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &AdminClaims{
		AdminID: AdminID,
		StandardClaims: jwt.StandardClaims{
			//发放时间等
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "danmu3.0",
			Subject:   "admin_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(adminJwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

/*********************************************************
** 函数功能: 解析管理员Token
** 日    期:2021/8/1
**********************************************************/
func ParseAdminToken(tokenString string) (*jwt.Token, *AdminClaims, error) {
	claims := &AdminClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, e error) {
		return adminJwtKey, nil
	})
	return token, claims, err
}
