package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExecuteCode(c *gin.Context) {
	c.String(http.StatusOK, "ExecuteCode")
}
