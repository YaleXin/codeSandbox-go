package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取支持的语言列表
func List(c *gin.Context) {
	c.String(http.StatusOK, "List")
}
