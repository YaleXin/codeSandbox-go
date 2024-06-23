package routes

import (
	"codeSandbox/docs"
	"codeSandbox/utils"
	"codeSandbox/utils/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 初始化 docker 相关信息
	_ "codeSandbox/service/sandboxDockerServices"
)

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

func Starter() {
	log.Info("init router...")
	if utils.Config.Server.AppMode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	// 注册 swagger
	registerSwagger(r)

	// 一些配置
	r.MaxMultipartMemory = 8 << 20                 // 8 MiB
	r.Use(middleware.Logger(log.StandardLogger())) // 使用Logger记录日志
	r.Use(gin.Recovery())                          // 恐慌恢复
	// TODO 开启限流
	// r.Use(middleware.RateMiddleware())             // 速率限制
	r.Use(middleware.Cors()) // 跨域处理
	// 绑定沙箱路由处理函数
	SandboxGroup(r)
	// 用户相关
	UserGroup(r)
	log.Info("init router run~")
	err := r.Run(fmt.Sprintf("%s:%s", utils.Config.Server.Host, utils.Config.Server.Port))
	if err != nil {
		log.Panic(fmt.Sprintf("Server startup failure, %v1", err))
	}
}
