package middleware

import (
	baseRes "codeSandbox/responses"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var START_TIME = utils.Config.Server.Maintenance.StartTime
var END_TIME = utils.Config.Server.Maintenance.EndTime

// MaintenanceMiddleware
//
//	@Description: 维护期中间件，维护期间不能使用的功能都可以添加该中间件
//	@return gin.HandlerFunc
func MaintenanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		nowTime := time.Now()
		nowTimeStr := nowTime.Format("2006-01-02 15:05:05")
		// 将日期和时间拼接在一起
		startTimeStr := nowTimeStr[:len("2006-01-02 ")] + START_TIME
		endTimeStr := nowTimeStr[:len("2006-01-02 ")] + END_TIME
		// 转为时间变量
		startTime, err := time.Parse("2006-01-02 15:05:05", startTimeStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, baseRes.ErrByCode(global.SYSTEM_ERROR))
			return
		}
		endTime, err := time.Parse("2006-01-02 15:05:05", endTimeStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, baseRes.ErrByCode(global.SYSTEM_ERROR))
			return
		}
		// 如果当前时刻在维护期间，则直接拒绝
		if nowTime.Unix() >= startTime.Unix() && nowTime.Unix() <= endTime.Unix() {
			c.AbortWithStatusJSON(http.StatusOK, baseRes.ErrByCode(global.MAINTAIN_ERROR))
			return
		}
	}
}
