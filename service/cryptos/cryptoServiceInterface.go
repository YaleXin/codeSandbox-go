package cryptos

import (
	"crypto/rsa"
)

type CryptoServiceInterface interface {
	// 生成密钥对
	generatePrivateKey() *rsa.PrivateKey
	GenerateKeyPairBase64() (string, string)

	// 从私钥字符串中还原密钥对
	getPrivateKeyFromBase64(base64PrivateKey string) *rsa.PrivateKey

	// 使用公钥加密数据
	EncryptedFromPulicStr(originData string, base64PublicKey string) string
	encryptedFromPulic(originData string, base64PublicKey string) string

	// 使用私钥解密数据
	DecryptedFromPrivateStr(encryptedData string, base64PrivateKey string)
	decryptedFromPrivate(encryptedData string, base64PrivateKey string)
}
