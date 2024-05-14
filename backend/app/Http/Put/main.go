package simple_account_http_put

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple_account/app/AccountStruct"
	"simple_account/app/Database"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
	"simple_account/app/Http/Url"
	"simple_account/app/Logger"
)

func Handler(author *AccountStruct.User, db *Database.API, writer http.ResponseWriter, req *http.Request) {
	url := Url.New(req.URL)

	var result string
	var errCode int
	var err error

	switch url.Shift() {
	case "user":
		result, errCode, err = User(author, db, writer)
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
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/json")
	writer.Write(responseBytes)
}
