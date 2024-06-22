package test

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserDaoFind(t *testing.T) {
	userDao := db.UserDao{}
	list, err := userDao.ListUser()
	if err != nil {
		t.Errorf("error : %v", err)
	}
	if len(list) == 0 {
		t.Error("empty")
	}
}

func TestUserDaoAdd(t *testing.T) {
	userDao := db.UserDao{}
	user := model.User{
		Username: "yalexin",
		Password: "password",
		Salt:     "salt",
		Role:     11,
	}
	add, err := userDao.UserAdd(&user)
	if err != nil {
		t.Errorf("add fail %v", err)
	}
	if add == 0 {
		t.Errorf("add fail add = %v", add)
	}
	t.Log("add count = ", add)

}
func TestUserUpdate(t *testing.T) {
	newUser := model.User{
		Username: "original",
		Password: "password",
		Salt:     "salt",
		Role:     11,
	}
	userDao := db.UserDao{}
	_, err := userDao.UserAdd(&newUser)
	id := newUser.ID
	if err != nil {
		t.Fatalf("add")
	}

	newUser.Username = "new-username"
	_, err = userDao.UpdateUserById(&newUser)
	if err != nil {
		t.Fatalf("update")
	}
	queryUser := model.User{}
	byId, err := userDao.GetUserById(&queryUser, id)
	if err != nil {
		t.Fatalf("query")
	}
	assert.NotEqual(t, byId.Username, "original")
}

func TestUserDelete(t *testing.T) {
	newUser := model.User{
		Username: "original",
		Password: "password",
		Salt:     "salt",
		Role:     11,
	}
	userDao := db.UserDao{}
	_, err := userDao.UserAdd(&newUser)
	id := newUser.ID
	t.Logf("insert id : %v", id)
	if err != nil {
		t.Fatalf("add")
	}
	assert.Equal(t, newUser.DeletedAt.Valid, false)

	err = userDao.DeleteUserById(id)
	if err != nil {
		t.Fatalf("delete")
	}
	assert.NotEqual(t, newUser.DeletedAt.Valid, true)
}
