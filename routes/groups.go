package routes

import (
	v1 "codeSandbox/api/v1"
	docs "codeSandbox/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SandboxGroup(r *gin.Engine) {
	registerSwagger(r)
	router := r.Group("api/v1/")
	{
		router.GET("languages", v1.List)
		router.POST("executeCode", v1.ExecuteCode)
	}
}
func registerSwagger(r gin.IRouter) {
	// API文档访问地址: http://host/swagger/index.html
	// 注解定义可参考 https://github.com/swaggo/swag#declarative-comments-format
	// 样例 https://github.com/swaggo/swag/blob/master/example/basic/api/api.go
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Title = "代码沙箱"
	docs.SwaggerInfo.Description = "代码沙箱，亦即代码远程执行器"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
