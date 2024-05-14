package Auths

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strings"
)

const defaultPrivateKeyPath = "./keys/default.pem"
const defaultPublicKeyPath = "./keys/default.pem"

func (auth *Auth) rsaKeyExist() bool {
	if auth.PrivateKeyFilePath == "" {
		auth.PrivateKeyFilePath = defaultPrivateKeyPath
	}

	if auth.PublicKeyFilePath == "" {
		auth.PublicKeyFilePath = defaultPublicKeyPath
	}
	_, privateKey := os.Stat(auth.PrivateKeyFilePath)
	_, publicKey := os.Stat(auth.PublicKeyFilePath)

	return !(os.IsNotExist(privateKey) || os.IsNotExist(publicKey))
}

func (auth *Auth) generateRsaKey() error {
	// 生成 RSA 私鑰
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	// 將私鑰編碼成 ASN.1 DER 格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	// 建立私鑰檔案
	privateKeyFile, err := os.Create(auth.PrivateKeyFilePath)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	// 將私鑰寫入檔案
	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		return err
	}

	// 生成 RSA 公鑰
	publicKey := &privateKey.PublicKey
	// 將公鑰編碼成 ASN.1 DER 格式
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	// 判定私鑰檔案與公鑰檔案是否相同
	if auth.PrivateKeyFilePath == auth.PublicKeyFilePath {
		// 相同
		err = pem.Encode(privateKeyFile, publicKeyPEM)
		if err != nil {
			return err
		}
	} else {
		// 不相同 建立公鑰檔案
		publicKeyFile, err := os.Create(auth.PublicKeyFilePath)
		if err != nil {
			return err
		}
		defer publicKeyFile.Close()
		err = pem.Encode(publicKeyFile, publicKeyPEM)
		if err != nil {
			return err
		}
	}
	// 將公鑰寫入檔案

	auth.privateKey = privateKey
	auth.publicKey = publicKey
	return nil
}

func (auth *Auth) readPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.ReadFile(auth.PrivateKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyFile)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (auth *Auth) readPublicKey() (*rsa.PublicKey, error) {
	publicKeyFile, err := os.ReadFile(auth.PublicKeyFilePath)
	if err != nil {
		return nil, err
	}

	publicKeyFileStr := string(publicKeyFile)
	index := strings.Index(publicKeyFileStr, "-----BEGIN RSA PUBLIC KEY-----")
	if index == -1 {
		return nil, errors.New("invalid public key")
	}
	publicKeyBytes := []byte(publicKeyFileStr[index:])

	block, _ := pem.Decode(publicKeyBytes)
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (auth *Auth) rsaEncode(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, auth.publicKey, data)
}

func (auth *Auth) RsaDecode(data []byte) (*[]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, auth.privateKey, data)
	if err != nil {
		return nil, err
	}

	return &decryptedData, nil
}
