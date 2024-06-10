package api_v1

import (
	baseRes "codeSandbox/responses"
	sandboxService "codeSandbox/service/service_v1"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetTodo
// @Summary 获取支持的语言列表
// @Description 获取支持的语言列表，只有在该列表中的语言代码才能运行
// @Tags Code
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/languages [get]
func List(c *gin.Context) {
	languages := sandboxService.GetSupportLanguages()
	log.Debug("languages = ", languages)
	c.JSON(http.StatusOK, baseRes.OK.WithData(languages))
}
