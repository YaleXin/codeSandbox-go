package test

import (
	"codeSandbox/model"
	"codeSandbox/service/userServices"
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
	if err != nil {
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
