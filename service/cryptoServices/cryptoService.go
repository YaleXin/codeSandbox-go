package cryptoServices

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const RSA_LEN int = 2048

type CryptoService struct {
}

// 生成RSA密钥对，并将公钥和私钥分别转为Base64字符串
func (service *CryptoService) GenerateRSAKeyPairBase64() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RSA_LEN)
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key: %v", err)
	}

	//私钥转为 pem 字节数组
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	})
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)

	// 拆分出私钥
	publicKey := &privateKey.PublicKey
	// 转为 pem 数组
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyPEM)

	return publicKeyBase64, privateKeyBase64, nil
}

// 使用公钥Base64字符串加密数据
func (service *CryptoService) EncryptWithPublicKeyBase64(publicKeyBase64, originData string) (string, error) {
	// 先解码为字节数组
	publicKeyPEM, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}
	// 从字节数组还原
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", fmt.Errorf("invalid public key PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing public key: %v", err)
	}

	// 断言是 rsa 公钥类型
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("not an RSA public key")
	}

	// 加密成字节数组
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(originData))
	if err != nil {
		return "", fmt.Errorf("error encrypting data: %v", err)
	}

	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedData)
	return encryptedBase64, nil
}

// 使用私钥Base64字符串解密数据
func (service *CryptoService) DecryptWithPrivateKeyBase64(privateKeyBase64, encryptedBase64 string) (string, error) {
	privateKeyPEM, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding private key: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("invalid private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing private key: %v", err)
	}

	// 将原先加密的base64数据转为加密的字节数组
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding encrypted data: %v", err)
	}
	// 解密成字节数组
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedData)
	if err != nil {
		return "", fmt.Errorf("error decrypting data: %v", err)
	}
	// 字节转为字符串
	return string(decryptedData), nil
}
