package Auths

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type AuthCache struct {
	Map              map[string]*Signature
	mux              sync.Mutex
	ticker           *time.Ticker
	ValidityDuration time.Duration
}

type CacheToken struct {
	Secret string
	Key    string
	Iat    int64
}

func (cache *AuthCache) Init(validityDuration time.Duration) {
	cache.ValidityDuration = validityDuration
	cache.Map = make(map[string]*Signature)
	ticker := time.NewTicker(time.Hour)

	go func() {
		for range ticker.C {
			cache.ClearExpired()
		}
	}()

	cache.ticker = ticker
}

func (cache *AuthCache) ClearExpired() {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	currentTimestamp := time.Now()
	for key, sign := range cache.Map {
		if sign.CreateTime.Add(cache.ValidityDuration).Before(currentTimestamp) {
			delete(cache.Map, key)
		}
	}
}

func (cache *AuthCache) GenerateKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	key := hex.EncodeToString(bytes)
	return key, nil
}

func (cache *AuthCache) Add(key string, sign Signature) error {

	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.Map[key] = &sign

	return nil
}

func (cache *AuthCache) Del(key string) {

	cache.mux.Lock()
	defer cache.mux.Unlock()

	_, ok := cache.Map[key]
	if ok {
		delete(cache.Map, key)
	}
}

func (cache *AuthCache) Verify(key string) *Signature {

	cache.mux.Lock()
	defer cache.mux.Unlock()

	sign, ok := cache.Map[key]
	if ok {
		return sign
	}

	return nil
}
