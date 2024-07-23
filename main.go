package main

import (
	"codeSandbox/db"
	"codeSandbox/routes"
	"codeSandbox/task"
	logPackage "codeSandbox/utils/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	db.InitRedis()
	logPackage.ConfigLog()
	taskCron := task.InitTask()
	routes.Starter()
	select {
	case <-signalChan:
		// 接收到信号，优雅地关闭 cron
		if taskCron != nil {
			taskCron.Stop()
		}
	}
}
