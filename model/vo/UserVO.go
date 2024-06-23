package vo

type UserVO struct {
	Id       uint   `json:"id"`
	Username string `json:"username"` //用户名
	Role     int    `json:"role"`     //权限 0管理员 1 普通用户
	Token    string `json:"token"`
}
