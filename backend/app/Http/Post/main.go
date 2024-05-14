package simple_account_http_post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
	"simple_account/app/Http/Url"
	"simple_account/app/Logger"
)

func Handler(auth *Auths.Auth, db *Database.API, writer http.ResponseWriter, req *http.Request, emailSender *Email.Sender) {
	url := Url.New(req.URL)

	var result string
	var errCode int
	var err error
	switch url.Shift() {
	case "login":
		errCode, err = Login(auth, db, writer, req)
	case "register":
		errCode, err = Register(db, req, emailSender)
	default:
		errCode = Error.CLIENT_INVALID_REQUEST
	}

	if err != nil {
		Logger.Error(err)
		fmt.Println(err)
	}

	response := Message.Response{
		Success: errCode == Error.NULL,
		Result:  result,
		Error:   errCode,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/json")
	writer.Write(responseBytes)
}

type UserData struct {
	Id       int
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     string
	Hash     string
}
