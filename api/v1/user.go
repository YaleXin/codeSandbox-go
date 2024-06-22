package api_v1

import (
	"codeSandbox/model"
	baseRes "codeSandbox/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 注册
func Register(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithData("error"))
	}

}

// Login 登陆
func Login(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, baseRes.Err.WithData("error"))
	}
}
