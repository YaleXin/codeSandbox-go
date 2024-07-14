package middleware

import (
	baseRes "codeSandbox/responses"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

var JWT_KEY = []byte(utils.Config.Server.JwtKey)
var JWT_EXPIRE_TIME = time.Duration(utils.Config.Server.JwtExpireTime) * time.Minute

// jwt 中要加密的内容
type MyClaims struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	Audit    bool   `json:"audit"`
	jwt.StandardClaims
}

// SetToken 生成token
func SetToken(userId uint, username string, role int, audit bool) (string, int) {
	expireTime := time.Now().Add(JWT_EXPIRE_TIME)
	SetClaims := MyClaims{
		userId,
		username,
		role,
		audit,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "codeSandbox",
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err := reqClaim.SignedString(JWT_KEY)
	if err != nil {
		return "", global.SYSTEM_ERROR
	}
	return token, global.SUCCESS
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, int) {
	setToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
	if err != nil {
		return nil, global.TOKEN_WRONG_ERROR
	}
	if key, ok := setToken.Claims.(*MyClaims); ok && setToken.Valid {
		// 如果 id 为 0 ，则是错的，因为我们返回的必然是一个正确的
		if key.UserId == 0 {
			return nil, global.TOKEN_WRONG_ERROR
		}
		return key, global.SUCCESS
	} else {
		return nil, global.TOKEN_WRONG_ERROR
	}
}

// JwtToken jwt中间件
// 参数： termination 是否中断(true:没权限直接静止访问,false:没权限只返回部分字段)

func JwtToken(termination bool, needRole int) gin.HandlerFunc {
	cRes := func(c *gin.Context, code int) {
		if termination {
			c.JSON(http.StatusOK, baseRes.ErrByCode(code))
			c.Abort()
		}
	}
	return func(c *gin.Context) {
		ckToken := c.GetHeader("Token")
		if ckToken == "" {
			//认证字符串判断 !没有认证字符串
			code := global.NOT_LOGIN_ERROR
			cRes(c, code)
			return
		}
		//认证字符串判断 !（token是否正确）
		keyData, tCode := CheckToken(ckToken)
		if tCode != global.SUCCESS {
			cRes(c, tCode)
			return
		}
		if keyData.Role > needRole {
			//认证字符串 权限不够
			code := global.LACK_AUTH_ERROR
			cRes(c, code)
			return
		}
		if time.Now().Unix() > keyData.ExpiresAt {
			//认证字符串时间判断 !过期
			code := global.TOKEN_RUNTIME_ERROR
			cRes(c, code)
			return
		}
		// 将 user 保存起来，后续用得到
		c.Set("user", keyData)
		c.Next()
	}
}
