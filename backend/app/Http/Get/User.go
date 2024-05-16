package simple_account_http_get

import (
	"encoding/json"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func User(context *Message.Context) (string, int, error) {
	if context.Author.User == nil {
		return "", Error.NOT_LOGGED_IN, nil
	}

	author := context.Author
	if author.User != nil {

		bytes, err := json.Marshal(author.User)
		if err != nil {
			return "", Error.SYSTEM, nil
		}
		return string(bytes), Error.NULL, nil
	}

	return "", Error.NOT_LOGGED_IN, nil
}
