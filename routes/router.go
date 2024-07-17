package routes

import (
	"codeSandbox/docs"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
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
	if utils.Config.Server.AppMode == global.APP_MODE_PROD {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	// 生产环境中必须要关闭 API 文档
	if utils.Config.Server.AppMode != global.APP_MODE_PROD {
		// 注册 swagger
		registerSwagger(r)
	}

	// 一些配置
	r.MaxMultipartMemory = 8 << 20                 // 8 MiB
	r.Use(middleware.Logger(log.StandardLogger())) // 使用Logger记录日志
	r.Use(gin.Recovery())                          // 恐慌恢复

	r.Use(middleware.GlobalRateMiddleware()) // 速率限制
	r.Use(middleware.Cors())                 // 跨域处理
	// 管理员路由注册
	AdminGroup(r)
	// 沙箱路由注册
	SandboxGroup(r)
	// 用户相关路由注册
	UserGroup(r)
	// 验证码路由注册
	CaptchaGroup(r)

	log.Info("init router run~")
	err := r.Run(fmt.Sprintf("%s:%s", utils.Config.Server.Host, utils.Config.Server.Port))
	if err != nil {
		log.Panic(fmt.Sprintf("Server startup failure, %v1", err))
	}
}
