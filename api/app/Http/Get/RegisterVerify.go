package simple_account_http_get

import (
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func RegisterVerify(context *Message.Context) (int, error) {
	db := context.Database
	verifyKey := context.Url.Shift()
	userData := context.Email.TimedKeys.Verify("register", verifyKey)
	if userData == nil {
		return Error.REGISTER_VERIFY_USER_NOT_EXIST, nil
	}

	errCode, err := db.UserExist(userData.Username, userData.Email, "")
	if err != nil {
		return Error.SYSTEM, err
	}
	if errCode != Error.NULL {
		return errCode, nil
	}

	connect := db.Connect()
	query := "INSERT INTO User(username, email, salt, hash) VALUE(?, ?, ?, ?)"
	stmt, err := connect.Prepare(query)
	if err != nil {
		return Error.SYSTEM, err
	}
	_, err = stmt.Exec(userData.Username, userData.Email, userData.Salt, userData.Hash)
	if err != nil {
		return Error.SYSTEM, err
	}
	context.Email.TimedKeys.Del("register", verifyKey)

	return Error.NULL, nil
}
