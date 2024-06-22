package tool

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// 通用分页获取
func PageParse(c *gin.Context) (int, int) {
	//每页数据
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	// 页码数
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))

	if pageSize <= 0 || pageSize > 50 {
		pageSize = 50
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	return pageSize, pageNum
}
