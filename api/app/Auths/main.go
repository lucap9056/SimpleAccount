package Auths

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"simple_account/app/Logger"
	"sync"
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
	KeyFilesPath          string
	keys                  map[string]Key
	expiredKeys           map[string]Key
	muxKeys               sync.RWMutex
	Cache                 AuthCache
	KeyValidityDuration   time.Duration
	ValidityDuration      time.Duration
	RenewTime             time.Duration
	CacheValidityDuration time.Duration
	Total                 int
	logger                *Logger.Manager
}

func (auth *Auth) Init(logger *Logger.Manager) error {
	auth.logger = logger

	auth.Cache.Init(auth.CacheValidityDuration)
	//確認 RSA256 公私鑰是否都存在
	keysId, err := auth.rsaKeysExist()
	if err != nil {
		return err
	}

	keys, err := auth.readKeys(keysId)
	if err != nil {
		return err
	}

	auth.keys = make(map[string]Key)
	auth.expiredKeys = make(map[string]Key)

	currentTime := time.Now()
	for _, key := range keys {
		if currentTime.Before(key.iat) {
			auth.expiredKeys[key.id] = key
		} else {
			auth.keys[key.id] = key
		}
	}

	for len(auth.keys) < auth.Total {
		key, err := auth.generateKey()
		if err != nil {
			return err
		}
		auth.keys[key.id] = key
	}

	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			auth.removeExpiredKeys()
		}
	}()

	return nil
}
