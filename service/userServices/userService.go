package userServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/utils/tool"
	"fmt"
)

var userDao db.UserDao

const SALT_LEN = 20
const NORMAL_USER_ROLE = int(1)
const ADMIN_USER_ROLE = int(2)

type UserService struct {
}

func (userService *UserService) UserLogin(user model.User) bool {
	return false
}

func (userService *UserService) UserLogout(user model.User) bool {
	return false
}

func checkUser(user *model.User) bool {
	return user != nil && !tool.IsBlankString(user.Username) && !tool.IsBlankString(user.Password)
}
func (userService *UserService) UserRegister(user *model.User) error {
	if !checkUser(user) {
		return fmt.Errorf("")
	}
	// 使用盐进行加密
	salt := tool.GenerateRandomVisibleString(SALT_LEN)
	md5Str := tool.MD5Str(user.Password + salt)

	user.Password = md5Str
	user.Salt = salt
	user.Role = NORMAL_USER_ROLE
	_, err := userDao.UserAdd(user)
	return err
}

func (userService *UserService) UserDelete(user *model.User) (*model.User, error) {
	return nil, nil
}

func (userService *UserService) GetUserById(id int) (*model.User, error) {
	return nil, nil
}
