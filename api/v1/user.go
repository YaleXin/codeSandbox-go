package api_v1

import (
	"codeSandbox/model/dto"
	baseRes "codeSandbox/responses"
	"codeSandbox/service/userServices"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 注册
// @Summary 用户注册
// @Description 提交用户名，邮箱和密码
// @Accept json
// @Produce json
// @Param userRegisterRequest body dto.UserRegisterRequest true "用户信息"
// @Success 200 {object} responses.Response "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/register [post]
func Register(c *gin.Context) {
	var data dto.UserRegisterRequest
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

// Login 登录
// @Summary 用户登录
// @Description 提交用户名和密码
// @Accept json
// @Produce json
// @Param userLoginRequest body dto.UserLoginRequest true "登录表单"
// @Success 200 {object} responses.Response "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	var data dto.UserLoginRequest
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

// NewKey 生成密钥对
// @Summary 生成密钥对
// @Description 生成密钥对，用于通过程序式提交代码
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/newKey [post]
func NewKey(c *gin.Context) {
	instance := &userServices.UserServiceInstance
	code, vo := instance.GenerateKeyPair(c)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithData(&global.CustomError{
			ErrorCode: code,
			Message:   global.GetErrMsg(code),
		}))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(vo))
		return
	}
}

// KeyList 展示用户的密钥对
// @Summary 展示用户的密钥对
// @Description 展示用户的密钥对
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/keys [get]
func KeyList(c *gin.Context) {
	instance := &userServices.UserServiceInstance
	code, pairs := instance.GetUserKeys(c)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithData(&global.CustomError{
			ErrorCode: code,
			Message:   global.GetErrMsg(code),
		}))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(pairs))
		return
	}
}
