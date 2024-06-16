package api_v1

import (
	"codeSandbox/model/dto"
	baseRes "codeSandbox/responses"
	"codeSandbox/service/sandbox"
	"codeSandbox/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ExecuteCodeHandler 处理执行代码的请求
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
		sandbox := sandbox.SandBox{
			DockerInfo: utils.DockerInfo{
				Language:       "Go",
				ImageName:      "golang:1.17",
				Filename:       "Main.go",
				CompileCmd:     "go build Main.go",
				RunCmd:         "./Main",
				ContainerCount: 1,
			},
		}
		code := sandbox.ExecuteCode(&exeRequest)
		c.JSON(http.StatusOK, baseRes.OK.WithData(code))
	} else {
		c.JSON(http.StatusOK, baseRes.Err.WithData("error"))
	}
}
