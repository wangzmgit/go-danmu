package manage

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
)

type UpdateVersion struct {
	Content string `json:"content"`
}

/*********************************************************
** 函数功能: 获取oss配置信息
** 日    期: 2022年2月24日16:11:27
**********************************************************/
func GetOssConfig(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"storage":         viper.GetBool("aliyunoss.storage"),
		"bucket":          viper.GetString("aliyunoss.bucket"),
		"endpoint":        viper.GetString("aliyunoss.endpoint"),
		"accesskeyId":     viper.GetString("aliyunoss.accesskey_id"),
		"accesskeySecret": viper.GetString("aliyunoss.accesskey_secret"),
		"domain":          viper.GetString("aliyunoss.domain"),
	}, response.OK)
}

/*********************************************************
** 函数功能: 获取邮箱配置信息
** 日    期:
**********************************************************/
func GetEmailConfig(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"name":     viper.GetString("email.name"),
		"host":     viper.GetString("email.host"),
		"port":     viper.GetInt("email.port"),
		"address":  viper.GetString("email.address"),
		"password": viper.GetString("email.password"),
	}, response.OK)
}

/*********************************************************
** 函数功能: 配置oss
** 日    期:
**********************************************************/
func SetOssConfig(ctx *gin.Context) {
	var oss dto.OssConfigDto
	err := ctx.Bind(&oss)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}

	viper.Set("aliyunoss.storage", oss.Storage)
	viper.Set("aliyunoss.bucket", oss.Bucket)
	viper.Set("aliyunoss.endpoint", oss.Endpoint)
	viper.Set("aliyunoss.accesskey_id", oss.AccesskeyId)
	viper.Set("aliyunoss.accesskey_secret", oss.AccesskeySecret)
	viper.Set("aliyunoss.domain", oss.Domain)

	viper.WriteConfig()
	response.Success(ctx, nil, response.OK)
}

/*********************************************************
** 函数功能: 配置邮箱
** 日    期:
**********************************************************/
func SetEmailConfig(ctx *gin.Context) {
	var email dto.EmailConfigDto
	err := ctx.Bind(&email)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}

	viper.Set("email.name", email.Name)
	viper.Set("email.host", email.Host)
	viper.Set("email.port", email.Port)
	viper.Set("email.address", email.Address)
	viper.Set("email.password", email.Password)

	viper.WriteConfig()
	response.Success(ctx, nil, response.OK)
}

/*********************************************************
** 函数功能: 设置管理员账号
** 日    期:
**********************************************************/
func SetAdminConfig(ctx *gin.Context) {
	var admin dto.AdminConfigDto
	err := ctx.Bind(&admin)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}

	viper.Set("admin.email", admin.Email)
	viper.Set("admin.password", admin.Password)

	viper.WriteConfig()
	response.Success(ctx, nil, response.OK)
}

/*********************************************************
** 函数功能: 获取其他配置信息
** 日    期: 2022年2月24日16:11:27
**********************************************************/
func GetOtherConfig(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"coding":     viper.GetString("transcoding.coding"),
		"max_res":    viper.GetInt("transcoding.max_res"),
		"video_user": viper.GetInt("user.video"),
	}, response.OK)
}

/*********************************************************
** 函数功能: 设置其他配置信息
** 日    期: 2022年2月24日16:27:14
**********************************************************/
func SetOtherConfig(ctx *gin.Context) {
	var other dto.OtherConfigDto
	err := ctx.Bind(&other)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}

	viper.Set("transcoding.coding", other.Coding)
	viper.Set("transcoding.max_res", other.MaxRes)
	viper.Set("user.video", other.VideoUser)

	viper.WriteConfig()
	response.Success(ctx, nil, response.OK)
}
