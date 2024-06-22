package test

import (
	"codeSandbox/service/cryptoServices"
	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"

	"testing"
)

func TestEncryptoAndDecrypto(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	origidata := "hello codesandbox"

	var cryptoService cryptoServices.CryptoServiceInterface = new(cryptoServices.CryptoService)
	publicKeyBase64, privateKeyBase64, err := cryptoService.GenerateRSAKeyPairBase64()
	if err != nil {
		t.Errorf("GenerateRSAKeyPairBase64 fail: %v", err)
	}

	base64, err := cryptoService.EncryptWithPublicKeyBase64(publicKeyBase64, origidata)
	if err != nil {
		t.Errorf("GenerateRSAKeyPairBase64 fail: %v", err)
	}

	keyBase64, err := cryptoService.DecryptWithPrivateKeyBase64(privateKeyBase64, base64)
	t.Log("keyBase64: ", keyBase64)
	if err != nil {
		t.Errorf("GenerateRSAKeyPairBase64 fail: %v", err)
	}

	assert.Equal(t, keyBase64, origidata)

}
