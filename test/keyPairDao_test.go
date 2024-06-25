package test

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/service/cryptoServices"
	"codeSandbox/service/userServices"
	"codeSandbox/utils/tool"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestKeyPairDaoAdd(t *testing.T) {
	passwd := tool.GenerateRandomVisibleString(10)
	username := tool.GenerateRandomVisibleString(10)
	user := model.User{
		Username: username,
		Password: passwd,
	}
	userService := userServices.UserService{}
	err := userService.UserRegister(&user)
	assert.Equal(t, err, nil)
	t.Logf("user:%v", user)

	keyPairDao := db.KeyPairDao{}
	service := cryptoServices.CryptoService{}
	pub, pri, errO := service.GenerateRSAKeyPairBase64()
	t.Logf("len pri:%v, len pub:%v", len(pri), len(pub))
	if errO != nil {
		t.Errorf("%v", err)
	}
	keyPair := model.KeyPair{
		SecretKey: pri,
		AccessKey: pub,
		User:      user,
		UserId:    user.ID,
	}
	_, errO = keyPairDao.KeyPairAdd(&keyPair)
	if errO != nil {
		t.Fatalf("%v", err)
	}

	keys, errO := keyPairDao.ListKeyPairByUserId(user.ID)
	if errO != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("keys %v", keys)
}
