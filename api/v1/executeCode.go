package api_v1

import (
	"codeSandbox/model/dto"
	baseRes "codeSandbox/responses"
	service "codeSandbox/service/sandboxServices"
	"codeSandbox/utils/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ExecuteCode 处理执行代码的请求
// @Summary 执行代码
// @Description 根据用户提交的代码和语言执行代码并返回结果
// @Tags Code Execution
// @Accept json
// @Produce json
// @Param executeCodeRequest body dto.ExecuteCodeRequest true "执行代码请求"
// @Success 200 {object} responses.Response "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/executeCode [post]
func ExecuteCode(c *gin.Context) {
	var exeRequest dto.ExecuteCodeRequest
	if c.ShouldBind(&exeRequest) == nil && exeRequest.Code != "" && exeRequest.Language != "" {
		sandboxService := service.SandboxService{}
		code, executeData := sandboxService.ExecuteCode(c, exeRequest, nil)
		if code != global.SUCCESS {
			c.JSON(http.StatusOK, baseRes.Err.WithMsg(global.GetErrMsg(code)))
		} else {
			c.JSON(http.StatusOK, baseRes.OK.WithData(executeData))
		}
	} else {
		c.JSON(http.StatusOK, baseRes.Err.WithData("error"))
	}
}

func ProgramExecuteCode(c *gin.Context) {
	var programExecuteCodeRequest dto.ProgramExecuteCodeRequest
	if c.ShouldBind(&programExecuteCodeRequest) == nil && programExecuteCodeRequest.PublicKey != "" && programExecuteCodeRequest.PublicKey != "" {
		sandboxService := service.SandboxService{}
		code, executeCode := sandboxService.ProgramExecuteCode(c, &programExecuteCodeRequest)
		if code != global.SUCCESS {
			c.JSON(http.StatusOK, baseRes.Err.WithData(global.GetErrMsg(code)))
			return
		} else {
			c.JSON(http.StatusOK, baseRes.OK.WithData(executeCode))
		}
	} else {
		c.JSON(http.StatusOK, baseRes.Err.WithData("error"))
	}
}
