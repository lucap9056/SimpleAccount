package simple_account_http_get

import (
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func ChangeEmailVerify(context *Message.Context) (int, error) {
	db := context.Database
	verifyKey := context.Url.Shift()
	user := context.Email.TimedKeys.Verify("change_email", verifyKey)
	if user == nil {
		return Error.EMAIL_VERIFY_NOT_EXIST, nil
	}

	connect := db.Connect()
	query := "UPDATE User SET email=?,last_edit=CURRENT_TIMESTAMP() WHERE id=?"
	_, err := connect.Exec(query, user.Email, user.Id)
	if err != nil {
		return Error.SYSTEM, err
	}
	context.Auth.Cache.ClearUser(user.Id)

	newUser, errCode, err := context.Database.GetUser(user.Id)
	if errCode != Error.NULL {
		return Error.SYSTEM, err
	}

	errCode, err = context.Author.GenerateToken(&newUser, context.Auth)
	if errCode != Error.NULL {
		return errCode, err
	}

	context.Email.TimedKeys.Del("change_email", verifyKey)

	return Error.NULL, nil
}
