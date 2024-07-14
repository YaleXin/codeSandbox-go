package routes

import (
	api_v1 "codeSandbox/api/v1"
	"codeSandbox/utils/global"
	"codeSandbox/utils/middleware"
	"github.com/gin-gonic/gin"
)

func AdminGroup(r *gin.Engine) {
	admin := r.Group("api/v1/admin")
	// 必须登录且要有管理员权限
	admin.Use(middleware.JwtToken(true, global.ADMIN_USER_ROLE))

	{
		admin.POST("/execution", api_v1.AdminPageExecution)
	}
}
