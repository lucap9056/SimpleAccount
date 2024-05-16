package simple_account_http_put

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
	"simple_account/app/Http/Url"
	"simple_account/app/Logger"
)

func Handler(context *Message.Context) {
	url := Url.New(context.Request.URL)

	var result string
	var errCode int
	var err error

	switch url.Shift() {
	case "user":
		switch url.Shift() {
		case "username":
			errCode, err = Username(context)
		case "email":
			errCode, err = Email(context)
		case "password":
			errCode, err = Password(context)
		}
	default:
		errCode = Error.SYSTEM_TEST
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
	writer := context.Writer
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if context.Author != nil {
		context.Author.UpdateToken()
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/json")
	writer.Write(responseBytes)
}
