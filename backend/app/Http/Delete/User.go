package simple_account_http_delete

import (
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func User(context *Message.Context) (string, int, error) {
	return "", Error.NULL, nil
}
