package routes

import (
	api_v1 "codeSandbox/api/v1"
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
}
