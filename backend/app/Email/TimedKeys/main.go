package TimedKeys

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"simple_account/app/AccountStruct"
	"strconv"
	"sync"
	"time"
)

type TimedKeys struct {
	mux         sync.Mutex
	Map         map[string]*AccountStruct.User
	ExpiredTime time.Duration
}

func New(duration time.Duration) *TimedKeys {
	verification := &TimedKeys{
		Map:         make(map[string]*AccountStruct.User),
		ExpiredTime: duration * time.Minute,
	}

	return verification
}

func (timedKeys *TimedKeys) ClearExpired() {
	timedKeys.mux.Lock()
	defer timedKeys.mux.Unlock()

	currentTimestamp := time.Now()

	for key, user := range timedKeys.Map {
		if user.CreateTime.Add(timedKeys.ExpiredTime).Before(currentTimestamp) {
			delete(timedKeys.Map, key)
		}
	}
}

func (timedKeys *TimedKeys) Add(pre string, k string, user AccountStruct.User) error {
	timedKeys.mux.Lock()
	defer timedKeys.mux.Unlock()
	key := pre + "-" + k
	timedKeys.Map[key] = &user
	return nil
}

func (timedKeys *TimedKeys) Del(pre string, k string) {

	timedKeys.mux.Lock()
	defer timedKeys.mux.Unlock()

	key := pre + "-" + k

	_, ok := timedKeys.Map[key]
	if ok {
		delete(timedKeys.Map, key)
	}
}

func (timedKeys *TimedKeys) Verify(pre string, k string) *AccountStruct.User {

	timedKeys.mux.Lock()
	defer timedKeys.mux.Unlock()

	key := pre + "-" + k

	user, ok := timedKeys.Map[key]
	if ok {
		return user
	}

	for key, u := range timedKeys.Map {
		if user.Username == u.Username || user.Email == u.Email {
			delete(timedKeys.Map, key)
		}
	}

	return nil
}

func GenerateKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	key := hex.EncodeToString(bytes)
	return key, nil
}

func GenerateShortKey() (string, error) {
	buf := make([]byte, 8)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	randomNumber := binary.BigEndian.Uint64(buf)

	num := strconv.FormatUint(randomNumber, 10)[:6]
	return num, nil
}
