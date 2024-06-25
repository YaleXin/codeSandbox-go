package test

import (
	"codeSandbox/model"
	"codeSandbox/service/userServices"
	"codeSandbox/utils/global"
	"codeSandbox/utils/tool"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestUserRegister(t *testing.T) {
	passwd := "password"
	user := model.User{
		Username: "gthtrh",
		Password: passwd,
	}
	userService := userServices.UserService{}
	err := userService.UserRegister(&user)
	if err != global.SUCCESS {
		t.Errorf("err:%v", err)
	}
	if tool.IsBlankString(user.Salt) {
		t.Errorf("salt is blank")
	}
	if tool.MD5Str(passwd+user.Salt) != user.Password {
		t.Errorf("md5(passwd + salt) != password")
	}
}
func TestExistUserRegister(t *testing.T) {
	passwd := "password"
	user := model.User{
		Username: "yalexin",
		Password: passwd,
	}
	userService := userServices.UserService{}
	err := userService.UserRegister(&user)
	assert.NotEqual(t, err, nil)

}

func TestUserLogin(t *testing.T) {
	passwd := tool.GenerateRandomVisibleString(10)
	username := tool.GenerateRandomVisibleString(10)
	user := model.User{
		Username: username,
		Password: passwd,
	}
	userService := userServices.UserService{}
	err := userService.UserRegister(&user)
	assert.Equal(t, err, nil)
	loginUser := model.User{
		Username: user.Username,
		Password: passwd,
	}
	login, _ := userService.UserLogin(&loginUser)
	assert.Equal(t, login, global.SUCCESS)
}
