package vo

import "time"

type UserDetailVO struct {
	Id           uint      `json:"id"`
	Username     string    `json:"username"`     //用户名
	Role         int       `json:"role"`         //权限 0管理员 1 普通用户
	Email        string    `json:"email"`        // 邮箱
	MonthLimit   int       `json:"monthLimit"`   // 每月限额
	CurrentUsage int       `json:"currentUsage"` // 每月已用额度
	CreateAt     time.Time `json:"createAt"`     //注册时间
}
