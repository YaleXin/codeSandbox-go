package db

import (
	"codeSandbox/model"
	"fmt"
)

type UserDao struct {
}

func (u *UserDao) UserAdd(user *model.User) (int64, error) {
	create := dBClinet.Create(user)
	err := create.Error
	if err != nil {
		return 0, err
	}
	affected := create.RowsAffected
	return affected, nil
}

// 找出没有被 “delete” 的
func (u *UserDao) ListUser() ([]model.User, error) {
	var users []model.User
	find := dBClinet.Find(&users)
	err := find.Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (u *UserDao) GetUserById(user *model.User, id uint) (*model.User, error) {
	first := dBClinet.First(user, id)
	err := first.Error
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return user, nil
}
func (u *UserDao) UpdateUserById(user *model.User) (*model.User, error) {
	// gorm 中，没有主键时会调用 create ，这是不符合我们的
	if user.ID == uint(0) {
		return nil, fmt.Errorf("user.id == 0")
	}
	// 根据 `struct` 更新属性，只会更新非零值的字段
	updates := dBClinet.Updates(user)
	err := updates.Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserDao) DeleteUserById(id uint) error {
	tx := dBClinet.Delete(&model.User{}, id)
	err := tx.Error
	if err != nil {
		return err
	}
	return nil
}
