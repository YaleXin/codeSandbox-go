package model

import "gorm.io/gorm"

type KeyPair struct {
	gorm.Model
	AccessKey string `gorm:"index;type:varchar(650);not null;comment:公钥" json:"accessKey"`
	SecretKey string `gorm:"type:varchar(2300);not null;comment:私钥" json:"secretKey"`
	UserId    uint   `json:"userId"`
	User      User
}

func (k KeyPair) TableName() string {
	return "keypairs"
}
