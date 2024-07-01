package dto

// @Summary 用户注册请求参数
// @Description 用户注册请求参数
// @Accept json
// @Param username body string true "用户名"
// @Param email body string true "邮箱"
// @Param password body string true "密码"
type UserRegisterRequest struct {
	Username string `bind:"required" json:"username"` // 用户名
	Email    string `bind:"required" json:"email"`    // 邮箱
	Password string `bind:"required" json:"password"` // 密码
}
