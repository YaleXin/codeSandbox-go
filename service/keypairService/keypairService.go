package keypairService

import (
	"codeSandbox/db"
	"codeSandbox/model"
	"codeSandbox/model/vo"
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

// 将密钥对切片转为脱敏后的信息
func getKeyPairsVO(keys []model.KeyPair) []vo.KeyPairVO {
	pairVO := make([]vo.KeyPairVO, 0, len(keys))
	for _, keypair := range keys {
		pairVO = append(pairVO, vo.KeyPairVO{
			ID:        keypair.ID,
			AccessKey: keypair.AccessKey,
			SecretKey: keypair.SecretKey,
			UserId:    keypair.UserId,
		})
	}
	return pairVO
}
