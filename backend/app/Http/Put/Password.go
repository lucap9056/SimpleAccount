package simple_account_http_put

import (
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func Password(context *Message.Context) (int, error) {
	if context.Author.User == nil {
		return Error.NOT_LOGGED_IN, nil
	}
	author := context.Author.User

	requestUser, errCode, err := context.ResquestBody()
	if errCode != Error.NULL {
		return errCode, err
	}

	if requestUser.Password == "" {
		return Error.PASSWORD_IS_EMPTY, nil
	}

	if requestUser.NewPassword == "" {
		return Error.NEW_PASSWORD_IS_EMPTY, nil
	}
	oldPassword := requestUser.Password
	newPassword := requestUser.NewPassword

	connect := context.Database.Connect()
	query := "SELECT id,salt,hash FROM User WHERE id=?"
	rows, err := connect.Query(query, author.Id)
	if err != nil {
		return Error.SYSTEM, err
	}

	var user AccountStruct.User
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Salt, &user.Hash)
		if err != nil {
			return Error.SYSTEM, err
		}
	} else {
		return Error.USER_NOT_EXIST, nil
	}

	oldHash := Auths.Hash(user.Salt, oldPassword)
	if oldHash != user.Hash {
		return Error.PASSWORD_NOT_MATCH, nil
	}

	newSalt := Auths.Salt()
	newHash := Auths.Hash(newSalt, newPassword)

	updateQuery := "UPDATE User SET salt=?,hash=?,last_edit=CURRENT_TIMESTAMP() WHERE id=?"
	_, err = connect.Exec(updateQuery, newSalt, newHash, user.Id)
	if err != nil {
		return Error.SYSTEM, nil
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

	return Error.NULL, nil
}
