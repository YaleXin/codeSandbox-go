package userServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/utils/global"
	"codeSandbox/utils/tool"
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
func encryptPwdWithSalt(pwd string, salt string) string {
	return tool.MD5Str(pwd + salt)
}

func (userService *UserService) UserRegister(user *model.User) error {
	if !checkUser(user) {
		return &global.CustomError{
			ErrorCode: global.PARAMS_ERROR,
			Message:   global.GetErrMsg(global.PARAMS_ERROR),
		}
	}
	_, err := userDao.GetUserByName(user)
	// 查不到时候会报 error
	if err == nil {
		return &global.CustomError{
			ErrorCode: global.DATA_REPEAT,
			Message:   global.GetErrMsg(global.DATA_REPEAT),
		}
	}

	// 使用盐进行加密
	salt := tool.GenerateRandomVisibleString(SALT_LEN)
	md5Str := encryptPwdWithSalt(user.Password, salt)

	user.Password = md5Str
	user.Salt = salt
	user.Role = NORMAL_USER_ROLE
	_, err = userDao.UserAdd(user)
	if err != nil {
		return &global.CustomError{
			ErrorCode: global.SYSTEM_ERROR,
			Message:   global.GetErrMsg(global.SYSTEM_ERROR),
		}
	}
	return nil
}

func (userService *UserService) UserDelete(user *model.User) (*model.User, error) {
	return nil, nil
}

func (userService *UserService) GetUserById(id int) (*model.User, error) {
	return nil, nil
}
