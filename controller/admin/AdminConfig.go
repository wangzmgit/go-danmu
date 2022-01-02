package admin_controller

import (
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetOssConfig(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"bucket":          viper.Get("aliyunoss.bucket"),
		"endpoint":        viper.Get("aliyunoss.endpoint"),
		"accesskeyId":     viper.Get("aliyunoss.accesskey_id"),
		"accesskeySecret": viper.Get("aliyunoss.accesskey_secret"),
		"domain":          viper.Get("aliyunoss.domain"),
	}, "ok")
}

func GetEmailConfig(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"name":     viper.Get("email.name"),
		"host":     viper.Get("email.host"),
		"port":     viper.Get("email.port"),
		"address":  viper.Get("email.address"),
		"password": viper.Get("email.password"),
	}, "ok")
}

func SetOssConfig(ctx *gin.Context) {
	var oss dto.OssConfigDto
	err := ctx.Bind(&oss)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	viper.Set("aliyunoss.bucket", oss.Bucket)
	viper.Set("aliyunoss.endpoint", oss.Endpoint)
	viper.Set("aliyunoss.accesskey_id", oss.AccesskeyId)
	viper.Set("aliyunoss.accesskey_secret", oss.AccesskeySecret)
	viper.Set("aliyunoss.domain", oss.Domain)

	viper.WriteConfig()
	response.Success(ctx, nil, "ok")
}

func SetEmailConfig(ctx *gin.Context) {
	var email dto.EmailConfigDto
	err := ctx.Bind(&email)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	viper.Set("email.name", email.Name)
	viper.Set("email.host", email.Host)
	viper.Set("email.port", email.Port)
	viper.Set("email.address", email.Address)
	viper.Set("email.password", email.Password)

	viper.WriteConfig()
	response.Success(ctx, nil, "ok")
}

func SetAdminConfig(ctx *gin.Context) {
	var admin dto.AdminConfigDto
	err := ctx.Bind(&admin)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	viper.Set("admin.email", admin.Email)
	viper.Set("admin.password", admin.Password)

	viper.WriteConfig()
	response.Success(ctx, nil, "ok")
}
