package dto

// @Summary 用户登录请求参数
// @Description 用户注册请求参数
// @Accept json
// @Param username body string true "用户名"
// @Param password body string true "密码"
type UserLoginRequest struct {
	Username string `bind:"required" json:"username"` // 用户名
	Password string `bind:"required" json:"password"` // 密码
}
