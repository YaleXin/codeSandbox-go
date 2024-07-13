package routes

import (
	api_v1 "codeSandbox/api/v1"
	"github.com/gin-gonic/gin"
)

func CaptchaGroup(r *gin.Engine) {
	captcha := r.Group("api/v1")
	captcha.GET("/captcha", api_v1.Captcha)
}
