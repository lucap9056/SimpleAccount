package Auths

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	MathRand "math/rand"
	"os"
	"path/filepath"
	"simple_account/app/Error"
	"strconv"
	"strings"
	"time"
)

const defaultKeyPath = "./keys/"

type Key struct {
	id      string
	iat     time.Time
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

func (auth *Auth) rsaKeysExist() ([]string, error) {
	keys := []string{}
	if auth.KeyFilesPath == "" {
		auth.KeyFilesPath = defaultKeyPath
	}

	_, Keys := os.Stat(auth.KeyFilesPath)
	if os.IsNotExist(Keys) {
		err := os.MkdirAll(auth.KeyFilesPath, 0644)
		return keys, err
	}

	files, err := os.ReadDir(auth.KeyFilesPath)
	if err != nil {
		return keys, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".key") {
			keyId := strings.TrimSuffix(file.Name(), ".key")
			keys = append(keys, keyId)
		}
	}

	return keys, nil
}

func (auth *Auth) removeExpiredKeys() {
	auth.muxKeys.Lock()
	defer auth.muxKeys.Unlock()
	currentTime := time.Now()
	for _, key := range auth.keys {
		expiredTime := key.iat.Add(auth.KeyValidityDuration)
		if currentTime.Before(expiredTime) {
			auth.expiredKeys[key.id] = key
			delete(auth.keys, key.id)
		}
	}

	for _, key := range auth.expiredKeys {
		expiredTime := key.iat.Add(auth.KeyValidityDuration * 2)
		if currentTime.Before(expiredTime) {
			delete(auth.expiredKeys, key.id)
			auth.removeExpiredKeyFile(key.id)
		}
	}
}

func (auth *Auth) removeExpiredKeyFile(keyId string) {
	filePath := filepath.Join(auth.KeyFilesPath, keyId+".key")
	err := os.Remove(filePath)
	if err != nil {
		auth.logger.Error.Write(err)
	}
}

func (auth *Auth) generateKey() (Key, error) {
	var key Key
	// 生成 RSA 私鑰
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return key, err
	}
	// 將私鑰編碼成 ASN.1 DER 格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	// 將私鑰編碼為Base64
	privateKeyStr := base64.StdEncoding.EncodeToString(privateKeyBytes)

	//生成公鑰
	publicKey := &privateKey.PublicKey
	// 將公鑰編碼成 ASN.1 DER 格式
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	// 將公鑰編碼為Base64
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKeyBytes)

	//生成ID
	keyIdBytes := make([]byte, 8)
	_, err = rand.Read(keyIdBytes)
	if err != nil {
		return key, err
	}
	keyIdStr := base64.RawURLEncoding.EncodeToString(keyIdBytes)
	//建立檔案
	filePath := filepath.Join(auth.KeyFilesPath, keyIdStr+".key")
	file, err := os.Create(filePath)
	if err != nil {
		return key, err
	}
	defer file.Close()

	currentTime := time.Now()
	key.id = keyIdStr
	key.iat = currentTime
	key.private = privateKey
	key.public = publicKey

	//寫入KEY生成時間
	timestamp := currentTime.Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	_, err = file.WriteString(timestampStr + "\n" + privateKeyStr + "\n" + publicKeyStr)
	if err != nil {
		return key, err
	}

	return key, nil
}

func (auth *Auth) readKeys(keysId []string) ([]Key, error) {
	auth.muxKeys.RLock()
	defer auth.muxKeys.RUnlock()
	keys := []Key{}
	for _, keyId := range keysId {
		filePath := filepath.Join(auth.KeyFilesPath, keyId+".key")
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		fileStr := string(fileBytes)
		lines := strings.Split(fileStr, "\n")
		if len(lines) != 3 {
			return nil, errors.New("invalid key file")
		}

		iatStr := lines[0]
		iatInt, err := strconv.ParseInt(iatStr, 10, 64)
		if err != nil {
			return nil, err
		}

		iat := time.Unix(iatInt, 0)

		privateStr := lines[1]
		privateKeyBytes, err := base64.StdEncoding.DecodeString(privateStr)
		if err != nil {
			return nil, err
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
		if err != nil {
			return nil, err
		}

		publicStr := lines[2]
		publicKeyBytes, err := base64.StdEncoding.DecodeString(publicStr)
		if err != nil {
			return nil, err
		}
		publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
		if err != nil {
			return nil, err
		}

		key := Key{
			id:      keyId,
			iat:     iat,
			private: privateKey,
			public:  publicKey,
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (auth *Auth) rsaRandomKeyId() string {
	keys := make([]string, 0, len(auth.keys))
	for key := range auth.keys {
		keys = append(keys, key)
	}

	// 隨機選擇一個鍵
	keyId := keys[MathRand.Intn(len(keys))]
	return keyId
}

func (auth *Auth) rsaEncode(keyId string, data []byte) ([]byte, error) {
	auth.muxKeys.RLock()
	defer auth.muxKeys.RUnlock()
	key := auth.keys[keyId]
	return rsa.EncryptPKCS1v15(rand.Reader, key.public, data)
}

func (auth *Auth) RsaDecode(keyId string, data []byte) (*[]byte, int, error) {
	auth.muxKeys.RLock()
	defer auth.muxKeys.RUnlock()
	key, exist := auth.keys[keyId]
	if !exist {
		key, exist = auth.expiredKeys[keyId]
		if !exist {
			return nil, Error.AUTHORIZATION_INVALID, nil
		}
	}
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, key.private, data)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, err
	}

	return &decryptedData, Error.NULL, nil
}
