package middleware

import (
	redisDb "codeSandbox/db"
	baseRes "codeSandbox/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果ip请求连接数在 2 秒内超过 5 次，返回 429 并抛出error
		if !redisDb.Allow(c.ClientIP(), 5, 2*time.Second) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, baseRes.ErrTooManyRequests)
			return
		}
		c.Next()
	}
}
