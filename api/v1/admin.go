package api_v1

import (
	"codeSandbox/model/dto"
	baseRes "codeSandbox/responses"
	"codeSandbox/service/adminServices"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminPageExecution 获取用户的代码执行记录
// @Summary 获取用户的代码执行记录
// @Description 获取用户的代码执行记录
// @Tags Admin
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param pageExecutionRequest body dto.PageExecutionRequest true "分页信息"
// @Success 200 {object} responses.Response{data=vo.PageDataVO{data=[]vo.AdminExecutionVO}} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/admin/execution [post]
func AdminPageExecution(c *gin.Context) {
	var data dto.PageExecutionRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.ErrByCode(global.PARAMS_ERROR))
		return
	}

	instance := &adminServices.AdminServiceInstance
	code, vo := instance.PageExecution(&data)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.ErrByCode(code))
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(vo))
	}

}

// AdminPageUser 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags Admin
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param pageExecutionRequest body dto.PageExecutionRequest true "分页信息"
// @Success 200 {object} responses.Response{data=vo.PageDataVO{data=[]vo.UserDetailVO}} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/admin/user [post]
func AdminPageUser(c *gin.Context) {
	var data dto.PageExecutionRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.ErrByCode(global.PARAMS_ERROR))
		return
	}

	instance := &adminServices.AdminServiceInstance
	code, vo := instance.PageUser(&data)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.ErrByCode(code))
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(vo))
	}
}

// AdminAuditUser 审核通过用户
// @Summary 审核通过用户
// @Description 审核通过用户
// @Tags Admin
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param pageExecutionRequest body dto.UserIdRequest true "用户id表单"
// @Success 200 {object} responses.Response{data=bool} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/admin/auditUser [post]
func AdminAuditUser(c *gin.Context) {
	var data dto.UserIdRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.ErrByCode(global.PARAMS_ERROR))
		return
	}

	instance := &adminServices.AdminServiceInstance
	code := instance.AuditUserByUserId(data.UserId)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.ErrByCode(code))
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(true))
	}
}

// AdminBanUser 审核禁用用户
// @Summary 审核禁用用户
// @Description 审核禁用用户
// @Tags Admin
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param pageExecutionRequest body dto.UserIdRequest true "用户id表单"
// @Success 200 {object} responses.Response{data=bool} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/admin/banUser [post]
func AdminBanUser(c *gin.Context) {
	var data dto.UserIdRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.ErrByCode(global.PARAMS_ERROR))
		return
	}

	instance := &adminServices.AdminServiceInstance
	code := instance.BanUserByUserId(data.UserId)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.ErrByCode(code))
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(true))
	}
}

// AdminCancelBanUser 审核从禁用状态中恢复用户
// @Summary 审核从禁用状态中恢复用户
// @Description 审核从禁用状态中恢复用户
// @Tags Admin
// @Accept json
// @Produce json
// @Param Token header string true "登录凭证，登录成功后会返回该凭证"
// @Param pageExecutionRequest body dto.UserIdRequest true "用户id表单"
// @Success 200 {object} responses.Response{data=bool} "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/admin/cancelBanUser [post]
func AdminCancelBanUser(c *gin.Context) {
	var data dto.UserIdRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.ErrByCode(global.PARAMS_ERROR))
		return
	}

	instance := &adminServices.AdminServiceInstance
	code := instance.CancelUserByUserId(data.UserId)
	if code != global.SUCCESS {
		c.JSON(http.StatusOK, baseRes.ErrByCode(code))
	} else {
		c.JSON(http.StatusOK, baseRes.OK.WithData(true))
	}
}
