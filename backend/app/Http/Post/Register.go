package simple_account_http_post

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	"simple_account/app/Error"
	"time"
)

func Register(db *Database.API, req *http.Request, email *Email.Sender) (int, error) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Error.CLIENT_INVALID_REQUEST, err
	}

	var userData UserData
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return Error.CLIENT_INVALID_REQUEST, err
	}

	if userData.Username == "" {
		return Error.USERNAME_IS_EMPTY, nil
	}
	usernameRegexp := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegexp.MatchString(userData.Username) {
		return Error.REGISTER_INVALID_USERNAME, nil
	}
	if len(userData.Username) < 4 {
		return Error.REGISTER_USERNAME_TOO_SHORT, nil
	}
	if len(userData.Username) > 24 {
		return Error.REGISTER_USERNAME_TOO_LENGTH, nil
	}

	if userData.Email == "" {
		return Error.EMAIL_IS_EMPTY, nil
	}
	emailRegexp := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegexp.MatchString(userData.Email) {
		return Error.REGISTER_EMAIL_EXISTED, nil
	}

	if userData.Password == "" {
		return Error.PASSWORD_IS_EMPTY, nil
	}
	passwordRegexp := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*_\-]+$`)
	if !passwordRegexp.MatchString(userData.Password) {
		return Error.REGISTER_INVALID_PASSWORD, nil
	}
	if len(userData.Password) < 8 {
		return Error.REGISTER_PASSWORD_TOO_SHORT, nil
	}

	errCode, err := db.UserExist(userData.Username, userData.Email)
	if err != nil {
		return Error.SYSTEM, err
	}
	if errCode != Error.NULL {
		return errCode, nil
	}

	userData.Salt = Auths.Salt()
	userData.Hash = Auths.Hash(userData.Salt, userData.Password)

	currentTime := time.Now()
	user := AccountStruct.User{
		Username:   userData.Username,
		Email:      userData.Email,
		Salt:       userData.Salt,
		Hash:       userData.Hash,
		CreateTime: &currentTime,
	}

	key, err := db.Verification.GenerateKey()
	if err != nil {
		return Error.SYSTEM, err
	}

	err = email.SendVerify(user.Email, key)
	if err != nil {
		return Error.SYSTEM, err
	}

	db.Verification.Add(key, user)

	return Error.NULL, nil
}
