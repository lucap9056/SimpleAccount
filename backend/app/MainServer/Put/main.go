package simple_account_http_put

import (
	"encoding/json"
	"net/http"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func Handler(context *Message.Context) {

	var result string
	var errCode int
	var err error

	switch context.Url.Shift() {
	case "user":
		switch context.Url.Shift() {
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
		context.Logs.Error.Write(err)
	}

	response := Message.Response{
		Success: errCode == Error.NULL,
		Result:  result,
		Error:   errCode,
	}

	if context.Author != nil {
		context.Author.UpdateToken()
	}

	responseBytes, _ := json.Marshal(response)
	writer := context.Writer

	if response.Success {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(responseBytes)
}
