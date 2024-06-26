package cryptoServices

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const RSA_LEN int = 2048
const AES_LEN int = 32

var KEY_DATA_SEPARATOR [3]byte = [3]byte{0x00, 0x00, 0x00}

type CryptoService struct {
}

var CryptoServiceInstance CryptoService

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

// 使用公钥Base64字符串加密数据（采用混合加密的方式，即结合AES对称加密，否则直接采用 RSA 会收到加密数据长度的限制，虽然也可以使用分组加密的方式）
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

	// 生成一个随机的对称密钥
	symmetricKey := make([]byte, AES_LEN)

	if _, err = rand.Read(symmetricKey); err != nil {
		return "", err
	}

	// 使用对称密钥和AES-GCM模式加密数据
	aesBlock, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	encryptedData := gcm.Seal(nonce, nonce, []byte(originData), nil)

	// 使用RSA公钥加密对称密钥
	encryptedSymmetricKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, symmetricKey, nil)
	if err != nil {
		return "", err
	}

	// 将加密后的对称密钥和加密后的数据组合(分隔符之前是加密密钥，分隔符后是加密数据)
	separator := KEY_DATA_SEPARATOR[:]
	dataBytes := append(append(encryptedSymmetricKey, separator...), encryptedData...)
	// 字节转 base64
	encryptedBase64 := base64.StdEncoding.EncodeToString(dataBytes)
	return encryptedBase64, nil
}

// 使用私钥Base64字符串解密数据
func (service *CryptoService) DecryptWithPrivateKeyBase64(privateKeyBase64, encryptedBase64 string) (string, error) {
	// 密钥先转为字节切片
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

	//加密数据接着转为字节切片
	decodeString, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("error DecodeString privateKeyBase64: %v", err)
	}
	// 分离出 AES 和 加密数据
	separator := KEY_DATA_SEPARATOR[:]
	separatorIndex := bytes.Index(decodeString, separator)
	if separatorIndex == -1 {
		return "", fmt.Errorf("separator not found")
	}
	// 分离加密后的对称密钥和加密后的数据
	encryptedSymmetricKey := decodeString[:separatorIndex]
	encryptedData := decodeString[separatorIndex+len(separator):]

	// 将AES key解密成字节数组
	symmetricKeyBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedSymmetricKey, nil)
	if err != nil {
		return "", err
	}
	// 使用解密后的对称密钥和AES-GCM模式解密数据
	aesBlock, aesErr := aes.NewCipher(symmetricKeyBytes)
	if aesErr != nil {
		return "", aesErr
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	open, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	// 字节转为字符串
	return string(open), nil
}
