package middleware

import (
	redisDb "codeSandbox/db"
	"codeSandbox/model/dto"
	baseRes "codeSandbox/responses"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var GLOBAL_TIME = time.Duration(utils.Config.Server.Rate.GlobalTime) * time.Second
var GLOBAL_COUNT = utils.Config.Server.Rate.GlobalCount
var LOGIN_TIME = time.Duration(utils.Config.Server.Rate.LoginTime) * time.Second
var LOGIN_COUNT = utils.Config.Server.Rate.LoginCount
var REGISTER_TIME = time.Duration(utils.Config.Server.Rate.RegisterTime) * time.Second
var REGISTER_COUNT = utils.Config.Server.Rate.RegisterCount
var EXECUTECODE_TIME = time.Duration(utils.Config.Server.Rate.ExecuteCodeTime) * time.Second
var EXECUTECODE_COUNT = utils.Config.Server.Rate.ExecuteCodeCount
var PROGRAM_EXECUTECODE_TIME = time.Duration(utils.Config.Server.Rate.ProgramExecuteCodeTime) * time.Second
var PROGRAM_EXECUTECODE_COUNT = utils.Config.Server.Rate.ProgramExecuteCodeCount

// 全局速率设置
func GlobalRateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果ip请求连接数在 GLOBAL_TIME 内超过 GLOBAL_COUNT 次，返回 429 并抛出error
		if !redisDb.Allow(c.ClientIP()+"-global", GLOBAL_COUNT, GLOBAL_TIME) {
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
			c.AbortWithStatusJSON(http.StatusOK, baseRes.ErrByCode(global.PARAMS_ERROR))
			return
		}

		// 为了防止被爆破密码，建立一个 ip-username-login 的 redis key 锁
		if !redisDb.Allow(c.ClientIP()+"-"+data.Username+"-login", LOGIN_COUNT, LOGIN_TIME) {
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
		if !redisDb.Allow(c.ClientIP()+"-register", REGISTER_COUNT, REGISTER_TIME) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, baseRes.ErrByCode(global.TOO_MANY_REQUESTS_ERROR))
			return
		}
		c.Next()
	}

}

// 执行代码防抖设置
func ExecuteCodeRateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !redisDb.Allow(c.ClientIP()+"-executeCode", EXECUTECODE_COUNT, EXECUTECODE_TIME) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, baseRes.ErrByCode(global.TOO_MANY_REQUESTS_ERROR))
			return
		}
		c.Next()
	}
}

// 执行代码防抖设置
func ProgramExecuteCodeRateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !redisDb.Allow(c.ClientIP()+"-programExecuteCode", PROGRAM_EXECUTECODE_COUNT, PROGRAM_EXECUTECODE_TIME) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, baseRes.ErrByCode(global.TOO_MANY_REQUESTS_ERROR))
			return
		}
		c.Next()
	}
}
