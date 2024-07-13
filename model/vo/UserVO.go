package vo

type UserVO struct {
	Id       uint   `json:"id"`
	Username string `json:"username"` // 用户名
	Role     int    `json:"role"`     // 权限 1:管理员 10: 普通用户
	Token    string `json:"token"`    // 登陆凭证
	Audit    bool   `json:"audit"`    // 是否已通过审核
}
