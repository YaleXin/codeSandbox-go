package cryptoServices

type CryptoServiceInterface interface {
	// 生成密钥对，并将私钥和公钥转为 base64
	GenerateRSAKeyPairBase64() (string, string, error)

	// 使用公钥加密数据
	EncryptWithPublicKeyBase64(publicKeyBase64, originData string) (string, error)

	// 使用私钥解密数据
	DecryptWithPrivateKeyBase64(privateKeyBase64, encryptedBase64 string) (string, error)
}
