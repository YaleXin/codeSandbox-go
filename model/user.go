package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"index;type:varchar(40);not null;comment:用户名" json:"username"` //用户名
	Password string `gorm:"type:varchar(100);not null;comment:密码" json:"password"`       //密码
	Salt     string `gorm:"type:varchar(20);not null;comment:加密盐" json:"salt"`           //加密盐
	Role     int    `gorm:"type:int;comment:权限" json:"role"`                             //权限 0管理员  10 普通用户 （值越低，权限越高）
}

func (u User) TableName() string {
	return "users"
}
