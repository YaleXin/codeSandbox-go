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
	// 常规方式执行（登录后即可提交代码）
	needLogin := r.Group("api/v1")
	needLogin.Use(middleware.JwtToken(true, global.NORMAL_USER_ROLE))
	{
		needLogin.POST("executeCode", v1.ExecuteCode)
	}
	// 程序方式执行（需要使用 secretKey 加密数据， 并提供 publicKey ）
	program := r.Group("api/v1")
	{
		program.POST("programExecuteCode", v1.ProgramExecuteCode)
	}
}
