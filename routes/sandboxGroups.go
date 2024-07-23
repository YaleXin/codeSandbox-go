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
		router.GET("languages", v1.LanguageList)

	}

	needLogin := r.Group("api/v1")
	// 维护期间不能提交代码
	needLogin.Use(middleware.MaintenanceMiddleware())
	// 需要登录
	needLogin.Use(middleware.JwtToken(true, global.NORMAL_USER_ROLE))
	{
		// 常规方式执行（登录后即可提交代码）
		needLogin.POST("executeCode", middleware.ExecuteCodeRateMiddleware(), v1.ExecuteCode)
	}

	program := r.Group("api/v1")
	// 维护期间不能提交代码
	program.Use(middleware.MaintenanceMiddleware())
	{
		// 程序方式执行（需要使用 secretKey 加密数据， 并提供 publicKey ）
		program.POST("programExecuteCode", middleware.ProgramExecuteCodeRateMiddleware(), v1.ProgramExecuteCode)
	}
}
