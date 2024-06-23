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
	if data.Username == "" || data.Password == "" {
		c.JSON(http.StatusOK, baseRes.Err.WithData(&global.CustomError{
			ErrorCode: global.PARAMS_ERROR,
			Message:   global.GetErrMsg(global.PARAMS_ERROR),
		}))
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
	}
}
