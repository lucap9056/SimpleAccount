package Auths

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const saltSize = 16

func Salt() string {
	saltBytes := make([]byte, saltSize)
	rand.Read(saltBytes)
	return hex.EncodeToString(saltBytes)
}

func Hash(salt string, value string) string {
	saltBytes := []byte(salt)
	hasher := sha256.New()
	combined := append([]byte(value), saltBytes...)
	hasher.Write(combined)
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

type Auth struct {
	PrivateKeyFilePath    string
	PublicKeyFilePath     string
	privateKey            *rsa.PrivateKey
	publicKey             *rsa.PublicKey
	Cache                 AuthCache
	ValidityDuration      time.Duration
	RenewTime             time.Duration
	CacheValidityDuration time.Duration
}

func (auth *Auth) Init() error {
	auth.Cache.Init(auth.CacheValidityDuration)
	//確認 RSA256 公私鑰是否都存在
	if auth.rsaKeyExist() {
		var err error
		// 讀取私鑰
		auth.privateKey, err = auth.readPrivateKey()
		if err != nil {
			return err
		}
		// 讀取公鑰
		auth.publicKey, err = auth.readPublicKey()
		if err != nil {
			return err
		}
	} else {
		//生成並讀取公私鑰
		err := auth.generateRsaKey()
		if err != nil {
			return err
		}
	}

	return nil
}
