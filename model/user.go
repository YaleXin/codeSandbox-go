package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"index;type:varchar(40);not null;comment:用户名" json:"username"` // 用户名
	Email        string `gorm:"index;type:varchar(20);not null;comment:邮箱" json:"email"`     // 邮箱
	Password     string `gorm:"type:varchar(100);not null;comment:密码" json:"password"`       // 密码
	Salt         string `gorm:"type:varchar(20);not null;comment:加密盐" json:"salt"`           // 加密盐
	Role         int    `gorm:"type:int;comment:权限" json:"role"`                             // 权限 @global.ADMIN_USER_ROLE :管理员 , @global.NORMAL_USER_ROLE : 普通用户 （值越低，权限越高）
	CurrentUsage int    `gorm:"type:int;comment:本月已用" json:"currentUsage"`                   // 每调用一次执行代码，该值加一，每月清空
	MonthLimit   int    `gorm:"type:int;comment:每月限额" json:"monthLimit"`                     // 每月限额
	Ban          bool   `gorm:"type:tinyint;default(0);comment:是否被禁用" json:"ban"`            // 是否被禁用（关进小黑屋），禁用后，直接不能登录，也无法使用程序执行代码
	Audit        bool   `gorm:"type:tinyint;default(0);comment:是否已激活" json:"audit"`          // 是否已激活，注册后，需要管理员审核，未激活前，只能登录
}

func (u User) TableName() string {
	return "users"
}
