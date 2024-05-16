package simple_account_http_get

import (
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func ChangeEmailVerify(context *Message.Context, verifyKey string) (int, error) {
	if context.Author.User == nil {
		return Error.NOT_LOGGED_IN, nil
	}

	db := context.Database
	userData := context.Email.TimedKeys.Verify("change_email", verifyKey)
	if userData == nil {
		return Error.EMAIL_VERIFY_NOT_EXIST, nil
	}

	connect := db.Connect()
	query := "UPDATE User SET email=? WHERE id=?"
	_, err := connect.Exec(query, userData.Email, userData.Id)
	if err != nil {
		return Error.SYSTEM, err
	}

	errCode, err := context.Author.GenerateToken(userData, context.Auth)
	if errCode != Error.NULL {
		return errCode, err
	}

	context.Email.TimedKeys.Del("change_email", verifyKey)

	return Error.NULL, nil
}
