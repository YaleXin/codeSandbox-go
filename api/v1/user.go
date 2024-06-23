package api_v1

import (
	"codeSandbox/model"
	baseRes "codeSandbox/responses"
	"codeSandbox/service/userServices"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 注册
func Register(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithData("error"))
		return
	}

	instance := &userServices.UserServiceInstance
	errCode := instance.UserRegister(&data)
	if errCode != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithData(&global.CustomError{
			ErrorCode: errCode,
			Message:   global.GetErrMsg(errCode),
		}))
	} else {
		c.JSON(http.StatusOK, baseRes.OK)
	}
}

// Login 登陆
func Login(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithData("error"))
		return
	}

	instance := &userServices.UserServiceInstance
	code, userVO := instance.UserLogin(&data)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithData(&global.CustomError{
			ErrorCode: code,
			Message:   global.GetErrMsg(code),
		}))
		return
	}

	c.JSON(http.StatusOK, baseRes.OK.WithData(userVO))
}
