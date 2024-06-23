package userServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/model/vo"
	"codeSandbox/utils/global"
	"codeSandbox/utils/middleware"
	"codeSandbox/utils/tool"
)

var userDao db.UserDao

const SALT_LEN = 20
const NORMAL_USER_ROLE = int(1)
const ADMIN_USER_ROLE = int(2)

type UserService struct {
}

var UserServiceInstance UserService

func verifyPwd(submitUser, databaseUser *model.User) bool {
	return encryptPwdWithSalt(submitUser.Password, databaseUser.Salt) == databaseUser.Password
}
func (userService *UserService) UserList() (int, []model.User) {
	user, err := userDao.ListUser()
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	return global.SUCCESS, user
}

// 返回执行结果，用户id， jwt token
func (userService *UserService) UserLogin(submitUser *model.User) (int, *vo.UserVO) {
	if !checkUser(submitUser) {
		return global.PARAMS_ERROR, nil
	}
	// 查询数据库中该用户信息
	databaseUser := model.User{
		Username: submitUser.Username,
	}
	_, err := userDao.GetUserByName(&databaseUser)
	if err != nil {
		return global.NOT_FOUND_USER_ERROR, nil
	}
	if !verifyPwd(submitUser, &databaseUser) {
		return global.PWD_ERROR, nil
	}
	token, tCode := middleware.SetToken(databaseUser.ID, databaseUser.Username, databaseUser.Role)
	if tCode != global.SUCCESS {
		return tCode, nil

	}
	var userVO vo.UserVO
	getUserVO(&databaseUser, token, &userVO)
	return global.SUCCESS, &userVO
}
func getUserVO(user *model.User, token string, userVO *vo.UserVO) {
	userVO.Id = user.ID
	userVO.Username = user.Username
	userVO.Role = user.Role
	userVO.Token = token
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

func (userService *UserService) UserRegister(user *model.User) int {
	if !checkUser(user) {
		return global.PARAMS_ERROR
	}
	_, err := userDao.GetUserByName(user)
	// 查不到时候会报 error
	if err == nil {
		return global.DATA_REPEAT_ERROR
	}

	// 使用盐进行加密
	salt := tool.GenerateRandomVisibleString(SALT_LEN)
	md5Str := encryptPwdWithSalt(user.Password, salt)

	user.Password = md5Str
	user.Salt = salt
	user.Role = NORMAL_USER_ROLE
	_, err = userDao.UserAdd(user)
	if err != nil {
		return global.SYSTEM_ERROR
	}
	return global.SUCCESS
}

func (userService *UserService) UserDelete(user *model.User) (*model.User, error) {
	return nil, nil
}

func (userService *UserService) GetUserById(id int) (*model.User, error) {
	return nil, nil
}
