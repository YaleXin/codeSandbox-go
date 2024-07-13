package middleware

import (
	baseRes "codeSandbox/responses"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"codeSandbox/utils/tool"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

var CAPTCHA_EXPIRE_TIME = time.Duration(utils.Config.Server.CaptchaExpireTime) * time.Minute
var CAPTCHA_KEY = utils.Config.Server.CaptchaKey

// jwt 中要加密的内容
type CaptchaClaims struct {
	CaptchaContentMd5 string `json:"captchaContentMd5"`
	jwt.StandardClaims
}

// SetToken 生成token
func SetCaptchaToken(CaptchaContentMd5 string) (string, int) {
	// 验证码过期时间即为 token 过期时间
	expireTime := time.Now().Add(CAPTCHA_EXPIRE_TIME)
	setClaims := CaptchaClaims{
		CaptchaContentMd5,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "codeSandbox",
		},
	}
	// 使用指定算法生成
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaims)
	// 使用密钥签名
	token, err := reqClaim.SignedString(JWT_KEY)
	if err != nil {
		return "", global.SYSTEM_ERROR
	}
	return token, global.SUCCESS
}

// CheckToken 验证token
func CheckCaptchaToken(token string) (*CaptchaClaims, int) {
	captchaToken, err := jwt.ParseWithClaims(token, &CaptchaClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
	if err != nil {
		return nil, global.SYSTEM_ERROR
	}
	if key, ok := captchaToken.Claims.(*CaptchaClaims); ok && captchaToken.Valid {
		return key, global.SUCCESS
	} else {
		return nil, global.SYSTEM_ERROR
	}
}

// JwtCaptchaToken jwt中间件
func JwtCaptchaToken() gin.HandlerFunc {
	cRes := func(c *gin.Context, code int) {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		c.Abort()
	}
	return func(c *gin.Context) {
		ckToken := c.GetHeader("CPT-Token")
		captcha := c.Query("captcha")
		if ckToken == "" {
			//认证字符串判断 !没有认证字符串
			errCode := global.CPT_TOKEN_NOT_FOUND_ERROR
			cRes(c, errCode)
			return
		}
		if captcha == "" {
			// 没有提交 captcha
			errCode := global.CAPTCHA_NOT_FOUND_ERROR
			cRes(c, errCode)
			return
		}
		//认证字符串判断 !（token是否是有效）
		keyData, tCode := CheckCaptchaToken(ckToken)
		if tCode == global.SYSTEM_ERROR {
			errCode := global.TOKEN_WRONG_ERROR
			cRes(c, errCode)
			return
		}
		if time.Now().Unix() > keyData.ExpiresAt {
			//认证字符串时间判断 !过期
			errCode := global.TOKEN_RUNTIME_ERROR
			cRes(c, errCode)
			return
		}

		// 判断提交的验证码对不对 计算  md5(key + ToLower(captcha))
		md5Str := tool.MD5Str(CAPTCHA_KEY + strings.ToLower(captcha))
		if md5Str != keyData.CaptchaContentMd5 {
			errCode := global.CAPTCHA_WRONG_ERROR
			cRes(c, errCode)
			return
		}
	}
}
