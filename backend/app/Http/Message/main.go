package Message

import (
	"encoding/json"
	"io"
	"net/http"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	"simple_account/app/Error"
	"simple_account/app/Http/Author"
)

type Response struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
	Error   int    `json:"error"`
}

type Context struct {
	Author   *Author.Author
	Auth     *Auths.Auth
	Database *Database.API
	Writer   http.ResponseWriter
	Request  *http.Request
	Email    Email.Manager
}

func (ctx *Context) ResquestBody() (*UserData, int, error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, Error.CLIENT_INVALID_REQUEST, err
	}

	user := &UserData{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, Error.CLIENT_INVALID_REQUEST, err
	}
	return user, Error.NULL, nil
}

type UserData struct {
	Id               int
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Salt             string
	Hash             string
	VerificationCode string `json:"code"`
	NewPassword      string `json:"new_password"`
}
