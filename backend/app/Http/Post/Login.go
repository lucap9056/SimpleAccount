package simple_account_http_post

import (
	"encoding/json"
	"io"
	"net/http"
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Error"
	"time"
)

func Login(auth *Auths.Auth, db *Database.API, writer http.ResponseWriter, req *http.Request) (int, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Error.CLIENT_INVALID_REQUEST, err
	}

	var user UserData
	err = json.Unmarshal(body, &user)
	if err != nil {
		return Error.CLIENT_INVALID_REQUEST, err
	}

	if user.Username == "" {
		return Error.USERNAME_IS_EMPTY, nil
	}

	if user.Password == "" {
		return Error.PASSWORD_IS_EMPTY, nil
	}

	connect := db.Connect()
	query := "SELECT id,username,email,salt,hash FROM User WHERE username=?"
	rows, err := connect.Query(query, user.Username)
	if err != nil {
		return Error.SYSTEM, err
	}

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Salt, &user.Hash)
		if err != nil {
			return Error.SYSTEM, err
		}
	} else {
		return Error.LOGIN_USER_NOT_EXIST, nil
	}

	hash := Auths.Hash(user.Salt, user.Password)
	if hash != user.Hash {
		return Error.LOGIN_PASSWORD_NOT_MATCH, nil
	}

	expiresTime := time.Now().Add(auth.ValidityDuration)

	playload := Auths.Playload{
		User: AccountStruct.User{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		},
		Iat: expiresTime.Unix(),
	}
	secret := Auths.Salt()
	token, err := auth.GenerateToken(playload, secret)
	if err != nil {
		return Error.SYSTEM, nil
	}

	cookie := http.Cookie{
		Name:     "secret",
		Value:    secret,
		HttpOnly: true,
		Expires:  expiresTime,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(writer, &cookie)
	writer.Header().Add("Authorization", token)

	return Error.NULL, nil
}
