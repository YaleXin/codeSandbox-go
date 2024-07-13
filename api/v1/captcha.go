package api_v1

import (
	baseRes "codeSandbox/responses"
	"codeSandbox/service/captchaServices"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Captcha 生成验证码
// @Summary 生成验证码
// @Tags Captcha
// @Description 生成验证码，需要用验证码的地方有登录和注册
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response{data=vo.CaptchaVO} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/captcha [get]
func Captcha(c *gin.Context) {
	instance := &captchaServices.CaptchaServiceInstance
	code, vo := instance.GenerateCaptchaToken()
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.ErrByCode(code))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(vo))
		return
	}
}
