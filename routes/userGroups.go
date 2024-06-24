package routes

import (
	api_v1 "codeSandbox/api/v1"
	"codeSandbox/utils/global"
	"codeSandbox/utils/middleware"
	"github.com/gin-gonic/gin"
)

func UserGroup(r *gin.Engine) {
	router := r.Group("api/v1")
	{
		// 注册
		router.POST("user/register", api_v1.Register)
		// 登录
		router.POST("user/login", api_v1.Login)
	}

	needLogin := r.Group("api/v1")
	needLogin.Use(middleware.JwtToken(true, global.NORMAL_USER_ROLE))
	{
		// 密钥对列表
		needLogin.GET("user/keys", api_v1.KeyList)
		// 生成密钥
	}
}
