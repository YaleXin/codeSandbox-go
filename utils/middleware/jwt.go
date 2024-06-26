package middleware

import (
	"codeSandbox/utils"
	"codeSandbox/utils/errmsg"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

var JWT_KEY = []byte(utils.Config.Server.JwtKey)
var JWT_EXPIRE_TIME = time.Duration(utils.Config.Server.JwtExpireTime) * time.Minute
var code int

// jwt 中要加密的内容
type MyClaims struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

// SetToken 生成token
func SetToken(userId uint, username string, role int) (string, int) {
	expireTime := time.Now().Add(JWT_EXPIRE_TIME)
	SetClaims := MyClaims{
		userId,
		username,
		role,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "codeSandbox",
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err := reqClaim.SignedString(JWT_KEY)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, global.SUCCESS
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, int) {
	setToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
	if err != nil {
		return nil, errmsg.ERROR
	}
	if key, ok := setToken.Claims.(*MyClaims); ok && setToken.Valid {
		return key, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR
	}
}

// JwtToken jwt中间件
// 参数： termination 是否中断(true:没权限直接静止访问,false:没权限只返回部分字段)

func JwtToken(termination bool, needRole int) gin.HandlerFunc {
	cRes := func(c *gin.Context, code int) {
		if termination {
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
		}
	}
	return func(c *gin.Context) {
		ckToken := c.GetHeader("Token")
		if ckToken == "" {
			//认证字符串判断 !没有认证字符串
			code = errmsg.ERROR_TOKEN_EXIST
			cRes(c, code)
			return
		}
		keyData, tCode := CheckToken(ckToken)
		if tCode == errmsg.ERROR {
			//认证字符串判断 !内容不对
			code = errmsg.ERROR_TOKEN_WRONG
			cRes(c, code)
			return
		}
		if keyData.Role > needRole {
			//认证字符串 权限不够
			code = errmsg.ERROR_ROLE_LOW
			cRes(c, code)
			return
		}
		if time.Now().Unix() > keyData.ExpiresAt {
			//认证字符串时间判断 !过期
			code = errmsg.ERROR_TOKEN_RUNTIME
			cRes(c, code)
			return
		}
		// 将 user 保存起来，后续用得到
		c.Set("user", keyData)
		c.Next()
	}
}
