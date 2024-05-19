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

type TempUser struct {
	User       AccountStruct.User
	CreateTime time.Time
}

type TimedKeys struct {
	mux         sync.Mutex
	Map         map[string]TempUser
	ExpiredTime time.Duration
}

func New(duration time.Duration) *TimedKeys {
	verification := &TimedKeys{
		Map:         make(map[string]TempUser),
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
	timedKeys.Map[key] = TempUser{
		User:       user,
		CreateTime: time.Now(),
	}
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

	tempUser, ok := timedKeys.Map[key]
	if ok {

		return &tempUser.User
	}

	for key, u := range timedKeys.Map {
		user := tempUser.User
		if user.Username == u.User.Username || user.Email == u.User.Email {
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
