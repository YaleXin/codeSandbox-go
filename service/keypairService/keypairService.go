package keypairService

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/model/vo"
	"codeSandbox/service/cryptoServices"
	"codeSandbox/utils/global"
)

var keyPairDao db.KeyPairDao

type KeyPairService struct {
}

var KeyPairServiceInstance KeyPairService

func (keyPairService *KeyPairService) GetUserKeys(loginUser *model.User) (int, []vo.KeyPairVO) {

	keys, err := keyPairDao.ListKeyPairByUserId(loginUser.ID)
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	return global.SUCCESS, getKeyPairsVO(keys)
}
func (keyPairService *KeyPairService) GenerateUserKey(user *model.User) (int, *model.KeyPair) {
	service := cryptoServices.CryptoService{}
	pub, pri, err := service.GenerateRSAKeyPairBase64()
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	keyPair := model.KeyPair{
		User:      *user,
		UserId:    user.ID,
		AccessKey: pub,
		SecretKey: pri,
	}
	_, err = keyPairDao.KeyPairAdd(&keyPair)
	if err != nil {
		return global.SYSTEM_ERROR, nil
	}
	return global.SUCCESS, &keyPair
}

// 将密钥对切片转为脱敏后的信息
func getKeyPairsVO(keys []model.KeyPair) []vo.KeyPairVO {
	pairVO := make([]vo.KeyPairVO, 0, len(keys))
	for _, keypair := range keys {
		pairVO = append(pairVO, vo.KeyPairVO{
			ID:        keypair.ID,
			AccessKey: keypair.AccessKey,
			SecretKey: keypair.SecretKey,
			UserId:    keypair.UserId,
			CreatedAt: keypair.CreatedAt,
		})
	}
	return pairVO
}
