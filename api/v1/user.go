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
// @Tags Register
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
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(global.PARAMS_ERROR)))
		return
	}

	instance := &userServices.UserServiceInstance
	errCode := instance.UserRegister(&data)
	if errCode != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(errCode)))
	} else {
		c.JSON(http.StatusOK, baseRes.OK)
	}
}

// Login 登录
// @Summary 用户登录
// @Description 提交用户名和密码
// @Tags Login
// @Accept json
// @Produce json
// @Param userLoginRequest body dto.UserLoginRequest true "登录表单"
// @Success 200 {object} responses.Response{data=vo.UserVO} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	var data dto.UserLoginRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(global.PARAMS_ERROR)))
		return
	}

	instance := &userServices.UserServiceInstance
	code, userVO := instance.UserLogin(&data)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		return
	}

	c.JSON(http.StatusOK, baseRes.OK.WithData(userVO))
}

// NewKey 生成密钥对
// @Summary 生成密钥对
// @Tags NewKey
// @Description 生成密钥对，用于通过程序式提交代码
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Success 200 {object} responses.Response{data=vo.KeyPairVO} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/newKey [post]
func NewKey(c *gin.Context) {
	instance := &userServices.UserServiceInstance
	code, vo := instance.GenerateKeyPair(c)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(vo))
		return
	}
}

// DeleteKey 删除指定 id 的密钥对
// @Summary 删除指定 id 的密钥对
// @Tags DeleteKey
// @Description 删除指定 id 的密钥对
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param deleteKeyRequest body dto.DeleteKeyRequest true "删除表单"
// @Success 200 {object} responses.Response{data=int} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/delKey [delete]
func DeleteKey(c *gin.Context) {
	var data dto.DeleteKeyRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(global.PARAMS_ERROR)))
		return
	}

	instance := &userServices.UserServiceInstance
	code, rowsAffected := instance.DeleteKeyPair(c, &data)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(rowsAffected))
		return
	}
}

// KeyList 展示用户的密钥对
// @Summary 展示用户的密钥对
// @Description 展示用户的密钥对
// @Tags KeyList
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Success 200 {object} responses.Response{data=[]vo.KeyPairVO} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/keys [get]
func KeyList(c *gin.Context) {
	instance := &userServices.UserServiceInstance
	code, pairs := instance.GetUserKeys(c)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(pairs))
		return
	}
}

// UserInfo 展示用户的信息
// @Summary 展示用户的信息
// @Description 展示用户的信息
// @Tags UserInfo
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Success 200 {object} responses.Response{data=[]vo.UserDetailVO} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/info [get]
func UserInfo(c *gin.Context) {
	instance := &userServices.UserServiceInstance
	code, info := instance.GetUserInfo(c)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(info))
		return
	}
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改密码
// @Tags ChangePassword
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param changePasswordRequest body dto.ChangePasswordRequest true "修改密码表单"
// @Success 200 {object} responses.Response{data=int} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/changePwd [put]
func ChangePassword(c *gin.Context) {
	var data dto.ChangePasswordRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(global.PARAMS_ERROR)))
		return
	}

	instance := &userServices.UserServiceInstance
	code := instance.ChangePassword(c, &data)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(1))
		return
	}
}

// PageExecution 获取用户的代码执行记录
// @Summary 获取用户的代码执行记录
// @Description 获取用户的代码执行记录
// @Tags PageExecution
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param pageExecutionRequest body dto.PageExecutionRequest true "分页信息"
// @Success 200 {object} responses.Response{data=[]vo.PageDataVO{data=[]vo.ExecutionVO}} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/user/execution [post]
func PageExecution(c *gin.Context) {

	var data dto.PageExecutionRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(global.PARAMS_ERROR)))
		return
	}

	instance := &userServices.UserServiceInstance
	code, info := instance.GetPageExecution(c, &data)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		return
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(info))
		return
	}
}
