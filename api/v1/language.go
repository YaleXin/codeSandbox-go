package api_v1

import (
	baseRes "codeSandbox/responses"
	sandboxService "codeSandbox/service/service_v1"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 获取支持的语言列表
func List(c *gin.Context) {
	languages := sandboxService.GetSupportLanguages()
	log.Debug("debug GetSupportLanguages")
	log.Info("info GetSupportLanguages")
	c.JSON(http.StatusOK, baseRes.OK.WithData(languages))
}
