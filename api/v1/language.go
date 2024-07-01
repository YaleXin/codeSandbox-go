package api_v1

import (
	baseRes "codeSandbox/responses"
	sandboxService "codeSandbox/service/sandboxServices"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// LanguageList
// @Summary 获取支持的语言列表
// @Description 获取支持的语言列表，只有在该列表中的语言代码才能运行
// @Tags Languages
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response "成功响应"
// @Failure 400 {object} responses.Response "错误响应"
// @Failure 500 {object} responses.Response "系统内部错误"
// @Router /api/v1/languages [get]
func LanguageList(c *gin.Context) {
	service := sandboxService.SandboxService{}
	languages := service.GetSupportLanguages()
	log.Debug("languages = ", languages)
	c.JSON(http.StatusOK, baseRes.OK.WithData(languages))
}
