package simple_account_http_get

import (
	"encoding/json"
	"net/http"
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"simple_account/app/Error"
)

func User(author *AccountStruct.User, auth *Auths.Auth, writer http.ResponseWriter, req *http.Request) (string, int, error) {
	if author != nil {

		bytes, err := json.Marshal(author)
		if err != nil {
			return "", Error.SYSTEM, nil
		}
		return string(bytes), Error.NULL, nil
	}
	return "", Error.NOT_LOGGED_IN, nil
}
