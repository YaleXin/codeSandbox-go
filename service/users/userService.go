package users

import "codeSandbox/model"

type UserService struct {
}

func (userService *UserService) UserLogin(user model.User) bool {
	return false
}

func (userService *UserService) UserLogout(user model.User) bool {
	return false
}

func (userService *UserService) UserRegister(user *model.User) (*model.User, error) {
	return nil, nil
}

func (userService *UserService) UserDelete(user *model.User) (*model.User, error) {
	return nil, nil
}

func (userService *UserService) GetUserById(id int) (*model.User, error) {
	return nil, nil
}
