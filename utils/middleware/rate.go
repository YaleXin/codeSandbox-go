package middleware

import (
	redisDb "codeSandbox/db"
	"codeSandbox/model/dto"
	baseRes "codeSandbox/responses"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 全局速率设置
func GlobalRateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果ip请求连接数在 per 内超过 events 次，返回 429 并抛出error
		// 这里是 1 秒内不能超过 20 次
		if !redisDb.Allow(c.ClientIP(), 20, 1*time.Second) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, baseRes.ErrByCode(global.TOO_MANY_REQUESTS_ERROR))
			return
		}
		c.Next()
	}
}

// 登录防抖设置
func LoginRateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data dto.UserLoginRequest
		// 由于 ShouldBindJSON 只能调用一次，因此后面我们要把它保存起来
		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(http.StatusOK, baseRes.ErrByCode(global.PARAMS_ERROR))
			return
		}

		// 为了防止被爆破密码，建立一个 ip-username-login 的 redis key 锁
		// 每 5 秒钟每个用户只能尝试登录一次
		if !redisDb.Allow(c.ClientIP()+"-"+data.Username+"-login", 1, 5*time.Second) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, baseRes.ErrByCode(global.TOO_MANY_REQUESTS_ERROR))
			return
		}
		c.Set("userLoginRequest", &data)
		c.Next()
	}

}

// 注册防抖设置
func RegisterRateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 为了防止被过度注册，建立一个 ip-register 的 redis key 锁
		// 每 60 秒钟每个IP只能尝试注册一次
		if !redisDb.Allow(c.ClientIP()+"-register", 1, 60*time.Second) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, baseRes.ErrByCode(global.TOO_MANY_REQUESTS_ERROR))
			return
		}
		c.Next()
	}

}
