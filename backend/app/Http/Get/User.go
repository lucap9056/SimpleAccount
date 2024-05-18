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

	user, errCode, err := context.Database.GetUser(author.User.Id)
	if errCode != Error.NULL {
		return "", errCode, err
	}

	if author.User.LastEditTime.Unix() != user.LastEditTime.Unix() {
		author.InvaildToken()
		return "", Error.USER_DATA_EDITED, nil
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		return "", Error.SYSTEM, nil
	}
	return string(bytes), Error.NULL, nil
}
