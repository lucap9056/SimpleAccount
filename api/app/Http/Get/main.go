package simple_account_http_get

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
		result, errCode, err = User(context)
	case "email":
		errCode, err = ChangeEmailVerify(context)
	case "register":
		errCode, err = RegisterVerify(context)
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
