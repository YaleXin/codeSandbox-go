package routes

import (
	"codeSandbox/utils"
	"codeSandbox/utils/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Starter() {
	log.Info("init router...")
	if utils.Config.Server.AppMode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	// 一些配置
	r.MaxMultipartMemory = 8 << 20                 // 8 MiB
	r.Use(middleware.Logger(log.StandardLogger())) // 使用Logger记录日志
	r.Use(gin.Recovery())                          // 恐慌恢复
	// TODO 开启限流
	// r.Use(middleware.RateMiddleware())             // 速率限制
	r.Use(middleware.Cors()) // 跨域处理
	// 绑定路由处理函数
	SandboxGroup(r)
	log.Info("init router run~")
	err := r.Run(fmt.Sprintf("%s:%s", utils.Config.Server.Host, utils.Config.Server.Port))
	if err != nil {
		log.Panic(fmt.Sprintf("Server startup failure, %v1", err))
	}
}
