package userServices

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/model/dto"
	"codeSandbox/model/vo"
	"codeSandbox/service/keypairService"
	"codeSandbox/utils/global"
	"codeSandbox/utils/middleware"
	"codeSandbox/utils/tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
	"time"
)

var userDao db.UserDao

const SALT_LEN = 20

type UserService struct {
}

var UserServiceInstance UserService

func verifyPwd(submitUser, databaseUser *model.User) bool {
	return encryptPwdWithSalt(submitUser.Password, databaseUser.Salt) == databaseUser.Password
}

// 检查该用户本地是否可以执行代码（是否还有额度）并增加相应调用次数
func (userService *UserService) CheckAndUpdateUserUsage(userId uint) int {
	user := model.User{
		Username: "aa",
		Model: gorm.Model{
			ID: userId,
		},
	}
	_, err := userDao.GetUserById(&user, userId)
	if err != nil {
		return global.NOT_FOUND_USER_ERROR
	}
	// TODO 把过期时间去掉
	createdAt := user.CreatedAt
	now := time.Now()
	if now.Sub(createdAt).Hours() > 24*global.USER_VALIDITY_PERIOD {
		return global.USER_TRY_TIME_EXPIRE
	}

	// 判断额度
	if user.CurrentUsage >= user.MonthLimit {
		return global.INSUFF_AMOUNT_ERROR
	}
	// 使用次数加一
	user.CurrentUsage += 1
	userDao.UpdateUserById(&user)
	return global.SUCCESS
}

func (userService *UserService) UserList() (int, []model.User) {
	user, err := userDao.ListUser()
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	return global.SUCCESS, user
}
func (userService *UserService) GetLoginUser(c *gin.Context) (int, *model.User) {
	get, exists := c.Get("user")
	if !exists {
		return global.NOT_LOGIN_ERROR, nil
	}
	// 使用类型断言转换为 MyClaims
	if myClaims, ok := get.(*middleware.MyClaims); ok {
		return global.SUCCESS, &model.User{
			Model: gorm.Model{
				ID: myClaims.UserId,
			},
			Username: myClaims.Username,
		}
	} else {
		// 断言失败，anyValue不是 MyClaims 类型
		return global.SYSTEM_ERROR, nil
	}

}
func (userService *UserService) GetUserKeys(c *gin.Context) (int, []vo.KeyPairVO) {
	// 获取当前登录用户
	code, loginUser := userService.GetLoginUser(c)
	if code != global.SUCCESS {
		return code, nil
	}
	// 查出该用户拥有的密钥对
	keypairServiceInstance := &keypairService.KeyPairServiceInstance
	code, keyPairVOs := keypairServiceInstance.GetUserKeys(loginUser)
	if code != global.SUCCESS {
		return code, nil
	}
	return global.SUCCESS, keyPairVOs
}

func checkLoginUser(user *model.User) bool {
	return user != nil && !tool.IsBlankString(user.Username) && !tool.IsBlankString(user.Password)
}

// 返回执行结果，用户id， jwt token
func (userService *UserService) UserLogin(submitUser *model.User) (int, *vo.UserVO) {
	if !checkLoginUser(submitUser) {
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

func checkRegisterUser(user *model.User) bool {
	return !tool.IsStructEmpty(user) && isValidPasswd(user.Password) && isValidEmail(user.Email)
}

func isValidPasswd(password string) bool {
	// 正则表达式，分别匹配小写字母、大写字母和数字
	pattern := `(?=.*[a-z])(?=.*[A-Z])(?=.*\d)`
	matched, _ := regexp.MatchString(pattern, password)
	return matched && len(password) >= 8 && len(password) <= 16
}

func isValidEmail(email string) bool {
	//开始是一个或多个字母、数字、点号、百分号、加号、减号或下划线。
	//然后是一个 "@" 符号。
	//接着是一个或多个字母、数字、点号或减号。
	//最后是一个点号，后面跟着两个或更多的字母（代表顶级域名）。
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	r := regexp.MustCompile(pattern)
	// 使用正则表达式的 MatchString 方法进行匹配
	return r.MatchString(email)
}
func encryptPwdWithSalt(pwd string, salt string) string {
	return tool.MD5Str(pwd + salt)
}

func (userService *UserService) UserRegister(userRegisterRequest *dto.UserRegisterRequest) int {
	user := model.User{
		Username: userRegisterRequest.Username,
		Password: userRegisterRequest.Password,
		Email:    userRegisterRequest.Email,
	}
	if !checkRegisterUser(&user) {
		return global.PARAMS_ERROR
	}
	_, err := userDao.GetUserByName(&user)
	// 查不到时候会报 error
	if err == nil {
		return global.DATA_REPEAT_ERROR
	}

	// 使用盐进行加密
	salt := tool.GenerateRandomVisibleString(SALT_LEN)
	md5Str := encryptPwdWithSalt(user.Password, salt)

	user.Password = md5Str
	user.Salt = salt
	user.Role = global.NORMAL_USER_ROLE
	_, err = userDao.UserAdd(&user)
	if err != nil {
		return global.SYSTEM_ERROR
	}
	return global.SUCCESS
}
func (userService *UserService) GenerateKeyPair(c *gin.Context) (int, *vo.KeyPairVO) {
	// 获取当前登录用户
	code, loginUser := userService.GetLoginUser(c)
	if code != global.SUCCESS {
		return code, nil
	}
	keypairServiceInstance := &keypairService.KeyPairServiceInstance
	code, keyPair := keypairServiceInstance.GenerateUserKey(loginUser)
	if code != global.SUCCESS {
		return code, nil
	}
	// 自行封装成 VO
	pairVO := vo.KeyPairVO{
		ID:        keyPair.ID,
		SecretKey: keyPair.SecretKey,
		AccessKey: keyPair.AccessKey,
		UserId:    loginUser.ID,
		CreatedAt: keyPair.CreatedAt,
	}
	return global.SUCCESS, &pairVO
}

func (userService *UserService) UserDelete(user *model.User) (*model.User, error) {
	return nil, nil
}

func (userService *UserService) GetUserById(id int) (*model.User, error) {
	return nil, nil
}
