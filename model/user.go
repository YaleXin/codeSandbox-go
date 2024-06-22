package model

type User struct {
	Model
	Username string `gorm:"type:varchar(40);not null;comment:用户名" json:"username"` //用户名
	Password string `gorm:"type:varchar(100);not null;comment:密码" json:"password"` //密码
	Salt     string `gorm:"type:varchar(20);not null;comment:加密盐" json:"salt"`     //加密盐
	Role     int    `gorm:"type:int;comment:权限" json:"role"`                       //权限 0管理员 1 普通用户
}

func (u User) TableName() string {
	return "users"
}
