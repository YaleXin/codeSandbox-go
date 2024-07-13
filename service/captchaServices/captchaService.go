package captchaServices

import (
	"codeSandbox/model/vo"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"codeSandbox/utils/middleware"
	"codeSandbox/utils/tool"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"strings"
)

var CAPTCHA_KEY = utils.Config.Server.CaptchaKey

type CaptchaService struct {
}

var CaptchaServiceInstance CaptchaService

// 返回状态码，验证码，图片base64
func generateCaptcha() (int, string, string) {
	captcha := tool.GenerateRandomEasyVisibleString(global.CAPTCHA_LEN)
	// 得到包含验证码的图片字节切片
	text := ImgText(30*global.CAPTCHA_LEN, 30, captcha)
	// 图片转base64
	imageBase64 := base64.StdEncoding.EncodeToString(text)
	return global.SUCCESS, captcha, imageBase64
}

// 获取验证码
func (c *CaptchaService) GenerateCaptchaToken() (int, *vo.CaptchaVO) {
	status, captcha, imageBase64 := generateCaptcha()
	if status != global.SUCCESS {
		return status, nil
	}
	// 不应直接在token中放明文，而是放 md5(key + ToLower(captcha))
	md5Str := tool.MD5Str(CAPTCHA_KEY + strings.ToLower(captcha))
	log.Debugf("captcha:%v, md5:%v", captcha, md5Str)
	token, tCode := middleware.SetCaptchaToken(md5Str)
	if tCode != global.SUCCESS {
		return tCode, nil
	}
	captchaVO := vo.CaptchaVO{
		ImageBase64: imageBase64,
		Token:       token,
	}
	return global.SUCCESS, &captchaVO
}

//
//// 验证验证码
//func (cp *CaptchaService) VerifyCaptchaToken(c *gin.Context, captcha string) int {
//	get, exists := c.Get("captcha")
//	if !exists {
//		return global.SYSTEM_ERROR
//	}
//
//	// 使用类型断言转换为 MyClaims
//	captchaClaims, ok := get.(*middleware.CaptchaClaims)
//	if !ok {
//		return global.TOKEN_WRONG_ERROR
//	}
//	// 判断验证码是否正确 计算 md5(key + captcha)
//	md5Str := tool.MD5Str(CAPTCHA_KEY + captcha)
//	if md5Str != captchaClaims.CaptchaContentMd5 {
//		return global.CAPTCHA_WRONG_ERROR
//	}
//	return global.SUCCESS
//}
