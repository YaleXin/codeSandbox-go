package task

import (
	"codeSandbox/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var TASK_TIME = utils.Config.Server.Maintenance.ExecuteTime
var START_TIME = utils.Config.Server.Maintenance.StartTime
var END_TIME = utils.Config.Server.Maintenance.EndTime

// 检查配置文件中的维护期是否是有效的
func checkTime() bool {
	nowTime := time.Now()
	nowTimeStr := nowTime.Format("2006-01-02 15:05:05")
	// 将日期和时间拼接在一起
	startTimeStr := nowTimeStr[:len("2006-01-02 ")] + START_TIME
	endTimeStr := nowTimeStr[:len("2006-01-02 ")] + END_TIME
	taskTimeStr := nowTimeStr[:len("2006-01-02 ")] + TASK_TIME
	// 转为时间变量
	startTime, err := time.Parse("2006-01-02 15:05:05", startTimeStr)
	if err != nil {
		return false
	}
	endTime, err := time.Parse("2006-01-02 15:05:05", endTimeStr)
	if err != nil {
		return false
	}
	taskTime, err := time.Parse("2006-01-02 15:05:05", taskTimeStr)
	if err != nil {
		return false
	}
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()
	taskUnix := taskTime.Unix()
	if endUnix <= startUnix {
		return false
	}
	// 维护任务必须在维护期期间执行
	if taskUnix > endUnix || taskUnix < startUnix {
		return false
	}
	return true
}

func InitTask() *cron.Cron {
	log.Infof("task init....")
	if checkTime() {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		// 创建一个新的 Cron 实例
		c := cron.New(cron.WithSeconds(), cron.WithLocation(loc))
		// TASK_TIME 形如 01:02:03
		split := strings.Split(TASK_TIME, ":")
		// 添加一个作业 ：秒 分 时
		_, err := c.AddFunc(fmt.Sprintf("%s %s %s * * *", split[2], split[1], split[0]), ResetCodeSandbox)
		if err != nil {
			log.Errorf("Create task ResetCodeSandbox fail: %v", err)
			return nil
		}
		// 开始调度
		c.Start()
		log.Infof("task init finish..")
		return c
	} else {
		log.Errorf("Task init fail")
		return nil
	}

}
