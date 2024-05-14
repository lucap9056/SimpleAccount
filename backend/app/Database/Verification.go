package Database

import (
	"crypto/rand"
	"encoding/hex"
	"simple_account/app/AccountStruct"
	"sync"
	"time"
)

type WaitingVerification struct {
	Map         map[string]*AccountStruct.User
	mux         sync.Mutex
	ticker      *time.Ticker
	ExpiredTime time.Duration
}

func (w *WaitingVerification) Init() {
	w.Map = make(map[string]*AccountStruct.User)
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for range ticker.C {
			w.ClearExpired()
		}
	}()

	w.ticker = ticker
}

func (w *WaitingVerification) ClearExpired() {
	w.mux.Lock()
	defer w.mux.Unlock()

	currentTimestamp := time.Now()

	for key, user := range w.Map {
		if user.CreateTime.Add(w.ExpiredTime).Before(currentTimestamp) {
			delete(w.Map, key)
		}
	}
}

func (w *WaitingVerification) GenerateKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	key := hex.EncodeToString(bytes)
	return key, nil
}

func (w *WaitingVerification) Add(key string, user AccountStruct.User) error {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.Map[key] = &user
	return nil
}

func (w *WaitingVerification) Del(key string) {

	w.mux.Lock()
	defer w.mux.Unlock()

	_, ok := w.Map[key]
	if ok {
		delete(w.Map, key)
	}
}

func (w *WaitingVerification) Verify(key string) *AccountStruct.User {

	w.mux.Lock()
	defer w.mux.Unlock()

	user, ok := w.Map[key]
	if ok {
		return user
	}

	for key, u := range w.Map {
		if user.Username == u.Username || user.Email == u.Email {
			delete(w.Map, key)
		}
	}

	return nil
}
