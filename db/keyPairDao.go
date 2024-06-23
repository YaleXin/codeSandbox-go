package db

import "codeSandbox/model"

type KeyPairDao struct {
}

// 找出没有被 “delete” 的
func (k *KeyPairDao) ListKeyPair() ([]model.KeyPair, error) {
	var keyPairs []model.KeyPair
	find := dBClinet.Find(&keyPairs)
	err := find.Error
	if err != nil {
		return nil, err
	}
	return keyPairs, nil
}

func (k *KeyPairDao) ListKeyPairByUserId(userId uint) ([]model.KeyPair, error) {
	var keyPairs []model.KeyPair
	// 使用Preload预加载User关系
	result := dBClinet.Preload("User").Where("user_id = ?", userId).Find(&keyPairs)
	if result.Error != nil {
		return nil, result.Error // 返回查询过程中可能遇到的错误
	}
	return keyPairs, nil
}

func (k *KeyPairDao) KeyPairAdd(keypair *model.KeyPair) (int64, error) {
	create := dBClinet.Create(keypair)
	err := create.Error
	if err != nil {
		return 0, err
	}
	affected := create.RowsAffected
	return affected, nil
}
