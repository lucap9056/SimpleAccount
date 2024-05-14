package simple_account_http_delete

import (
	"net/http"
	"simple_account/app/AccountStruct"
	"simple_account/app/Database"
	"simple_account/app/Error"
)

func User(author *AccountStruct.User, db *Database.API, writer http.ResponseWriter) (string, int, error) {
	return "", Error.NULL, nil
}
