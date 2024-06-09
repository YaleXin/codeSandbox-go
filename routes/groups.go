package routes

import (
	v1 "codeSandbox/api/v1"
	"github.com/gin-gonic/gin"
)

func SandboxGroup(r *gin.Engine) {
	router := r.Group("api/v1/")
	{
		router.GET("languages", v1.List)
		router.POST("executeCode", v1.ExecuteCode)
	}
}
