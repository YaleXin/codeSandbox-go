package routes

import (
	v1 "codeSandbox/api/v1"
	"codeSandbox/utils/global"
	"codeSandbox/utils/middleware"
	"github.com/gin-gonic/gin"
)

func SandboxGroup(r *gin.Engine) {

	router := r.Group("api/v1/")
	{
		router.GET("languages", v1.List)

	}
	needLogin := r.Group("api/v1")
	needLogin.Use(middleware.JwtToken(true, global.NORMAL_USER_ROLE))
	{
		needLogin.POST("executeCode", v1.ExecuteCode)
	}
}
